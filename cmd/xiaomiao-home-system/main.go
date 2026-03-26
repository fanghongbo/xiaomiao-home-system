package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"
	"xiaomiao-home-system/internal/conf"
	"xiaomiao-home-system/internal/task"

	register "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "xiaomiao-home-system"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// grpcRegistrar gRPC服务的注册器包装
type grpcRegistrar struct {
	registry.Registrar
}

// Register 注册服务
func (r *grpcRegistrar) Register(ctx context.Context, service *registry.ServiceInstance) error {
	var grpcEndpoints []string
	for _, endpoint := range service.Endpoints {
		if strings.HasPrefix(endpoint, "grpc://") || strings.HasPrefix(endpoint, "grpc+") {
			grpcEndpoints = append(grpcEndpoints, endpoint)
		}
	}

	if len(grpcEndpoints) == 0 {
		return nil
	}

	grpcService := &registry.ServiceInstance{
		ID:        service.ID,
		Name:      Name,
		Version:   service.Version,
		Metadata:  service.Metadata,
		Endpoints: grpcEndpoints,
	}

	return r.Registrar.Register(ctx, grpcService)
}

// Deregister 注销服务
func (r *grpcRegistrar) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return r.Registrar.Deregister(ctx, service)
}

func newApp(logger log.Logger, task *task.TaskManager, gs *grpc.Server, hs *http.Server, r *register.Registry) *kratos.App {
	if err := task.Start(); err != nil {
		panic(err)
	}

	// 创建只注册 gRPC 服务的注册器包装
	// gr := &grpcRegistrar{Registrar: r}

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{
			"id":      id,
			"name":    Name,
			"version": Version,
			"date":    time.Now().Format("2006-01-02 15:04:05"),
		}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		// kratos.Registrar(r),
		// kratos.Registrar(gr),
	)
}

func main() {
	flag.Parse()

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	// 方案1: 从本地配置文件读取配置
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)

	r := &register.Registry{}

	// 方案2: 从环境变量中读取nacos配置
	// rc, err := utils.GetNacosConfigFromEnv()
	// if err != nil {
	// 	panic(err)
	// }

	// c := config.New(
	// 	config.WithSource(
	// 		nacosConf.NewConfigSource(utils.NewNacosConfigClient(rc, Name), nacosConf.WithGroup(rc.Nacos.Group), nacosConf.WithDataID(rc.Nacos.DataId)),
	// 	),
	// )

	// r := register.New(utils.NewNacosNamingClient(rc, Name))

	defer c.Close()

	var err error

	if err = c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err = c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Config, bc.Jwt, bc.Static, r, logger)
	if err != nil {
		panic(err)
	}

	defer cleanup()

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		panic(err)
	}
}
