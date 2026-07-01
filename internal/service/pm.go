package service

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type PMService interface {
	Send(senderID, receiverID uint64, content string) (*model.PrivateMessage, error)
	GetConversation(userID, otherID uint64, page, pageSize int) ([]model.PrivateMessage, int64, error)
	GetConversationList(userID uint64) ([]map[string]interface{}, error)
	MarkRead(userID, otherID uint64) error
	GetUnreadCount(userID uint64) (int64, error)
	SoftDelete(userID, msgID uint64) error
}

type pmService struct {
	repo    repository.PMRepo
}

func NewPMService(repo repository.PMRepo) PMService {
	return &pmService{repo: repo}
}

func (s *pmService) Send(senderID, receiverID uint64, content string) (*model.PrivateMessage, error) {
	return s.repo.Send(senderID, receiverID, content)
}

func (s *pmService) GetConversation(userID, otherID uint64, page, pageSize int) ([]model.PrivateMessage, int64, error) {
	return s.repo.GetConversation(userID, otherID, page, pageSize)
}

func (s *pmService) GetConversationList(userID uint64) ([]map[string]interface{}, error) {
	return s.repo.GetConversations(userID)
}

func (s *pmService) MarkRead(userID, otherID uint64) error {
	return s.repo.MarkRead(userID, otherID)
}

func (s *pmService) GetUnreadCount(userID uint64) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

func (s *pmService) SoftDelete(userID, msgID uint64) error {
	// Check if user is sender or receiver of this message
	msg, err := s.repo.GetByID(msgID)
	if err != nil {
		return err
	}
	if msg.SenderID == userID {
		return s.repo.SoftDeleteSender(userID, msgID)
	}
	if msg.ReceiverID == userID {
		return s.repo.SoftDeleteReceiver(userID, msgID)
	}
	return nil
}
