package service

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type ForumService struct {
	repo *repository.ForumRepo
}

func NewForumService(repo *repository.ForumRepo) *ForumService {
	return &ForumService{repo: repo}
}

func (s *ForumService) Create(name, description string, sortOrder int) (*model.Forum, error) {
	forum := &model.Forum{
		Name:        name,
		Description: description,
		SortOrder:   sortOrder,
	}
	if err := s.repo.Create(forum); err != nil {
		return nil, err
	}
	return forum, nil
}

func (s *ForumService) GetByID(id uint) (*model.Forum, error) {
	return s.repo.GetByID(id)
}

func (s *ForumService) Update(id uint, updates map[string]interface{}) error {
	return s.repo.Update(id, updates)
}

func (s *ForumService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ForumService) List() ([]model.Forum, error) {
	return s.repo.List()
}
