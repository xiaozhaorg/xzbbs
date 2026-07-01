package repository

import (
	"time"
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type PMRepo interface {
	Send(senderID, receiverID uint64, content string) (*model.PrivateMessage, error)
	GetConversation(userID, otherID uint64, page, pageSize int) ([]model.PrivateMessage, int64, error)
	GetConversations(userID uint64) ([]map[string]interface{}, error)
	MarkRead(userID, otherID uint64) error
	GetUnreadCount(userID uint64) (int64, error)
	GetByID(id uint64) (*model.PrivateMessage, error)
	SoftDeleteSender(userID, msgID uint64) error
	SoftDeleteReceiver(userID, msgID uint64) error
}

type pmRepo struct {
	db *gorm.DB
}

func NewPMRepo(db *gorm.DB) PMRepo {
	return &pmRepo{db: db}
}

func (r *pmRepo) Send(senderID, receiverID uint64, content string) (*model.PrivateMessage, error) {
	pm := &model.PrivateMessage{
		SenderID: senderID, ReceiverID: receiverID,
		Content: content, CreatedAt: time.Now(),
	}
	if err := r.db.Create(pm).Error; err != nil {
		return nil, err
	}
	return pm, nil
}

func (r *pmRepo) GetConversation(userID, otherID uint64, page, pageSize int) ([]model.PrivateMessage, int64, error) {
	var total int64
	var items []model.PrivateMessage
	q := r.db.Model(&model.PrivateMessage{}).Where(
		"((sender_id = ? AND receiver_id = ? AND is_deleted = ?) OR (sender_id = ? AND receiver_id = ? AND is_received_deleted = ?))",
		userID, otherID, false, otherID, userID, false,
	)
	q.Count(&total)
	q.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items)
	return items, total, q.Error
}

func (r *pmRepo) GetConversations(userID uint64) ([]map[string]interface{}, error) {
	// Find the most recent message for each conversation partner
	var results []map[string]interface{}
	r.db.Model(&model.PrivateMessage{}).
		Select("CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as other_id", userID).
		Where("((sender_id = ? AND is_deleted = ?) OR (receiver_id = ? AND is_received_deleted = ?))", userID, false, userID, false).
		Group("other_id").
		Order("MAX(created_at) DESC").
		Find(&results)
	return results, nil
}

func (r *pmRepo) MarkRead(userID, otherID uint64) error {
	return r.db.Model(&model.PrivateMessage{}).Where(
		"sender_id = ? AND receiver_id = ? AND is_read = ?", otherID, userID, false,
	).Update("is_read", true).Error
}

func (r *pmRepo) GetUnreadCount(userID uint64) (int64, error) {
	var count int64
	return count, r.db.Model(&model.PrivateMessage{}).Where(
		"receiver_id = ? AND is_read = ? AND is_received_deleted = ?", userID, false, false,
	).Count(&count).Error
}

func (r *pmRepo) SoftDeleteSender(userID, msgID uint64) error {
	return r.db.Model(&model.PrivateMessage{}).Where("id = ? AND sender_id = ?", msgID, userID).Update("is_deleted", true).Error
}

func (r *pmRepo) SoftDeleteReceiver(userID, msgID uint64) error {
	return r.db.Model(&model.PrivateMessage{}).Where("id = ? AND receiver_id = ?", msgID, userID).Update("is_received_deleted", true).Error
}

func (r *pmRepo) GetByID(id uint64) (*model.PrivateMessage, error) {
	var msg model.PrivateMessage
	err := r.db.First(&msg, id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
