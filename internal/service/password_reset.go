package service

import (
	"errors"
	"time"

	"github.com/xiaozhaorg/xzbbs/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetService struct {
	repo    *repository.PasswordResetRepo
	userSvc *UserService
}

func NewPasswordResetService(repo *repository.PasswordResetRepo, userSvc *UserService) *PasswordResetService {
	return &PasswordResetService{repo: repo, userSvc: userSvc}
}

func (s *PasswordResetService) RequestReset(email string) (string, error) {
	user, err := s.userSvc.GetByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	token := generateRandomToken(48)
	expiresAt := time.Now().Add(30 * time.Minute)
	err = s.repo.Create(user.ID, email, token, expiresAt)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *PasswordResetService) ConfirmReset(token, newPassword string) error {
	pr, err := s.repo.FindByToken(token)
	if err != nil {
		return errors.New("invalid or expired token")
	}
	if time.Now().After(pr.ExpiresAt) {
		return errors.New("token expired")
	}

	if err := ValidatePassword(newPassword); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.userSvc.Update(pr.UserID, map[string]interface{}{"password": string(hash)}); err != nil {
		return err
	}

	return s.repo.MarkUsed(pr.ID)
}
