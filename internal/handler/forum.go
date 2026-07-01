package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/dto"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type ForumHandler struct {
	forumSvc *service.ForumService
}

func NewForumHandler(svc *service.ForumService) *ForumHandler {
	return &ForumHandler{forumSvc: svc}
}

func (h *ForumHandler) List(c *gin.Context) {
	forums, err := h.forumSvc.List()
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, forums)
}

func (h *ForumHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	forum, err := h.forumSvc.GetByID(uint(id))
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrForumNotFound)
		return
	}
	errcode.OK(c, forum)
}

func (h *ForumHandler) Create(c *gin.Context) {
	var req dto.CreateForumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	forum, err := h.forumSvc.Create(req.Name, req.Description, req.SortOrder)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, forum)
}

func (h *ForumHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req dto.UpdateForumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}

	if err := h.forumSvc.Update(uint(id), updates); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "updated")
}

func (h *ForumHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.forumSvc.Delete(uint(id)); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "deleted")
}
