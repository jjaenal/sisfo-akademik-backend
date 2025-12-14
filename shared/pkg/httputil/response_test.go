package httputil

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestSuccess(t *testing.T) {
	rr := httptest.NewRecorder()
	Success(rr, map[string]string{"a": "b"})
	if rr.Code != 200 {
		t.Fatalf("status=%d want 200", rr.Code)
	}
	var m map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &m); err != nil {
		t.Fatalf("json error: %v", err)
	}
	if m["success"] != true {
		t.Fatalf("success not true")
	}
}

func TestError(t *testing.T) {
	rr := httptest.NewRecorder()
	Error(rr, 400, "VALIDATION_ERROR", "invalid", nil)
	if rr.Code != 400 {
		t.Fatalf("status=%d want 400", rr.Code)
	}
}
