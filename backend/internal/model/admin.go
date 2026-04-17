package model

// DashboardStats 后台仪表盘数据结构
type DashboardStats struct {
	ArticleCount int64  `json:"article_count"`
	UserCount    int64  `json:"user_count"`
	CommentCount int64  `json:"comment_count"`
	ViewCount    uint64 `json:"view_count"`
}
