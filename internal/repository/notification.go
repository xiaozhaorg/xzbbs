package repository

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type NotificationRepo interface {
	Create(n *model.Notification) error
	GetUserNotifications(userID uint64, unreadOnly bool, page, pageSize int) ([]model.Notification, int64, error)
	GetUnreadCount(userID uint64) (int64, error)
	MarkRead(userID uint64, ids []uint64) error
	MarkAllRead(userID uint64) error
	Delete(userID uint64, ids []uint64) error
	DeleteByThread(userID, threadID uint64) error
	DeleteAllByThread(threadID uint64) error
}

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) NotificationRepo {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}

func (r *notificationRepo) GetUserNotifications(userID uint64, unreadOnly bool, page, pageSize int) ([]model.Notification, int64, error) {
	var total int64
	var items []model.Notification
	q := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	if unreadOnly {
		q = q.Where("is_read = ?", false)
	}
	q.Count(&total)
	q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	return items, total, q.Error
}

func (r *notificationRepo) GetUnreadCount(userID uint64) (int64, error) {
	var count int64
	return count, r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
}

func (r *notificationRepo) MarkRead(userID uint64, ids []uint64) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND id IN ?", userID, ids).Update("is_read", true).Error
}

func (r *notificationRepo) MarkAllRead(userID uint64) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *notificationRepo) Delete(userID uint64, ids []uint64) error {
	return r.db.Where("user_id = ? AND id IN ?", userID, ids).Delete(&model.Notification{}).Error
}

func (r *notificationRepo) DeleteByThread(userID, threadID uint64) error {
	return r.db.Where("user_id = ? AND thread_id = ?", userID, threadID).Delete(&model.Notification{}).Error
}

func (r *notificationRepo) DeleteAllByThread(threadID uint64) error {
	return r.db.Where("thread_id = ?", threadID).Delete(&model.Notification{}).Error
}
