package middleware

import (
	"GoWork_9/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方法
		method := c.Request.Method

		// 请求路由
		path := c.Request.URL.Path

		// 状态码
		statusCode := c.Writer.Status()

		// 客户端IP
		clientIP := c.ClientIP()

		// 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 获取日志实例
		logger := utils.GetLogger()

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			// 服务器错误
			logger.Error("HTTP Request",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.String("ip", clientIP),
				zap.Duration("latency", latency),
				zap.String("error", errorMessage),
			)
		} else if statusCode >= 400 {
			// 客户端错误
			logger.Warn("HTTP Request",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.String("ip", clientIP),
				zap.Duration("latency", latency),
				zap.String("error", errorMessage),
			)
		} else {
			// 正常请求
			logger.Info("HTTP Request",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.String("ip", clientIP),
				zap.Duration("latency", latency),
			)
		}
	}
}
