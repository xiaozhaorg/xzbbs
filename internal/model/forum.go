package model

import "time"

type Forum struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(64);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"type:varchar(255);default:''" json:"icon"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	Threads     uint      `gorm:"default:0" json:"threads"`
	Posts       uint      `gorm:"default:0" json:"posts"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Forum) TableName() string {
	return "forums"
}

type ForumPermission struct {
	ForumID     uint `gorm:"primaryKey" json:"forum_id"`
	GroupID     uint `gorm:"primaryKey" json:"group_id"`
	AllowRead   bool `gorm:"default:true" json:"allow_read"`
	AllowThread bool `gorm:"default:false" json:"allow_thread"`
	AllowPost   bool `gorm:"default:false" json:"allow_post"`
	AllowAttach bool `gorm:"default:false" json:"allow_attach"`
}

func (ForumPermission) TableName() string {
	return "forum_permissions"
}
