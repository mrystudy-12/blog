package repository

import (
	"GoWork_9/backend/internal/model"
	"context"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id uint64) (*model.Comment, error) // 新增
	GetByArticleID(ctx context.Context, articleID uint64) ([]model.Comment, error)
	AdminList(ctx context.Context, page, pageSize int, keyword string) ([]model.Comment, int64, error)
	UpdateStatus(ctx context.Context, id uint64, status int8) error // 新增
	Delete(ctx context.Context, id uint64) error
}

type commentRepoImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepoImpl{db: db}
}

func (r *commentRepoImpl) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// GetByID 用于 Service 层校验和状态对比
func (r *commentRepoImpl) GetByID(ctx context.Context, id uint64) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = 0", id).First(&comment).Error
	return &comment, err
}

func (r *commentRepoImpl) GetByArticleID(ctx context.Context, articleID uint64) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.WithContext(ctx).
		Where("article_id = ? AND is_deleted = 0 AND status = 1", articleID).
		Order("created_at DESC").
		Preload("User").
		Find(&comments).Error
	return comments, err
}

// 实现 AdminList
func (r *commentRepoImpl) AdminList(ctx context.Context, page, pageSize int, keyword string) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Comment{}).Where("is_deleted = 0")

	// 如果有关键词，匹配评论内容
	if keyword != "" {
		db = db.Where("content LIKE ?", "%"+keyword+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询并加载关联用户
	err := db.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Preload("User").
		Find(&comments).Error

	return comments, total, err
}

// UpdateStatus 实现审核状态切换
func (r *commentRepoImpl) UpdateStatus(ctx context.Context, id uint64, status int8) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *commentRepoImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", id).
		Update("is_deleted", 1).Error
}
