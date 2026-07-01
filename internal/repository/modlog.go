package repository

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type ModLogRepo struct {
	db *gorm.DB
}

func NewModLogRepo(db *gorm.DB) *ModLogRepo {
	return &ModLogRepo{db: db}
}

func (r *ModLogRepo) Create(userID uint64, action, detail string, threadID uint64) error {
	return r.db.Create(&model.ModLog{
		UserID:   userID,
		Action:   action,
		Detail:   detail,
		ThreadID: threadID,
	}).Error
}

func (r *ModLogRepo) CreateBatch(userID uint64, threadIDs []uint64, action string) error {
	var logs []model.ModLog
	for _, tid := range threadIDs {
		logs = append(logs, model.ModLog{
			UserID:   userID,
			Action:   action,
			ThreadID: tid,
		})
	}
	return r.db.Create(&logs).Error
}

func (r *ModLogRepo) List(page, pageSize int) ([]model.ModLog, int64, error) {
	var logs []model.ModLog
	var total int64
	r.db.Model(&model.ModLog{}).Count(&total)
	err := r.db.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&logs).Error
	return logs, total, err
}
