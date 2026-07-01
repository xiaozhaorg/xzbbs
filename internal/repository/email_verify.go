package repository

import (
	"time"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type EmailVerifyRepo struct {
	db *gorm.DB
}

func NewEmailVerifyRepo(db *gorm.DB) *EmailVerifyRepo {
	return &EmailVerifyRepo{db: db}
}

func (r *EmailVerifyRepo) Create(userID uint64, email, token string, expiresAt time.Time) error {
	return r.db.Create(&model.EmailVerify{
		UserID: userID, Email: email, Token: token, ExpiresAt: expiresAt,
	}).Error
}

func (r *EmailVerifyRepo) FindByToken(token string) (*model.EmailVerify, error) {
	var ev model.EmailVerify
	err := r.db.Where("token = ? AND used = ?", token, false).First(&ev).Error
	return &ev, err
}

func (r *EmailVerifyRepo) MarkUsed(id uint64) error {
	return r.db.Model(&model.EmailVerify{}).Where("id = ?", id).Update("used", true).Error
}
