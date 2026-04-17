package db

import (
	"GoWork_9/backend/internal/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 确保配置已初始化
	if config.GlobalConfig == nil {
		config.InitDefaultConfig()
	}

	cfg := config.GlobalConfig.Database

	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		// 禁用外键约束，符合用户要求 "不用外键"
		DisableForeignKeyConstraintWhenMigrating: true,
		// 配置日志，方便调试
		Logger: logger.New(

			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get sqlDB: %v", err))
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.LifeTime) * time.Hour)

	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("Database initialized but Ping failed: %v", err))
	}

	fmt.Println("Database connection established successfully")
}
