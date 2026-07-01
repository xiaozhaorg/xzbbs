package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xiaozhaorg/xzbbs/internal/config"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
)

type Claims struct {
	UserID  uint64 `json:"user_id"`
	GroupID uint   `json:"group_id"`
	jwt.RegisteredClaims
}

// Auth requires a valid JWT token
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			errcode.Fail(c, http.StatusUnauthorized, errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			errcode.Fail(c, http.StatusUnauthorized, errcode.ErrTokenInvalid)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("group_id", claims.GroupID)
		c.Next()
	}
}

// OptionalAuth sets user info if token exists, but doesn't block
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			claims, err := ParseToken(token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("group_id", claims.GroupID)
			}
		}
		c.Next()
	}
}

// AdminOnly requires admin group
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		gid, exists := c.Get("group_id")
		if !exists {
			errcode.Fail(c, http.StatusUnauthorized, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		if gid.(uint) != 1 {
			errcode.Fail(c, http.StatusForbidden, errcode.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

// ModOnly requires moderator group (gid <= 5)
func ModOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		gid, exists := c.Get("group_id")
		if !exists {
			errcode.Fail(c, http.StatusUnauthorized, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		if gid.(uint) > 5 {
			errcode.Fail(c, http.StatusForbidden, errcode.ErrNoPermission)
			c.Abort()
			return
		}
		c.Next()
	}
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Global.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func extractToken(c *gin.Context) string {
	// Bearer token
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	// Query param fallback
	return c.Query("token")
}

// GetUserID returns current user ID from context
func GetUserID(c *gin.Context) uint64 {
	v, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return v.(uint64)
}

// GetGroupID returns current user group from context
func GetGroupID(c *gin.Context) uint {
	v, exists := c.Get("group_id")
	if !exists {
		return 100 // guest
	}
	return v.(uint)
}
