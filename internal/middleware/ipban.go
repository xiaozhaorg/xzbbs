package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

func IPBan(ipBanSvc service.IPBanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			c.Next()
			return
		}

		banned, err := ipBanSvc.Check(ip)
		if err != nil {
			c.Next()
			return
		}
		if banned {
			c.Abort()
			errcode.Fail(c, http.StatusForbidden, errcode.ErrIPBanned)
			return
		}
		c.Next()
	}
}
