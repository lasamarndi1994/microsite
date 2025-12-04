package request

type UserRequest struct {
	UserName         string `json:"user_name"  binding:"required"`
	AboutMe          string `json:"about_me"`
	BusinessName     string `json:"business_name"`
	BusinessLocation string `json:"business_location"`
	UserAvatar       string `json:"user_avatar"`
}
