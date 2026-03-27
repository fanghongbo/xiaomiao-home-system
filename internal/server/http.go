package server

import (
	"context"
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
func NewHTTPServer(c *conf.Server, config *conf.Jwt, user *service.UserService, role *service.RoleService, userNotification *service.UserNotificationService, userSetting *service.UserSettingService, file *service.FileService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.ProtoValidate(),
			logging.Server(logger),
			ratelimit.Server(),
			selector.Server(
				token.Server(config),
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

	userV1.RegisterUserHTTPServer(srv, user)
	roleV1.RegisterRoleHTTPServer(srv, role)
	userNotificationV1.RegisterUserNotificationHTTPServer(srv, userNotification)
	userSettingV1.RegisterUserSettingHTTPServer(srv, userSetting)
	fileV1.RegisterFileHTTPServer(srv, file)
	return srv
}
