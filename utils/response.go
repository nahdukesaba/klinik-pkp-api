package utils

// API RESPONSE STRUCTURE
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// SUCCESS RESPONSE
func SuccessResponse(message string, data any) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ERROR RESPONSE
func ErrorResponse(err error) Response {
	return Response{
		Success: false,
		Error:   err.Error(),
	}
}
