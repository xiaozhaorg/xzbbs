package service

import (
	"time"

	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/plugin"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
	"gorm.io/gorm"
)

type ThreadService struct {
	db           *gorm.DB
	threadRepo   *repository.ThreadRepo
	postRepo     *repository.PostRepo
	forumRepo    *repository.ForumRepo
	userRepo     *repository.UserRepo
	notifService NotificationService
	creditSvc    *CreditService
}

func NewThreadService(db *gorm.DB, tr *repository.ThreadRepo, pr *repository.PostRepo, fr *repository.ForumRepo, ur *repository.UserRepo, notif NotificationService, credit *CreditService) *ThreadService {
	return &ThreadService{
		db: db, threadRepo: tr, postRepo: pr, forumRepo: fr,
		userRepo: ur, notifService: notif, creditSvc: credit,
	}
}

func (s *ThreadService) Create(userID uint64, forumID uint, title, content string, contentType uint8, ip string) (*model.Thread, *model.Post, error) {
	now := time.Now()

	var thread *model.Thread
	var post *model.Post

	err := s.db.Transaction(func(tx *gorm.DB) error {
		thread = &model.Thread{
			ForumID: forumID, UserID: userID, Title: title,
			LastReplyAt: &now, LastReplyUID: userID,
		}
		if err := tx.Create(thread).Error; err != nil {
			return err
		}

		post = &model.Post{
			ThreadID: thread.ID, UserID: userID, IsFirst: true,
			Content: content, ContentType: contentType, IP: ip,
		}
		if err := tx.Create(post).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Forum{}).Where("id = ?", forumID).
			Update("threads", gorm.Expr("threads + 1")).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.User{}).Where("id = ?", userID).
			Update("threads", gorm.Expr("threads + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	if s.creditSvc != nil {
		_ = s.creditSvc.AwardThreadCreation(userID, thread.ID)
		_ = s.creditSvc.UpdateLevel(userID)
	}

	plugin.Global().Dispatch(plugin.EventThreadCreated, thread.ID, userID, title)
	return thread, post, nil
}

func (s *ThreadService) GetByID(id uint64) (*model.Thread, error) {
	return s.threadRepo.GetByID(id)
}

func (s *ThreadService) ListByForum(forumID uint, orderBy string, page, pageSize int) ([]model.Thread, int64, error) {
	return s.threadRepo.ListByForum(forumID, orderBy, page, pageSize)
}

func (s *ThreadService) ListByUser(userID uint64, page, pageSize int) ([]model.Thread, int64, error) {
	return s.threadRepo.ListByUser(userID, page, pageSize)
}

func (s *ThreadService) Update(id uint64, updates map[string]interface{}) error {
	return s.threadRepo.Update(id, updates)
}

func (s *ThreadService) Delete(id uint64) error {
	thread, err := s.threadRepo.GetByID(id)
	if err != nil {
		return err
	}
	s.postRepo.DeleteByThread(id)
	s.threadRepo.Delete(id)
	s.forumRepo.IncrThreads(thread.ForumID, -1)
	s.userRepo.IncrThreads(thread.UserID, -1)

	if s.notifService != nil {
		_ = s.notifService.DeleteAllByThread(id)
	}
	return nil
}

func (s *ThreadService) IncrViews(id uint64) {
	s.threadRepo.IncrViews(id)
}

// Moderation
func (s *ThreadService) SetTop(ids []uint64, top uint8) error {
	return s.threadRepo.UpdateBatch(ids, map[string]interface{}{"is_top": top})
}

func (s *ThreadService) SetClosed(ids []uint64, closed bool) error {
	return s.threadRepo.UpdateBatch(ids, map[string]interface{}{"is_closed": closed})
}

func (s *ThreadService) Move(ids []uint64, forumID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Get old forum IDs and thread counts
		type threadInfo struct {
			ID       uint64
			ForumID  uint
			PostCount int
		}
		var threads []threadInfo
		if err := tx.Model(&model.Thread{}).Where("id IN ?", ids).
			Select("id, forum_id, posts").Scan(&threads).Error; err != nil {
			return err
		}

		// Group by old forum and count threads/posts
		forumStats := make(map[uint]struct { threads int; posts int })
		for _, t := range threads {
			stats := forumStats[t.ForumID]
			stats.threads++
			stats.posts += t.PostCount + 1 // +1 for the first post
			forumStats[t.ForumID] = stats
		}

		// Update old forums (decrement)
		for oldForumID, stats := range forumStats {
			if err := tx.Model(&model.Forum{}).Where("id = ?", oldForumID).
				Update("threads", gorm.Expr("threads - ?", stats.threads)).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Forum{}).Where("id = ?", oldForumID).
				Update("posts", gorm.Expr("posts - ?", stats.posts)).Error; err != nil {
				return err
			}
		}

		// Update new forum (increment)
		var totalThreads int
		var totalPosts int
		for _, stats := range forumStats {
			totalThreads += stats.threads
			totalPosts += stats.posts
		}
		if err := tx.Model(&model.Forum{}).Where("id = ?", forumID).
			Update("threads", gorm.Expr("threads + ?", totalThreads)).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Forum{}).Where("id = ?", forumID).
			Update("posts", gorm.Expr("posts + ?", totalPosts)).Error; err != nil {
			return err
		}

		// Move threads
		return tx.Model(&model.Thread{}).Where("id IN ?", ids).
			Update("forum_id", forumID).Error
	})
}

func (s *ThreadService) Count() (int64, error) {
	return s.threadRepo.Count()
}

// NotifyReply creates a notification for the thread author when someone replies
func (s *ThreadService) NotifyReply(threadID, replierID uint64) error {
	thread, err := s.threadRepo.GetByID(threadID)
	if err != nil {
		return err
	}
	// Don't notify yourself
	if thread.UserID == replierID {
		return nil
	}
	replier, _ := s.userRepo.GetByID(replierID)
	msg := ""
	if replier != nil {
		msg = replier.Username + " 回复了你的帖子"
	} else {
		msg = "有人回复了你的帖子"
	}
	return s.notifService.Create(thread.UserID, 0, replierID, threadID, 0, msg)
}
