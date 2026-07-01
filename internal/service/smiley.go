package service

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type SmileyService interface {
	List() ([]model.Smiley, error)
}

type smileyService struct {
	repo repository.SmileyRepo
}

func NewSmileyService(repo repository.SmileyRepo) SmileyService {
	return &smileyService{repo: repo}
}

func (s *smileyService) List() ([]model.Smiley, error) {
	return s.repo.List()
}
