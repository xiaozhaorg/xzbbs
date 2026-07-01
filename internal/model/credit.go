package model

import "time"

// CreditLog records credit transactions
type CreditLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"index;not null" json:"user_id"`
	Amount    int       `gorm:"not null" json:"amount"`        // positive=earn, negative=spend
	Reason    string    `gorm:"type:varchar(100);not null" json:"reason"`
	RelatedID uint64    `gorm:"default:0" json:"related_id"`   // thread/post id
	CreatedAt time.Time `json:"created_at"`
}

func (CreditLog) TableName() string {
	return "credit_logs"
}
