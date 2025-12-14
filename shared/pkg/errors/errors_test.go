package errors

import (
	"fmt"
	"net/http"
	"testing"
)

func TestStatusFromCode(t *testing.T) {
	tests := []struct {
		code   string
		status int
	}{
		{"VALIDATION_ERROR", http.StatusBadRequest},
		{"2001", http.StatusUnauthorized},
		{"FORBIDDEN", http.StatusForbidden},
		{"NOT_FOUND", http.StatusNotFound},
		{"DUPLICATE_ENTRY", http.StatusConflict},
		{"THIRD_PARTY_ERROR", http.StatusBadGateway},
		{"TIMEOUT", http.StatusGatewayTimeout},
		{"INTERNAL_SERVER_ERROR", http.StatusInternalServerError},
		{"UNKNOWN", http.StatusInternalServerError},
	}
	for _, tt := range tests {
		if got := statusFromCode(tt.code); got != tt.status {
			t.Fatalf("statusFromCode(%s)=%d want %d", tt.code, got, tt.status)
		}
	}
}

func TestToHTTP(t *testing.T) {
	e := New("VALIDATION_ERROR", "Invalid data")
	e = WithDetails(e, []FieldError{{Field: "email", Message: "invalid"}})
	status, body := ToHTTP(e)
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d want %d", status, http.StatusBadRequest)
	}
	errMap := body["error"].(map[string]any)
	if errMap["code"] != "VALIDATION_ERROR" || errMap["message"] != "Invalid data" {
		t.Fatalf("unexpected body error content")
	}
}

func TestWrapAndDetails(t *testing.T) {
	inner := fmt.Errorf("db error")
	e := Wrap("5001", "duplicate", inner)
	e = WithDetails(e, []FieldError{{Field: "id", Message: "exists"}})
	if e.Status != http.StatusConflict {
		t.Fatalf("status mismatch")
	}
	_, body := ToHTTP(e)
	dets := body["error"].(map[string]any)["details"].([]FieldError)
	if len(dets) != 1 || dets[0].Field != "id" {
		t.Fatalf("details missing")
	}
}

func TestToHTTPUnknownCode(t *testing.T) {
	e := New("SOME_UNKNOWN", "unknown")
	status, body := ToHTTP(e)
	if status != http.StatusInternalServerError {
		t.Fatalf("expected 500 for unknown code")
	}
	errMap := body["error"].(map[string]any)
	if errMap["code"] != "SOME_UNKNOWN" {
		t.Fatalf("code passthrough missing")
	}
}

func TestErrorStringIncludesInner(t *testing.T) {
	inner := fmt.Errorf("inner")
	e := Wrap("1001", "internal", inner)
	s := e.Error()
	if s == "" || (s != "" && s == "internal") {
		t.Fatalf("error string should include inner context")
	}
}
