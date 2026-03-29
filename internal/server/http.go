package server

import (
	"context"
	publishV1 "xiaomiao-home-system/api/publish/v1"
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
	whiteList := map[string]bool{
		"/api.user.v1.User/GetSecretKey":  true,
		"/api.user.v1.User/WebLogin":      true,
		"/api.user.v1.User/AppLogin":      true,
		"/api.user.v1.User/MpLogin":       true,
		"/api.file.v1.File/GetStaticFile": true,
	}

	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, jwtConfig *conf.Jwt, staticConfig *conf.Static, user *service.UserService, role *service.RoleService, userNotification *service.UserNotificationService, userSetting *service.UserSettingService, file *service.FileService, publish *service.PublishService, logger log.Logger) *http.Server {
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

	route.GET("/static/{path:.*}", file.StaticFileHandler)
	route.POST("/api/v1/user/avatar/upload", file.UploadAvatarHandler)

	userV1.RegisterUserHTTPServer(srv, user)
	roleV1.RegisterRoleHTTPServer(srv, role)
	userNotificationV1.RegisterUserNotificationHTTPServer(srv, userNotification)
	userSettingV1.RegisterUserSettingHTTPServer(srv, userSetting)
	publishV1.RegisterPublishHTTPServer(srv, publish)
	return srv
}
