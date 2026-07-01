package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/pagination"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type SearchHandler struct {
	svc service.SearchService
}

func NewSearchHandler(svc service.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

// Search threads
// GET /api/search?q=keyword&forum_id=0&page=1
func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		errcode.Fail(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}
	forumIDStr := c.DefaultQuery("forum_id", "0")
	forumID, _ := strconv.ParseUint(forumIDStr, 10, 64)
	forumIDUint := uint(forumID)

	page, pageSize := pagination.Normalize(c.Query("page"), c.Query("page_size"))

	threads, total, err := h.svc.SearchThreads(q, forumIDUint, page, pageSize)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	pages := (total + int64(pageSize) - 1) / int64(pageSize)
	errcode.OK(c, gin.H{
		"items": threads, "total": total, "page": page,
		"page_size": pageSize, "pages": pages,
	})
}
