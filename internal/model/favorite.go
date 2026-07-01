package model

import "time"

// Favorite 用户收藏/书签
type Favorite struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index:idx_user_thread;not null" json:"user_id"`
	ThreadID  uint64    `gorm:"index:idx_user_thread;not null" json:"thread_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Favorite) TableName() string {
	return "favorites"
}
