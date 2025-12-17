package entity

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationDocument struct {
	ID            uuid.UUID `json:"id"`
	ApplicationID uuid.UUID `json:"application_id"`
	DocumentType  string    `json:"document_type"`
	FileURL       string    `json:"file_url"`
	FileName      string    `json:"file_name"`
	FileSize      int64     `json:"file_size"`
	UploadedAt    time.Time `json:"uploaded_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ApplicationDocument) TableName() string {
	return "application_documents"
}
