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

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func Success(w http.ResponseWriter, data any) {
	m := Meta{Timestamp: time.Now().UTC(), RequestID: uuid.NewString()}
	write(w, http.StatusOK, map[string]any{"success": true, "data": data, "meta": m})
}

func Error(w http.ResponseWriter, status int, code string, message string, details any) {
	m := Meta{Timestamp: time.Now().UTC(), RequestID: uuid.NewString()}
	write(w, status, map[string]any{
		"success": false,
		"error": map[string]any{
			"code":    code,
			"message": message,
			"details": details,
		},
		"meta": m,
	})
}
