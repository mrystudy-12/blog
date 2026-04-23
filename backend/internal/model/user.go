package model

// User 代表用户模型，简化后的 4 个核心字段
type User struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password  string `gorm:"size:255;not null" json:"-"`
	Email     string `gorm:"size:255;" json:"email"`
	AvatarUrl string `gorm:"size:100;" json:"avatar_url"`
	Role      string `gorm:"type:enum('admin', 'user');default:'user'" json:"role"`
}
