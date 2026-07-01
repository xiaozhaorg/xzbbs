package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaozhaorg/xzbbs/internal/config"
	"github.com/xiaozhaorg/xzbbs/internal/handler"
	"github.com/xiaozhaorg/xzbbs/internal/middleware"
	"github.com/xiaozhaorg/xzbbs/internal/model"
	"github.com/xiaozhaorg/xzbbs/internal/plugin"
	"github.com/xiaozhaorg/xzbbs/internal/repository"
	"github.com/xiaozhaorg/xzbbs/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfgPath := config.FindConfig()
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if !filepath.IsAbs(cfg.Upload.Path) {
		cfg.Upload.Path = config.ResolvePath(cfg.Upload.Path)
	}
	if cfg.Database.Driver == "sqlite" && !filepath.IsAbs(cfg.Database.DSN) {
		cfg.Database.DSN = config.ResolvePath(cfg.Database.DSN)
	}

	// Connect database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto migrate & seed
	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	if err := model.Seed(db); err != nil {
		log.Fatalf("Failed to seed: %v", err)
	}

	// Create upload directories
	os.MkdirAll(cfg.Upload.Path+"/avatars", 0755)

	// Initialize layers
	userRepo := repository.NewUserRepo(db)
	forumRepo := repository.NewForumRepo(db)
	threadRepo := repository.NewThreadRepo(db)
	postRepo := repository.NewPostRepo(db)
	attachRepo := repository.NewAttachRepo(db)
	favRepo := repository.NewFavoriteRepo(db)
	notifRepo := repository.NewNotificationRepo(db)
	postEditRepo := repository.NewPostEditRepo(db)
	searchRepo := repository.NewSearchRepo(db)
	onlineRepo := repository.NewOnlineRepo(db)
	pmRepo := repository.NewPMRepo(db)
	creditRepo := repository.NewCreditRepo(db)
	ipBanRepo := repository.NewIPBanRepo(db)
	modLogRepo := repository.NewModLogRepo(db)
	emailVerifyRepo := repository.NewEmailVerifyRepo(db)
	passwordResetRepo := repository.NewPasswordResetRepo(db)

	userSvc := service.NewUserService(userRepo, onlineRepo)
	forumSvc := service.NewForumService(forumRepo)
	creditSvc := service.NewCreditService(creditRepo, userRepo)
	emailVerifySvc := service.NewEmailVerifyService(emailVerifyRepo, userSvc)
	passwordResetSvc := service.NewPasswordResetService(passwordResetRepo, userSvc)
	notifSvc := service.NewNotificationService(notifRepo)
	ipBanSvc := service.NewIPBanService(ipBanRepo)
	modLogSvc := service.NewModLogService(modLogRepo)
	attachSvc := service.NewAttachService(attachRepo)
	threadSvc := service.NewThreadService(db, threadRepo, postRepo, forumRepo, userRepo, notifSvc, creditSvc)
	postSvc := service.NewPostService(db, postRepo, threadRepo, forumRepo, userRepo, threadSvc, postEditRepo, creditSvc, attachSvc)
	favSvc := service.NewFavoriteService(favRepo, threadRepo)
	searchSvc := service.NewSearchService(searchRepo)

	authH := handler.NewAuthHandler(userSvc, emailVerifySvc, passwordResetSvc)
	userH := handler.NewUserHandler(userSvc, threadSvc, postSvc)
	forumH := handler.NewForumHandler(forumSvc)
	creditH := handler.NewCreditHandler(creditSvc)
	threadH := handler.NewThreadHandler(threadSvc, postSvc)
	postH := handler.NewPostHandler(postSvc, threadSvc)
	rssH := handler.NewRSSHandler(threadSvc)
	postEditH := handler.NewPostEditHandler(postEditRepo)
	attachH := handler.NewAttachHandler(attachSvc)
	modH := handler.NewModHandler(threadSvc, userSvc, modLogSvc)
	adminH := handler.NewAdminHandler(userSvc, forumSvc, threadSvc, postSvc, db)
	favH := handler.NewFavoriteHandler(favSvc)
	notifH := handler.NewNotificationHandler(notifSvc)
	searchH := handler.NewSearchHandler(searchSvc)
	pmSvc := service.NewPMService(pmRepo)
	pmH := handler.NewPMHandler(pmSvc)

	// Boot plugins
	(&plugin.LoggingPlugin{}).Boot()
	(&plugin.EditorPlugin{}).Boot()

	// Setup Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Configure trusted proxies for correct ClientIP
	if len(cfg.Server.TrustedProxies) > 0 {
		_ = r.SetTrustedProxies(cfg.Server.TrustedProxies)
	} else {
		_ = r.SetTrustedProxies(nil)
	}

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.IPBan(ipBanSvc))
	r.Use(middleware.RateLimit(10, 30)) // 10 req/s, burst 30

	// Static files (uploads)
	r.Static("/uploads", cfg.Upload.Path)

	// RSS (public)
	r.GET("/rss.xml", rssH.Feed)

	// SEO-friendly .htm routes (public)
	r.GET("/index.htm", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/") })
	r.GET("/forum-:id.htm", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, "/forums/"+id)
	})
	r.GET("/thread-:id.htm", func(c *gin.Context) {
		id := c.Param("id")
		c.Redirect(http.StatusMovedPermanently, "/threads/"+id)
	})

	// API routes
	api := r.Group("/api")
	{
		// Auth (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
		}
		// Email verification
		email := api.Group("/email")
		email.Use(middleware.Auth())
		{
			email.POST("/verify-request", authH.RequestEmailVerification)
			email.POST("/verify-confirm", authH.ConfirmEmailVerification)
		}

		// Password reset (public)
		password := api.Group("/password")
		{
			password.POST("/reset-request", authH.RequestPasswordReset)
			password.POST("/reset-confirm", authH.ConfirmPasswordReset)
		}

		// Users
		users := api.Group("/users")
		{
			users.GET("/:id", userH.GetProfile)
			users.GET("/:id/threads", userH.GetUserThreads)
			users.GET("/:id/posts", userH.GetUserPosts)
			users.PUT("/:id", middleware.Auth(), userH.UpdateProfile)
			users.PUT("/:id/password", middleware.Auth(), userH.ChangePassword)
			users.POST("/:id/avatar", middleware.Auth(), userH.UploadAvatar)
		}

		// Credits
		credits := api.Group("/credits")
		credits.Use(middleware.Auth())
		{
			credits.GET("/logs", creditH.ListLogs)
		}

		// Forums (public read)
		forums := api.Group("/forums")
		{
			forums.GET("", forumH.List)
			forums.GET("/:id", forumH.Get)
			forums.GET("/:id/threads", threadH.ListByForum)
		}

		// Threads (public list, auth for create)
		threads := api.Group("/threads")
		{
			threads.GET("", threadH.ListAll)
			threads.POST("", middleware.Auth(), threadH.Create)
			threads.GET("/:id", middleware.OptionalAuth(), threadH.Get)
			threads.PUT("/:id", middleware.Auth(), threadH.Update)
			threads.DELETE("/:id", middleware.Auth(), threadH.Delete)
			// Replies under thread
			threads.POST("/:tid/posts", middleware.Auth(), postH.Create)
		}

		// Posts
		posts := api.Group("/posts")
		posts.Use(middleware.Auth())
		{
			posts.PUT("/:id", postH.Update)
			posts.DELETE("/:id", postH.Delete)
			posts.GET("/:id/edits", postEditH.GetPostEdits)
		}

		// Attachments
		attach := api.Group("/attachments")
		{
			attach.POST("", middleware.Auth(), attachH.Upload)
			attach.GET("/:id", attachH.Download)
		}

		// Moderation
		mod := api.Group("/mod")
		mod.Use(middleware.Auth(), middleware.ModOnly())
		{
			mod.PUT("/threads/top", modH.SetTop)
			mod.PUT("/threads/close", modH.SetClosed)
			mod.PUT("/threads/move", modH.Move)
			mod.DELETE("/users/:id", modH.BanUser)
		}

		// Admin
		admin := api.Group("/admin")
		admin.Use(middleware.Auth(), middleware.AdminOnly())
		{
			admin.GET("/stats", adminH.Stats)
			admin.GET("/groups", adminH.ListGroups)
			admin.PUT("/groups/:id", adminH.UpdateGroup)
			admin.GET("/users", adminH.ListUsers)
			admin.PUT("/users/:id", adminH.UpdateUser)
			admin.DELETE("/users/:id", adminH.DeleteUser)
			admin.POST("/forums", forumH.Create)
			admin.PUT("/forums/:id", forumH.Update)
			admin.DELETE("/forums/:id", forumH.Delete)
			admin.GET("/settings", adminH.GetSettings)
			admin.PUT("/settings", adminH.UpdateSettings)
			// IP Bans
			admin.POST("/ip-bans", adminH.BanIP)
			admin.DELETE("/ip-bans/:id", adminH.UnbanIP)
			admin.GET("/ip-bans", adminH.ListIPBans)
			admin.GET("/ip-bans/check", adminH.CheckIPBan)
			// Smilies
			admin.GET("/smilies", adminH.ListSmilies)
			// Email verification
			admin.POST("/users/:id/verify-email", adminH.VerifyEmail)
		}

		// Auth protected: get current user
		api.GET("/auth/me", middleware.Auth(), authH.Me)

		// Search
		api.GET("/search", searchH.Search)

		// Favorites
		fav := api.Group("/favorites")
		fav.Use(middleware.Auth())
		{
			fav.POST("/threads/:id", favH.Toggle)
			fav.GET("/threads", favH.List)
			fav.GET("/threads/:id/check", favH.Check)
		}

		// Notifications
		notifs := api.Group("/notifications")
		notifs.Use(middleware.Auth())
		{
			notifs.GET("", notifH.List)
			notifs.GET("/unread-count", notifH.UnreadCount)
			notifs.POST("/read", notifH.MarkRead)
			notifs.POST("/read-all", notifH.MarkAllRead)
			notifs.DELETE("/:id", notifH.Delete)
		}

		// Private Messages
		pms := api.Group("/pms")
		pms.Use(middleware.Auth())
		{
			pms.POST("", pmH.Send)
			pms.GET("", pmH.List)
			pms.GET("/conversations/:otherId", pmH.Conversation)
			pms.POST("/read/:otherId", pmH.MarkRead)
			pms.GET("/unread-count", pmH.UnreadCount)
			pms.DELETE("/:id", pmH.Delete)
		}

		// Smilies (public)
		api.GET("/smilies", adminH.ListSmilies)

		// Online users
		api.GET("/online", middleware.OptionalAuth(), adminH.OnlineUsers)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start online cleanup goroutine (every 5 min)
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			_ = onlineRepo.Cleanup()
		}
	}()

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 XzBBS server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch cfg.Database.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.Database.DSN)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.DSN)
	default:
		return nil, fmt.Errorf("unsupported db driver: %s", cfg.Database.Driver)
	}

	logLevel := logger.Warn
	if cfg.Server.Mode == "debug" {
		logLevel = logger.Info
	}

	return gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
}
