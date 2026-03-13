package utils

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// GetCurrentUser 查询当前用户
func GetCurrentUser(ctx context.Context) (string, error) {
	user := ctx.Value("x-md-global-user")
	if user == nil {
		return "", fmt.Errorf("failed to get current user from ctx")
	}

	return fmt.Sprintf("%v", user), nil
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
