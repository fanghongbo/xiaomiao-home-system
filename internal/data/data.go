package data

import (
	"context"
	"os"
	"time"
	"xiaomiao-home-system/internal/conf"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	snowflake "github.com/sony/sonyflake"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGID, NewDB, NewCache, NewUserRepo, NewRoleRepo, NewUserPostRepo, NewUserNotificationRepo, NewUserSettingRepo, NewFileRepo, NewUserCollectRepo, NewDiscoverRepo, NewUserCatRepo)

// Data .
type Data struct {
	db       *gorm.DB
	cache    *redis.Client
	log      *klog.Helper
	config   *conf.Config
	dbConfig *conf.Data
	static   *conf.Static
	jwt      *conf.Jwt
	gid      *snowflake.Sonyflake
}

func NewGID(logger klog.Logger) *snowflake.Sonyflake {
	log := klog.NewHelper(klog.With(logger, "module", "xiaomiao-home-system/data/snowflake"))

	// 使用主机名和进程ID的组合来生成机器ID
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// 将主机名转换为数字
	var machineId uint16
	for _, c := range hostname {
		machineId = machineId*31 + uint16(c)
	}

	// 加上进程ID
	machineId += uint16(os.Getpid())

	sf := snowflake.NewSonyflake(snowflake.Settings{
		MachineID: func() (uint16, error) {
			return machineId, nil
		},
	})

	if sf == nil {
		log.Fatalf("failed to new snowflake instance")
	}

	return sf
}

// NewDB 初始化数据库
func NewDB(conf *conf.Data, logger klog.Logger) *gorm.DB {
	log := klog.NewHelper(klog.With(logger, "module", "xiaomiao-home-system/data/gorm"))

	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{Logger: nil})
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed getting the database connection: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

// NewCache 初始化redis
func NewCache(conf *conf.Data, logger klog.Logger) *redis.Client {
	log := klog.NewHelper(klog.With(logger, "module", "xiaomiao-home-system/data/redis"))

	client := redis.NewClient(&redis.Options{
		Addr:            conf.Redis.Addr,
		Password:        conf.Redis.Password,
		DB:              int(conf.Redis.Db),
		MaxRetries:      3,                      // 设置最大重试次数为3
		MinRetryBackoff: 8 * time.Millisecond,   // 最小重试间隔时间
		MaxRetryBackoff: 512 * time.Millisecond, // 最大重试间隔时间
		PoolSize:        10,                     // 连接池大小
		MinIdleConns:    5,                      // 最小空闲连接数
		DialTimeout:     conf.Redis.DialTimeout.AsDuration(),
		WriteTimeout:    conf.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:     conf.Redis.ReadTimeout.AsDuration(),
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed opening connection to redis: %v", err)
	}

	client.AddHook(redisotel.TracingHook{})

	return client
}

// NewData .
func NewData(dbConfig *conf.Data, config *conf.Config, static *conf.Static, jwt *conf.Jwt, logger klog.Logger) (*Data, func(), error) {
	log := klog.NewHelper(klog.With(logger, "module", "xiaomiao-home-system/data"))

	d := &Data{
		db:       NewDB(dbConfig, logger),
		cache:    NewCache(dbConfig, logger),
		gid:      NewGID(logger),
		config:   config,
		dbConfig: dbConfig,
		static:   static,
		jwt:      jwt,
		log:      log,
	}

	return d, func() {
		log.Info("message", "closing the data resources")

		if db, err := d.db.DB(); err != nil {
			log.Error(err)
		} else {
			if err = db.Close(); err != nil {
				log.Error(err)
			}
		}

		if err := d.cache.Close(); err != nil {
			log.Error(err)
		}
	}, nil
}
