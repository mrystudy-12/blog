package middleware

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthJWT JWT 认证中间件
func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.Result{
				Code:    401,
				Message: "请求未携带有效的 Authorization 头部",
			})
			c.Abort()
			return
		}

		// 按空格分割，期望格式为 "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization 头部格式错误"})
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			// 如果 token 过期，返回特定 code 告知前端刷新 token
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "登录已过期，请重新登录",
			})
			c.Abort()
			return
		}

		// 将解析出的 Claims 存入上下文，方便后续 Handler 使用
		c.Set("userID", claims.ID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// IsAdmin 校验是否为管理员身份
func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取 AuthJWT 存入的 role
		role, exists := c.Get("role")
		fmt.Println("-----------------------------", role)
		// 逻辑判断：如果角色不是管理员（假设 1 为管理员，或者判断字符串 "admin"）
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限不足：该操作仅限管理员",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUID  从上下文中获取当前用户 ID
func GetUID(c *gin.Context) uint64 {
	uidVal, exists := c.Get("userID")
	if !exists {
		return 0
	}
	switch v := uidVal.(type) {
	case uint64:
		return v
	case int64:
		return uint64(v)
	case float64: // JWT 解析数字时有时会默认转为 float64
		return uint64(v)
	default:
		return 0
	}
}
