package repository

import (
	"time"
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type OnlineRepo interface {
	Upsert(userID uint64, username, ip string) error
	List(limit int) ([]model.OnlineUserWithInfo, error)
	Cleanup() error
}

type onlineRepo struct {
	db *gorm.DB
}

func NewOnlineRepo(db *gorm.DB) OnlineRepo {
	return &onlineRepo{db: db}
}

func (r *onlineRepo) Upsert(userID uint64, username, ip string) error {
	now := time.Now()
	return r.db.Transaction(func(tx *gorm.DB) error {
		var ou model.OnlineUser
		err := tx.Where("user_id = ?", userID).First(&ou).Error
		if err == gorm.ErrRecordNotFound {
			return tx.Create(&model.OnlineUser{
				UserID: userID, Username: username, IP: ip, LastActive: now,
			}).Error
		}
		if err != nil {
			return err
		}
		return tx.Model(&ou).Updates(map[string]interface{}{
			"username": username, "ip": ip, "last_active": now,
		}).Error
	})
}

func (r *onlineRepo) List(limit int) ([]model.OnlineUserWithInfo, error) {
	var list []model.OnlineUser
	err := r.db.Order("last_active DESC").Limit(limit).Find(&list).Error
	if err != nil {
		return nil, err
	}
	result := make([]model.OnlineUserWithInfo, len(list))
	for i, ou := range list {
		result[i] = model.OnlineUserWithInfo{
			ID: ou.ID, UserID: ou.UserID, Username: ou.Username,
			LastActive: ou.LastActive, IP: ou.IP,
		}
	}
	return result, nil
}

func (r *onlineRepo) Cleanup() error {
	return model.CleanupExpiredOnlineUsers(r.db)
}
