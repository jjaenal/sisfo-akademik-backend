package usecase

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/usecase"
)

type applicationDocumentUseCase struct {
	repo            repository.ApplicationDocumentRepository
	applicationRepo repository.ApplicationRepository
	uploadDir       string
}

func NewApplicationDocumentUseCase(
	repo repository.ApplicationDocumentRepository,
	appRepo repository.ApplicationRepository,
	uploadDir string,
) domainUseCase.ApplicationDocumentUseCase {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	// Ensure directory exists
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		// Panic on startup if cannot create directory
		panic(err)
	}
	
	return &applicationDocumentUseCase{
		repo:            repo,
		applicationRepo: appRepo,
		uploadDir:       uploadDir,
	}
}

func (u *applicationDocumentUseCase) Upload(ctx context.Context, applicationID uuid.UUID, documentType string, file *multipart.FileHeader) (*entity.ApplicationDocument, error) {
	// 1. Check Application exists
	app, err := u.applicationRepo.GetByID(ctx, applicationID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, errors.New("application not found")
	}

	// 2. Validate File
	// Max size 5MB
	if file.Size > 5*1024*1024 {
		return nil, errors.New("file too large (max 5MB)")
	}
	
	// 3. Save File
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	dst := filepath.Join(u.uploadDir, filename)
	
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	out, err := os.Create(dst) // #nosec G304 -- dst is constructed from trusted config and safe UUID
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return nil, err
	}

	// 4. Create Record
	doc := &entity.ApplicationDocument{
		ID:            uuid.New(),
		ApplicationID: applicationID,
		DocumentType:  documentType,
		FileURL:       "/uploads/" + filename, // In real app, this would be S3 URL
		FileName:      file.Filename,
		FileSize:      file.Size,
		UploadedAt:    time.Now(),
	}

	if err := u.repo.Create(ctx, doc); err != nil {
		// Cleanup file if DB fail
		_ = os.Remove(dst)
		return nil, err
	}

	return doc, nil
}

func (u *applicationDocumentUseCase) GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*entity.ApplicationDocument, error) {
	return u.repo.GetByApplicationID(ctx, applicationID)
}

func (u *applicationDocumentUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	doc, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if doc == nil {
		return errors.New("document not found")
	}

	// Delete from DB
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Delete file (best effort)
	// Extract filename from URL or store path in DB
	// For now, assuming URL format "/uploads/filename"
	if len(doc.FileURL) > 9 {
		filename := doc.FileURL[9:]
		filepath := filepath.Join(u.uploadDir, filename)
		_ = os.Remove(filepath)
	}

	return nil
}
