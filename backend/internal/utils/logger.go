package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger 初始化zap日志
func InitLogger() {
	// 创建基本的日志配置
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}

	// 初始化日志
	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// 测试日志输出
	Logger.Info("Logger initialized successfully")
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}
