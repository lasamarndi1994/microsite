package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id               uint64    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Uuid             uuid.UUID `json:"uuid" gorm:"size:256;not nul;unique"`
	UserName         string    `json:"user_name" gorm:"size:256;not null"`
	Email            string    `json:"email" gorm:"unique;unique_email"`
	MobileNumber     int       `json:"mobile_number" gorm:"unique;default null"`
	Status           string    `json:"status" gorm:"type:enum('Active','Pending','Approved','Deactive');default:'Active';not null"`
	AboutMe          string    `json:"about_me"`
	BusinessName     string    `json:"business_name"  gorm:"size:256"`
	BusinessLocation string    `json:"business_location" gorm:"size:256"`
	UserAvatar       string    `json:"user_avatar" gorm:"size:256"`
	TimeStamp
}

/*
* BeforeCreate hook to generate UUID
* @param tx *gorm.DB
* @return error
 */
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Uuid = uuid.New() // NOT uuid.UUID{}
	return
}

/*
* Set table name for User
* @return string
 */
func (User) TableName() string {
	return "users"
}
