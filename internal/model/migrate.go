package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AutoMigrate creates all tables
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Group{},
		&Forum{},
		&ForumPermission{},
		&Thread{},
		&Post{},
		&Attachment{},
		&ModLog{},
		&Favorite{},
		&Notification{},
		&PostEdit{},
		&OnlineUser{},
		&PrivateMessage{},
		&IPBan{},
		&Smiley{},
		&EmailVerify{},
		&CreditLog{},
	)
}

// Seed inserts default data if tables are empty
func Seed(db *gorm.DB) error {
	// Seed groups
	var count int64
	db.Model(&Group{}).Count(&count)
	if count == 0 {
		groups := []Group{
			{ID: 1, Name: "管理员", AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true, AllowTop: true, AllowUpdate: true, AllowDelete: true, AllowMove: true, AllowBanUser: true, IsAdmin: true},
			{ID: 2, Name: "超级版主", AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true, AllowTop: true, AllowUpdate: true, AllowDelete: true, AllowMove: true, AllowBanUser: true},
			{ID: 4, Name: "版主", AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true, AllowTop: true, AllowUpdate: true, AllowDelete: true, AllowMove: true, AllowBanUser: true},
			{ID: 5, Name: "实习版主", AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true, AllowTop: true, AllowUpdate: true},
			{ID: 6, Name: "待验证", AllowRead: true, AllowPost: true, AllowDown: true},
			{ID: 7, Name: "禁止", AllowRead: false},
			{ID: 100, Name: "游客", AllowRead: true, AllowDown: true},
			{ID: 101, Name: "一级用户", CreditsFrom: 0, CreditsTo: 50, AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true},
			{ID: 102, Name: "二级用户", CreditsFrom: 50, CreditsTo: 200, AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true},
			{ID: 103, Name: "三级用户", CreditsFrom: 200, CreditsTo: 1000, AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true},
			{ID: 104, Name: "四级用户", CreditsFrom: 1000, CreditsTo: 10000, AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true},
			{ID: 105, Name: "五级用户", CreditsFrom: 10000, CreditsTo: 10000000, AllowRead: true, AllowThread: true, AllowPost: true, AllowAttach: true, AllowDown: true},
		}
		if err := db.Create(&groups).Error; err != nil {
			return err
		}
	}

	// Seed admin user
	db.Model(&User{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := User{
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(hash),
			GroupID:  1,
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	}

	// Seed default forum
	db.Model(&Forum{}).Count(&count)
	if count == 0 {
		forum := Forum{
			Name:        "默认版块",
			Description: "这是默认版块，可在管理后台修改",
			SortOrder:   1,
		}
		if err := db.Create(&forum).Error; err != nil {
			return err
		}
	}

	// Seed default smilies
	db.Model(&Smiley{}).Count(&count)
	if count == 0 {
		smilies := []Smiley{
			{Code: ":)", Image: "smiley.gif", Sort: 1},
			{Code: ":(", Image: "sad.gif", Sort: 2},
			{Code: ":D", Image: "grin.gif", Sort: 3},
			{Code: ":o", Image: "shocked.gif", Sort: 4},
			{Code: ":P", Image: "tongue.gif", Sort: 5},
			{Code: ":love:", Image: "heart.gif", Sort: 6},
			{Code: ":think:", Image: "think.gif", Sort: 7},
			{Code: ":angry:", Image: "angry.gif", Sort: 8},
			{Code: ":cool:", Image: "cool.gif", Sort: 9},
			{Code: ":cry:", Image: "cry.gif", Sort: 10},
			{Code: ":ok:", Image: "ok.gif", Sort: 11},
			{Code: ":no:", Image: "no.gif", Sort: 12},
		}
		if err := db.Create(&smilies).Error; err != nil {
			return err
		}
	}

	return nil
}
