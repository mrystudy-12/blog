package controller

import (
	"GoWork_9/backend/internal/middleware"
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"GoWork_9/backend/internal/service"
	"GoWork_9/backend/internal/utils"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
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
		c.JSON(http.StatusBadRequest, model.Result{
			Code:    400,
			Message: "输入格式不正确，请检查后提交",
		})
		return
	}

	user, err := ctrl.userService.Register(c.Request.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Result{
			Code:    500,
			Message: "注册失败：" + err.Error(), // 封装一个翻译函数，隐藏底层细节
		})
		return
	}

	c.JSON(http.StatusOK, model.Result{
		Code:    200,
		Message: "注册成功",
		Data:    user,
	})
}

// Login 处理用户登录请求
func (ctrl *AuthController) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Result{
			Code:    400,
			Message: "用户名或密码格式不正确",
		})
		return
	}

	user, err := ctrl.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Result{
			Code:    401,
			Message: "用户名或密码错误",
		})
		return
	}

	// 生成 JWT Token，按照新格式传递参数
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Result{
			Code:    500,
			Message: "登录系统繁忙，请稍后再试",
		})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		Code:    200,
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

// GetProfile 获取当前用户资料
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, model.Result{
			Code:    401,
			Message: "登录已失效，请重新登录",
		})
		return
	}

	user, err := ctrl.userService.GetProfile(c.Request.Context(), uid.(uint64))
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, model.Result{
		Code:    200,
		Message: "获取成功",
		Data:    user,
	})
}

// UpdateProfile 更新当前用户资料
func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, model.Result{
			Code:    401,
			Message: "登录已失效，请重新登录",
		})
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.Result{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	user, err := ctrl.userService.UpdateProfile(c.Request.Context(), uid.(uint64), req)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Code:    500,
			Message: "更新失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Result{
		Code:    200,
		Message: "更新成功",
		Data:    user,
	})
}

// UploadImage 处理用户头像上传
func (ctrl *AuthController) UploadImage(c *gin.Context) {
	// 1. 获取用户ID
	userID := middleware.GetUID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, model.Result{
			Code:    401,
			Message: "未授权，请先登录",
		})
		return
	}

	// 2. 解析文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Result{
			Code:    400,
			Message: "获取上传文件失败：" + err.Error(),
		})
		return
	}

	// 3. 调用 Service 层处理上传逻辑
	avatarURL, err := ctrl.userService.UploadAvatar(userID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Result{
			Code:    500,
			Message: "上传头像失败：" + err.Error(),
		})
		return
	}
	// 4. 返回成功响应
	c.JSON(http.StatusOK, model.Result{
		Code:    200,
		Message: "头像上传成功",
		Data:    gin.H{"avatar_url": avatarURL},
	})
}

func (ctrl *AuthController) GetUserList(c *gin.Context) {
	// 1. 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 2. 调用 Service 层获取用户列表
	users, total, err := ctrl.userService.GetUserList(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Code:    500,
			Message: "获取用户列表失败: " + err.Error(),
		})
		return
	}

	// 3. 返回结果
	c.JSON(http.StatusOK, model.Result{
		Code:    200,
		Message: "获取成功",
		Data: gin.H{
			"list":      users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (ctrl *AuthController) UpdateUserStatus(c *gin.Context) {
	// 1. 解析用户ID
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Code:    400,
			Message: "无效的用户ID",
		})
		return
	}

	// 2. 获取当前管理员ID
	adminID := middleware.GetUID(c)
	if adminID == 0 {
		c.JSON(http.StatusOK, model.Result{
			Code:    401,
			Message: "未授权",
		})
		return
	}

	// 3. 防止管理员禁用自己
	if userID == adminID {
		c.JSON(http.StatusOK, model.Result{
			Code:    400,
			Message: "不能修改自己的账号状态",
		})
		return
	}

	// 4. 解析请求体
	fmt.Println("========== 开始调试 ==========")
	fmt.Printf("Content-Type: %s\n", c.ContentType())
	fmt.Printf("Content-Length: %d\n", c.Request.ContentLength)

	// 读取原始请求体
	bodyBytes, _ := c.GetRawData()
	fmt.Printf("原始请求体内容: %s\n", string(bodyBytes))
	fmt.Printf("原始请求体长度: %d\n", len(bodyBytes))

	// 重新设置请求体，因为 GetRawData 会消耗它
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req struct {
		Status int `json:"status"`
	}

	err = c.ShouldBindJSON(&req)
	fmt.Printf("绑定错误信息: %v\n", err)
	fmt.Printf("解析后的 status 值: %d\n", req.Status)
	fmt.Println("========== 结束调试 ==========")
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Code:    400,
			Message: "参数错误：请提供 status 字段（0: 禁用, 1: 正常）",
		})
		return
	}

	// 5. 验证状态值
	if req.Status != 0 && req.Status != 1 {
		c.JSON(http.StatusOK, model.Result{
			Code:    400,
			Message: "无效的状态值，只能是 0（禁用）或 1（正常）",
		})
		return
	}

	// 6. 调用 Service 层更新状态
	err = ctrl.userService.UpdateUserStatus(c.Request.Context(), userID, req.Status)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusOK, model.Result{
				Code:    404,
				Message: "用户不存在",
			})
			return
		}
		c.JSON(http.StatusOK, model.Result{
			Code:    500,
			Message: "更新状态失败: " + err.Error(),
		})
		return
	}

	// 7. 返回成功响应
	statusText := "已启用"
	if req.Status == 0 {
		statusText = "已禁用"
	}
	c.JSON(http.StatusOK, model.Result{
		Code:    200,
		Message: statusText,
		Data: gin.H{
			"user_id": userID,
			"status":  req.Status,
		},
	})
}
