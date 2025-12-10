package model

import (
	"micro-site/internal/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MicroSite struct {
	Id              uint64       `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Uuid            uuid.UUID    `json:"uuid" gorm:"size:256;not nul;unique"`
	UserId          uint64       `json:"user_id" gorm:"index;foreignKey:user_id;constraint:OnDelete:CASCADE;size:256"`
	FullName        string       `json:"full_name" gorm:"size:256"`
	Title           string       `json:"title" gorm:"size:256;not nul"`
	Slug            string       `json:"slug" gorm:"size:256;not nul;unique"`
	Description     string       `json:"description"`
	Status          string       `json:"status" gorm:"type:enum('Active','Pending','Approved','Rejected');default:'Pending';not null"`
	IsDraft         string       `json:"is_draft" gorm:"type:boolean;default:0"`
	AvatarIcon      string       `json:"avatar_icon" gorm:"size:256"`
	BannerImage     string       `json:"banner_image" gorm:"size:256"`
	AdminId         uint64       `json:"admin_id"`
	RejectionReason string       `json:"rejection_reason"`
	Services        []Service    `gorm:"foreignKey:MicroSiteID;constraint:OnDelete:CASCADE;"`
	SocialLinks     []SocialLink `gorm:"foreignKey:MicroSiteID;constraint:OnDelete:CASCADE;"`
	TimeStamp
}

/*
* BeforeCreate hook to generate UUID and slug
* @param tx *gorm.DB
* @return error
 */
func (u *MicroSite) BeforeCreate(tx *gorm.DB) (err error) {
	u.Uuid = uuid.New() // NOT uuid.UUID{}
	u.Slug = helper.GenerateUniqueSlug(tx, u.Title, &MicroSite{}, "slug")
	return
}

/*
* Set table name for MicroSite
* @return string
 */
func (MicroSite) TableName() string {
	return "micro_sites"
}
