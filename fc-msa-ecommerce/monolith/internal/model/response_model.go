package model

import "net/http"

// SuccessResponse defines the standard response structure for successful API calls.
type SuccessResponse struct {
	Code    int         `json:"code"`           // HTTP status code (e.g. 200)
	Message string      `json:"message"`        // Human-readable success message
	Data    interface{} `json:"data,omitempty"` // Optional response payload
}

// NewSuccessResponse creates a standard success response object.
func NewSuccessResponse(code int, message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse represents the standard structure for API error responses.
type ErrorResponse struct {
	Code    int    `json:"code"`            // HTTP status code (e.g. 400, 404)
	Message string `json:"message"`         // Human-readable error message
	Error   string `json:"error,omitempty"` // Optional technical error details
}

// NewErrorResponse constructs a new error response with optional internal error details.
func NewErrorResponse(statusCode int, message string, err error) *ErrorResponse {
	response := &ErrorResponse{
		Code:    statusCode,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	return response
}

// Predefined error helpers (can be reused across handlers)
func ErrNotFound() *ErrorResponse {
	return NewErrorResponse(http.StatusNotFound, "Resource not found", nil)
}

func ErrUnauthorized() *ErrorResponse {
	return NewErrorResponse(http.StatusUnauthorized, "Unauthorized", nil)
}

func ErrForbidden() *ErrorResponse {
	return NewErrorResponse(http.StatusForbidden, "Forbidden", nil)
}

func ErrBadRequest() *ErrorResponse {
	return NewErrorResponse(http.StatusBadRequest, "Bad request", nil)
}

func ErrInternalServer() *ErrorResponse {
	return NewErrorResponse(http.StatusInternalServerError, "Internal server error", nil)
}
