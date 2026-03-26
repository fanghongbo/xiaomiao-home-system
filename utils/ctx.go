package utils

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// GetCurrentUserId 查询当前用户Id
func GetCurrentUserId(ctx context.Context) (int64, error) {
	userId := ctx.Value("x-md-global-user")
	if userId == nil {
		return 0, errors.New("user id not found in ctx")
	}

	id, ok := userId.(int64)
	if !ok {
		return 0, errors.New("user id is not int64")
	}

	return id, nil
}

// GetUserIP 获取用户IP
func GetUserIP(ctx context.Context) (string, error) {
	request, isOk := http.RequestFromServerContext(ctx)
	if isOk {
		xForwardedFor := request.Header.Get("X-FORWARDED-FOR")
		ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
		if ip != "" {
			return ip, nil
		}

		ip = strings.TrimSpace(request.Header.Get("X-Real-Ip"))
		if ip != "" {
			return ip, nil
		}

		if ip, _, err := net.SplitHostPort(strings.TrimSpace(request.RemoteAddr)); err == nil {
			return ip, nil
		}
	}

	return "", fmt.Errorf("failed to get user ip from ctx")
}
