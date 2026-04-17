package model

// Categories 分类模型
type Categories struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"size:50;not null" json:"name"`
	Sort        int    `gorm:"default:0" json:"sort"`       // 对应图片中的 sort
	Description string `gorm:"size:255" json:"description"` // 对应图片中的 description
}

// CreateCategoryRequest 创建分类
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1,max=50"`
	Sort int    `json:"sort"`
}

// UpdateCategoryRequest 更新分类
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"max=50"`
	Sort int    `json:"sort"`
}
