package model

type Admin struct {
	Id       uint64 `json:"id" gorm:"unique;autoIncrement;primaryKey"`
	Email    string `json:"email" gorm:"unique;unique_email;not null"`
	Password string `json:"password" gorm:"unique;not null"`
	Status   bool   `json:"status" gorm:"default: 1"`
	TimeStamp
}

/*
* Set table name for Admin
* @return string
 */
func (Admin) TableName() string {
	return "admins"
}
