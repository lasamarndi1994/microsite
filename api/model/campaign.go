package model

import (
	"time"

	"github.com/google/uuid"
)

type Campaign struct {
	Id          uint64     `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Uuid        uuid.UUID  `json:"uuid" gorm:"size:256;not nul;unique"`
	UserId      []User     `json:"user_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title       string     `json:"titile" gorm:"size:256;not nul;unique"`
	Slug        string     `json:"slug" gorm:"size:256;not nul;unique"`
	Status      string     `json:"status" gorm:"type:enum('Active','Pending','Approved','Deactive');default:'Pending';not null"`
	BannerImage string     `json:"banner_image"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
