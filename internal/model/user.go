package model

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(32);uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(72);not null" json:"-"`
	GroupID   uint      `gorm:"default:101;index" json:"group_id"`
	Avatar    string    `gorm:"type:varchar(255);default:''" json:"avatar"`
	Threads   uint      `gorm:"default:0" json:"threads"`
	Posts     uint      `gorm:"default:0" json:"posts"`
	Credits   int       `gorm:"default:0" json:"credits"`
	Level     uint      `gorm:"default:1" json:"level"`
	LastLogin *time.Time `json:"last_login"`
	LastIP    string    `gorm:"type:varchar(45);default:''" json:"last_ip,omitempty"`
	Signature string    `gorm:"type:varchar(255);default:''" json:"signature"`
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
