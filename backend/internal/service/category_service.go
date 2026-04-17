package service

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"context"
	"errors"
)

var (
	ErrCategoryNameEmpty = errors.New("分类名称不能为空")
	ErrCategoryNotFound  = errors.New("分类不存在")
)

// CategoryService 定义分类业务接口
type CategoryService interface {
	// GetAll 获取全部分类（用于前台展示和后台列表）
	GetAll(ctx context.Context) ([]model.Categories, error)
	// GetByID 根据ID获取单个分类
	GetByID(ctx context.Context, id uint64) (*model.Categories, error)
	// Create 创建分类
	Create(ctx context.Context, category *model.Categories) error
	// Update 更新分类
	Update(ctx context.Context, category *model.Categories) error
	// Delete 删除分类（物理删除）
	Delete(ctx context.Context, id uint64) error
}

type categoryServiceImpl struct {
	repo repository.CategoryRepository
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{repo: repo}
}

// GetAll 获取所有分类
func (s *categoryServiceImpl) GetAll(ctx context.Context) ([]model.Categories, error) {
	// 调用 Repo 层，Repo 内部已实现按 sort 排序
	return s.repo.GetAll(ctx)
}

// GetByID 根据 ID 获取详情
func (s *categoryServiceImpl) GetByID(ctx context.Context, id uint64) (*model.Categories, error) {
	return s.repo.GetByID(ctx, id)
}

// Create 创建新分类
func (s *categoryServiceImpl) Create(ctx context.Context, category *model.Categories) error {
	if category.Name == "" {
		return ErrCategoryNameEmpty
	}
	// 可以在这里增加逻辑：检查分类名是否已存在
	return s.repo.Create(ctx, category)
}

// Update 更新分类信息
func (s *categoryServiceImpl) Update(ctx context.Context, category *model.Categories) error {
	// 1. 先校验分类是否存在
	if _, err := s.repo.GetByID(ctx, category.ID); err != nil {
		return ErrCategoryNotFound
	}

	if category.Name == "" {
		return ErrCategoryNameEmpty
	}

	return s.repo.Update(ctx, category)
}

// Delete 删除分类
func (s *categoryServiceImpl) Delete(ctx context.Context, id uint64) error {
	// 1. 校验是否存在
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return ErrCategoryNotFound
	}

	// 2. 执行物理删除
	return s.repo.Delete(ctx, id)
}
