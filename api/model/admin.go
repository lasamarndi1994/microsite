package model

type Admin struct {
	Id       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   bool   `json:"status"`
	TimeStamp
}

type AdminMigration struct {
	Id       uint64 `json:"id" gorm:"unique;autoIncrement;primaryKey"`
	Email    string `json:"email" gorm:"unique;unique_email;not null"`
	Password string `json:"password" gorm:"unique;not null"`
	Status   bool   `json:"status" gorm:"default: 1"`
	TimeStamp
}

/*
* Set table name for AdminMigration
* @return string
 */
func (AdminMigration) TableName() string {
	return "admins"
}
