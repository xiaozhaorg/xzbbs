package repository

import (
	"time"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepo) GetByID(id uint64) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepo) Update(id uint64, updates map[string]interface{}) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Post{}, id).Error
}

func (r *PostRepo) DeleteByThread(threadID uint64) error {
	return r.db.Where("thread_id = ?", threadID).Delete(&model.Post{}).Error
}

func (r *PostRepo) ListByThread(threadID uint64, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("thread_id = ?", threadID)
	query.Count(&total)

	err := query.Preload("User").
		Offset((page-1)*pageSize).Limit(pageSize).
		Order("id ASC").
		Find(&posts).Error
	return posts, total, err
}

func (r *PostRepo) ListByUser(userID uint64, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("user_id = ?", userID)
	query.Count(&total)

	err := query.Offset((page-1)*pageSize).Limit(pageSize).
		Order("id DESC").
		Find(&posts).Error
	return posts, total, err
}

func (r *PostRepo) GetFirstPost(threadID uint64) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("thread_id = ? AND is_first = ?", threadID, true).
		Preload("User").First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Post{}).Count(&count).Error
	return count, err
}

func (r *PostRepo) TodayCount() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&model.Post{}).
		Where("created_at >= ?", today).Count(&count).Error
	return count, err
}
