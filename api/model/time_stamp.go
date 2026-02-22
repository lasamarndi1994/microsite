package model

import "time"

type TimeStamp struct {
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
