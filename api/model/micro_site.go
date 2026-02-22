package model

import (
	"micro-site/internal/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MicroSite struct {
	Id              uint64    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Uuid            uuid.UUID `json:"uuid" gorm:"size:256;not nul;unique"`
	UserId          uint64    `json:"user_id" gorm:"not null"`
	User            User      `json:"user" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	FullName        string    `json:"full_name" gorm:"size:256"`
	Title           string    `json:"title" gorm:"size:256;not nul"`
	SubTitle        string    `json:"sub_title" gorm:"size:256"`
	Slug            string    `json:"slug" gorm:"size:256;not nul;unique"`
	BusinessName    string    `json:"business_name" gorm:"size:256"`
	Location        string    `json:"location" gorm:"size:256"`
	Description     string    `json:"description"`
	Status          string    `json:"status" gorm:"type:enum('Active','Pending','Approved','Rejected','Draft');default:'Pending';not null"`
	ViewCount       uint64    `json:"view_count" gorm:"default:0"`
	EngagementCount uint64    `json:"engagement_count" gorm:"default:0"`

	// IsDraft         string       `json:"is_draft" gorm:"type:boolean;default:0"`
	AvatarIcon      string       `json:"avatar_icon" gorm:"size:256"`
	BannerImage     string       `json:"banner_image" gorm:"size:256"`
	AdminId         uint64       `json:"admin_id"`
	RejectionReason string       `json:"rejection_reason"`
	ApproveMessage  string       `json:"approve_message"`
	Services        []Service    `json:"services" gorm:"foreignKey:MicroSiteID;constraint:OnDelete:CASCADE;"`
	SocialLinks     []SocialLink `json:"social_links" gorm:"foreignKey:MicroSiteID;constraint:OnDelete:CASCADE;"`
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
