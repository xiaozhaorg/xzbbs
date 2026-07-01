package model

type Group struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `gorm:"type:varchar(32);not null" json:"name"`
	CreditsFrom  int    `gorm:"default:0" json:"credits_from"`
	CreditsTo    int    `gorm:"default:0" json:"credits_to"`
	AllowRead    bool   `gorm:"default:true" json:"allow_read"`
	AllowThread  bool   `gorm:"default:false" json:"allow_thread"`
	AllowPost    bool   `gorm:"default:false" json:"allow_post"`
	AllowAttach  bool   `gorm:"default:false" json:"allow_attach"`
	AllowDown    bool   `gorm:"default:true" json:"allow_down"`
	AllowTop     bool   `gorm:"default:false" json:"allow_top"`
	AllowUpdate  bool   `gorm:"default:false" json:"allow_update"`
	AllowDelete  bool   `gorm:"default:false" json:"allow_delete"`
	AllowMove    bool   `gorm:"default:false" json:"allow_move"`
	AllowBanUser bool   `gorm:"default:false" json:"allow_ban_user"`
	IsAdmin      bool   `gorm:"default:false" json:"is_admin"`
}

func (Group) TableName() string {
	return "groups"
}
