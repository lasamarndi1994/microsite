package model

import "time"

type SocialLink struct {
	Id          uint64     `json:"id" gorm:"unique;autoIncrement;primaryKey"`
	Name        string     `json:"name" gorm:"size:256"`
	Url         string     `json:"url" gorm:"size:256"`
	Type        string     `json:"link" gorm:"size:256"`
	UserId      uint64     `json:"user_id" gorm:"not null;size:256"`
	MicroSiteId uint64     `json:"micro_site_id"`
	CreatedAt   *time.Time `json:"-"`
	UpdatedAt   *time.Time `json:"-"`
}

/*
* Set table name for SocialLink
* @return string
 */
func (SocialLink) TableName() string {
	return "social_medias"
}
