package service

import (
	"fmt"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
	"gorm.io/gorm"
)

// CreditService handles credit earning and spending
type CreditService struct {
	repo     repository.CreditRepo
	userRepo *repository.UserRepo
}

func NewCreditService(cr repository.CreditRepo, ur *repository.UserRepo) *CreditService {
	return &CreditService{repo: cr, userRepo: ur}
}

// CreditRules defines how many credits to award per action
var CreditRules = map[string]int{
	"new_thread":  5,
	"new_post":    2,
	"daily_login": 1,
}

// Earn adds credits to a user and logs the transaction
func (s *CreditService) Earn(userID uint64, reason string, amount int, relatedID uint64) error {
	if amount <= 0 {
		return nil
	}
	if err := s.repo.Log(userID, amount, reason, relatedID); err != nil {
		return err
	}
	return s.userRepo.DB().Model(&model.User{}).Where("id = ?", userID).
		Update("credits", gorm.Expr("credits + ?", amount)).Error
}

// Spend deducts credits from a user
func (s *CreditService) Spend(userID uint64, reason string, amount int, relatedID uint64) error {
	if amount <= 0 {
		return nil
	}
	db := s.userRepo.DB()
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.User{}).
			Where("id = ? AND credits >= ?", userID, amount).
			Update("credits", gorm.Expr("credits - ?", amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			var credits int
			tx.Model(&model.User{}).Where("id = ?", userID).Pluck("credits", &credits)
			return fmt.Errorf("insufficient credits: have %d, need %d", credits, amount)
		}
		return s.repo.LogWithTx(tx, userID, -amount, reason, relatedID)
	})
}

// AwardThreadCreation credits for creating a thread
func (s *CreditService) AwardThreadCreation(userID uint64, threadID uint64) error {
	return s.Earn(userID, "new_thread", CreditRules["new_thread"], threadID)
}

// AwardPostCreation credits for creating a post
func (s *CreditService) AwardPostCreation(userID uint64, postID uint64) error {
	return s.Earn(userID, "new_post", CreditRules["new_post"], postID)
}

// UpdateLevel recalculates user level based on credits
func (s *CreditService) UpdateLevel(userID uint64) error {
	var credits int
	if err := s.userRepo.DB().Model(&model.User{}).Where("id = ?", userID).Pluck("credits", &credits).Error; err != nil {
		return err
	}

	var level uint = 1
	switch {
	case credits >= 1000:
		level = 10
	case credits >= 500:
		level = 9
	case credits >= 200:
		level = 8
	case credits >= 100:
		level = 7
	case credits >= 50:
		level = 6
	case credits >= 20:
		level = 5
	case credits >= 10:
		level = 4
	case credits >= 5:
		level = 3
	}
	return s.userRepo.DB().Model(&model.User{}).Where("id = ?", userID).Update("level", level).Error
}

func (s *CreditService) ListLogs(userID uint64, page, pageSize int) ([]model.CreditLog, int64, error) {
	return s.repo.ListByUser(userID, page, pageSize)
}
