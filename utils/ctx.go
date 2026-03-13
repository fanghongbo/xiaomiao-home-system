package utils

import (
	"context"
	"fmt"
)

// GetCurrentUser 查询当前用户
func GetCurrentUser(ctx context.Context) (string, error) {
	user := ctx.Value("x-md-global-user")
	if user == nil {
		return "", fmt.Errorf("failed to get current user from ctx")
	}

	return fmt.Sprintf("%v", user), nil
}
