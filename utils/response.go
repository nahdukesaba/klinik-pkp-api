package utils

// Response represents standard API response structure
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(data any) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// SuccessMessageResponse creates a success response with message
func SuccessMessageResponse(message string, data any) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(error string) Response {
	return Response{
		Success: false,
		Error:   error,
	}
}
