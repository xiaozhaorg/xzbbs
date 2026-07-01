package dto

// ========== Auth ==========

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"` // email or username
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"` // seconds
}

// ========== User ==========

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=2,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=64"`
}

// ========== Forum ==========

type CreateForumRequest struct {
	Name        string `json:"name" binding:"required,max=64"`
	Description string `json:"description" binding:"omitempty"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateForumRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=64"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sort_order"`
	Icon        *string `json:"icon"`
}

// ========== Thread ==========

type CreateThreadRequest struct {
	ForumID     uint   `json:"forum_id" binding:"required"`
	Title       string `json:"title" binding:"required,max=200"`
	Content     string `json:"content" binding:"required"`
	ContentType uint8  `json:"content_type"` // 0=md, 1=html
}

type UpdateThreadRequest struct {
	Title   *string `json:"title" binding:"omitempty,max=200"`
	ForumID *uint   `json:"forum_id"`
}

// ========== Post ==========

type CreatePostRequest struct {
	Content     string `json:"content" binding:"required"`
	ContentType uint8  `json:"content_type"`
	ReplyTo     uint64 `json:"reply_to"`
}

type UpdatePostRequest struct {
	Content     string `json:"content" binding:"required"`
	ContentType uint8  `json:"content_type"`
}

// ========== Moderation ==========

type ModTopRequest struct {
	ThreadIDs []uint64 `json:"thread_ids" binding:"required"`
	Top       uint8    `json:"top"` // 0=cancel, 1=forum, 2=global
}

type ModCloseRequest struct {
	ThreadIDs []uint64 `json:"thread_ids" binding:"required"`
	Closed    bool     `json:"closed"`
}

type ModMoveRequest struct {
	ThreadIDs []uint64 `json:"thread_ids" binding:"required"`
	ForumID   uint     `json:"forum_id" binding:"required"`
}

// ========== Admin ==========

type UpdateGroupRequest struct {
	Name         *string `json:"name"`
	AllowRead    *bool   `json:"allow_read"`
	AllowThread  *bool   `json:"allow_thread"`
	AllowPost    *bool   `json:"allow_post"`
	AllowAttach  *bool   `json:"allow_attach"`
	AllowDown    *bool   `json:"allow_down"`
	AllowTop     *bool   `json:"allow_top"`
	AllowUpdate  *bool   `json:"allow_update"`
	AllowDelete  *bool   `json:"allow_delete"`
	AllowMove    *bool   `json:"allow_move"`
	AllowBanUser *bool   `json:"allow_ban_user"`
}

type UpdateSettingsRequest struct {
	SiteName string `json:"site_name"`
	Brief    string `json:"brief"`
	PageSize int    `json:"page_size"`
}
