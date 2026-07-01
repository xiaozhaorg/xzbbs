package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type AttachHandler struct {
	attachSvc *service.AttachService
}

func NewAttachHandler(svc *service.AttachService) *AttachHandler {
	return &AttachHandler{attachSvc: svc}
}

func (h *AttachHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		errcode.FailValidation(c, "file required")
		return
	}

	userID := middleware.GetUserID(c)
	attach, err := h.attachSvc.Upload(file, userID)
	if err != nil {
		errcode.FailMsg(c, http.StatusBadRequest, errcode.ErrBadRequest, err.Error())
		return
	}

	errcode.OK(c, gin.H{
		"id":       attach.ID,
		"url":      "/uploads/" + attach.Filename,
		"filename": attach.OriginalName,
		"size":     attach.FileSize,
		"is_image": attach.IsImage,
	})
}

func (h *AttachHandler) Download(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	attach, err := h.attachSvc.GetByID(id)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrNotFound)
		return
	}

	// Redirect to static file
	c.Redirect(http.StatusFound, "/uploads/"+attach.Filename)
}
