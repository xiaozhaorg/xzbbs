package model

import "time"

type Attachment struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID       uint64    `gorm:"index;not null" json:"post_id"`
	UserID       uint64    `gorm:"index;not null" json:"user_id"`
	Filename     string    `gorm:"type:varchar(255);not null" json:"filename"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"original_name"`
	FileSize     uint      `gorm:"default:0" json:"file_size"`
	MimeType     string    `gorm:"type:varchar(64);default:''" json:"mime_type"`
	IsImage      bool      `gorm:"default:false" json:"is_image"`
	Width        uint      `gorm:"default:0" json:"width"`
	Height       uint      `gorm:"default:0" json:"height"`
	Downloads    uint      `gorm:"default:0" json:"downloads"`
	CreatedAt    time.Time `json:"created_at"`
}

func (Attachment) TableName() string {
	return "attachments"
}
