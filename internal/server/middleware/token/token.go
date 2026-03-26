package token

import (
	"context"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

const (
	// reason holds the error reason.
	reason string = "UNAUTHORIZED"

	userKey  string = "x-md-global-user"
	tokenKey string = "x-md-global-token"
)

var (
	ErrInvalidCookie      = errors.Unauthorized(reason, "登录凭证无效")
	ErrInvalidSession     = errors.Unauthorized(reason, "登录失效, 请重新登录")
	ErrWrongCookieContext = errors.BadRequest(reason, "登录凭证无效")
)

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(c *conf.Config) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			var (
				header transport.Header
				token  string
				tr     transport.Transporter
				ok     bool
			)

			if tr, ok = transport.FromServerContext(ctx); ok {
				header = tr.RequestHeader()
				token = header.Get("token")

				if token == "" {
					return nil, ErrInvalidSession
				}

				ctx = context.WithValue(ctx, tokenKey, token)

				return handler(ctx, req)
			}

			return nil, ErrWrongCookieContext
		}
	}
}
