package repository

import (
	"GoWork_9/backend/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// ArticleRepository 定义文章操作接口
type ArticleRepository interface {
	Create(ctx context.Context, tx *gorm.DB, article *model.Article) error
	GetByID(ctx context.Context, id uint64) (*model.Article, error)
	List(ctx context.Context, page, pageSize int, keyword string, status int) ([]model.Article, int64, error)
	Update(ctx context.Context, tx *gorm.DB, article *model.Article) error
	Delete(ctx context.Context, tx *gorm.DB, id uint64) error

	// CreateImage 创建图片记录
	CreateImage(ctx context.Context, image *model.Image) error
	// UpdateImagesArticleID 批量更新图片的关联文章ID
	UpdateImagesArticleID(ctx context.Context, tx *gorm.DB, articleID uint64, paths []string) error
	// GetImagesByArticleID 获取文章关联的所有图片
	GetImagesByArticleID(ctx context.Context, articleID uint64) ([]model.Image, error)
	// UnbindImages 将图片的 article_id 置为 0
	UnbindImages(ctx context.Context, tx *gorm.DB, paths []string) error
	// DeleteImagesByArticleID 根据文章ID删除对应的图片记录
	DeleteImagesByArticleID(ctx context.Context, articleID uint64) error
	// Transaction 开启事务
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type articleRepoImpl struct {
	db *gorm.DB
}

// NewArticleRepository 创建新的 Repository 实例
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepoImpl{db: db}
}

func (r *articleRepoImpl) Create(ctx context.Context, tx *gorm.DB, article *model.Article) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Create(article).Error
}

// CreateImage 创建图片记录
func (r *articleRepoImpl) CreateImage(ctx context.Context, image *model.Image) error {
	return r.db.WithContext(ctx).Create(image).Error
}

// UpdateImagesArticleID 负责将图片记录与文章ID关联。该函数需支持事务上下文。
func (r *articleRepoImpl) UpdateImagesArticleID(ctx context.Context, tx *gorm.DB, articleID uint64, paths []string) error {
	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	// 更新 article_id 为 0 的记录（确保只关联本次上传且未被关联过的图片）
	return tx.Model(&model.Image{}).
		Where("url IN ? AND article_id = ?", paths, 0).
		Update("article_id", articleID).Error
}

// GetImagesByArticleID 获取文章关联的所有图片
func (r *articleRepoImpl) GetImagesByArticleID(ctx context.Context, articleID uint64) ([]model.Image, error) {
	var images []model.Image
	err := r.db.WithContext(ctx).Where("article_id = ?", articleID).Find(&images).Error
	return images, err
}

// UnbindImages 将图片的 article_id 置为 0
func (r *articleRepoImpl) UnbindImages(ctx context.Context, tx *gorm.DB, paths []string) error {
	if tx == nil {
		tx = r.db.WithContext(ctx)
	}
	return tx.Model(&model.Image{}).
		Where("url IN ?", paths).
		Update("article_id", 0).Error
}

// Transaction 执行事务
func (r *articleRepoImpl) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *articleRepoImpl) GetByID(ctx context.Context, id uint64) (*model.Article, error) {
	var article model.Article
	// 1. 保持基础查询：不带状态过滤，这样后台回显也能用
	// 2. Preload 保持不变，用于关联查询作者和分类信息
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Category").
		Preload("Images").
		First(&article, id).Error

	if err != nil {
		fmt.Println("文章消失了")
		return nil, err // 发生错误（如记录不存在）返回 nil
	}
	return &article, nil
}

func (r *articleRepoImpl) List(ctx context.Context, page, pageSize int, keyword string, status int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Article{})

	// 处理状态过滤：status 为 1(发布) 或 0(草稿)；传入 -1 时表示查询全部(后台用)
	if status != -1 {
		db = db.Where("status = ?", status)
	}

	if keyword != "" {
		// 使用 GORM 的 Group 条件，避免与 status 产生逻辑冲突
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 1. 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 2. 分页查询
	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").
		Preload("User").   // 预加载作者，防止前端收到 null
		Preload("Images"). // 预加载图片列表
		Preload("Category").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error
	return articles, total, err
}

func (r *articleRepoImpl) Update(ctx context.Context, tx *gorm.DB, article *model.Article) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", article.ID).
		Select("Title", "Content", "CategoryID", "Status").
		Updates(article).Error
}

// DeleteImagesByArticleID 根据文章ID删除对应的图片记录
func (r *articleRepoImpl) DeleteImagesByArticleID(ctx context.Context, articleID uint64) error {
	return r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&model.Image{}).Error
}

func (r *articleRepoImpl) Delete(ctx context.Context, tx *gorm.DB, id uint64) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Delete(&model.Article{}, id).Error
}
