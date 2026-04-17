package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint64         `gorm:"column:article_id;index;not null" json:"article_id"`
	UserID    uint32         `gorm:"column:user_id;index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ParentID  uint64         `gorm:"column:parent_id;default:0;index" json:"parent_id"`
	Status    int8           `gorm:"column:status;default:1" json:"status"`
	IsDeleted int8           `gorm:"column:is_deleted;default:0" json:"is_deleted"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateCommentRequest 提交评论
type CreateCommentRequest struct {
	ArticleID uint64 `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required,max=500"`
	ParentID  uint64 `json:"parent_id"`
}

// TableName 指定表名为 comments
func (Comment) TableName() string {
	return "comments"
}
