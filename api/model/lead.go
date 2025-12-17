package model

import "time"

type Lead struct {
	Id           uint64     `json:"id" gorm:"unique;autoIncrement;primaryKey"`
	Name         string     `json:"name" gorm:"size:256;not null"`
	Email        string     `json:"email" gorm:"size:256;not null"`
	MobileNumber string     `json:"mobile_number" gorm:"size:20;not null"`
	MicrositeId  uint64     `json:"microsite_id" gorm:"not null"`
	UserId       uint64     `json:"user_id" gorm:"not null"`
	LeadStatus   bool       `json:"lead_status" gorm:"boolean;default:true"`
	LeadMessage  string     `json:"lead_message" gorm:"size:256;default: null"`
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
