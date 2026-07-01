package service

import (
	"time"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/plugin"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
	"gorm.io/gorm"
)

type PostService struct {
	db           *gorm.DB
	postRepo     *repository.PostRepo
	threadRepo   *repository.ThreadRepo
	forumRepo    *repository.ForumRepo
	userRepo     *repository.UserRepo
	threadSvc    *ThreadService
	postEditRepo repository.PostEditRepo
	creditSvc    *CreditService
	attachSvc    *AttachService
}

func NewPostService(db *gorm.DB, pr *repository.PostRepo, tr *repository.ThreadRepo, fr *repository.ForumRepo, ur *repository.UserRepo, ts *ThreadService, per repository.PostEditRepo, credit *CreditService, attachSvc *AttachService) *PostService {
	return &PostService{
		db: db, postRepo: pr, threadRepo: tr, forumRepo: fr,
		userRepo: ur, threadSvc: ts, postEditRepo: per, creditSvc: credit, attachSvc: attachSvc,
	}
}

func (s *PostService) Create(threadID, userID uint64, content string, contentType uint8, replyTo uint64, ip string) (*model.Post, error) {
	if r := plugin.Global().Dispatch(plugin.EventBeforePostSave, content); len(r) > 0 {
		if v, ok := r[0].(string); ok && v != "" {
			content = v
		}
	}

	var post *model.Post
	var forumID uint

	err := s.db.Transaction(func(tx *gorm.DB) error {
		post = &model.Post{
			ThreadID: threadID, UserID: userID, IsFirst: false,
			Content: content, ContentType: contentType, ReplyTo: replyTo, IP: ip,
		}
		if err := tx.Create(post).Error; err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&model.Thread{}).Where("id = ?", threadID).Updates(map[string]interface{}{
			"last_reply_at": now, "last_reply_uid": userID,
			"posts": gorm.Expr("posts + 1"),
		}).Error; err != nil {
			return err
		}

		var thread model.Thread
		if err := tx.Select("forum_id").First(&thread, threadID).Error; err != nil {
			return err
		}
		forumID = thread.ForumID

		if err := tx.Model(&model.Forum{}).Where("id = ?", forumID).
			Update("posts", gorm.Expr("posts + 1")).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.User{}).Where("id = ?", userID).
			Update("posts", gorm.Expr("posts + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	plugin.Global().Dispatch(plugin.EventAfterPostSave, post.ID)

	_ = s.threadSvc.NotifyReply(threadID, userID)

	if s.creditSvc != nil {
		_ = s.creditSvc.AwardPostCreation(userID, post.ID)
		_ = s.creditSvc.UpdateLevel(userID)
	}

	plugin.Global().Dispatch(plugin.EventPostCreated, post.ID, userID, threadID)
	return post, nil
}

func (s *PostService) GetByID(id uint64) (*model.Post, error) {
	return s.postRepo.GetByID(id)
}

func (s *PostService) Update(id uint64, updates map[string]interface{}) error {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return err
	}

	if newContent, ok := updates["content"].(string); ok && s.postEditRepo != nil {
		if newContent != post.Content {
			_ = s.postEditRepo.Create(id, post.UserID, post.Content, newContent)
		}
	}

	return s.postRepo.Update(id, updates)
}

func (s *PostService) Delete(id uint64) error {
	post, err := s.postRepo.GetByID(id)
	if err != nil {
		return err
	}

	var threadID uint64
	var isFirst bool

	err = s.db.Transaction(func(tx *gorm.DB) error {
		threadID = post.ThreadID
		isFirst = post.IsFirst

		if post.IsFirst {
			var thread model.Thread
			if err := tx.First(&thread, post.ThreadID).Error; err != nil {
				return err
			}

			if err := tx.Delete(&model.Thread{}, post.ThreadID).Error; err != nil {
				return err
			}
			if err := tx.Where("thread_id = ?", post.ThreadID).
				Delete(&model.Post{}).Error; err != nil {
				return err
			}
			if err := tx.Where("thread_id = ?", post.ThreadID).
				Delete(&model.Attachment{}).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.Forum{}).Where("id = ?", thread.ForumID).
				Update("threads", gorm.Expr("threads - 1")).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Forum{}).Where("id = ?", thread.ForumID).
				Update("posts", gorm.Expr("posts - ?", thread.Posts+1)).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.User{}).Where("id = ?", post.UserID).
				Update("threads", gorm.Expr("threads - 1")).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.User{}).Where("id = ?", post.UserID).
				Update("posts", gorm.Expr("posts - ?", thread.Posts+1)).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Delete(&model.Post{}, id).Error; err != nil {
				return err
			}
			if err := tx.Where("post_id = ?", id).
				Delete(&model.Attachment{}).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.Thread{}).Where("id = ?", post.ThreadID).
				Update("posts", gorm.Expr("posts - 1")).Error; err != nil {
				return err
			}

			var thread model.Thread
			if err := tx.Select("forum_id").First(&thread, post.ThreadID).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.Forum{}).Where("id = ?", thread.ForumID).
				Update("posts", gorm.Expr("posts - 1")).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.User{}).Where("id = ?", post.UserID).
				Update("posts", gorm.Expr("posts - 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Clean up physical files after successful DB transaction
	if s.attachSvc != nil {
		if isFirst {
			_ = s.attachSvc.DeleteByThread(threadID)
		} else {
			_ = s.attachSvc.DeleteByPost(id)
		}
	}

	return nil
}

func (s *PostService) ListByThread(threadID uint64, page, pageSize int) ([]model.Post, int64, error) {
	return s.postRepo.ListByThread(threadID, page, pageSize)
}

func (s *PostService) ListByUser(userID uint64, page, pageSize int) ([]model.Post, int64, error) {
	return s.postRepo.ListByUser(userID, page, pageSize)
}

func (s *PostService) GetFirstPost(threadID uint64) (*model.Post, error) {
	return s.postRepo.GetFirstPost(threadID)
}

func (s *PostService) Count() (int64, error) {
	return s.postRepo.Count()
}

func (s *PostService) TodayCount() (int64, error) {
	return s.postRepo.TodayCount()
}
