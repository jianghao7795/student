package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewStudentRepo)

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
