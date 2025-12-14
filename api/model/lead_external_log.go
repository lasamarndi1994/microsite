package model

type LeadExternalLog struct {
	Id              uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	LeadId          uint64 `json:"lead_id" gorm:"not null"`
	ApiUrl          string `json:"api_url" gorm:"size:512"`
	RequestPayload  string `json:"request_payload" gorm:"type:text"`
	ResponsePayload string `json:"response_payload" gorm:"type:text"`
	HttpStatusCode  int    `json:"http_status_code"`
	Status          string `json:"status" gorm:"size:50"` // Success, Failed
	ErrorMessage    string `json:"error_message" gorm:"type:text"`
	TimeStamp
}

/*
* Set table name for LeadExternalLog
* @return string
 */
func (LeadExternalLog) TableName() string {
	return "lead_external_logs"
}
