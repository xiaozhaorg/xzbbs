package service

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type SearchService interface {
	SearchThreads(query string, forumID uint, page, pageSize int) ([]model.Thread, int64, error)
}

type searchService struct {
	repo repository.SearchRepo
}

func NewSearchService(repo repository.SearchRepo) SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) SearchThreads(query string, forumID uint, page, pageSize int) ([]model.Thread, int64, error) {
	return s.repo.SearchThreads(query, forumID, page, pageSize)
}
