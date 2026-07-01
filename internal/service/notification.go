package service

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type NotificationService interface {
	Create(userID uint64, notifType uint8, actorID, threadID, postID uint64, message string) error
	List(userID uint64, unreadOnly bool, page, pageSize int) ([]model.Notification, int64, error)
	UnreadCount(userID uint64) (int64, error)
	MarkRead(userID uint64, ids []uint64) error
	MarkAllRead(userID uint64) error
	Delete(userID uint64, ids []uint64) error
	DeleteAllByThread(threadID uint64) error
}

type notificationService struct {
	repo repository.NotificationRepo
}

func NewNotificationService(repo repository.NotificationRepo) NotificationService {
	return &notificationService{repo: repo}
}

func (s *notificationService) Create(userID uint64, notifType uint8, actorID, threadID, postID uint64, message string) error {
	return s.repo.Create(&model.Notification{
		UserID: userID, Type: notifType, ActorID: actorID,
		ThreadID: threadID, PostID: postID, Message: message,
		CreatedAt: model.Now(),
	})
}

func (s *notificationService) List(userID uint64, unreadOnly bool, page, pageSize int) ([]model.Notification, int64, error) {
	return s.repo.GetUserNotifications(userID, unreadOnly, page, pageSize)
}

func (s *notificationService) UnreadCount(userID uint64) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

func (s *notificationService) MarkRead(userID uint64, ids []uint64) error {
	return s.repo.MarkRead(userID, ids)
}

func (s *notificationService) MarkAllRead(userID uint64) error {
	return s.repo.MarkAllRead(userID)
}

func (s *notificationService) Delete(userID uint64, ids []uint64) error {
	return s.repo.Delete(userID, ids)
}

func (s *notificationService) DeleteAllByThread(threadID uint64) error {
	return s.repo.DeleteAllByThread(threadID)
}
