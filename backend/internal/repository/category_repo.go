package repository

import (
	"GoWork_9/backend/internal/model"
	"context"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Categories) error
	GetAll(ctx context.Context) ([]model.Categories, error)
	GetByID(ctx context.Context, id uint32) (*model.Categories, error)
	Update(ctx context.Context, category *model.Categories) error
	Delete(ctx context.Context, id uint32) error // 注意：这里将执行物理删除
}

type categoryRepoImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepoImpl{db: db}
}

func (r *categoryRepoImpl) Create(ctx context.Context, category *model.Categories) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetAll 获取所有分类，按 sort 升序排列
func (r *categoryRepoImpl) GetAll(ctx context.Context) ([]model.Categories, error) {
	var categories []model.Categories
	err := r.db.WithContext(ctx).Order("sort asc").Find(&categories).Error
	return categories, err
}

func (r *categoryRepoImpl) GetByID(ctx context.Context, id uint32) (*model.Categories, error) {
	var category model.Categories
	err := r.db.WithContext(ctx).First(&category, id).Error
	return &category, err
}

func (r *categoryRepoImpl) Update(ctx context.Context, category *model.Categories) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete 物理删除
func (r *categoryRepoImpl) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&model.Categories{}, id).Error
}
