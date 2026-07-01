package repository

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type SmileyRepo interface {
	List() ([]model.Smiley, error)
}

type smileyRepo struct {
	db *gorm.DB
}

func NewSmileyRepo(db *gorm.DB) SmileyRepo {
	return &smileyRepo{db: db}
}

func (r *smileyRepo) List() ([]model.Smiley, error) {
	var smilies []model.Smiley
	err := r.db.Order("sort ASC").Find(&smilies).Error
	return smilies, err
}
