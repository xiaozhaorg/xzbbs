package model

import "time"

// PostEdit 帖子编辑历史
type PostEdit struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PostID    uint64    `gorm:"index;not null" json:"post_id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`       // 编辑者
	OldContent string   `gorm:"type:text;not null" json:"old_content"`
	NewContent string   `gorm:"type:text;not null" json:"new_content"`
	CreatedAt time.Time `json:"created_at"`
}

func (PostEdit) TableName() string {
	return "post_edits"
}
