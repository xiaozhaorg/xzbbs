package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/pagination"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

// PMHandler handles private message endpoints
type PMHandler struct {
	svc service.PMService
}

func NewPMHandler(svc service.PMService) *PMHandler {
	return &PMHandler{svc: svc}
}

// Send PM
// POST /api/pms
func (h *PMHandler) Send(c *gin.Context) {
	var req struct {
		ReceiverID uint64 `json:"receiver_id" binding:"required"`
		Content    string `json:"content" binding:"required,max=5000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}
	senderID := middleware.GetUserID(c)
	if senderID == req.ReceiverID {
		errcode.Fail(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}
	pm, err := h.svc.Send(senderID, req.ReceiverID, req.Content)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, pm)
}

// Get conversation with a user
// GET /api/pms/conversations/:otherId?page=1
func (h *PMHandler) Conversation(c *gin.Context) {
	uid := middleware.GetUserID(c)
	otherID, _ := strconv.ParseUint(c.Param("otherId"), 10, 64)
	page, pageSize := pagination.Normalize(c.Query("page"), c.Query("page_size"))
	items, total, err := h.svc.GetConversation(uid, otherID, page, pageSize)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, pagination.NewResult(items, total, page, pageSize))
}

// Get conversation list
// GET /api/pms
func (h *PMHandler) List(c *gin.Context) {
	uid := middleware.GetUserID(c)
	convs, err := h.svc.GetConversationList(uid)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, convs)
}

// Get unread count
// GET /api/pms/unread-count
func (h *PMHandler) UnreadCount(c *gin.Context) {
	uid := middleware.GetUserID(c)
	count, err := h.svc.GetUnreadCount(uid)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, gin.H{"count": count})
}

// Mark conversation as read
// POST /api/pms/read/:otherId
func (h *PMHandler) MarkRead(c *gin.Context) {
	uid := middleware.GetUserID(c)
	otherID, _ := strconv.ParseUint(c.Param("otherId"), 10, 64)
	if err := h.svc.MarkRead(uid, otherID); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "marked as read")
}

// Delete PM
// DELETE /api/pms/:id
func (h *PMHandler) Delete(c *gin.Context) {
	uid := middleware.GetUserID(c)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.SoftDelete(uid, id); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "deleted")
}
