package repository

import (
	"time"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type PasswordResetRepo struct {
	db *gorm.DB
}

func NewPasswordResetRepo(db *gorm.DB) *PasswordResetRepo {
	return &PasswordResetRepo{db: db}
}

func (r *PasswordResetRepo) Create(userID uint64, email, token string, expiresAt time.Time) error {
	return r.db.Create(&model.PasswordReset{
		UserID: userID, Email: email, Token: token, ExpiresAt: expiresAt,
	}).Error
}

func (r *PasswordResetRepo) FindByToken(token string) (*model.PasswordReset, error) {
	var pr model.PasswordReset
	err := r.db.Where("token = ? AND used = ?", token, false).First(&pr).Error
	return &pr, err
}

func (r *PasswordResetRepo) MarkUsed(id uint64) error {
	return r.db.Model(&model.PasswordReset{}).Where("id = ?", id).Update("used", true).Error
}
