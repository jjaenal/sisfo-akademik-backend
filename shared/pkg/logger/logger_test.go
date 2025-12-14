package logger

import "testing"

func TestNew(t *testing.T) {
	l, err := New("development")
	if err != nil || l == nil {
		t.Fatalf("logger not created")
	}
}
