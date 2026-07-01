package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/config"
	"github.com/xiaozhaorg/xzbbs/internal/dto"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/pagination"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type ThreadHandler struct {
	threadSvc *service.ThreadService
	postSvc   *service.PostService
}

func NewThreadHandler(ts *service.ThreadService, ps *service.PostService) *ThreadHandler {
	return &ThreadHandler{threadSvc: ts, postSvc: ps}
}

func (h *ThreadHandler) Create(c *gin.Context) {
	var req dto.CreateThreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	thread, _, err := h.threadSvc.Create(userID, req.ForumID, req.Title, req.Content, req.ContentType, c.ClientIP())
	if err != nil {
		errcode.FailMsg(c, http.StatusInternalServerError, errcode.ErrInternal, err.Error())
		return
	}

	errcode.OK(c, thread)
}

func (h *ThreadHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	thread, err := h.threadSvc.GetByID(id)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrThreadNotFound)
		return
	}

	// Increment views (fire and forget)
	go h.threadSvc.IncrViews(id)

	// Get first post
	firstPost, _ := h.postSvc.GetFirstPost(id)

	// Get paginated replies
	var pq pagination.PageQuery
	c.ShouldBindQuery(&pq)
	pq.Normalize(config.Global.Site.PageSize)

	posts, total, _ := h.postSvc.ListByThread(id, pq.Page, pq.PageSize)

	errcode.OK(c, gin.H{
		"thread":     thread,
		"first_post": firstPost,
		"replies":    pagination.NewResult(posts, total, pq.Page, pq.PageSize),
	})
}

func (h *ThreadHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	groupID := middleware.GetGroupID(c)

	thread, err := h.threadSvc.GetByID(id)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrThreadNotFound)
		return
	}

	// Only author or mod can edit
	if thread.UserID != userID && groupID > 5 {
		errcode.Fail(c, http.StatusForbidden, errcode.ErrNoPermission)
		return
	}

	var req dto.UpdateThreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.ForumID != nil {
		updates["forum_id"] = *req.ForumID
	}

	if err := h.threadSvc.Update(id, updates); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "updated")
}

func (h *ThreadHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)
	groupID := middleware.GetGroupID(c)

	thread, err := h.threadSvc.GetByID(id)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrThreadNotFound)
		return
	}

	if thread.UserID != userID && groupID > 5 {
		errcode.Fail(c, http.StatusForbidden, errcode.ErrNoPermission)
		return
	}

	if err := h.threadSvc.Delete(id); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "deleted")
}

func (h *ThreadHandler) ListByForum(c *gin.Context) {
	forumID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	orderBy := c.DefaultQuery("order", "reply")

	var pq pagination.PageQuery
	c.ShouldBindQuery(&pq)
	pq.Normalize(config.Global.Site.PageSize)

	threads, total, err := h.threadSvc.ListByForum(uint(forumID), orderBy, pq.Page, pq.PageSize)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, pagination.NewResult(threads, total, pq.Page, pq.PageSize))
}

// ListAll lists threads across all forums
// GET /api/threads?order=views&page=1
func (h *ThreadHandler) ListAll(c *gin.Context) {
	orderBy := c.DefaultQuery("order", "reply")
	var pq pagination.PageQuery
	c.ShouldBindQuery(&pq)
	pq.Normalize(config.Global.Site.PageSize)

	threads, total, err := h.threadSvc.ListByForum(0, orderBy, pq.Page, pq.PageSize)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, pagination.NewResult(threads, total, pq.Page, pq.PageSize))
}
