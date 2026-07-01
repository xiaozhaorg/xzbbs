package repository

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type AttachRepo struct {
	db *gorm.DB
}

func NewAttachRepo(db *gorm.DB) *AttachRepo {
	return &AttachRepo{db: db}
}

func (r *AttachRepo) Create(attach *model.Attachment) error {
	return r.db.Create(attach).Error
}

func (r *AttachRepo) GetByID(id uint64) (*model.Attachment, error) {
	var attach model.Attachment
	err := r.db.First(&attach, id).Error
	if err != nil {
		return nil, err
	}
	return &attach, nil
}

func (r *AttachRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Attachment{}, id).Error
}

func (r *AttachRepo) DeleteByPost(postID uint64) error {
	return r.db.Where("post_id = ?", postID).Delete(&model.Attachment{}).Error
}

func (r *AttachRepo) ListByPost(postID uint64) ([]model.Attachment, error) {
	var attachments []model.Attachment
	err := r.db.Where("post_id = ?", postID).Find(&attachments).Error
	return attachments, err
}

func (r *AttachRepo) IncrDownloads(id uint64) error {
	return r.db.Model(&model.Attachment{}).Where("id = ?", id).
		Update("downloads", gorm.Expr("downloads + 1")).Error
}

func (r *AttachRepo) UpdatePostID(id, postID uint64) error {
	return r.db.Model(&model.Attachment{}).Where("id = ?", id).
		Update("post_id", postID).Error
}

func (r *AttachRepo) ListByThread(threadID uint64) ([]model.Attachment, error) {
	var attachments []model.Attachment
	err := r.db.Where("post_id IN (?)",
		r.db.Model(&model.Post{}).Select("id").Where("thread_id = ?", threadID),
	).Find(&attachments).Error
	return attachments, err
}
