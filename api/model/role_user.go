package model

import (
	"time"

	"gorm.io/gorm"
)

type RoleUser struct {
	RoleId int `gorm:"primaryKey"`
	UserId int `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}