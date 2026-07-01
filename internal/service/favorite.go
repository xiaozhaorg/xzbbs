package service

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type FavoriteService interface {
	Toggle(userID, threadID uint64) (bool, error)
	IsFavorited(userID, threadID uint64) (bool, error)
	ListByUser(userID uint64, page, pageSize int) ([]model.Thread, int, error)
	CountByThread(threadID uint64) (int64, error)
}

type favoriteService struct {
	favRepo    repository.FavoriteRepo
	threadRepo *repository.ThreadRepo
}

func NewFavoriteService(favRepo repository.FavoriteRepo, threadRepo *repository.ThreadRepo) FavoriteService {
	return &favoriteService{favRepo: favRepo, threadRepo: threadRepo}
}

func (s *favoriteService) Toggle(userID, threadID uint64) (bool, error) {
	return s.favRepo.Toggle(userID, threadID)
}

func (s *favoriteService) IsFavorited(userID, threadID uint64) (bool, error) {
	return s.favRepo.IsFavorited(userID, threadID)
}

func (s *favoriteService) ListByUser(userID uint64, page, pageSize int) ([]model.Thread, int, error) {
	favs, total, err := s.favRepo.ListByUser(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	threadIDs := make([]uint64, len(favs))
	for i, f := range favs {
		threadIDs[i] = f.ThreadID
	}
	threads, _ := s.threadRepo.FindByIDs(threadIDs)
	return threads, int(total), nil
}

func (s *favoriteService) CountByThread(threadID uint64) (int64, error) {
	return s.favRepo.CountByThread(threadID)
}
