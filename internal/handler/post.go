package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/dto"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type PostHandler struct {
	postSvc   *service.PostService
	threadSvc *service.ThreadService
}

func NewPostHandler(ps *service.PostService, ts *service.ThreadService) *PostHandler {
	return &PostHandler{postSvc: ps, threadSvc: ts}
}

func (h *PostHandler) Create(c *gin.Context) {
	threadID, _ := strconv.ParseUint(c.Param("tid"), 10, 64)

	// Check thread exists and not closed
	thread, err := h.threadSvc.GetByID(threadID)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrThreadNotFound)
		return
	}
	if thread.IsClosed {
		errcode.Fail(c, http.StatusForbidden, errcode.ErrThreadClosed)
		return
	}

	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	post, err := h.postSvc.Create(threadID, userID, req.Content, req.ContentType, req.ReplyTo, c.ClientIP())
	if err != nil {
		errcode.FailMsg(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}
	errcode.OK(c, post)
}

func (h *PostHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	groupID := middleware.GetGroupID(c)

	post, err := h.postSvc.GetByID(id)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrPostNotFound)
		return
	}

	if post.UserID != userID && groupID > 5 {
		errcode.Fail(c, http.StatusForbidden, errcode.ErrNoPermission)
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"content":      req.Content,
		"content_type": req.ContentType,
	}
	if err := h.postSvc.Update(id, updates); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "updated")
}

func (h *PostHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	groupID := middleware.GetGroupID(c)

	post, err := h.postSvc.GetByID(id)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrPostNotFound)
		return
	}

	if post.UserID != userID && groupID > 5 {
		errcode.Fail(c, http.StatusForbidden, errcode.ErrNoPermission)
		return
	}

	// If first post, delete thread instead
	if post.IsFirst {
		h.threadSvc.Delete(post.ThreadID)
	} else {
		h.postSvc.Delete(id)
	}
	errcode.OKMsg(c, "deleted")
}
