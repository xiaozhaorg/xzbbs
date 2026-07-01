package model

import (
	"time"

	"gorm.io/gorm"
)

// OnlineUser 在线用户（临时表，定期清理）
type OnlineUser struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"uniqueIndex;not null" json:"user_id"`
	Username  string    `gorm:"type:varchar(32);not null" json:"username"`
	LastActive time.Time `gorm:"index;not null" json:"last_active"`
	IP        string    `gorm:"type:varchar(45);not null" json:"ip"`
}

func (OnlineUser) TableName() string {
	return "online_users"
}

// OnlineUserWithInfo 带统计信息的在线用户
type OnlineUserWithInfo struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	Username   string    `json:"username"`
	Avatar     string    `json:"avatar"`
	LastActive time.Time `json:"last_active"`
	IP         string    `json:"ip"`
}

// CleanupExpiredOnlineUsers 清理过期在线用户
func CleanupExpiredOnlineUsers(db *gorm.DB) error {
	return db.Where("last_active < ?", time.Now().Add(-30*time.Minute)).Delete(&OnlineUser{}).Error
}
