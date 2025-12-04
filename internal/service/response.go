package service

type FieldResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func SuccessResponse(message string) interface{} {
	return &FieldResponse{
		Status:  true,
		Message: message,
	}
}

func ErrorResponse(message string) interface{} {
	return &FieldResponse{
		Status:  false,
		Message: message,
	}
}
