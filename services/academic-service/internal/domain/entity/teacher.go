package entity

import (
	"time"

	"github.com/google/uuid"
)

type Teacher struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	TenantID   string     `json:"tenant_id" db:"tenant_id"`
	UserID     *uuid.UUID `json:"user_id" db:"user_id"`
	NIP        string     `json:"nip" db:"nip"`
	Name       string     `json:"name" db:"name"`
	Gender     string     `json:"gender" db:"gender"`
	TitleFront string     `json:"title_front" db:"title_front"`
	TitleBack  string     `json:"title_back" db:"title_back"`
	Phone      string     `json:"phone" db:"phone"`
	Email      string     `json:"email" db:"email"`
	Status     string     `json:"status" db:"status"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy  *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy  *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (t *Teacher) Validate() map[string]string {
	errors := make(map[string]string)
	if t.Name == "" {
		errors["name"] = "Name is required"
	}
	if t.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	return errors
}
