package service

import (
	"context"
	stdhttp "net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	pb "xiaomiao-home-system/api/file/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type FileService struct {
	pb.UnimplementedFileServer

	file   *biz.FileUsecase
	log    *log.Helper
	config *conf.Config
	static *conf.Static
}

func NewFileService(file *biz.FileUsecase, config *conf.Config, static *conf.Static, logger log.Logger) *FileService {
	return &FileService{
		file:   file,
		config: config,
		static: static,
		log:    log.NewHelper(log.With(logger, "service", "FileService")),
	}
}

// DownloadFile 下载文件
func (s *FileService) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileReply, error) {
	return s.file.DownloadFile(ctx, req)
}

// UpdateAvatar 上传头像
func (s *FileService) UploadAvatar(ctx context.Context, req *pb.UploadAvatarRequest) (*pb.UploadAvatarReply, error) {
	return s.file.UploadAvatar(ctx, req)
}

// StaticFile 获取静态文件
func (s *FileService) StaticFileHandler(httpCtx http.Context) error {
	http.SetOperation(httpCtx, pb.OperationFileGetStaticFile)
	logPath := ""
	if req, ok := http.RequestFromServerContext(httpCtx); ok && req != nil {
		logPath = req.Method + " " + req.URL.RequestURI()
	}
	h := httpCtx.Middleware(func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, s.serveStaticFile(httpCtx, s.static)
	})
	_, err := h(httpCtx, logPath)
	return err
}

func (s *FileService) serveStaticFile(httpCtx http.Context, staticConfig *conf.Static) error {
	req, ok := http.RequestFromServerContext(httpCtx)
	if !ok || req == nil {
		return errors.BadRequest("ERR_INVALID_REQUEST", "无效请求")
	}

	relPath := strings.TrimSpace(strings.TrimPrefix(req.URL.Path, "/static/"))
	if relPath == "" {
		return errors.BadRequest("ERR_INVALID_REQUEST", "无效文件路径")
	}

	// 替换掉目录中的所有 ".." 字符后，再做规范化处理
	relPath = strings.ReplaceAll(relPath, "..", "")
	cleanRelPath := strings.TrimPrefix(path.Clean("/"+relPath), "/")
	if cleanRelPath == "." || cleanRelPath == "" {
		return errors.BadRequest("ERR_INVALID_REQUEST", "无效文件路径")
	}

	// 不允许请求隐藏文件/目录
	for _, seg := range strings.Split(cleanRelPath, "/") {
		if seg == "" {
			continue
		}
		if strings.HasPrefix(seg, ".") {
			return errors.BadRequest("ERR_INVALID_REQUEST", "无效文件路径")
		}
	}

	baseDirAbs, err := filepath.Abs(staticConfig.BaseDir)
	if err != nil {
		return errors.InternalServer("ERR_SYSTEM_EXCEPTION", "系统错误, 请稍后再试")
	}
	fullPath := filepath.Join(baseDirAbs, filepath.FromSlash(cleanRelPath))
	fullPathAbs, err := filepath.Abs(fullPath)
	if err != nil {
		return errors.InternalServer("ERR_SYSTEM_EXCEPTION", "系统错误, 请稍后再试")
	}

	// 最终路径前缀必须是 BaseDir
	basePrefix := strings.TrimRight(baseDirAbs, string(os.PathSeparator)) + string(os.PathSeparator)
	if fullPathAbs != baseDirAbs && !strings.HasPrefix(fullPathAbs, basePrefix) {
		return errors.BadRequest("ERR_INVALID_REQUEST", "无效文件路径")
	}

	// 禁止路径上任意层级使用符号链接
	hasSymlink, err := s.hasSymlinkInPath(baseDirAbs, fullPathAbs)
	if err != nil {
		return errors.InternalServer("ERR_SYSTEM_EXCEPTION", "系统错误, 请稍后再试")
	}
	if hasSymlink {
		return errors.BadRequest("ERR_INVALID_REQUEST", "无效文件路径")
	}

	fullPath = fullPathAbs
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return errors.NotFound("ERR_BAD_REQUEST", "文件不存在")
		}
		return errors.InternalServer("ERR_SYSTEM_EXCEPTION", "系统错误, 请稍后再试")
	}
	// 直接让 net/http 处理 Content-Type、Range 等
	stdhttp.ServeFile(httpCtx.Response(), req, fullPath)
	return nil
}

func (s *FileService) hasSymlinkInPath(baseDirAbs string, fullPathAbs string) (bool, error) {
	rel, err := filepath.Rel(baseDirAbs, fullPathAbs)
	if err != nil {
		return false, err
	}
	if rel == "." {
		return false, nil
	}

	cur := baseDirAbs
	for _, seg := range strings.Split(rel, string(os.PathSeparator)) {
		if seg == "" || seg == "." {
			continue
		}
		cur = filepath.Join(cur, seg)
		info, err := os.Lstat(cur)
		if err != nil {
			// 由后续 os.Stat 统一处理不存在等情况
			if os.IsNotExist(err) {
				return false, nil
			}
			return false, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return true, nil
		}
	}
	return false, nil
}

func (s *FileService) UploadAvatarHandler(httpCtx http.Context) error {
	http.SetOperation(httpCtx, pb.OperationFileUploadAvatar)
	h := httpCtx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
		var in pb.UploadAvatarRequest
		if req, ok := http.RequestFromServerContext(httpCtx); ok && req != nil {
			// multipart 由业务层自行解析，这里跳过默认 body bind（默认 codec 不支持 multipart）
			if !strings.HasPrefix(strings.ToLower(req.Header.Get("Content-Type")), "multipart/form-data") {
				if err := httpCtx.Bind(&in); err != nil {
					return nil, err
				}
			}
		} else {
			if err := httpCtx.Bind(&in); err != nil {
				return nil, err
			}
		}
		if err := httpCtx.BindQuery(&in); err != nil {
			return nil, err
		}
		return s.UploadAvatar(ctx, &in)
	})
	avatarLogArg := "POST /api/v1/user/avatar/upload"
	if req, ok := http.RequestFromServerContext(httpCtx); ok && req != nil {
		avatarLogArg = req.Method + " " + req.URL.RequestURI()
	}
	out, err := h(httpCtx, avatarLogArg)
	if err != nil {
		return err
	}
	reply := out.(*pb.UploadAvatarReply)
	return httpCtx.Result(200, reply)
}
