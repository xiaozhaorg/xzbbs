package model

import "time"

type PasswordReset struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Email     string    `gorm:"type:varchar(255);not null" json:"email"`
	Token     string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"token"`
	Used      bool      `gorm:"default:false" json:"used"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}
