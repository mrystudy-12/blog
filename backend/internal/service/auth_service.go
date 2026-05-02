package service

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"GoWork_9/backend/internal/utils"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	GetProfile(ctx context.Context, id uint64) (*model.User, error)
	UpdateProfile(ctx context.Context, id uint64, req model.UpdateProfileRequest) (*model.User, error)
	UploadAvatar(id uint64, file *multipart.FileHeader) (string, error)
	GetUserList(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	UpdateUserStatus(ctx context.Context, userID uint64, status int) error
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

func (s *userServiceImpl) GetProfile(ctx context.Context, id uint64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) UpdateProfile(ctx context.Context, id uint64, req model.UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.AvatarUrl != "" && req.AvatarUrl != user.AvatarUrl {
		oldAvatar := user.AvatarUrl
		user.AvatarUrl = req.AvatarUrl
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
		if oldAvatar != "" {
			go s.deleteOldAvatar(oldAvatar)
		}
	} else {
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

// UploadAvatar 处理用户头像上传，仅保存文件并返回URL
// 数据库更新由 profile 接口统一管理
func (s *userServiceImpl) UploadAvatar(userID uint64, file *multipart.FileHeader) (string, error) {
	if err := s.validateImage(file); err != nil {
		return "", err
	}

	saveDir, err := utils.GetSavePath("avatars")
	if err != nil {
		return "", fmt.Errorf("无法准备存储目录: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	newFilename := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), userID, ext)
	fullFilePath := filepath.Join(saveDir, newFilename)

	if err := s.saveFileToDisk(file, fullFilePath); err != nil {
		return "", err
	}

	avatarURL := utils.GetAccessURL("avatars", newFilename)
	return avatarURL, nil
}

func (s *userServiceImpl) GetUserList(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := s.userRepo.ListUsers(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// deleteOldAvatar 删除旧的头像文件
func (s *userServiceImpl) deleteOldAvatar(avatarURL string) {
	filename := filepath.Base(avatarURL)
	if filename == "" || filename == "." {
		return
	}

	// 直接复用 utils，获取 avatars 文件夹的绝对路径
	avatarsDir, err := utils.GetSavePath("avatars")
	if err != nil {
		utils.GetLogger().Error("获取头像目录失败", zap.Error(err))
		return
	}

	oldFilePath := filepath.Join(avatarsDir, filename)

	// 检查并删除
	if _, err := os.Stat(oldFilePath); err == nil {
		if err := os.Remove(oldFilePath); err != nil {
			utils.GetLogger().Warn("删除旧头像失败", zap.String("path", oldFilePath), zap.Error(err))
		}
	}
}

// validateImage 校验图片大小和格式
func (s *userServiceImpl) validateImage(file *multipart.FileHeader) error {
	// 限制 2MB
	const maxFileSize = 2 << 20
	if file.Size > maxFileSize {
		return fmt.Errorf("头像大小不能超过 2MB (当前 %.2f MB)", float64(file.Size)/1024/1024)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowExts[ext] {
		return fmt.Errorf("不支持格式: %s", ext)
	}
	return nil
}

// saveFileToDisk 封装具体的物理写入逻辑
func (s *userServiceImpl) saveFileToDisk(file *multipart.FileHeader, destPath string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			fmt.Println("关闭文件失败")
		}
	}(src)

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		err := dst.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf("写入文件失败: %w", err)
	}

	_ = dst.Sync() // 确保落盘
	return dst.Close()
}

func (s *userServiceImpl) UpdateUserStatus(ctx context.Context, userID uint64, status int) error {
	// 1. 验证用户是否存在
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// 2. 验证状态值是否合法（0: 禁用, 1: 正常）
	if status != 0 && status != 1 {
		return fmt.Errorf("无效的状态值: %d", status)
	}

	// 3. 不允许禁用自己
	// 注意：这里需要在 Controller 层传入当前管理员ID进行判断

	// 4. 更新状态
	return s.userRepo.UpdateStatus(ctx, userID, status)
}
