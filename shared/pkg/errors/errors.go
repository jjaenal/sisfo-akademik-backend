package errors

import (
	"fmt"
	"net/http"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []FieldError  `json:"details,omitempty"`
	Inner   error         `json:"-"`
	Status  int           `json:"-"`
}

func (e *AppError) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Inner)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code string, message string) *AppError {
	return &AppError{Code: code, Message: message, Status: statusFromCode(code)}
}

func Wrap(code string, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Inner: err, Status: statusFromCode(code)}
}

func WithDetails(e *AppError, details []FieldError) *AppError {
	e.Details = details
	return e
}

func statusFromCode(code string) int {
	switch code {
	case "VALIDATION_ERROR", "4001", "4002":
		return http.StatusBadRequest
	case "UNAUTHORIZED", "2001", "2002":
		return http.StatusUnauthorized
	case "FORBIDDEN", "3001", "3002":
		return http.StatusForbidden
	case "NOT_FOUND", "5002":
		return http.StatusNotFound
	case "DUPLICATE_ENTRY", "5001":
		return http.StatusConflict
	case "THIRD_PARTY_ERROR", "6001":
		return http.StatusBadGateway
	case "TIMEOUT", "6002":
		return http.StatusGatewayTimeout
	case "INTERNAL_SERVER_ERROR", "1001":
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func ToHTTP(e *AppError) (int, map[string]any) {
	status := e.Status
	if status == 0 {
		status = statusFromCode(e.Code)
	}
	body := map[string]any{
		"success": false,
		"error": map[string]any{
			"code":    e.Code,
			"message": e.Message,
			"details": e.Details,
		},
	}
	return status, body
}
