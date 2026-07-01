package model

import "time"

// PrivateMessage 私信
type PrivateMessage struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID   uint64    `gorm:"index:idx_pm_conversation;not null" json:"sender_id"`
	ReceiverID uint64    `gorm:"index:idx_pm_conversation;not null" json:"receiver_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsRead     bool      `gorm:"default:false;index" json:"is_read"`
	IsDeleted  bool      `gorm:"default:false" json:"is_deleted"`       // sender soft delete
	IsReceivedDeleted bool `gorm:"default:false" json:"is_received_deleted"` // receiver soft delete
	CreatedAt  time.Time `gorm:"index;not null" json:"created_at"`
}

func (PrivateMessage) TableName() string {
	return "private_messages"
}

// IPBan IP封禁
type IPBan struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	IP        string    `gorm:"type:varchar(45);uniqueIndex;not null" json:"ip"`
	Reason    string    `gorm:"type:varchar(255);not null" json:"reason"`
	ExpireAt  *time.Time `gorm:"index" json:"expire_at"`       // nil = permanent
	CreatedBy uint64    `gorm:"not null" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (IPBan) TableName() string {
	return "ip_bans"
}

// Smiley 表情
type Smiley struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Code     string `gorm:"type:varchar(32);uniqueIndex;not null" json:"code"`
	Image    string `gorm:"type:varchar(255);not null" json:"image"`
	Sort     int    `gorm:"default:0" json:"sort"`
}

func (Smiley) TableName() string {
	return "smilies"
}

// EmailVerify 邮箱验证
type EmailVerify struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Email     string    `gorm:"type:varchar(128);not null" json:"email"`
	Token     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"`
	Used      bool      `gorm:"default:false" json:"used"`
	ExpiresAt time.Time `gorm:"index;not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (EmailVerify) TableName() string {
	return "email_verifies"
}
