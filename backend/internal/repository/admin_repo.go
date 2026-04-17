package repository

import (
	"GoWork_9/backend/internal/model"
	"context"
	"gorm.io/gorm"
)

type AdminRepository interface {
	// GetGlobalStats 将四个数据库统计操作封装为一个函数
	GetGlobalStats(ctx context.Context) (*model.DashboardStats, error)
}

type adminRepoImpl struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepoImpl{db: db}
}

// GetGlobalStats 提取出的数据库操作函数
func (r *adminRepoImpl) GetGlobalStats(ctx context.Context) (*model.DashboardStats, error) {
	var stats model.DashboardStats

	// 使用事务（可选）或单一 DB 实例执行多个统计
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 文章总数
		if err := tx.Model(&model.Article{}).Count(&stats.ArticleCount).Error; err != nil {
			return err
		}
		// 2. 用户总数
		if err := tx.Model(&model.User{}).Count(&stats.UserCount).Error; err != nil {
			return err
		}
		// 3. 待审核评论 (status=0)
		if err := tx.Model(&model.Comment{}).Where("status = ?", 0).Count(&stats.CommentCount).Error; err != nil {
			return err
		}
		// 4. 总阅读量 (处理 SUM 为空的情况)
		var totalViews int64
		tx.Model(&model.Article{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews)
		stats.ViewCount = uint64(totalViews)

		return nil
	})

	return &stats, err
}
