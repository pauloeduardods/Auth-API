package utils

import "fmt"

type ApiError struct {
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
	StatusCode  int    `json:"-"`
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("message: %s, description: %s, status_code: %d", e.Message, e.Description, e.StatusCode)
}

func NewApiError(statusCode int, message, description string) *ApiError {
	return &ApiError{
		Message:     message,
		StatusCode:  statusCode,
		Description: description,
	}
}
