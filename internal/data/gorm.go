package data

import (
	"log"
	"os"
	"student/internal/conf"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDB(c *conf.Bootstrap) (*gorm.DB, error) {
	// config mysql
	mysqlConfig := mysql.Config{
		DSN:                       c.Data.Database.Source,
		DefaultStringSize:         256,
		SkipInitializeWithVersion: false,
	}
	loggerConfig := logger.New(NewWriter(log.New(os.Stdout, "\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      false,
	})
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields:                              true,
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   false,
		Logger:                                   loggerConfig.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(int(c.Data.Database.MaxIdleConns))
	sqlDB.SetMaxOpenConns(int(c.Data.Database.MaxOpenConns))
	sqlDB.SetConnMaxLifetime(time.Second * 25)
	if c.Data.Database.Debug {
		db = db.Debug()
	}
	return db, err
}

type Writer struct {
	logger.Writer
}

// NewWriter writer 构造函数

func NewWriter(w logger.Writer) *Writer {
	return &Writer{Writer: w}
}

// Printf 格式化打印日志

func (w *Writer) Printf(message string, data ...any) {
	w.Writer.Printf(message, data...)
}
