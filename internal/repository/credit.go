package repository

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type CreditRepo interface {
	Log(userID uint64, amount int, reason string, relatedID uint64) error
	LogWithTx(tx *gorm.DB, userID uint64, amount int, reason string, relatedID uint64) error
	ListByUser(userID uint64, page, pageSize int) ([]model.CreditLog, int64, error)
}

type creditRepo struct {
	db *gorm.DB
}

func NewCreditRepo(db *gorm.DB) CreditRepo {
	return &creditRepo{db: db}
}

func (r *creditRepo) Log(userID uint64, amount int, reason string, relatedID uint64) error {
	return r.db.Create(&model.CreditLog{
		UserID: userID, Amount: amount, Reason: reason, RelatedID: relatedID,
	}).Error
}

func (r *creditRepo) LogWithTx(tx *gorm.DB, userID uint64, amount int, reason string, relatedID uint64) error {
	return tx.Create(&model.CreditLog{
		UserID: userID, Amount: amount, Reason: reason, RelatedID: relatedID,
	}).Error
}

func (r *creditRepo) ListByUser(userID uint64, page, pageSize int) ([]model.CreditLog, int64, error) {
	var logs []model.CreditLog
	var total int64
	r.db.Model(&model.CreditLog{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("id DESC").Offset((page-1)*pageSize).Limit(pageSize).
		Find(&logs).Error
	return logs, total, err
}
