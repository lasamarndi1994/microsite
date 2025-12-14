package model

type MicrositeVisitor struct {
	Id          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId      uint64    `json:"user_id" gorm:"not null"`
	User        User      `json:"user" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
	MicroSiteId uint64    `json:"micro_site_id" gorm:"not null"`
	MicroSite   MicroSite `json:"micro_site" gorm:"foreignKey:MicroSiteId;constraint:OnDelete:CASCADE;"`
	IpAddress   string    `json:"ip_address" gorm:"size:256;not null"`

	TimeStamp
}

func (MicrositeVisitor) TableName() string {
	return "microsite_visitors"
}
