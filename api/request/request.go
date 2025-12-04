package request

type UserAvatarRequest struct {
	UserAvatar string `json:"user_avatar" binding:"required"`
}
