package main

import (
	"GoWork_9/backend/internal/config"
	"GoWork_9/backend/internal/controller"
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/router"
	"GoWork_9/backend/internal/utils"
	"GoWork_9/backend/pkg/db"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化日志
	utils.InitLogger()
	logger := utils.GetLogger()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal("日志初始化，失败")
		}
	}(logger)

	// 2. 初始化配置
	if err := config.LoadConfig("backend/internal/config/config.yaml"); err != nil {
		logger.Warn("Load config failed, using default", zap.Error(err))
		config.InitDefaultConfig()
	}

	// 3. 初始化数据库
	db.InitDB()

	// 执行自动迁移（包含 Image 模型）
	if err := db.DB.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Image{},
		&model.Categories{}, // 对应分类表
		&model.Comment{},    // 对应评论表
	); err != nil {
		logger.Fatal("AutoMigrate failed", zap.Error(err))
	}
	controller.InitAuth(db.DB)
	controller.InitArticle(db.DB)
	controller.InitCommentModule(db.DB)
	controller.InitCategoryModule(db.DB)
	controller.InitAdmin(db.DB)
	// 4. 设置 Gin 路由
	r := router.SetupRouter()

	// 5. 配置 HTTP Server
	srv := &http.Server{
		Addr:    ":" + config.GlobalConfig.Server.Port,
		Handler: r,
	}

	// 6. 开启协程启动服务 (非阻塞)
	go func() {
		logger.Info("Server is running", zap.String("port", config.GlobalConfig.Server.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Listen failed", zap.Error(err))
		}
	}()

	// 7. 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so no need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 8. 设置 5 秒超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
