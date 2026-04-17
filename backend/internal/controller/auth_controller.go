package controller

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"GoWork_9/backend/internal/service"
	"GoWork_9/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var (
	AuthCtrl *AuthController
)

// InitAuth 处理认证组件的初始化
func InitAuth(db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	AuthCtrl = NewAuthController(userService)
}

// AuthController 处理身份验证相关的 HTTP 请求
type AuthController struct {
	userService service.UserService
}

// NewAuthController 创建新的 AuthController 实例
func NewAuthController(userService service.UserService) *AuthController {
	return &AuthController{userService: userService}
}

// Register 处理用户注册请求
func (ctrl *AuthController) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	user, err := ctrl.userService.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.RegisterResponse{
		Message: "注册成功",
		Data:    user,
	})
}

// Login 处理用户登录请求
func (ctrl *AuthController) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	user, err := ctrl.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "登录失败: " + err.Error()})
		return
	}

	// 生成 JWT Token，按照新格式传递参数
	token, err := utils.GenerateToken(user.ID, user.Username, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		Message: "登录成功",
		Token:   token,
		Data:    user,
	})
}

func (ctrl *AuthController) GetMe(c *gin.Context) {
	// 1. 从 Context 中提取中间件 (AuthJWT) 存入的 userID
	// 注意：c.Get 返回的是 interface{}，需要类型断言
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, model.Result{
			Code:    401,
			Message: "登录已失效，请重新登录",
		})
		return
	}

	// 2. 类型转换并调用 Service
	// 假设你 Service 接受的是 uint，这里断言后需转换
	// 如果中间件存的是 uint32，则 uint(uid.(uint32))
	res, err := ctrl.userService.GetUserInfo(c.Request.Context(), uid.(uint64))
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Code:    500,
			Message: "服务器内部错误: " + err.Error(),
		})
		return
	}

	// 3. 返回 Service 封装好的 model.Result
	c.JSON(http.StatusOK, res)
}
