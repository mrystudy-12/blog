package model

import (
	"gorm.io/gorm"
	"time"
)

// Article 代表文章模型
type Article struct {
	// ID: 自增主键
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// Title: 文章标题
	Title string `gorm:"size:255;not null" json:"title"`

	// Content: 文章正文
	Content string `gorm:"type:longtext;not null" json:"content"`

	// AuthorID:用户id
	AuthorID uint32 `gorm:"column:user_id;index" json:"user_id"`

	// Author: 作者信息
	Author User `gorm:"foreignKey:AuthorID" json:"author"`

	// CategoryID: 分类id
	CategoryID uint `gorm:"column:category_id;index" json:"category_id"`

	// Status: 文章状态 0:草稿/未发布，1:已发布
	Status int8 `gorm:"column:status;default:0" json:"status"`

	// ViewCount: 阅读量
	ViewCount uint `gorm:"column:view_count;default:0" json:"view_count"`

	// CreatedAt:创建文件时间
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt: 修改文件时间
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt: 软删除的标签
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// Images: 文章关联的图片列表
	Images []Image `gorm:"foreignKey:ArticleID" json:"images"`

	// ImageURLs: 自定义数组，用于接收前端传递的图片 URL 列表进行关联更新
	ImageURLs []string `gorm:"-" json:"image_urls"`
}

// Image 代表文章关联的图片模型
type Image struct {
	// ID: 自增主键 (bigint unsigned)
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// ArticleID: 关联文章ID (bigint unsigned)
	ArticleID uint64 `gorm:"column:article_id;index" json:"article_id"`

	// UserID: 上传用户ID (int unsigned)
	UserID uint32 `gorm:"column:user_id;index" json:"user_id"`

	// URL: 图片访问地址 (varchar(255))
	URL string `gorm:"size:255;not null" json:"url"`

	// CreatedAt: 创建时间
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// Result 统一响应体（成功和失败都用这个）
// 这样可以解决你感觉“响应体东西被删了”的问题，因为结构固定了
type Result struct {
	Code    int         `json:"code"`    // 业务状态码 (200:成功, 400:前端错误, 500:后端错误)
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 实际业务数据
}

// --- 请求参数部分 ---

// CreateArticleRequest 新增文章请求
type CreateArticleRequest struct {
	Title      string   `json:"title" binding:"required,min=1,max=255"`
	Content    string   `json:"content" binding:"required"`
	CategoryID uint     `json:"category_id"`
	Status     int8     `json:"status"`     // 0: 草稿 1: 发布
	ImageURLs  []string `json:"image_urls"` // 接收上传后的图片URL列表
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title      string   `json:"title" binding:"max=255"`
	Content    string   `json:"content"`
	CategoryID uint     `json:"category_id"`
	Status     int8     `json:"status"`     // 修改状态（发布或改回草稿）
	ImageURLs  []string `json:"image_urls"` // 更新后的图片URL列表
}

// ArticleListResponse 文章列表分页数据
type ArticleListResponse struct {
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	Keyword  string    `json:"keyword"` // 新增：回传搜索词，方便前端显示搜索状态
	Data     []Article `json:"data"`
}

// TableName 指定映射的数据库表名
func (Image) TableName() string {
	return "images"
}
