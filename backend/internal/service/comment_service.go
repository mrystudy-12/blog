package service

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrCommentContentEmpty = errors.New("评论内容不能为空")
	ErrCommentNotFound     = errors.New("评论不存在")
)

// CommentService 定义评论业务接口
type CommentService interface {
	// Create 前台接口
	Create(ctx context.Context, userID uint64, req model.CreateCommentRequest) error
	GetByArticle(ctx context.Context, articleID uint64) ([]model.Comment, error)

	// AdminList 后台管理接口
	AdminList(ctx context.Context, page, pageSize int, keyword string) (*model.Result, error)
	Delete(ctx context.Context, id uint64) error
	Audit(ctx context.Context, id uint64, pass bool) error // 审核评论（显示/隐藏）
}

type commentServiceImpl struct {
	repo repository.CommentRepository
}

// NewCommentService 创建实例
func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentServiceImpl{repo: repo}
}

// Create 发表评论
func (s *commentServiceImpl) Create(ctx context.Context, userID uint64, req model.CreateCommentRequest) error {
	if req.Content == "" {
		return ErrCommentContentEmpty
	}

	comment := &model.Comment{
		ArticleID: req.ArticleID,
		UserID:    userID,
		Content:   req.Content,
		ParentID:  req.ParentID, // 支持回复功能
		Status:    1,            // 默认正常显示
		IsDeleted: 0,
	}

	return s.repo.Create(ctx, comment)
}

// GetByArticle 获取文章评论
// 前台调用：Repo 层必须过滤 status=1 和 is_deleted=0
func (s *commentServiceImpl) GetByArticle(ctx context.Context, articleID uint64) ([]model.Comment, error) {
	// 调用 Repo 层，Repo 内部应执行 Preload("User") 以便前端显示头像
	return s.repo.GetByArticleID(ctx, articleID)
}

// Delete 管理员或用户执行删除
// 业务逻辑：将 is_deleted 标记为 1
func (s *commentServiceImpl) Delete(ctx context.Context, id uint64) error {
	// 建议增加：先检查该评论是否存在
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return ErrCommentNotFound
	}
	return s.repo.Delete(ctx, id)
}

// Audit 审核评论逻辑
// 模式：先发布(Status=1)，后审核。违规时 pass 传 false 隐藏
func (s *commentServiceImpl) Audit(ctx context.Context, id uint64, pass bool) error {
	// 1. 检查评论是否存在
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrCommentNotFound
	}

	// 2. 确定状态值：通过为 1，屏蔽为 0
	var status int8 = 0
	if pass {
		status = 1
	}

	// 3. 如果状态没变，直接返回节省数据库开销
	if comment.Status == status {
		return nil
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// AdminList 后台管理列表 (待实现具体 Repo 分页逻辑)
func (s *commentServiceImpl) AdminList(ctx context.Context, page, pageSize int, keyword string) (*model.Result, error) {
	list, total, err := s.repo.AdminList(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, err
	}

	return &model.Result{
		Code:    200,
		Message: "success",
		Data: gin.H{
			"list":  list,
			"total": total,
		},
	}, nil
}
