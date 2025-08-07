package data

import (
	"student/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewGormDB, NewData, NewRedis, NewStudentRepo)

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
