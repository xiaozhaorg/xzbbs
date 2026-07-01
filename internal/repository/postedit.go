package repository

import (
	"time"
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type PostEditRepo interface {
	Create(postID, userID uint64, oldContent, newContent string) error
	ListByPost(postID uint64) ([]model.PostEdit, error)
}

type postEditRepo struct {
	db *gorm.DB
}

func NewPostEditRepo(db *gorm.DB) PostEditRepo {
	return &postEditRepo{db: db}
}

func (r *postEditRepo) Create(postID, userID uint64, oldContent, newContent string) error {
	return r.db.Create(&model.PostEdit{
		PostID:    postID,
		UserID:    userID,
		OldContent: oldContent,
		NewContent: newContent,
		CreatedAt: time.Now(),
	}).Error
}

func (r *postEditRepo) ListByPost(postID uint64) ([]model.PostEdit, error) {
	var edits []model.PostEdit
	err := r.db.Where("post_id = ?", postID).Order("created_at DESC").Find(&edits).Error
	return edits, err
}
