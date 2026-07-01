package repository

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type ForumRepo struct {
	db *gorm.DB
}

func NewForumRepo(db *gorm.DB) *ForumRepo {
	return &ForumRepo{db: db}
}

func (r *ForumRepo) Create(forum *model.Forum) error {
	return r.db.Create(forum).Error
}

func (r *ForumRepo) GetByID(id uint) (*model.Forum, error) {
	var forum model.Forum
	err := r.db.First(&forum, id).Error
	if err != nil {
		return nil, err
	}
	return &forum, nil
}

func (r *ForumRepo) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.Forum{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ForumRepo) Delete(id uint) error {
	return r.db.Delete(&model.Forum{}, id).Error
}

func (r *ForumRepo) List() ([]model.Forum, error) {
	var forums []model.Forum
	err := r.db.Order("sort_order DESC, id ASC").Find(&forums).Error
	return forums, err
}

func (r *ForumRepo) IncrThreads(id uint, delta int) error {
	return r.db.Model(&model.Forum{}).Where("id = ?", id).
		Update("threads", gorm.Expr("threads + ?", delta)).Error
}

func (r *ForumRepo) IncrPosts(id uint, delta int) error {
	return r.db.Model(&model.Forum{}).Where("id = ?", id).
		Update("posts", gorm.Expr("posts + ?", delta)).Error
}

// Permissions
func (r *ForumRepo) GetPermissions(forumID uint) ([]model.ForumPermission, error) {
	var perms []model.ForumPermission
	err := r.db.Where("forum_id = ?", forumID).Find(&perms).Error
	return perms, err
}

func (r *ForumRepo) SetPermission(perm *model.ForumPermission) error {
	return r.db.Save(perm).Error
}

func (r *ForumRepo) DeletePermissions(forumID uint) error {
	return r.db.Where("forum_id = ?", forumID).Delete(&model.ForumPermission{}).Error
}
