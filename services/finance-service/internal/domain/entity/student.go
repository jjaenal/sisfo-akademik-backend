package entity

import (
	"time"

	"github.com/google/uuid"
)

type StudentStatus string

const (
	StudentStatusActive   StudentStatus = "active"
	StudentStatusInactive StudentStatus = "inactive"
)

type Student struct {
	ID        uuid.UUID     `json:"id"`
	TenantID  uuid.UUID     `json:"tenant_id"`
	Name      string        `json:"name"`
	Status    StudentStatus `json:"status"`
	ClassID   *uuid.UUID    `json:"class_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (Student) TableName() string {
	return "students"
}
