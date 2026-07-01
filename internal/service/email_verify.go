package service

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

func generateRandomToken(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	return string(b)
}

type EmailVerifyService struct {
	repo    *repository.EmailVerifyRepo
	userSvc *UserService
}

func NewEmailVerifyService(repo *repository.EmailVerifyRepo, userSvc *UserService) *EmailVerifyService {
	return &EmailVerifyService{repo: repo, userSvc: userSvc}
}

func (s *EmailVerifyService) CreateToken(userID uint64, email string) (string, time.Time, error) {
	token := generateRandomToken(32)
	expiresAt := time.Now().Add(24 * time.Hour)
	err := s.repo.Create(userID, email, token, expiresAt)
	return token, expiresAt, err
}

func (s *EmailVerifyService) ConfirmToken(token string) error {
	ev, err := s.repo.FindByToken(token)
	if err != nil {
		return err
	}
	if time.Now().After(ev.ExpiresAt) {
		return err
	}
	if err := s.repo.MarkUsed(ev.ID); err != nil {
		return err
	}
	return s.userSvc.Update(ev.UserID, map[string]interface{}{"email_verified": true})
}

func (s *EmailVerifyService) GetUser(userID uint64) (*model.User, error) {
	return s.userSvc.GetByID(userID)
}

func (s *EmailVerifyService) IsEmailVerified(userID uint64) (bool, error) {
	user, err := s.userSvc.GetByID(userID)
	if err != nil {
		return false, err
	}
	return user.EmailVerified, nil
}
