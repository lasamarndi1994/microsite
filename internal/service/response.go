package service

type FieldResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

/*
* Create success response
* @param message string
* @param data ...interface{}
* @return interface{}
 */
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

/*
* Create error response
* @param message string
* @return interface{}
 */
func ErrorResponse(message string) interface{} {
	return &FieldResponse{
		Status:  false,
		Message: message,
	}
}
