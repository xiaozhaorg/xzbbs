package repository

import (
	"time"
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type IPBanRepo interface {
	Create(ip, reason string, createdBy uint64, expireAt *time.Time) error
	List(page, pageSize int) ([]model.IPBan, int64, error)
	Delete(id uint64) error
	Check(ip string) (bool, error)
}

type ipBanRepo struct {
	db *gorm.DB
}

func NewIPBanRepo(db *gorm.DB) IPBanRepo {
	return &ipBanRepo{db: db}
}

func (r *ipBanRepo) Create(ip, reason string, createdBy uint64, expireAt *time.Time) error {
	ban := &model.IPBan{
		IP: ip, Reason: reason, CreatedBy: createdBy,
		ExpireAt: expireAt, CreatedAt: model.Now(),
	}
	return r.db.Create(ban).Error
}

func (r *ipBanRepo) List(page, pageSize int) ([]model.IPBan, int64, error) {
	var items []model.IPBan
	var total int64
	r.db.Model(&model.IPBan{}).Count(&total)
	err := r.db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *ipBanRepo) Delete(id uint64) error {
	return r.db.Delete(&model.IPBan{}, id).Error
}

func (r *ipBanRepo) Check(ip string) (bool, error) {
	var count int64
	err := r.db.Model(&model.IPBan{}).Where("ip = ? AND (expire_at IS NULL OR expire_at > ?)", ip, model.Now()).Count(&count).Error
	return count > 0, err
}
