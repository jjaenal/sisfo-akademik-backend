package validator

import "testing"

type sample struct {
	Email string `validate:"required,email"`
}

func TestValidate(t *testing.T) {
	s := sample{Email: "invalid"}
	if err := Validate(s); err == nil {
		t.Fatalf("expected error")
	}
	s.Email = "a@b.com"
	if err := Validate(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
