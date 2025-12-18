package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/domain"
)

type FileUseCase struct {
	repo    domain.FileRepository
	storage domain.StorageProvider
}

func NewFileUseCase(repo domain.FileRepository, storage domain.StorageProvider) *FileUseCase {
	return &FileUseCase{
		repo:    repo,
		storage: storage,
	}
}

func (u *FileUseCase) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, tenantID, uploadedBy uuid.UUID, bucket string) (*domain.File, error) {
	// Generate unique path
	// Structure: bucket/tenant_id/year/month/uuid.ext
	ext := filepath.Ext(header.Filename)
	fileID := uuid.New()
	now := time.Now()
	path := fmt.Sprintf("%s/%s/%d/%02d/%s%s", bucket, tenantID.String(), now.Year(), now.Month(), fileID.String(), ext)

	// Upload to storage
	storedPath, err := u.storage.Upload(ctx, file, header, path)
	if err != nil {
		return nil, err
	}

	// Create file entity
	f := &domain.File{
		ID:           fileID,
		TenantID:     tenantID,
		Name:         header.Filename, // Or generated name? Let's keep original name as display name for now, but usually we might want to sanitize it.
		OriginalName: header.Filename,
		MimeType:     header.Header.Get("Content-Type"),
		Size:         header.Size,
		Path:         storedPath,
		Bucket:       bucket,
		UploadedBy:   uploadedBy,
	}

	// Save metadata
	if err := u.repo.Create(ctx, f); err != nil {
		// Rollback storage if db fails
		_ = u.storage.Delete(ctx, storedPath)
		return nil, err
	}

	// Generate URL
	url, err := u.storage.GetURL(ctx, storedPath)
	if err != nil {
		return nil, err
	}
	f.URL = url

	return f, nil
}

func (u *FileUseCase) Get(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	f, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, nil
	}

	url, err := u.storage.GetURL(ctx, f.Path)
	if err != nil {
		return nil, err
	}
	f.URL = url

	return f, nil
}

func (u *FileUseCase) Download(ctx context.Context, id uuid.UUID) (*domain.File, io.ReadCloser, error) {
	f, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if f == nil {
		return nil, nil, nil
	}

	reader, err := u.storage.Get(ctx, f.Path)
	if err != nil {
		return nil, nil, err
	}

	return f, reader, nil
}

func (u *FileUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	f, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if f == nil {
		return nil
	}

	// Delete from storage
	if err := u.storage.Delete(ctx, f.Path); err != nil {
		// Log error but continue to delete metadata? Or fail?
		// For now, let's return error
		return err
	}

	// Delete metadata
	return u.repo.Delete(ctx, id)
}

func (u *FileUseCase) List(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]*domain.File, int64, error) {
	offset := (page - 1) * limit
	files, err := u.repo.List(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for _, f := range files {
		url, err := u.storage.GetURL(ctx, f.Path)
		if err == nil {
			f.URL = url
		}
	}

	// TODO: Get total count for pagination
	return files, 0, nil
}
