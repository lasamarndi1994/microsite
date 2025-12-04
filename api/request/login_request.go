package request


type LoginRequest struct {
	MobileNumber         string      `json:"mobile_number" binding:"required"`
	MobileOtp string `json:"mobile_otp" binding:"required"`
}