package service

import (
	"time"
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type IPBanService interface {
	Ban(ip, reason string, createdBy uint64, expireAt *time.Time) error
	Unban(id uint64) error
	List(page, pageSize int) ([]model.IPBan, int64, error)
	Check(ip string) (bool, error)
}

type ipBanService struct {
	repo repository.IPBanRepo
}

func NewIPBanService(repo repository.IPBanRepo) IPBanService {
	return &ipBanService{repo: repo}
}

func (s *ipBanService) Ban(ip, reason string, createdBy uint64, expireAt *time.Time) error {
	return s.repo.Create(ip, reason, createdBy, expireAt)
}

func (s *ipBanService) Unban(id uint64) error {
	return s.repo.Delete(id)
}

func (s *ipBanService) List(page, pageSize int) ([]model.IPBan, int64, error) {
	return s.repo.List(page, pageSize)
}

func (s *ipBanService) Check(ip string) (bool, error) {
	return s.repo.Check(ip)
}
