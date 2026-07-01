package model

import "time"

type Post struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ThreadID    uint64    `gorm:"index:idx_thread;not null" json:"thread_id"`
	UserID      uint64    `gorm:"index;not null" json:"user_id"`
	IsFirst     bool      `gorm:"default:false" json:"is_first"`
	Content     string    `gorm:"type:longtext;not null" json:"content"`
	ContentType uint8     `gorm:"default:0" json:"content_type"` // 0=markdown, 1=html
	IP          string    `gorm:"type:varchar(45);default:''" json:"-"`
	ReplyTo     uint64    `gorm:"default:0" json:"reply_to"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Post) TableName() string {
	return "posts"
}
