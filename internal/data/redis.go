package data

import (
	"context"
	"student/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

func NewRedis(c *conf.Bootstrap) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            c.Data.Redis.Addr, // use default Addr
		Password:        "",
		DialTimeout:     c.Data.Redis.DialTimeout.AsDuration(),
		ReadTimeout:     c.Data.Redis.ReadTimeout.AsDuration(),
		WriteTimeout:    c.Data.Redis.WriteTimeout.AsDuration(),
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
