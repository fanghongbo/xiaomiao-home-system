package server

import (
	collectV1 "xiaomiao-home-system/api/collect/v1"
	fileV1 "xiaomiao-home-system/api/file/v1"
	publishV1 "xiaomiao-home-system/api/publish/v1"
	roleV1 "xiaomiao-home-system/api/role/v1"
	userNotificationV1 "xiaomiao-home-system/api/user/notification/v1"
	userSettingV1 "xiaomiao-home-system/api/user/setting/v1"
	userV1 "xiaomiao-home-system/api/user/v1"
	"xiaomiao-home-system/internal/conf"
	"xiaomiao-home-system/internal/service"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, user *service.UserService, role *service.RoleService, userNotification *service.UserNotificationService, userSetting *service.UserSettingService, file *service.FileService, publish *service.PublishService, collect *service.CollectService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			validate.ProtoValidate(),
			logging.Server(logger),
			ratelimit.Server(),
		),
	}

	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}

	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}

	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	userV1.RegisterUserServer(srv, user)
	roleV1.RegisterRoleServer(srv, role)
	userNotificationV1.RegisterUserNotificationServer(srv, userNotification)
	userSettingV1.RegisterUserSettingServer(srv, userSetting)
	fileV1.RegisterFileServer(srv, file)
	publishV1.RegisterPublishServer(srv, publish)
	collectV1.RegisterCollectServer(srv, collect)
	return srv
}
