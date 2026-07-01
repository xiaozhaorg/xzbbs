package service

import (
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type ModLogService struct {
	repo *repository.ModLogRepo
}

func NewModLogService(repo *repository.ModLogRepo) *ModLogService {
	return &ModLogService{repo: repo}
}

func (s *ModLogService) Create(userID uint64, action, detail string, threadID uint64) error {
	return s.repo.Create(userID, action, detail, threadID)
}

func (s *ModLogService) CreateBatch(userID uint64, threadIDs []uint64, action string) error {
	return s.repo.CreateBatch(userID, threadIDs, action)
}
