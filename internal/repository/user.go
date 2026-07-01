package repository

import (
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) GetByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByAccount(account string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? OR username = ?", account, account).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(id uint64, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepo) Delete(id uint64) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepo) DeleteWithCleanup(userID uint64, groupID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete related data in order
		if err := tx.Where("user_id = ?", userID).Delete(&model.OnlineUser{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&model.Notification{}).Error; err != nil {
			return err
		}
		if err := tx.Where("sender_id = ? OR receiver_id = ?", userID, userID).Delete(&model.PrivateMessage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&model.Favorite{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&model.CreditLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&model.EmailVerify{}).Error; err != nil {
			return err
		}

		return tx.Delete(&model.User{}, userID).Error
	})
}

func (r *UserRepo) List(page, pageSize int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepo) IncrThreads(id uint64, delta int) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("threads", gorm.Expr("threads + ?", delta)).Error
}

func (r *UserRepo) IncrPosts(id uint64, delta int) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("posts", gorm.Expr("posts + ?", delta)).Error
}

func (r *UserRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepo) DB() *gorm.DB {
	return r.db
}

func (r *UserRepo) ListGroups() ([]model.Group, error) {
	var groups []model.Group
	err := r.db.Order("id ASC").Find(&groups).Error
	return groups, err
}

func (r *UserRepo) UpdateGroup(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.Group{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepo) VerifyEmail(id uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("email_verified", true).Error
}
