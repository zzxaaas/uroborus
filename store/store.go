package store

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"uroborus/common/logging"
	settings "uroborus/common/setting"
)

// DB 存储
type DB struct {
	*gorm.DB
}

// NewPgDB 从配置中新建 Postgres 存储
func NewPgDB(config *settings.Config, logger *logging.ZapLogger) *DB {
	postgresConfig := config.Postgres
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password='%s' sslmode=disable",
		postgresConfig.Host, postgresConfig.Port, postgresConfig.Username, postgresConfig.DBName, postgresConfig.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: postgresConfig.BatchSize,
		//Logger:          gl.Default.LogMode(gl.Info),
	})
	if err != nil {
		logger.Sugar().Fatal("open dsn", dsn, "failed!", err)
	}
	return &DB{
		DB: db,
	}
}
