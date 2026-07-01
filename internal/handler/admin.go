package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/config"
	"github.com/xiaozhaorg/xzbbs/internal/dto"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/pagination"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
	"github.com/xiaozhaorg/xzbbs/internal/service"
	"gorm.io/gorm"
)

type AdminHandler struct {
	userSvc   *service.UserService
	forumSvc  *service.ForumService
	threadSvc *service.ThreadService
	postSvc   *service.PostService
	db        *gorm.DB
}

func NewAdminHandler(us *service.UserService, fs *service.ForumService, ts *service.ThreadService, ps *service.PostService, db *gorm.DB) *AdminHandler {
	return &AdminHandler{userSvc: us, forumSvc: fs, threadSvc: ts, postSvc: ps, db: db}
}

func (h *AdminHandler) Stats(c *gin.Context) {
	users, _ := h.userSvc.Count()
	threads, _ := h.threadSvc.Count()
	posts, _ := h.postSvc.Count()
	todayPosts, _ := h.postSvc.TodayCount()
	errcode.OK(c, gin.H{
		"users": users, "threads": threads, "posts": posts, "today_posts": todayPosts,
	})
}

func (h *AdminHandler) ListGroups(c *gin.Context) {
	groups, err := h.userSvc.ListGroups()
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, groups)
}

func (h *AdminHandler) UpdateGroup(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req dto.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}
	updates := make(map[string]interface{})
	if req.Name != nil { updates["name"] = *req.Name }
	if req.AllowRead != nil { updates["allow_read"] = *req.AllowRead }
	if req.AllowThread != nil { updates["allow_thread"] = *req.AllowThread }
	if req.AllowPost != nil { updates["allow_post"] = *req.AllowPost }
	if req.AllowAttach != nil { updates["allow_attach"] = *req.AllowAttach }
	if req.AllowDown != nil { updates["allow_down"] = *req.AllowDown }
	if req.AllowTop != nil { updates["allow_top"] = *req.AllowTop }
	if req.AllowUpdate != nil { updates["allow_update"] = *req.AllowUpdate }
	if req.AllowDelete != nil { updates["allow_delete"] = *req.AllowDelete }
	if req.AllowMove != nil { updates["allow_move"] = *req.AllowMove }
	if req.AllowBanUser != nil { updates["allow_ban_user"] = *req.AllowBanUser }
	if err := h.userSvc.UpdateGroup(uint(id), updates); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "updated")
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var pq pagination.PageQuery
	c.ShouldBindQuery(&pq)
	pq.Normalize(20)
	search := c.Query("search")
	users, total, _ := h.userSvc.List(pq.Page, pq.PageSize, search)
	errcode.OK(c, pagination.NewResult(users, total, pq.Page, pq.PageSize))
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}
	allowed := map[string]bool{"group_id": true, "username": true, "email": true, "credits": true, "email_verified": true}
	updates := make(map[string]interface{})
	for k, v := range body {
		if allowed[k] { updates[k] = v }
	}
	h.userSvc.Update(id, updates)
	errcode.OKMsg(c, "updated")
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	user, err := h.userSvc.GetByID(id)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrUserNotFound)
		return
	}
	if user.GroupID == 1 {
		errcode.FailMsg(c, http.StatusForbidden, errcode.ErrForbidden, "cannot delete admin")
		return
	}
	h.userSvc.Delete(id)
	errcode.OKMsg(c, "deleted")
}

func (h *AdminHandler) GetSettings(c *gin.Context) {
	cfg := config.Global
	errcode.OK(c, gin.H{"site_name": cfg.Site.Name, "brief": cfg.Site.Brief, "page_size": cfg.Site.PageSize})
}

func (h *AdminHandler) UpdateSettings(c *gin.Context) {
	var req dto.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}
	cfg := config.Global
	if req.SiteName != "" { cfg.Site.Name = req.SiteName }
	if req.Brief != "" { cfg.Site.Brief = req.Brief }
	if req.PageSize > 0 { cfg.Site.PageSize = req.PageSize }
	if err := config.Save(); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "updated")
}

// OnlineUsers returns currently online users
// GET /api/online
func (h *AdminHandler) OnlineUsers(c *gin.Context) {
	users, err := h.userSvc.ListOnlineUsers(50)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, users)
}

// IP Ban Management
func (h *AdminHandler) BanIP(c *gin.Context) {
	var req struct {
		IP      string `json:"ip" binding:"required"`
		Reason  string `json:"reason" binding:"required,max=255"`
		ExpireH int    `json:"expire_hours"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}
	banSvc := service.NewIPBanService(repository.NewIPBanRepo(h.db))
	var expireAt *time.Time
	if req.ExpireH > 0 {
		t := time.Now().Add(time.Duration(req.ExpireH) * time.Hour)
		expireAt = &t
	}
	if err := banSvc.Ban(req.IP, req.Reason, middleware.GetUserID(c), expireAt); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "banned")
}

func (h *AdminHandler) UnbanIP(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	banSvc := service.NewIPBanService(repository.NewIPBanRepo(h.db))
	if err := banSvc.Unban(id); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "unbanned")
}

func (h *AdminHandler) ListIPBans(c *gin.Context) {
	page, pageSize := pagination.Normalize(c.Query("page"), c.Query("page_size"))
	banSvc := service.NewIPBanService(repository.NewIPBanRepo(h.db))
	items, total, err := banSvc.List(page, pageSize)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, pagination.NewResult(items, total, page, pageSize))
}

func (h *AdminHandler) CheckIPBan(c *gin.Context) {
	ip := c.Query("ip")
	banSvc := service.NewIPBanService(repository.NewIPBanRepo(h.db))
	banned, err := banSvc.Check(ip)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, gin.H{"banned": banned})
}

func (h *AdminHandler) ListSmilies(c *gin.Context) {
	smileySvc := service.NewSmileyService(repository.NewSmileyRepo(h.db))
	smilies, err := smileySvc.List()
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OK(c, smilies)
}

func (h *AdminHandler) VerifyEmail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.userSvc.VerifyEmail(id); err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}
	errcode.OKMsg(c, "email verified")
}
