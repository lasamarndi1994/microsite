package request

type LeadRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	MobileNumber string `json:"mobile_number" binding:"required"`
	MicrositeId  uint64 `json:"microsite_id" binding:"required"`
}
