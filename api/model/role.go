package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Id        uint   `json:"id" gorm:"unique;autoIncrement;primaryKey"`
	Name      string `json:"name" gorm:"unique"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt
}