package data

import (
	stdlog "log"
	"student/internal/conf"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewGormDB, NewData, NewRedis, NewRBACRepo, NewRBACConfig, NewRBACModelPath, NewEnforcer)

// Data
type Data struct {
	// TODO wrapped database client
	gormDB *gorm.DB
	// TODO redis
	redis    *redis.Client
	enforcer *casbin.Enforcer
}

// NewData .
func NewData(logger log.Logger, db *gorm.DB, redis *redis.Client, enforcer *casbin.Enforcer) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{gormDB: db, redis: redis, enforcer: enforcer}, cleanup, nil
}

// NewGormDB 创建数据库连接
func NewGormDB(c *conf.Bootstrap) *gorm.DB {
	config := &gorm.Config{}
	if c.Data.Database.Debug {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(c.Data.Database.Source), config)
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get sql.DB")
	}

	sqlDB.SetMaxIdleConns(int(c.Data.Database.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(c.Data.Database.MaxOpenConns))

	return db
}

// NewRedis 创建Redis客户端
func NewRedis(c *conf.Bootstrap) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Data.Redis.Addr,
		Password:     c.Data.Redis.Password,
		DialTimeout:  c.Data.Redis.DialTimeout.AsDuration(),
		ReadTimeout:  c.Data.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Data.Redis.WriteTimeout.AsDuration(),
	})
	return rdb
}

// NewRBACConfig 创建RBAC配置
func NewRBACConfig(c *conf.Bootstrap) *conf.RBAC {
	return c.Rbac
}

// NewRBACModelPath 获取RBAC模型路径
func NewRBACModelPath(c *conf.Bootstrap) string {
	return c.Rbac.ModelPath
}

// NewEnforcer 创建Casbin执行器
func NewEnforcer(db *gorm.DB, modelPath string) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		stdlog.Printf("NewEnforcer error: %v", err)
		panic(err)
	}

	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		stdlog.Printf("NewEnforcer error: %v", err)
		panic(err)
	}

	// 加载策略
	enforcer.LoadPolicy()

	return enforcer
}
