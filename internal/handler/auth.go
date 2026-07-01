package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xiaozhaorg/xzbbs/internal/config"
	"github.com/xiaozhaorg/xzbbs/internal/dto"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type AuthHandler struct {
	userSvc          *service.UserService
	emailVerifySvc   *service.EmailVerifyService
	passwordResetSvc *service.PasswordResetService
}

func NewAuthHandler(svc *service.UserService, emailVerifySvc *service.EmailVerifyService, passwordResetSvc *service.PasswordResetService) *AuthHandler {
	return &AuthHandler{userSvc: svc, emailVerifySvc: emailVerifySvc, passwordResetSvc: passwordResetSvc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	user, err := h.userSvc.Register(req.Username, req.Email, req.Password)
	if err != nil {
		errcode.FailMsg(c, http.StatusConflict, errcode.ErrUserExists, err.Error())
		return
	}

	token := issueToken(user)
	errcode.OK(c, gin.H{
		"user":       safeUser(user),
		"token":      token,
		"expires_in": config.Global.JWT.ExpireHour * 3600,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	user, token, err := h.userSvc.Login(req.Account, req.Password, c.ClientIP())
	if err != nil {
		errcode.Fail(c, http.StatusUnauthorized, errcode.ErrInvalidCredentials)
		return
	}

	errcode.OK(c, gin.H{
		"user":       safeUser(user),
		"token":      token,
		"expires_in": config.Global.JWT.ExpireHour * 3600,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrUserNotFound)
		return
	}
	errcode.OK(c, safeUser(user))
}

// issueToken generates a JWT for the given user
func issueToken(user *model.User) string {
	cfg := config.Global
	claims := middleware.Claims{
		UserID:  user.ID,
		GroupID: user.GroupID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWT.ExpireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(cfg.JWT.Secret))
	return signed
}

// safeUser strips sensitive fields
func safeUser(u *model.User) gin.H {
	return gin.H{
		"id":             u.ID,
		"username":       u.Username,
		"email":          u.Email,
		"group_id":       u.GroupID,
		"avatar":         u.Avatar,
		"threads":        u.Threads,
		"posts":          u.Posts,
		"credits":        u.Credits,
		"level":          u.Level,
		"signature":      u.Signature,
		"email_verified": u.EmailVerified,
		"last_login":     u.LastLogin,
		"created_at":     u.CreatedAt,
	}
}

// RequestEmailVerification creates a verification token and sends email
// POST /api/email/verify-request
func (h *AuthHandler) RequestEmailVerification(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		errcode.Fail(c, http.StatusNotFound, errcode.ErrUserNotFound)
		return
	}
	if user.EmailVerified {
		errcode.FailMsg(c, http.StatusBadRequest, errcode.ErrBadRequest, "email already verified")
		return
	}

	token, _, err := h.emailVerifySvc.CreateToken(userID, user.Email)
	if err != nil {
		errcode.Fail(c, http.StatusInternalServerError, errcode.ErrInternal)
		return
	}

	emailSvc := service.NewEmailService()
	if err := emailSvc.SendVerificationEmail(user.Email, token); err != nil {
		errcode.OK(c, gin.H{"token": token, "email": user.Email, "dev_mode": true, "note": "email not configured, token shown for dev"})
		return
	}
	errcode.OK(c, gin.H{"token": token, "email": user.Email, "dev_mode": false, "note": "verification email sent"})
}

// ConfirmEmailVerification marks email as verified via token
// POST /api/email/verify-confirm
func (h *AuthHandler) ConfirmEmailVerification(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	if err := h.emailVerifySvc.ConfirmToken(req.Token); err != nil {
		errcode.Fail(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}
	errcode.OKMsg(c, "verified")
}

// RequestPasswordReset sends a password reset link to email
// POST /api/password/reset-request
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	token, err := h.passwordResetSvc.RequestReset(req.Email)
	if err != nil {
		errcode.OKMsg(c, "if the email exists, a reset link has been sent")
		return
	}

	emailSvc := service.NewEmailService()
	if err := emailSvc.SendPasswordResetEmail(req.Email, token); err != nil {
		// Dev mode: return token for testing
		errcode.OK(c, gin.H{"token": token, "dev_mode": true, "note": "email not configured, token shown for dev"})
		return
	}
	errcode.OKMsg(c, "if the email exists, a reset link has been sent")
}

// ConfirmPasswordReset resets password via token
// POST /api/password/reset-confirm
func (h *AuthHandler) ConfirmPasswordReset(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errcode.FailValidation(c, err.Error())
		return
	}

	if err := h.passwordResetSvc.ConfirmReset(req.Token, req.NewPassword); err != nil {
		errcode.FailMsg(c, http.StatusBadRequest, errcode.ErrBadRequest, err.Error())
		return
	}
	errcode.OKMsg(c, "password reset successful")
}
