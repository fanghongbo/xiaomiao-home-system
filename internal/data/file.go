package data

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	v1 "xiaomiao-home-system/api/file/v1"
	"xiaomiao-home-system/internal/biz"
	"xiaomiao-home-system/utils"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

const maxAvatarSize = 5 << 20 // 5MB
const maxAvatarUploadCount = 10
const avatarUploadCountTTL = 8 * time.Hour

type UploadImageMeta struct {
	fileName    string
	contentType string
	size        int64
}

type fileRepo struct {
	data *Data
	log  *log.Helper
}

// NewFileRepo .
func NewFileRepo(data *Data, logger log.Logger) biz.FileRepo {
	return &fileRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "FileRepo")),
	}
}

// DownloadFile 下载文件
func (u *fileRepo) DownloadFile(ctx context.Context, req *v1.DownloadFileRequest) (*v1.DownloadFileReply, error) {
	return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
}

// UpdateAvatar 上传头像
func (u *fileRepo) UploadAvatar(ctx context.Context, req *v1.UploadAvatarRequest) (*v1.UploadAvatarReply, error) {
	fileReader, fileHeader, err := u.pickUploadImage(ctx, req)
	if err != nil {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_INVALID_REQUEST.String(), "请上传头像文件")
	}
	defer fileReader.Close()
	if fileHeader.size > maxAvatarSize {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_INVALID_REQUEST.String(), "头像文件不能超过 5MB")
	}

	if !strings.HasPrefix(strings.ToLower(fileHeader.contentType), "image/") {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_INVALID_REQUEST.String(), "仅支持图片类型")
	}

	userId, err := utils.GetCurrentUserId(ctx)
	if err != nil {
		return nil, errors.Unauthorized(v1.ErrorReason_ERR_INVALID_SESSION.String(), "登录失效, 请重新登录")
	}

	if err := u.CheckAvatarUploadCountLimit(ctx, userId); err != nil {
		return nil, err
	}

	fileSuffix, err := u.resolveAvatarSuffix(fileHeader)
	if err != nil {
		return nil, errors.BadRequest(v1.ErrorReason_ERR_INVALID_REQUEST.String(), err.Error())
	}

	relativeDir := filepath.Join("avatar", strconv.FormatInt(userId, 10))
	absoluteDir := filepath.Join(u.data.static.BaseDir, relativeDir)
	if err := os.MkdirAll(absoluteDir, 0o755); err != nil {
		u.log.Errorf("mkdir upload dir failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	finalName := "user" + fileSuffix
	absolutePath := filepath.Join(absoluteDir, finalName)

	// 清理该用户目录下旧头像文件，避免残留 avatar.* 历史文件
	if err := u.removeOldAvatarFiles(absoluteDir, absolutePath); err != nil {
		u.log.Errorf("remove old avatar files failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	dst, err := os.Create(absolutePath)
	if err != nil {
		u.log.Errorf("create avatar file failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	defer dst.Close()

	written, err := io.Copy(dst, io.LimitReader(fileReader, maxAvatarSize+1))
	if err != nil {
		u.log.Errorf("save avatar file failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if written > maxAvatarSize {
		_ = os.Remove(absolutePath)
		return nil, errors.BadRequest(v1.ErrorReason_ERR_INVALID_REQUEST.String(), "头像文件不能超过 5MB")
	}

	avatarPath := strings.ReplaceAll(filepath.ToSlash(filepath.Join(relativeDir, finalName)), "\\", "/")
	if err := u.data.db.Model(&User{}).
		Where("id = ?", userId).
		Where("deleted_flag = ?", 0).
		Update("avatar", avatarPath).Error; err != nil {
		u.log.Errorf("update user avatar failed: %v", err)
		return nil, errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}

	frontendAvatarURL := strings.TrimRight(u.data.static.BaseUrl, "/") + "/static/" + avatarPath

	return &v1.UploadAvatarReply{
		Code:    200,
		Message: "操作成功",
		Success: true,
		Data: &v1.ImageInfo{
			Url: frontendAvatarURL,
		},
	}, nil
}

func (u *fileRepo) CheckAvatarUploadCountLimit(ctx context.Context, userId int64) error {
	uploadCountKey := fmt.Sprintf("user:avatar:upload:count:%d", userId)
	uploadCount, err := u.data.cache.Incr(ctx, uploadCountKey).Result()
	if err != nil {
		u.log.Errorf("increase avatar upload count failed: %v", err)
		return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
	}
	if uploadCount == 1 {
		if err := u.data.cache.Expire(ctx, uploadCountKey, avatarUploadCountTTL).Err(); err != nil {
			u.log.Errorf("set avatar upload count ttl failed: %v", err)
			return errors.InternalServer(v1.ErrorReason_ERR_SYSTEM_EXCEPTION.String(), "系统错误, 请稍后再试")
		}
	}
	if uploadCount > maxAvatarUploadCount {
		return errors.BadRequest(v1.ErrorReason_ERR_TOO_MANY_REQUEST.String(), "今日已达最大修改次数, 请明日再试")
	}
	return nil
}

func (u *fileRepo) resolveAvatarSuffix(meta *UploadImageMeta) (string, error) {
	contentType := strings.ToLower(strings.TrimSpace(meta.contentType))
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/webp":
		return ".webp", nil
	case "image/gif":
		return ".gif", nil
	}

	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(meta.fileName)))
	switch ext {
	case ".jpg", ".jpeg":
		return ".jpg", nil
	case ".png", ".webp", ".gif":
		return ext, nil
	default:
		return "", fmt.Errorf("仅支持 jpg/png/webp/gif 格式")
	}
}

func (u *fileRepo) removeOldAvatarFiles(dir string, keepPath string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if !strings.HasPrefix(name, "avatar.") {
			continue
		}
		fullPath := filepath.Join(dir, entry.Name())
		if fullPath == keepPath {
			continue
		}
		if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (u *fileRepo) pickUploadImage(ctx context.Context, req *v1.UploadAvatarRequest) (io.ReadCloser, *UploadImageMeta, error) {
	if req != nil && len(req.Content) > 0 {
		if int64(len(req.Content)) > maxAvatarSize {
			return nil, nil, errors.BadRequest(v1.ErrorReason_ERR_INVALID_REQUEST.String(), "头像文件不能超过 5MB")
		}
		fileName := strings.TrimSpace(req.FileName)
		if fileName == "" {
			fileName = "avatar.png"
		}
		contentType := strings.TrimSpace(req.ContentType)
		if contentType == "" {
			contentType = "image/*"
		}
		return io.NopCloser(bytes.NewReader(req.Content)), &UploadImageMeta{
			fileName:    fileName,
			contentType: contentType,
			size:        int64(len(req.Content)),
		}, nil
	}

	httpReq, ok := http.RequestFromServerContext(ctx)
	if !ok || httpReq == nil {
		return nil, nil, os.ErrInvalid
	}
	if err := httpReq.ParseMultipartForm(maxAvatarSize); err != nil {
		return nil, nil, err
	}
	file, fileHeader, err := u.pickAvatarFile(httpReq.MultipartForm)
	if err != nil {
		return nil, nil, err
	}
	if fileHeader.Size > maxAvatarSize {
		_ = file.Close()
		return nil, nil, errors.BadRequest(v1.ErrorReason_ERR_INVALID_REQUEST.String(), "头像文件不能超过 5MB")
	}
	return file, &UploadImageMeta{
		fileName:    fileHeader.Filename,
		contentType: fileHeader.Header.Get("Content-Type"),
		size:        fileHeader.Size,
	}, nil
}

func (u *fileRepo) pickAvatarFile(form *multipart.Form) (multipart.File, *multipart.FileHeader, error) {
	if form == nil {
		return nil, nil, os.ErrInvalid
	}
	for _, key := range []string{"file", "avatar", "image"} {
		files := form.File[key]
		if len(files) == 0 {
			continue
		}
		fd, err := files[0].Open()
		if err != nil {
			return nil, nil, err
		}
		return fd, files[0], nil
	}
	for _, files := range form.File {
		if len(files) == 0 {
			continue
		}
		fd, err := files[0].Open()
		if err != nil {
			return nil, nil, err
		}
		return fd, files[0], nil
	}
	return nil, nil, os.ErrNotExist
}
