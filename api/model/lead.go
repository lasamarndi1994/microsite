package model

import "time"

type Lead struct {
	Id           uint64     `json:"id" gorm:"unique;autoIncrement;primaryKey"`
	Name         string     `json:"name" gorm:"size:256;not null"`
	Email        string     `json:"email" gorm:"size:256;not null"`
	MobileNumber string     `json:"mobile_number" gorm:"size:20;not null"`
	MicrositeId  uint64     `json:"microsite_id" gorm:"not null"`
	Slug         string     `json:"slug" gorm:"size:256"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

/*
* Set table name for Lead
* @return string
 */
func (Lead) TableName() string {
	return "leads"
}
