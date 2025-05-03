package data

import (
	"context"
	"student/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Data
type Data struct {
	// TODO wrapped database client
	gormDB *gorm.DB
	// TODO redis
	redis *redis.Client
}

func NewGormDB(c *conf.Bootstrap) (*gorm.DB, error) {
	dsn := c.Data.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: &logger.Recorderr.New(log.With(log.NewHelper(log.DefaultLogger), "gorm", "")),

	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(150)
	sqlDB.SetConnMaxLifetime(time.Second * 25)
	if c.Data.Database.Debug {
		db = db.Debug()
	}
	return db, err
}

func NewRedis(c *conf.Bootstrap) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            c.Data.Redis.Addr, // use default Addr
		Password:        "",
		ConnMaxIdleTime: c.Data.Redis.ReadTimeout.AsDuration(),
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("redis erros is: ", err.Error())
	} else {
		log.Infow(pong)
	}
	return client, err
}

// NewData .
func NewData(logger log.Logger, db *gorm.DB, redis *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{gormDB: db, redis: redis}, cleanup, nil
}
