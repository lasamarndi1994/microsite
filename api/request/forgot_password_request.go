package request

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Otp         int32  `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
