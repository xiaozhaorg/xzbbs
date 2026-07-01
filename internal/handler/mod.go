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

type ModHandler struct {
	threadSvc   *service.ThreadService
	userSvc     *service.UserService
	modLogSvc   *service.ModLogService
}

func NewModHandler(ts *service.ThreadService, us *service.UserService, modLogSvc *service.ModLogService) *ModHandler {
	return &ModHandler{threadSvc: ts, userSvc: us, modLogSvc: modLogSvc}
}

func (h *ModHandler) SetTop(c *gin.Context) {
	var req dto.ModTopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	if err := h.threadSvc.SetTop(req.ThreadIDs, req.Top); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	h.logAction(c, req.ThreadIDs, "top")
	errcode.OKMsg(c, "ok")
}

func (h *ModHandler) SetClosed(c *gin.Context) {
	var req dto.ModCloseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	if err := h.threadSvc.SetClosed(req.ThreadIDs, req.Closed); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	h.logAction(c, req.ThreadIDs, "close")
	errcode.OKMsg(c, "ok")
}

func (h *ModHandler) Move(c *gin.Context) {
	var req dto.ModMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	if err := h.threadSvc.Move(req.ThreadIDs, req.ForumID); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	h.logAction(c, req.ThreadIDs, "move")
	errcode.OKMsg(c, "ok")
}

func (h *ModHandler) BanUser(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user, err := h.userSvc.GetByID(uid)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrUserNotFound)
		return
	}
	if user.GroupID <= 5 {
		errcode.FailMsg(c, http.StatusForbidden, errcode.ErrForbidden, "cannot ban admin/mod")
		return
	}

	// Set to banned group (7)
	h.userSvc.Update(uid, map[string]interface{}{"group_id": 7})

	modUID := middleware.GetUserID(c)
	h.modLogSvc.Create(modUID, "ban_user", "banned user "+strconv.FormatUint(uid, 10), 0)
	errcode.OKMsg(c, "user banned")
}

func (h *ModHandler) logAction(c *gin.Context, threadIDs []uint64, action string) {
	uid := middleware.GetUserID(c)
	h.modLogSvc.CreateBatch(uid, threadIDs, action)
}
