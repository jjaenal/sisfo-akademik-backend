package httputil

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Meta struct {
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Error   ErrorDetail `json:"error"`
	Meta    Meta        `json:"meta"`
}

type SuccessResponse struct {
	Success bool `json:"success" example:"true"`
	Data    any  `json:"data"`
	Meta    Meta `json:"meta"`
}

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func Success(w http.ResponseWriter, data any) {
	m := Meta{Timestamp: time.Now().UTC(), RequestID: uuid.NewString()}
	resp := SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    m,
	}
	write(w, http.StatusOK, resp)
}

func Error(w http.ResponseWriter, status int, code string, message string, details any) {
	m := Meta{Timestamp: time.Now().UTC(), RequestID: uuid.NewString()}
	resp := ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: m,
	}
	write(w, status, resp)
}
