package model

import "time"

type ModLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	ThreadID  uint64    `gorm:"index;default:0" json:"thread_id"`
	Action    string    `gorm:"type:varchar(32);not null" json:"action"`
	Detail    string    `gorm:"type:varchar(255);default:''" json:"detail"`
	CreatedAt time.Time `json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ModLog) TableName() string {
	return "mod_logs"
}
