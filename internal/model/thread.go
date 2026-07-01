package model

import "time"

type Thread struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ForumID      uint      `gorm:"index:idx_forum_last;index:idx_forum_created;not null" json:"forum_id"`
	UserID       uint64    `gorm:"index;not null" json:"user_id"`
	Title        string    `gorm:"type:varchar(200);not null" json:"title"`
	IsTop        uint8     `gorm:"default:0" json:"is_top"`       // 0=normal, 1=forum top, 2=global top
	IsClosed     bool      `gorm:"default:false" json:"is_closed"`
	Views        uint      `gorm:"default:0" json:"views"`
	Posts        uint      `gorm:"default:0" json:"posts"`        // reply count (excluding first post)
	LastReplyAt  *time.Time `gorm:"index:idx_forum_last" json:"last_reply_at"`
	LastReplyUID uint64    `gorm:"default:0" json:"last_reply_uid"`
	CreatedAt    time.Time `gorm:"index:idx_forum_created" json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations (not stored, for preload)
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Forum *Forum `gorm:"foreignKey:ForumID" json:"forum,omitempty"`
}

func (Thread) TableName() string {
	return "threads"
}
