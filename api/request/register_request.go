package request

type RegisterRequest struct {
	UserName     string `json:"user_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	MobileNumber int    `json:"mobile_number" binding:"required"`
	UserCode     string `json:"user_code" binding:"required"`
}
