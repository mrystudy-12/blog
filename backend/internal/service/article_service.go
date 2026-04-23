package service

import (
	"GoWork_9/backend/internal/model"
	"GoWork_9/backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrInvalidArticleTitle   = errors.New("标题不能为空")
	ErrInvalidArticleContent = errors.New("内容不能为空")
	ErrArticleNotFound       = errors.New("文章不存在")
)

// ArticleService 定义文章业务接口
type ArticleService interface {
	Create(ctx context.Context, authorID uint64, req model.CreateArticleRequest) (*model.Article, error)
	// GetAdminDetail 后台获取文章详情 (用于编辑回显，不校验发布状态)
	GetAdminDetail(ctx context.Context, id uint64) (*model.Article, error)

	// GetPortalDetail 前台获取文章详情 (仅限已发布 Status=1，用于展示)
	GetPortalDetail(ctx context.Context, id uint64) (*model.Article, error)

	List(ctx context.Context, page, pageSize int, keyword string, status int) (*model.ArticleListResponse, error)
	Update(ctx context.Context, id uint64, req model.UpdateArticleRequest) (*model.Article, error)
	Delete(ctx context.Context, id uint64) error

	// HandleImageUpload 处理图片上传
	HandleImageUpload(ctx context.Context, file *multipart.FileHeader, userID uint64, baseURl string) (string, error)
}

type articleServiceImpl struct {
	repo repository.ArticleRepository
	cop  repository.CommentRepository
}

const articleImagePath = "frontend/static/images/articles"

// 获取项目根目录并构建图片存储路径
func getImagesPath() string {
	// 获取当前工作目录 (假设在 backend 目录下运行)
	wd, err := os.Getwd()
	fmt.Println(">>> 当前工作目录 WD:", wd)
	if err != nil {
		return articleImagePath
	}
	path := filepath.Join(wd, "frontend", "static", "images", "articles")

	if _, err := os.Stat(path); err == nil {
		return path
	}
	// 如果在 backend 目录下运行，只需要向上跳一级到项目根目录
	// D:\GoWork_9\backend -> D:\GoWork_9 -> D:\GoWork_9\frontend\...
	//path := filepath.Join(wd, "..", "frontend", "static", "images", "articles")
	fmt.Println(">>> 尝试保存图片的绝对路径为:", path) // 在后端控制台查看输出
	return path
}

// NewArticleService 创建新的 Article Service 实例
func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleServiceImpl{repo: repo}
}

// Create 新增文章业务逻辑
func (s *articleServiceImpl) Create(ctx context.Context, authorID uint64, req model.CreateArticleRequest) (*model.Article, error) {
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if title == "" {
		return nil, ErrInvalidArticleTitle
	}
	if content == "" {
		return nil, ErrInvalidArticleContent
	}

	status := req.Status
	if status != 0 && status != 1 {
		status = 0
	}

	article := &model.Article{
		Title:      title,
		Content:    content,
		AuthorID:   authorID,
		CategoryID: req.CategoryID,
		Status:     status,
		ImageURLs:  req.ImageURLs, // 将自定义数组传入模型
	}

	// 开启事务，确保文章创建和图片关联原子化
	err := s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		// 1. 创建文章
		if err := s.repo.Create(ctx, tx, article); err != nil {
			return err
		}

		// 2. 获取并关联图片 (从自定义数组中读取)
		if len(article.ImageURLs) > 0 {
			if err := s.repo.UpdateImagesArticleID(ctx, tx, article.ID, article.ImageURLs); err != nil {
				return fmt.Errorf("关联图片记录失败: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return article, nil
}

// GetAdminDetail 后台回显：直接从仓库获取原始数据
func (s *articleServiceImpl) GetAdminDetail(ctx context.Context, id uint64) (*model.Article, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrArticleNotFound // 使用文件定义的错误变量
	}
	return article, nil
}

// GetPortalDetail 前台展示：增加状态校验逻辑
func (s *articleServiceImpl) GetPortalDetail(ctx context.Context, id uint64) (*model.Article, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		fmt.Println("这个问题在这里有问题", err)
		return nil, ErrArticleNotFound
	}

	// 核心区别：前台必须校验文章是否为发布状态
	// 根据你的 Create 逻辑，Status 1 为发布
	if article.Status != 1 {
		return nil, errors.New("该文章尚未发布或已被下架")
	}

	// 此处可扩展：例如增加阅读量统计逻辑
	return article, nil
}

func (s *articleServiceImpl) List(ctx context.Context, page, pageSize int, keyword string, status int) (*model.ArticleListResponse, error) {
	if pageSize > 100 {
		pageSize = 100
	}
	articles, total, err := s.repo.List(ctx, page, pageSize, keyword, status)
	if err != nil {
		return nil, err
	}

	return &model.ArticleListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Data:     articles,
	}, nil
}

func (s *articleServiceImpl) Update(ctx context.Context, id uint64, req model.UpdateArticleRequest) (*model.Article, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrArticleNotFound
	}

	if req.Title != "" {
		article.Title = strings.TrimSpace(req.Title)
	}
	if req.Content != "" {
		article.Content = strings.TrimSpace(req.Content)
	}
	if req.CategoryID != 0 {
		article.CategoryID = req.CategoryID
	}
	if req.Status == 0 || req.Status == 1 {
		article.Status = req.Status
	}
	article.ImageURLs = req.ImageURLs // 将自定义数组传入模型

	// 开启事务进行更新和关联
	err = s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		// 1. 获取更新前的所有关联图片记录
		oldImages, err := s.repo.GetImagesByArticleID(ctx, id)
		if err != nil {
			return fmt.Errorf("获取旧图片关联失败: %w", err)
		}

		// 2. 更新文章基本信息
		if err := s.repo.Update(ctx, tx, article); err != nil {
			return err
		}

		// 3. 计算需要解除关联并删除的图片（差集：oldImages - req.ImageURLs）
		newURLMap := make(map[string]bool)
		for _, url := range req.ImageURLs {
			newURLMap[url] = true
		}

		var toDeletePaths []string
		for _, img := range oldImages {
			if !newURLMap[img.URL] {
				toDeletePaths = append(toDeletePaths, img.URL)
			}
		}

		// 4. 执行物理删除和解绑
		if len(toDeletePaths) > 0 {
			// A. 物理删除文件
			basePath := getImagesPath()
			for _, url := range toDeletePaths {
				// 从 URL 提取文件名 (例如 /static/images/articles/xxx.jpg -> xxx.jpg)
				filename := filepath.Base(url)
				fullPath := filepath.Join(basePath, filename)
				if err := os.Remove(fullPath); err != nil {
					fmt.Printf("警告：物理删除文件失败 [%s]: %v\n", fullPath, err)
				}
			}
			// B. 数据库解绑 (article_id 置 0)
			if err := s.repo.UnbindImages(ctx, tx, toDeletePaths); err != nil {
				return fmt.Errorf("数据库解绑图片失败: %w", err)
			}
		}

		// 5. 绑定新图片 (article_id = 0 -> article_id = id)
		if len(req.ImageURLs) > 0 {
			if err := s.repo.UpdateImagesArticleID(ctx, tx, article.ID, req.ImageURLs); err != nil {
				return fmt.Errorf("关联新图片记录失败: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return article, nil
}

func (s *articleServiceImpl) Delete(ctx context.Context, id uint64) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrArticleNotFound
	}

	images, err := s.repo.GetImagesByArticleID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取文章关联图片失败: %w", err)
	}

	// 开启事务，确保文章和关联图片的删除原子化
	err = s.repo.Transaction(ctx, func(tx *gorm.DB) error {

		// 3. 删除数据库中的图片记录
		if err := s.repo.DeleteImagesByArticleID(ctx, id); err != nil {
			return fmt.Errorf("删除图片记录失败: %w", err)
		}
		//4. 删除评论
		if err := s.cop.Delete(ctx, id); err != nil {
			return err
		}
		// 5. 删除文章
		if err := s.repo.Delete(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		basePath := getImagesPath()
		for _, img := range images {
			// filepath.Base 提取文件名，防止路径穿越风险
			fullPath := filepath.Join(basePath, filepath.Base(img.URL))
			if err := os.Remove(fullPath); err != nil {
				// 物理删除失败仅记录警告，不回滚已提交的数据库事务
				if !os.IsNotExist(err) {
					fmt.Printf("警告：删除文章图片物理文件失败 [%s]: %v\n", fullPath, err)
				}
			}
		}
	}
	return err
}

// HandleImageUpload 处理图片上传的存储和记录逻辑
func (s *articleServiceImpl) HandleImageUpload(ctx context.Context, file *multipart.FileHeader, userID uint64, baseURL string) (string, error) {
	// 1. 业务层精确校验：限制 5MB
	const maxFileSize = 5 << 20
	if file.Size > maxFileSize {
		return "", fmt.Errorf("图片大小不能超过 5MB (当前为 %.2f MB)", float64(file.Size)/1024/1024)
	}

	// 2. 格式校验
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowExts[ext] {
		return "", fmt.Errorf("不支持的文件格式: %s", ext)
	}

	// 3. 存储准备
	articlesImagePath := getImagesPath() // 确保该函数返回 frontend/static/images/articles
	if err := os.MkdirAll(articlesImagePath, 0755); err != nil {
		return "", fmt.Errorf("创建存储目录失败: %w", err)
	}

	newFilename := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), userID, ext)
	filePath := filepath.Join(articlesImagePath, newFilename)

	// 4. 保存物理文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			fmt.Printf("关闭失败1")
		}
	}(src)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	// 执行写入并显式落盘
	if _, err := io.Copy(dst, src); err != nil {
		err := dst.Close()
		if err != nil {
			return "关闭失败2", err
		}
		return "写入失败", err
	}
	_ = dst.Sync() // 强制刷入磁盘
	err = dst.Close()
	if err != nil {
		return "关闭失败3", err
	} // 写入完成立即关闭，释放句柄

	// 5. 数据库记录入库
	urlPath := baseURL + "/static/images/articles/" + newFilename
	imgRecord := &model.Image{
		ArticleID: 0,
		UserID:    userID,
		URL:       urlPath,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateImage(ctx, imgRecord); err != nil {
		// 关键：如果数据库记录失败，立即清理刚存好的物理文件
		_ = os.Remove(filePath)
		return "", fmt.Errorf("图片记录入库失败，已清理残留文件: %w", err)
	}

	return urlPath, nil
}
