package controller

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"GoWork_9/backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

var CategoryCtrl *CategoryController

type CategoryController struct {
	categoryService service.CategoryService
}

// InitCategoryModule 统一初始化逻辑
func InitCategoryModule(db *gorm.DB) {
	repo := repository.NewCategoryRepository(db)
	svc := service.NewCategoryService(repo)
	CategoryCtrl = &CategoryController{categoryService: svc}
}

func (ctrl *CategoryController) AdminList(c *gin.Context) {
	list, err := ctrl.categoryService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 500, Message: "获取失败"})
		return
	}
	c.JSON(http.StatusOK, model.Result{Code: 200, Message: "success", Data: list})
}

// GetByID 用于编辑回显
func (ctrl *CategoryController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 400, Message: "无效的分类ID"})
		return
	}

	category, err := ctrl.categoryService.GetByID(c.Request.Context(), uint32(id))
	if err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 404, Message: "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Result{Code: 200, Message: "success", Data: category})
}

func (ctrl *CategoryController) Create(c *gin.Context) {
	var category model.Categories
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 400, Message: "参数错误"})
		return
	}
	if err := ctrl.categoryService.Create(c.Request.Context(), &category); err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Result{Code: 200, Message: "创建成功"})
}

// Update 更新分类
func (ctrl *CategoryController) Update(c *gin.Context) {
	var category model.Categories // 统一使用单数
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 400, Message: "参数错误"})
		return
	}

	if err := ctrl.categoryService.Update(c.Request.Context(), &category); err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 500, Message: "更新失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Result{Code: 200, Message: "更新成功"})
}

func (ctrl *CategoryController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := ctrl.categoryService.Delete(c.Request.Context(), uint32(id)); err != nil {
		c.JSON(http.StatusOK, model.Result{Code: 500, Message: "删除失败"})
		return
	}
	c.JSON(http.StatusOK, model.Result{Code: 200, Message: "删除成功"})
}
