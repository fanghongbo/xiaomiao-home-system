package server

import (
	"context"
	stdhttp "net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	fileV1 "xiaomiao-home-system/api/file/v1"
	roleV1 "xiaomiao-home-system/api/role/v1"
	userNotificationV1 "xiaomiao-home-system/api/user/notification/v1"
	userSettingV1 "xiaomiao-home-system/api/user/setting/v1"
	userV1 "xiaomiao-home-system/api/user/v1"
	"xiaomiao-home-system/internal/conf"
	"xiaomiao-home-system/internal/server/encoder"
	"xiaomiao-home-system/internal/server/middleware/token"
	"xiaomiao-home-system/internal/service"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/selector"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := map[string]interface{}{
		"/api.user.v1.User/GetSecretKey": struct{}{},
		"/api.user.v1.User/WebLogin":     struct{}{},
		"/api.user.v1.User/AppLogin":     struct{}{},
		"/api.user.v1.User/MpLogin":      struct{}{},
	}

	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, jwtConfig *conf.Jwt, staticConfig *conf.Static, user *service.UserService, role *service.RoleService, userNotification *service.UserNotificationService, userSetting *service.UserSettingService, file *service.FileService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.ProtoValidate(),
			logging.Server(logger),
			ratelimit.Server(),
			selector.Server(
				token.Server(jwtConfig),
			).
				Match(NewWhiteListMatcher()).
				Build(),
		),
		http.ErrorEncoder(encoder.DefaultHttpServerErrorEncoder),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}

	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}

	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	route := srv.Route("/")

	route.GET("/static/{path:.*}", func(ctx http.Context) error {
		req, ok := http.RequestFromServerContext(ctx)
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
		hasSymlink, err := hasSymlinkInPath(baseDirAbs, fullPathAbs)
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
		stdhttp.ServeFile(ctx.Response(), req, fullPath)
		return nil
	})

	// 自定义文件路由，避免 multipart 处理逻辑被代码生成覆盖
	route.POST("/api/v1/user/avatar/upload", func(ctx http.Context) error {
		var in fileV1.UploadAvatarRequest
		if req, ok := http.RequestFromServerContext(ctx); ok && req != nil {
			// multipart 由业务层自行解析，这里跳过默认 body bind（默认 codec 不支持 multipart）
			if !strings.HasPrefix(strings.ToLower(req.Header.Get("Content-Type")), "multipart/form-data") {
				if err := ctx.Bind(&in); err != nil {
					return err
				}
			}
		} else {
			if err := ctx.Bind(&in); err != nil {
				return err
			}
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, fileV1.OperationFileUploadAvatar)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return file.UploadAvatar(ctx, req.(*fileV1.UploadAvatarRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*fileV1.UploadAvatarReply)
		return ctx.Result(200, reply)
	})

	userV1.RegisterUserHTTPServer(srv, user)
	roleV1.RegisterRoleHTTPServer(srv, role)
	userNotificationV1.RegisterUserNotificationHTTPServer(srv, userNotification)
	userSettingV1.RegisterUserSettingHTTPServer(srv, userSetting)
	return srv
}

func hasSymlinkInPath(baseDirAbs string, fullPathAbs string) (bool, error) {
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
