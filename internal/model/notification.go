package model

import "time"

// Notification 用户通知
type Notification struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Type      uint8     `gorm:"not null;default:0" json:"type"`       // 0=reply, 1=at, 2=system
	ActorID   uint64    `gorm:"not null;default:0" json:"actor_id"`   // 触发者UID, 0=系统
	ThreadID  uint64    `gorm:"index;not null;default:0" json:"thread_id"`
	PostID    uint64    `gorm:"index;not null;default:0" json:"post_id"`
	Message   string    `gorm:"type:varchar(200);not null" json:"message"`
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time `gorm:"index;not null" json:"created_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
