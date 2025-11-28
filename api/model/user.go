package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id               uint64     `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Uuid             uuid.UUID  `json:"uuid" gorm:"size:256;not nul;unique"`
	FirstName        string     `json:"first_name" gorm:"size:256"`
	LastName         string     `json:"last_name" gorm:"size:256;not null"`
	Email            string     `json:"email" gorm:"unique;not null;unique;unique_email"`
	MobileNumber     int        `json:"mobile_number" gorm:"unique"`
	Status           string     `json:"status" gorm:"type:enum('Active','Pending','Approved','Deactive');default:'Active';not null"`
	AboutMe          string     `json:"about_me"`
	BusinessName     string     `json:"business_name"`
	BusinessLocation string     `json:"business_location" gorm:"size:256"`
	UserAvatar       string     `json:"user_avatar" gorm:"size:256"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Uuid = uuid.New() // NOT uuid.UUID{}
	//u.EmployeeId = strconv.Itoa(100000 + rand.Intn(900000))
	return
}
