package request

type UpdatePasswordRequest struct {
	MobileNumber    string `json:"mobile_number" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
