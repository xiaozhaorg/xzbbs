package repository

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type FavoriteRepo interface {
	Toggle(userID, threadID uint64) (bool, error)
	IsFavorited(userID, threadID uint64) (bool, error)
	ListByUser(userID uint64, page, pageSize int) ([]model.Favorite, int64, error)
	CountByThread(threadID uint64) (int64, error)
	Delete(userID, threadID uint64) error
}

type favoriteRepo struct {
	db *gorm.DB
}

func NewFavoriteRepo(db *gorm.DB) FavoriteRepo {
	return &favoriteRepo{db: db}
}

func (r *favoriteRepo) Toggle(userID, threadID uint64) (bool, error) {
	var fav model.Favorite
	err := r.db.Where("user_id = ? AND thread_id = ?", userID, threadID).First(&fav).Error
	if err == nil {
		// Already favorited, unfavorite
		r.db.Delete(&fav)
		return false, nil
	}
	if err == gorm.ErrRecordNotFound {
		fav = model.Favorite{UserID: userID, ThreadID: threadID}
		return true, r.db.Create(&fav).Error
	}
	return false, err
}

func (r *favoriteRepo) IsFavorited(userID, threadID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).Where("user_id = ? AND thread_id = ?", userID, threadID).Count(&count).Error
	return count > 0, err
}

func (r *favoriteRepo) ListByUser(userID uint64, page, pageSize int) ([]model.Favorite, int64, error) {
	var total int64
	var items []model.Favorite
	r.db.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *favoriteRepo) CountByThread(threadID uint64) (int64, error) {
	var count int64
	return count, r.db.Model(&model.Favorite{}).Where("thread_id = ?", threadID).Count(&count).Error
}

func (r *favoriteRepo) Delete(userID, threadID uint64) error {
	return r.db.Where("user_id = ? AND thread_id = ?", userID, threadID).Delete(&model.Favorite{}).Error
}
