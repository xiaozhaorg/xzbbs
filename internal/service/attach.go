package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xiaozhaorg/xzbbs/internal/config"
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type AttachService struct {
	repo *repository.AttachRepo
}

func NewAttachService(repo *repository.AttachRepo) *AttachService {
	return &AttachService{repo: repo}
}

func (s *AttachService) Upload(file *multipart.FileHeader, userID uint64) (*model.Attachment, error) {
	cfg := config.Global

	// Check size
	if file.Size > cfg.Upload.MaxSize*1024*1024 {
		return nil, fmt.Errorf("file too large (max %dMB)", cfg.Upload.MaxSize)
	}

	// Check extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, e := range cfg.Upload.AllowExt {
		if ext == e {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, fmt.Errorf("file type %s not allowed", ext)
	}

	// Generate storage path: uploads/2024/06/uuid.ext
	now := time.Now()
	dir := filepath.Join(cfg.Upload.Path, now.Format("2006"), now.Format("01"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	filename := uuid.New().String() + ext
	destPath := filepath.Join(dir, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Determine if image
	isImage := false
	mimeType := file.Header.Get("Content-Type")
	if strings.HasPrefix(mimeType, "image/") {
		isImage = true
	}

	// Relative path for storage
	relPath := filepath.Join(now.Format("2006"), now.Format("01"), filename)

	attach := &model.Attachment{
		UserID:       userID,
		Filename:     filepath.ToSlash(relPath),
		OriginalName: file.Filename,
		FileSize:     uint(file.Size),
		MimeType:     mimeType,
		IsImage:      isImage,
	}
	if err := s.repo.Create(attach); err != nil {
		os.Remove(destPath)
		return nil, err
	}

	return attach, nil
}

func (s *AttachService) GetByID(id uint64) (*model.Attachment, error) {
	return s.repo.GetByID(id)
}

func (s *AttachService) Delete(id uint64) error {
	attach, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	// Remove file
	path := filepath.Join(config.Global.Upload.Path, attach.Filename)
	os.Remove(path)
	return s.repo.Delete(id)
}

func (s *AttachService) AssocPost(id, postID uint64) error {
	return s.repo.UpdatePostID(id, postID)
}

func (s *AttachService) DeleteByPost(postID uint64) error {
	attachments, err := s.repo.ListByPost(postID)
	if err != nil {
		return err
	}
	for _, a := range attachments {
		path := filepath.Join(config.Global.Upload.Path, a.Filename)
		_ = os.Remove(path)
	}
	return s.repo.DeleteByPost(postID)
}

func (s *AttachService) ListByPost(postID uint64) ([]model.Attachment, error) {
	return s.repo.ListByPost(postID)
}

func (s *AttachService) DeleteByThread(threadID uint64) error {
	attachments, err := s.repo.ListByThread(threadID)
	if err != nil {
		return err
	}
	for _, a := range attachments {
		path := filepath.Join(config.Global.Upload.Path, a.Filename)
		_ = os.Remove(path)
	}
	return nil
}
