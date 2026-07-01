package repository

import (
	"strings"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type SearchRepo interface {
	SearchThreads(query string, forumID uint, page, pageSize int) ([]model.Thread, int64, error)
}

type searchRepo struct {
	db *gorm.DB
}

func NewSearchRepo(db *gorm.DB) SearchRepo {
	return &searchRepo{db: db}
}

func (r *searchRepo) SearchThreads(query string, forumID uint, page, pageSize int) ([]model.Thread, int64, error) {
	// Escape LIKE special characters
	escapedQuery := strings.ReplaceAll(query, "%", "\\%")
	escapedQuery = strings.ReplaceAll(escapedQuery, "_", "\\_")
	like := "%" + escapedQuery + "%"

	var threads []model.Thread
	var total int64

	// Search in thread titles
	titleQ := r.db.Model(&model.Thread{}).Where("title LIKE ? ESCAPE '\\'", like)
	if forumID > 0 {
		titleQ = titleQ.Where("forum_id = ?", forumID)
	}
	titleQ.Count(&total)

	err := titleQ.Preload("User").Preload("Forum").
		Order("is_top DESC, last_reply_at DESC, created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&threads).Error

	return threads, total, err
}
