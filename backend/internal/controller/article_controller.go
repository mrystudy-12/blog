package controller

import (
	"GoWork_9/backend/internal/config"
	"GoWork_9/backend/internal/middleware"
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"GoWork_9/backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var ArticleCtrl *ArticleController

// InitArticle 初始化文章组件
func InitArticle(db *gorm.DB) {
	repo := repository.NewArticleRepository(db)
	svc := service.NewArticleService(repo)
	ArticleCtrl = NewArticleController(svc)
}

type ArticleController struct {
	articleService service.ArticleService
}

func NewArticleController(articleService service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

// sendResponse 统一响应处理
func (ctrl *ArticleController) sendResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, model.Result{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Create 处理新增文章请求
func (ctrl *ArticleController) Create(c *gin.Context) {
	var req model.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.sendResponse(c, 400, "参数验证失败: "+err.Error(), nil)
		return
	}

	authorID := middleware.GetUID(c)
	if authorID == 0 {
		ctrl.sendResponse(c, 401, "无效的用户身份", nil)
		return
	}

	article, err := ctrl.articleService.Create(c.Request.Context(), authorID, req)
	if err != nil {
		ctrl.sendResponse(c, 500, "新增失败: "+err.Error(), nil)
		return
	}

	ctrl.sendResponse(c, 200, "发布成功", article)
}

func (ctrl *ArticleController) GetByID(c *gin.Context) {
	// 1. 解析 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctrl.sendResponse(c, 400, "无效的文章ID", nil)
		return
	}

	// 2. 调用后台专用的 Service 方法 (不校验 Status，支持获取草稿)
	article, err := ctrl.articleService.GetAdminDetail(c.Request.Context(), id)
	if err != nil {
		ctrl.sendResponse(c, 500, "获取文章失败: "+err.Error(), nil)
		return
	}

	// 3. 返回完整数据
	ctrl.sendResponse(c, 200, "success", article)
}

// UploadImage 处理图片上传请求
func (ctrl *ArticleController) UploadImage(c *gin.Context) {
	// 1. 权限拦截
	authorID := middleware.GetUID(c)
	if authorID == 0 {
		ctrl.sendResponse(c, 401, "请先登录", nil)
		return
	}

	// 2. Controller 层轻量拦截：根据请求头判断大小
	// 这样不需要解析 Body 就能直接拒绝超大请求
	const maxBodySize = 5 << 20
	if c.Request.ContentLength > maxBodySize {
		ctrl.sendResponse(c, 413, "上传的文件太大，系统拒绝处理", nil)
		return
	}

	// 3. 解析表单文件
	file, err := c.FormFile("image")
	if err != nil {
		ctrl.sendResponse(c, 400, "文件获取失败", nil)
		return
	}

	// 1. 优先从请求头获取真实的协议 (Nginx/LB 转发)
	baseURL := config.GlobalConfig.Server.BaseURL

	// 4. 调用 Service
	url, err := ctrl.articleService.HandleImageUpload(c.Request.Context(), file, authorID, baseURL)
	if err != nil {
		// 如果是业务校验错误返回 400，系统错误返回 500
		if strings.Contains(err.Error(), "格式") || strings.Contains(err.Error(), "大小") {
			ctrl.sendResponse(c, 400, err.Error(), nil)
		} else {
			ctrl.sendResponse(c, 500, "服务器保存图片出错", nil)
		}
		return
	}

	// 5. 成功返回
	ctrl.sendResponse(c, 200, "上传成功", gin.H{"url": url})
}

// AdminList 获取文章列表
func (ctrl *ArticleController) AdminList(c *gin.Context) {
	// 获取分页与搜索参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword") // 模糊查询关键字

	// 此时需要你的 Service 层 List 方法支持 keyword 参数
	// 如果 Service 还没改，请务必在 Service 层增加对 keyword 的处理
	result, err := ctrl.articleService.List(c.Request.Context(), page, pageSize, keyword, -1)
	if err != nil {
		ctrl.sendResponse(c, 500, "获取列表失败", nil)
		return
	}

	ctrl.sendResponse(c, 200, "success", result)
}

// Update 更新文章
func (ctrl *ArticleController) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req model.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.sendResponse(c, 400, "数据格式错误", nil)
		return
	}

	article, err := ctrl.articleService.Update(c.Request.Context(), id, req)
	if err != nil {
		ctrl.sendResponse(c, 500, "更新失败", nil)
		return
	}

	ctrl.sendResponse(c, 200, "更新成功", article)
}

// Delete 删除文章
func (ctrl *ArticleController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := ctrl.articleService.Delete(c.Request.Context(), id); err != nil {
		ctrl.sendResponse(c, 500, "删除失败", nil)
		return
	}
	ctrl.sendResponse(c, 200, "删除成功", nil)
}

// =================================================================前台门户接口(Portal)============================================================================================

func (ctrl *ArticleController) PortalGet(c *gin.Context) {
	// 1. 解析 ID (建议使用 64 位解析以防溢出)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctrl.sendResponse(c, 400, "无效的文章ID", nil)
		return
	}

	// 2. 调用专门的前台 Service 方法 (内部已包含 Status=1 的逻辑检查)
	// 注意：根据你的 service 接口定义，需要转换为 uint
	article, err := ctrl.articleService.GetPortalDetail(c.Request.Context(), id)
	if err != nil {
		// 如果文章未发布或不存在，Service 会返回相应错误
		ctrl.sendResponse(c, 404, err.Error(), nil)
		return
	}

	// 3. 返回数据
	ctrl.sendResponse(c, 200, "success", article)
}

func (ctrl *ArticleController) PortalList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword") // 前台模糊查询

	// 逻辑：强制过滤 status = 1 (已发布)
	result, err := ctrl.articleService.List(c.Request.Context(), page, pageSize, keyword, 1)
	if err != nil {
		ctrl.sendResponse(c, 500, "获取列表失败", nil)
		return
	}
	ctrl.sendResponse(c, 200, "success", result)
}
