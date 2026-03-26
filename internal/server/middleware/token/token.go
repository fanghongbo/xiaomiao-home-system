package token

import (
	"context"
	"xiaomiao-home-system/internal/conf"
	"xiaomiao-home-system/third_party/jwt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

const (
	// reason holds the error reason.
	reason string = "UNAUTHORIZED"

	userKey string = "x-md-global-user"
)

var (
	ErrTokenInvalid     = errors.Unauthorized(reason, "登录失效, 请重新登录")
	ErrWrongTokenFormat = errors.Unauthorized(reason, "登录凭证无效")
	ErrRequestInvalid   = errors.BadRequest(reason, "无效请求, 请检查请求参数")
)

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(c *conf.Jwt) middleware.Middleware {
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
					return nil, ErrTokenInvalid
				}

				claims, err := jwt.ParseToken(c.SecretKey, jwt.BearerTokenFromAuthorization(token))
				if err != nil {
					return nil, ErrTokenInvalid
				}

				ctx = context.WithValue(ctx, userKey, claims.UserId)

				return handler(ctx, req)
			}

			return nil, ErrRequestInvalid
		}
	}
}
