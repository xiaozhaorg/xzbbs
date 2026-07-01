package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type FavoriteHandler struct {
	svc service.FavoriteService
}

func NewFavoriteHandler(svc service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{svc: svc}
}

// Toggle favorite/unfavorite
// POST /api/favorites/threads/:id
func (h *FavoriteHandler) Toggle(c *gin.Context) {
	uid := middleware.GetUserID(c)
	tidStr := c.Param("id")
	tid, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		errcode.Fail(c, http.StatusOK, errcode.ErrBadRequest)
		return
	}
	favorited, err := h.svc.Toggle(uid, tid)
	if err != nil {
		errcode.Fail(c, http.StatusOK, errcode.ErrBadRequest)
		return
	}
	errcode.OK(c, gin.H{"favorited": favorited})
}

// List user's favorites
// GET /api/favorites/threads?page=1
func (h *FavoriteHandler) List(c *gin.Context) {
	uid := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	threads, total, err := h.svc.ListByUser(uid, page, 20)
	if err != nil {
		errcode.Fail(c, http.StatusOK, errcode.ErrBadRequest)
		return
	}
	pages := (total + 19) / 20
	errcode.OK(c, gin.H{
		"items": threads, "total": total, "page": page, "pages": pages,
	})
}

// Check if current user favorited a thread
// GET /api/favorites/threads/:id/check
func (h *FavoriteHandler) Check(c *gin.Context) {
	uid := middleware.GetUserID(c)
	tidStr := c.Param("id")
	tid, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil {
		errcode.Fail(c, http.StatusOK, errcode.ErrBadRequest)
		return
	}
	favorited, err := h.svc.IsFavorited(uid, tid)
	if err != nil {
		errcode.Fail(c, http.StatusOK, errcode.ErrBadRequest)
		return
	}
	errcode.OK(c, gin.H{"favorited": favorited})
}
