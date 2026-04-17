package controller

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"GoWork_9/backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

var AdminCtrl *AdminController

func InitAdmin(db *gorm.DB) {
	// 1. 初始化 Repository (负责执行数据库统计 SQL)
	repo := repository.NewAdminRepository(db)

	// 2. 初始化 Service (负责业务逻辑，注入 AdminRepository)
	svc := service.NewAdminService(repo)

	// 3. 初始化 Controller (注入 AdminService 并赋值给全局变量 AdminCtrl)
	// 注意：这里调用的是我们之前定义的 NewAdminController
	NewAdminController(svc)
}

type AdminController struct {
	adminService service.AdminService // 需在 service 层定义对应的统计接口
}

// NewAdminController 初始化控制器并注入 Service
func NewAdminController(adminService service.AdminService) {
	AdminCtrl = &AdminController{
		adminService: adminService,
	}
}

// GetStats 获取后台仪表盘统计数据
// GET /api/v1/admin/dashboard
func (ctrl *AdminController) GetStats(c *gin.Context) {
	// 调用 service 层获取汇总数据
	stats, err := ctrl.adminService.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, model.Result{
			Code:    500,
			Message: "获取统计数据失败",
		})
		return
	}

	c.JSON(http.StatusOK, model.Result{
		Code:    200,
		Message: "success",
		Data:    stats,
	})
}
