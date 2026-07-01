package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/errcode"
	"github.com/xiaozhaorg/xzbbs/internal/pkg/pagination"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type CreditHandler struct {
	creditSvc *service.CreditService
}

func NewCreditHandler(creditSvc *service.CreditService) *CreditHandler {
	return &CreditHandler{creditSvc: creditSvc}
}

// ListLogs returns credit transaction history for current user
// GET /api/credits/logs
func (h *CreditHandler) ListLogs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, pageSize := pagination.Normalize(c.Query("page"), c.Query("page_size"))

	logs, total, err := h.creditSvc.ListLogs(userID, page, pageSize)
	if err != nil {
		errcode.Fail(c, 500, errcode.ErrInternal)
		return
	}

	result := map[string]interface{}{
		"items": logs, "total": total, "page": page,
		"pages": (total + int64(pageSize) - 1) / int64(pageSize),
	}
	errcode.OK(c, result)
}
