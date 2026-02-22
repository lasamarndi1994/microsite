package request

type MobileRequest struct {
	MobileNumber string `json:"mobile_number" binding:"required"`
}
