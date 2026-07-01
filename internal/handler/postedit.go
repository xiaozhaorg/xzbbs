package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
)

type PostEditHandler struct {
	postEditRepo repository.PostEditRepo
}

func NewPostEditHandler(postEditRepo repository.PostEditRepo) *PostEditHandler {
	return &PostEditHandler{postEditRepo: postEditRepo}
}

// GetPostEdits returns edit history for a post
// GET /api/posts/:id/edits
func (h *PostEditHandler) GetPostEdits(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	edits, err := h.postEditRepo.ListByPost(id)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, edits)
}
