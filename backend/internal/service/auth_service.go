package service

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists = errors.New("用户名已存在")
	ErrUserNotFound      = errors.New("用户不存在")
	ErrInvalidPassword   = errors.New("密码错误")
)

// UserService 定义用户业务接口
type UserService interface {
	Register(ctx context.Context, username, password, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (*model.User, error)
	GetUserInfo(ctx context.Context, id uint64) (*model.Result, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService 创建新的 User Service 实例
func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: repo}
}

// Register 处理用户注册逻辑
func (s *userServiceImpl) Register(ctx context.Context, username, password, email string) (*model.User, error) {
	// 1. 检查用户是否已存在
	existingUser, err := s.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 2. 密码加密 (bcrypt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. 创建用户模型
	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		// Avatar 字段留空，后续再处理
	}

	// 4. 持久化到数据库
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 处理用户登录逻辑
func (s *userServiceImpl) Login(ctx context.Context, username, password string) (*model.User, error) {
	// 1. 根据用户名查询用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 2. 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

func (s *userServiceImpl) GetUserInfo(ctx context.Context, id uint64) (*model.Result, error) {
	// 1. 调用 Repository 层获取数据
	user, err := s.userRepo.GetByID(ctx, id)

	// 2. 逻辑处理
	if err != nil {
		// 如果是 GORM 的“未找到记录”错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.Result{
				Code:    404,
				Message: "该用户不存在",
			}, nil
		}
		// 其他数据库错误
		return nil, err
	}

	// 3. 封装统一的返回格式
	return &model.Result{
		Code:    200,
		Message: "获取成功",
		Data:    user,
	}, nil
}
