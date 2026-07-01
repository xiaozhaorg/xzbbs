package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/config"
	"github.com/xiaozhaorg/xzbbs/internal/dto"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/pagination"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type UserHandler struct {
	userSvc   *service.UserService
	threadSvc *service.ThreadService
	postSvc   *service.PostService
}

func NewUserHandler(us *service.UserService, ts *service.ThreadService, ps *service.PostService) *UserHandler {
	return &UserHandler{userSvc: us, threadSvc: ts, postSvc: ps}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user, err := h.userSvc.GetByID(id)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrUserNotFound)
		return
	}
	errcode.OK(c, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"group_id":       user.GroupID,
		"avatar":         user.Avatar,
		"threads":        user.Threads,
		"posts":          user.Posts,
		"credits":        user.Credits,
		"level":          user.Level,
		"signature":      user.Signature,
		"email_verified": user.EmailVerified,
		"created_at":     user.CreatedAt,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	currentUID := middleware.GetUserID(c)
	if id != currentUID {
		errcode.Fail(c, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if len(updates) == 0 {
		errcode.OKMsg(c, "nothing to update")
		return
	}

	if err := h.userSvc.Update(id, updates); err != nil {
		errcode.FailMsg(c, http.StatusConflict, errcode.ErrConflict, err.Error())
		return
	}
	errcode.OKMsg(c, "updated")
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	currentUID := middleware.GetUserID(c)
	if id != currentUID {
		errcode.Fail(c, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	if err := h.userSvc.ChangePassword(id, req.OldPassword, req.NewPassword); err != nil {
		errcode.FailMsg(c, http.StatusBadRequest, errcode.ErrBadRequest, err.Error())
		return
	}
	errcode.OKMsg(c, "password changed")
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	currentUID := middleware.GetUserID(c)
	if id != currentUID {
		errcode.Fail(c, http.StatusForbidden, errcode.ErrForbidden)
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		errcode.FailValidation(c, "avatar file required")
		return
	}

	maxSize := config.Global.Upload.MaxSize * 1024 * 1024
	if file.Size > maxSize {
		errcode.Fail(c, http.StatusBadRequest, errcode.ErrFileTooLarge)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, e := range config.Global.Upload.AllowExt {
		if ext == e {
			allowed = true
			break
		}
	}
	if !allowed {
		errcode.Fail(c, http.StatusBadRequest, errcode.ErrFileTypeNotAllowed)
		return
	}

	dst := config.Global.Upload.Path + "/avatars/" + strconv.FormatUint(id, 10) + ext
	if err := c.SaveUploadedFile(file, dst); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	avatarURL := "/uploads/avatars/" + strconv.FormatUint(id, 10) + ext
	h.userSvc.Update(id, map[string]interface{}{"avatar": avatarURL})
	errcode.OK(c, gin.H{"avatar": avatarURL})
}

func (h *UserHandler) GetUserThreads(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var pq pagination.PageQuery
	c.ShouldBindQuery(&pq)
	pq.Normalize(config.Global.Site.PageSize)

	threads, total, err := h.threadSvc.ListByUser(id, pq.Page, pq.PageSize)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, pagination.NewResult(threads, total, pq.Page, pq.PageSize))
}

func (h *UserHandler) GetUserPosts(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var pq pagination.PageQuery
	c.ShouldBindQuery(&pq)
	pq.Normalize(config.Global.Site.PageSize)

	posts, total, err := h.postSvc.ListByUser(id, pq.Page, pq.PageSize)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, pagination.NewResult(posts, total, pq.Page, pq.PageSize))
}
