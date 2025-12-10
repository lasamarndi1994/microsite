package request

type MicrositeRequest struct {
	FullName     string           `json:"full_name" binding:"required"`
	Title        string           `json:"title" binding:"required"`
	SubTitle     string           `json:"sub_title"`
	BusinessName string           `json:"business_name" binding:"required"`
	Location     string           `json:"location" binding:"required"`
	Description  string           `json:"description" binding:"required"`
	IsDraft      string           `json:"is_draft"`
	AvatarIcon   string           `json:"avatar_icon"`
	BannerImage  string           `json:"banner_image"`
	ServicesName []ServiceRequest `json:"services_name" binding:"required"`
	SocialLink   []SocialRequest  `json:"social_link" binding:"required"`
}

type ServiceRequest struct {
	MicroSiteId uint   `json:"micro_site_id"`
	Name        string `json:"name"`
}

type SocialRequest struct {
	MicroSiteId uint   `json:"micro_site_id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
}
