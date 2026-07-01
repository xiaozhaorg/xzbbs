package pagination

import (
	"math"
	"strconv"
)

type PageQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type PageResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Pages    int         `json:"pages"`
}

func (q *PageQuery) Normalize(defaultSize int) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = defaultSize
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
}

func (q *PageQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

func Normalize(pageStr, pageSizeStr string) (int, int) {
	page := 1
	pageSize := 20
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
		pageSize = ps
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func NewResult(items interface{}, total int64, page, pageSize int) *PageResult {
	pages := int(math.Ceil(float64(total) / float64(pageSize)))
	return &PageResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}
}
