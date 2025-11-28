package model

import "time"

type SocialLink struct {
	Name      string     `json:"name"`
	Link      string     `json:"link"`
	UserId    []User     `json:"user_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
