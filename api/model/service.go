package model

type Service struct {
	Id          uint64 `json:"id" gorm:"unique;autoIncrement;primaryKey"`
	Name        string `json:"name" gorm:"not null;size:256"`
	UserId      uint64 `json:"user_id" gorm:"not null"`
	MicroSiteId uint64 `json:"micro_site_id"`
	Status      bool   `json:"status" gorm:"default: 1"`
	TimeStamp
}

func (Service) TableName() string {
	return "services"
}
