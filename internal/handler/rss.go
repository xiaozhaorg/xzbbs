package handler

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/config"
	"github.com/xiaozhaorg/xzbbs/internal/service"
)

type RSSHandler struct {
	threadSvc *service.ThreadService
}

func NewRSSHandler(ts *service.ThreadService) *RSSHandler {
	return &RSSHandler{threadSvc: ts}
}

// Feed returns RSS 2.0 feed of latest threads
// GET /rss.xml
func (h *RSSHandler) Feed(c *gin.Context) {
	forumIDStr := c.DefaultQuery("forum_id", "0")
	var forumID uint
	fmt.Sscanf(forumIDStr, "%d", &forumID)

	threads, _, _ := h.threadSvc.ListByForum(forumID, "reply", 1, 20)

	siteURL := "http://" + c.Request.Host
	siteName := config.Global.Site.Name

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(`<rss version="2.0"><channel>`)
	sb.WriteString(`<title>` + html.EscapeString(siteName) + `</title>`)
	sb.WriteString(`<link>` + siteURL + `</link>`)
	sb.WriteString(`<description>` + html.EscapeString(siteName) + ` - Latest threads</description>`)
	sb.WriteString(`<lastBuildDate>` + time.Now().Format(time.RFC1123) + `</lastBuildDate>`)

	for _, t := range threads {
		link := fmt.Sprintf("%s/threads/%d", siteURL, t.ID)
		sb.WriteString(`<item>`)
		sb.WriteString(`<title>` + html.EscapeString(t.Title) + `</title>`)
		sb.WriteString(`<link>` + link + `</link>`)
		sb.WriteString(`<guid>` + link + `</guid>`)
		if t.LastReplyAt != nil {
			sb.WriteString(`<pubDate>` + t.LastReplyAt.Format(time.RFC1123) + `</pubDate>`)
		}
		sb.WriteString(`</item>`)
	}
	sb.WriteString(`</channel></rss>`)

	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=300")
	c.String(200, sb.String())
}
