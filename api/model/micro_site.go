package model

import (
	"micro-site/internal/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MicroSite struct {
	Id          uint64       `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Uuid        uuid.UUID    `json:"uuid" gorm:"size:256;not nul;unique"`
	UserId      uint64       `json:"user_id" gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE;size:256"`
	FullName    string       `json:"full_name" gorm:"size:256"`
	Title       string       `json:"titile" gorm:"size:256;not nul"`
	Slug        string       `json:"slug" gorm:"size:256;not nul;unique"`
	About       string       `json:"about"`
	Status      string       `json:"status" gorm:"type:enum('Active','Pending','Approved','Rejected');default:'Pending';not null"`
	IsDraft     string       `json:"is_draft" gorm:"type:boolean;default:0"`
	AvatarIcon  string       `json:"avatar_icon" gorm:"size:256"`
	BannerImage string       `json:"banner_image" gorm:"size:256"`
	AdminId     uint64       `json:"admin_id"`
	Comment     string       `json:"comment"`
	Services    []Service    `gorm:"foreignKey:MicroSiteID;constraint:OnDelete:CASCADE;"`
	SocialLinks []SocialLink `gorm:"foreignKey:MicroSiteID;constraint:OnDelete:CASCADE;"`
	TimeStamp
}

func (u *MicroSite) BeforeCreate(tx *gorm.DB) (err error) {
	u.Uuid = uuid.New() // NOT uuid.UUID{}
	u.Slug = helper.GenerateUniqueSlug(tx, u.Title, &MicroSite{}, "slug")
	return
}

func (MicroSite) TableName() string {
	return "micro_sites"
}
