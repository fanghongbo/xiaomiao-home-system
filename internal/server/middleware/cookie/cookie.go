package cookie

import (
	"context"
	"time"
	"xiaomiao-home-system/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/patrickmn/go-cache"
)

type authKey struct{}

const (
	// reason holds the error reason.
	reason string = "UNAUTHORIZED"

	userKey   string = "x-md-global-user"
	dayuKey   string = "x-md-global-dayu-url"
	cookieKey string = "x-md-global-cookie"
)

var (
	// 本地缓存
	localCache *cache.Cache

	ErrMissingCookie      = errors.Unauthorized(reason, "cookie has expired, please login again")
	ErrInvalidCookie      = errors.Unauthorized(reason, "cookie is invalid")
	ErrInvalidSession     = errors.Unauthorized(reason, "session has expired, please login again")
	ErrWrongCookieContext = errors.BadRequest(reason, "wrong cookie context for middleware")
)

func init() {
	localCache = cache.New(10*time.Minute, 30*time.Minute)
}

// Server is a server auth middleware. Check the token and extract the info from token.
func Server(c *conf.Config) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			var (
				header transport.Header
				cookie string
				tr     transport.Transporter
				ok     bool
			)

			if tr, ok = transport.FromServerContext(ctx); ok {
				header = tr.RequestHeader()
				cookie = header.Get("cookie")

				if cookie == "" {
					return nil, ErrMissingCookie
				}

				ctx = context.WithValue(ctx, cookieKey, cookie)

				// 默认从本地缓存获取
				// if res, exist = localCache.Get(cookie); exist && res != nil {
				// 	userInfo = res.(*dayu.GetLoginUserReply)
				// } else {
				// 	if userInfo, err = dayu.GetLoginUser(ctx); err != nil {
				// 		if errors.Is(err, dayu.ErrInvalidSession) {
				// 			return nil, ErrInvalidSession
				// 		}

				// 		return nil, errors.BadRequest(reason, err.Error())
				// 	}
				// }

				// if userInfo != nil && userInfo.Data != nil {
				// 	// 本地缓存用户信息
				// 	localCache.Set(cookie, userInfo, 10*time.Minute)
				// 	ctx = context.WithValue(ctx, userKey, userInfo.Data.Username)
				// } else {
				// 	return nil, ErrInvalidCookie
				// }

				return handler(ctx, req)
			}

			return nil, ErrWrongCookieContext
		}
	}
}
