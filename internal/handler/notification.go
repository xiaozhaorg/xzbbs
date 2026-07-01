package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type NotificationHandler struct {
	svc service.NotificationService
}

func NewNotificationHandler(svc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// List notifications
// GET /api/notifications?unread=1&page=1
func (h *NotificationHandler) List(c *gin.Context) {
	uid := middleware.GetUserID(c)
	unreadOnly := c.Query("unread") == "1"
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	items, total, err := h.svc.List(uid, unreadOnly, page, 20)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	pages := (total + 19) / 20
	errcode.OK(c, gin.H{
		"items": items, "total": total, "page": page, "pages": pages,
	})
}

// Unread count
// GET /api/notifications/unread-count
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	uid := middleware.GetUserID(c)
	count, err := h.svc.UnreadCount(uid)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, gin.H{"count": count})
}

// Mark as read
// POST /api/notifications/read  body: {ids: [1,2,3]}
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	uid := middleware.GetUserID(c)
	var req struct {
		IDs []uint64 `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.IDs) == 0 {
		errcode.FailValidation(c, "invalid ids")
		return
	}
	if err := h.svc.MarkRead(uid, req.IDs); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "marked as read")
}

// Mark all as read
// POST /api/notifications/read-all
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	uid := middleware.GetUserID(c)
	if err := h.svc.MarkAllRead(uid); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "all marked as read")
}

// Delete notification
// DELETE /api/notifications/:id
func (h *NotificationHandler) Delete(c *gin.Context) {
	uid := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(uid, []uint64{id}); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "deleted")
}
