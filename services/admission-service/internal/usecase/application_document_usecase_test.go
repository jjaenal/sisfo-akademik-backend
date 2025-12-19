package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func createMultipartFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile("file", filename)
	assert.NoError(t, err)
	_, err = part.Write(content)
	assert.NoError(t, err)
	err = w.Close()
	assert.NoError(t, err)

	r := multipart.NewReader(&b, w.Boundary())
	form, err := r.ReadForm(int64(len(content)) + 1024)
	assert.NoError(t, err)
	
	files := form.File["file"]
	assert.NotEmpty(t, files)
	return files[0]
}

func TestApplicationDocumentUseCase_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDocRepo := mocks.NewMockApplicationDocumentRepository(ctrl)
	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)

	// Create temp dir for uploads
	tmpDir, err := os.MkdirTemp("", "uploads")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	u := usecase.NewApplicationDocumentUseCase(mockDocRepo, mockAppRepo, tmpDir)

	ctx := context.Background()
	appID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Mock Application Exists
		mockAppRepo.EXPECT().GetByID(ctx, appID).Return(&entity.Application{ID: appID}, nil)

		// Create file header
		fileContent := []byte("test content")
		fh := createMultipartFileHeader(t, "test.pdf", fileContent)

		// Mock Repo Create
		mockDocRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, doc *entity.ApplicationDocument) error {
			assert.Equal(t, appID, doc.ApplicationID)
			assert.Equal(t, "transcript", doc.DocumentType)
			assert.Equal(t, "test.pdf", doc.FileName)
			assert.NotEmpty(t, doc.FileURL)
			return nil
		})

		doc, err := u.Upload(ctx, appID, "transcript", fh)
		assert.NoError(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, "test.pdf", doc.FileName)
	})

	t.Run("ApplicationNotFound", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(ctx, appID).Return(nil, nil)

		fh := createMultipartFileHeader(t, "test.pdf", []byte("content"))
		doc, err := u.Upload(ctx, appID, "transcript", fh)
		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Equal(t, "application not found", err.Error())
	})

	t.Run("FileTooLarge", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(ctx, appID).Return(&entity.Application{ID: appID}, nil)

		// Fake size
		fh := createMultipartFileHeader(t, "large.pdf", []byte("content"))
		fh.Size = 6 * 1024 * 1024 // 6MB

		doc, err := u.Upload(ctx, appID, "transcript", fh)
		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Contains(t, err.Error(), "file too large")
	})

	t.Run("RepoError", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(ctx, appID).Return(&entity.Application{ID: appID}, nil)
		fh := createMultipartFileHeader(t, "test.pdf", []byte("content"))

		mockDocRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("db error"))

		doc, err := u.Upload(ctx, appID, "transcript", fh)
		assert.Error(t, err)
		assert.Nil(t, doc)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestApplicationDocumentUseCase_GetByApplicationID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDocRepo := mocks.NewMockApplicationDocumentRepository(ctrl)
	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	tmpDir := os.TempDir()

	u := usecase.NewApplicationDocumentUseCase(mockDocRepo, mockAppRepo, tmpDir)
	ctx := context.Background()
	appID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		expectedDocs := []*entity.ApplicationDocument{
			{ID: uuid.New(), ApplicationID: appID},
		}
		mockDocRepo.EXPECT().GetByApplicationID(ctx, appID).Return(expectedDocs, nil)

		docs, err := u.GetByApplicationID(ctx, appID)
		assert.NoError(t, err)
		assert.Equal(t, expectedDocs, docs)
	})

	t.Run("Error", func(t *testing.T) {
		mockDocRepo.EXPECT().GetByApplicationID(ctx, appID).Return(nil, errors.New("db error"))

		docs, err := u.GetByApplicationID(ctx, appID)
		assert.Error(t, err)
		assert.Nil(t, docs)
	})
}

func TestApplicationDocumentUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDocRepo := mocks.NewMockApplicationDocumentRepository(ctrl)
	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "uploads")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	u := usecase.NewApplicationDocumentUseCase(mockDocRepo, mockAppRepo, tmpDir)
	ctx := context.Background()
	docID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Create a dummy file to be deleted
		filename := "test.pdf"
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte("content"), 0644)
		assert.NoError(t, err)

		doc := &entity.ApplicationDocument{
			ID:      docID,
			FileURL: "/uploads/" + filename,
		}
		
		mockDocRepo.EXPECT().GetByID(ctx, docID).Return(doc, nil)
		mockDocRepo.EXPECT().Delete(ctx, docID).Return(nil)

		err = u.Delete(ctx, docID)
		assert.NoError(t, err)
		
		// Verify file is deleted
		_, err = os.Stat(filePath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("NotFound", func(t *testing.T) {
		mockDocRepo.EXPECT().GetByID(ctx, docID).Return(nil, nil)

		err := u.Delete(ctx, docID)
		assert.Error(t, err)
		assert.Equal(t, "document not found", err.Error())
	})

	t.Run("DBError", func(t *testing.T) {
		doc := &entity.ApplicationDocument{ID: docID}
		mockDocRepo.EXPECT().GetByID(ctx, docID).Return(doc, nil)
		mockDocRepo.EXPECT().Delete(ctx, docID).Return(errors.New("db error"))

		err := u.Delete(ctx, docID)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}
