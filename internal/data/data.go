package data

import (
	stdlog "log"
	"student/internal/conf"
	"student/internal/pkg/jwt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewGormDB, NewData, NewRedis, NewStudentRepo, NewUserRepo, NewRBACRepo, NewErrorRepo, NewJWTConfig, NewRBACConfig, NewRBACModelPath)

// Data
type Data struct {
	// TODO wrapped database client
	gormDB *gorm.DB
	// TODO redis
	redis *redis.Client
}

// NewData .
func NewData(logger log.Logger, db *gorm.DB, redis *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{gormDB: db, redis: redis}, cleanup, nil
}

// NewJWTConfig 创建JWT配置
func NewJWTConfig(c *conf.Bootstrap) *jwt.Config {
	config := &jwt.Config{
		SecretKey: c.Jwt.SecretKey,
		Expire:    c.Jwt.Expire.AsDuration(),
	}

	// 添加调试日志
	stdlog.Printf("JWT配置: SecretKey长度=%d, SecretKey=%s, Expire=%v",
		len(config.SecretKey), config.SecretKey, config.Expire)

	return config
}

// NewRBACConfig 创建RBAC配置
func NewRBACConfig(c *conf.Bootstrap) *conf.RBAC {
	return c.Rbac
}

// NewRBACModelPath 获取RBAC模型路径
func NewRBACModelPath(c *conf.Bootstrap) string {
	return c.Rbac.ModelPath
}
