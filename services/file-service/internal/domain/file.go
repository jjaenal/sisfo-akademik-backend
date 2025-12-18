package domain

import (
	"context"
	"io"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TenantID    uuid.UUID `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	OriginalName string   `json:"original_name" db:"original_name"`
	MimeType    string    `json:"mime_type" db:"mime_type"`
	Size        int64     `json:"size" db:"size"`
	Path        string    `json:"path" db:"path"`
	Bucket      string    `json:"bucket" db:"bucket"` // e.g. "uploads", "reports"
	URL         string    `json:"url" db:"-"`         // Generated on read
	UploadedBy  uuid.UUID  `json:"uploaded_by" db:"uploaded_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type FileRepository interface {
	Create(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id uuid.UUID) (*File, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*File, error)
}

type StorageProvider interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, path string) (string, error)
	Delete(ctx context.Context, path string) error
	GetURL(ctx context.Context, path string) (string, error)
	Get(ctx context.Context, path string) (io.ReadCloser, error)
}

type FileUseCase interface {
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, tenantID, uploadedBy uuid.UUID, bucket string) (*File, error)
	Get(ctx context.Context, id uuid.UUID) (*File, error)
	Download(ctx context.Context, id uuid.UUID) (*File, io.ReadCloser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]*File, int64, error)
}
