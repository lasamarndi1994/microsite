package service

type FieldResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(message string, data ...interface{}) interface{} {
	response := &FieldResponse{
		Status:  true,
		Message: message,
	}
	if len(data) > 0 {
		response.Data = data[0]
	}
	return response
}

func ErrorResponse(message string) interface{} {
	return &FieldResponse{
		Status:  false,
		Message: message,
	}
}
