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
