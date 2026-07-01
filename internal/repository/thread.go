package repository

import (
	"time"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type ThreadRepo struct {
	db *gorm.DB
}

func NewThreadRepo(db *gorm.DB) *ThreadRepo {
	return &ThreadRepo{db: db}
}

func (r *ThreadRepo) Create(thread *model.Thread) error {
	return r.db.Create(thread).Error
}

func (r *ThreadRepo) GetByID(id uint64) (*model.Thread, error) {
	var thread model.Thread
	err := r.db.Preload("User").First(&thread, id).Error
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r *ThreadRepo) Update(id uint64, updates map[string]interface{}) error {
	return r.db.Model(&model.Thread{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ThreadRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Thread{}, id).Error
}

func (r *ThreadRepo) ListByForum(forumID uint, orderBy string, page, pageSize int) ([]model.Thread, int64, error) {
	var threads []model.Thread
	var total int64

	query := r.db.Model(&model.Thread{})
	if forumID > 0 {
		query = query.Where("forum_id = ?", forumID)
	}
	query.Count(&total)

	order := "is_top DESC, "
	switch orderBy {
	case "created":
		order += "id DESC"
	default: // "reply" or default
		order += "last_reply_at DESC, id DESC"
	}

	err := query.Preload("User").
		Offset((page-1)*pageSize).Limit(pageSize).
		Order(order).
		Find(&threads).Error
	return threads, total, err
}

func (r *ThreadRepo) ListByUser(userID uint64, page, pageSize int) ([]model.Thread, int64, error) {
	var threads []model.Thread
	var total int64

	query := r.db.Model(&model.Thread{}).Where("user_id = ?", userID)
	query.Count(&total)
	err := query.Preload("Forum").
		Offset((page-1)*pageSize).Limit(pageSize).
		Order("id DESC").
		Find(&threads).Error
	return threads, total, err
}

func (r *ThreadRepo) IncrViews(id uint64) error {
	return r.db.Model(&model.Thread{}).Where("id = ?", id).
		Update("views", gorm.Expr("views + 1")).Error
}

func (r *ThreadRepo) IncrPosts(id uint64, delta int) error {
	return r.db.Model(&model.Thread{}).Where("id = ?", id).
		Update("posts", gorm.Expr("posts + ?", delta)).Error
}

func (r *ThreadRepo) UpdateBatch(ids []uint64, updates map[string]interface{}) error {
	return r.db.Model(&model.Thread{}).Where("id IN ?", ids).Updates(updates).Error
}

func (r *ThreadRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Thread{}).Count(&count).Error
	return count, err
}

func (r *ThreadRepo) TodayCount() (int64, error) {
	var count int64
	today := time.Now().Truncate(24 * time.Hour)
	err := r.db.Model(&model.Thread{}).
		Where("created_at >= ?", today).Count(&count).Error
	return count, err
}

func (r *ThreadRepo) FindByIDs(ids []uint64) ([]model.Thread, error) {
	var threads []model.Thread
	if len(ids) == 0 {
		return threads, nil
	}
	err := r.db.Where("id IN ?", ids).Preload("User").Preload("Forum").Find(&threads).Error
	return threads, err
}
