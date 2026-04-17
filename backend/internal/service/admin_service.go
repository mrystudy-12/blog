package service

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"context"
)

// AdminService 定义后台管理业务接口
type AdminService interface {
	// GetDashboardStats 获取仪表盘统计数据
	GetDashboardStats(ctx context.Context) (*model.DashboardStats, error)
}

type adminServiceImpl struct {
	// 修正：确保字段名与 Repo 接口对应
	adminRepo repository.AdminRepository
}

// NewAdminService 创建 Admin Service 实例
// 修正：参数类型应改为 AdminRepository
func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminServiceImpl{
		adminRepo: adminRepo,
	}
}

// GetDashboardStats 实现统计逻辑
func (s *adminServiceImpl) GetDashboardStats(ctx context.Context) (*model.DashboardStats, error) {
	// 逻辑层只需简单调用 Repo 层封装好的数据库操作
	return s.adminRepo.GetGlobalStats(ctx)
}
