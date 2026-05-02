package controller

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"GoWork_9/backend/internal/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var CommentCtrl *CommentController

type CommentController struct {
	commentService service.CommentService
}

// InitCommentModule 统一初始化逻辑
func InitCommentModule(db *gorm.DB) {
	repo := repository.NewCommentRepository(db)
	svc := service.NewCommentService(repo)
	CommentCtrl = &CommentController{commentService: svc}
}

// ======================== 前台功能 (Portal) ========================

// Create 发表评论
// 适配 Service: Create(ctx, userID uint32, req model.CreateCommentRequest)
func (ctrl *CommentController) Create(c *gin.Context) {
	var req model.CreateCommentRequest
	// 1. 绑定到指定的 Request 结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(req)
		c.JSON(http.StatusOK, model.Result{Code: 400, Message: "参数错误"})
		return
	}

	// 2. 从中间件获取 UserID (uint32)
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, model.Result{Code: 401, Message: "请登录后发表评论"})
		return
	}

	// 3. 调用 Service
	if err := ctrl.commentService.Create(c.Request.Context(), uid.(uint64), req); err != nil {
		code := 500
		if errors.Is(err, service.ErrCommentContentEmpty) {
			code = http.StatusBadRequest // 400
		}
		c.JSON(http.StatusOK, model.Result{Code: code, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Result{Code: 200, Message: "评论成功"})
}

// GetByArticleID 查看文章评论列表（分页）
func (ctrl *CommentController) GetByArticleID(c *gin.Context) {
	idStr := c.Param("aid")
	articleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 400, Message: "无效的文章ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	res, err := ctrl.commentService.GetByArticle(c.Request.Context(), articleID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 500, Message: "获取评论失败"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ======================== 后台管理 (Admin) ========================

// AdminList 后台管理列表
// 适配 Service: AdminList(ctx, page, pageSize, keyword)
func (ctrl *CommentController) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	res, err := ctrl.commentService.AdminList(c.Request.Context(), page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 500, Message: "获取失败"})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Audit 审核评论
// 适配 Service: Audit(ctx, id uint64, pass bool)
func (ctrl *CommentController) Audit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64) // Service 要求 uint64
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 400, Message: "无效的评论ID"})
		return
	}

	// 假设前端传 ?pass=true 或 ?pass=false
	pass, _ := strconv.ParseBool(c.Query("pass"))

	if err := ctrl.commentService.Audit(c.Request.Context(), id, pass); err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 500, Message: err.Error()})
		return
	}
	message := "审核通过"
	if !pass {
		message = "评论已屏蔽"
	}
	c.JSON(http.StatusOK, model.Result{Code: 200, Message: message})
}

// Delete 物理删除/逻辑删除
// 适配 Service: Delete(ctx, id uint64)
func (ctrl *CommentController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64) // Service 要求 uint64
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 400, Message: "无效的评论ID"})
		return
	}

	if err := ctrl.commentService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Result{Code: 200, Message: "删除成功"})
}
