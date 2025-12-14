package model

type UserActivity struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       *uint  `gorm:"index"` // nullable for guests
	Method       string `gorm:"size:10"`
	Path         string `gorm:"size:255"`
	IP           string `gorm:"size:45"`
	UserAgent    string `gorm:"size:255"`
	StatusCode   int
	ResponseTime int64 // milliseconds
	TimeStamp
}

/*
* Set table name for UserActivity
* @return string
 */
func (UserActivity) TableName() string {
	return "user_activities"
}
