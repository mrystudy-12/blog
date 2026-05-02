package repository

import (
	"GoWork_9/backend/internal/model"
	"context"
	"gorm.io/gorm"
)

// UserRepository 定义用户操作接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdateAvatar(ctx context.Context, userID uint64, avatarURL string) error
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	UpdateStatus(ctx context.Context, userID uint64, status int) error
}

type userRepoImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建新的 User Repository 实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepoImpl{db: db}
}

// Create 注册用户
func (r *userRepoImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByUsername 用于登录验证
func (r *userRepoImpl) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepoImpl) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *userRepoImpl) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
func (r *userRepoImpl) UpdateAvatar(ctx context.Context, userID uint64, avatarURL string) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("avatar_url", avatarURL).Error
}

func (r *userRepoImpl) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按 ID 降序排列
	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Order("id DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateStatus 更新用户状态
func (r *userRepoImpl) UpdateStatus(ctx context.Context, userID uint64, status int) error {
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Update("status", status).Error
}
