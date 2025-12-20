package usecase

import (
	"context"
	"io"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFileRepository ...
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(ctx context.Context, file *domain.File) error {
	args := m.Called(ctx, file)
	return args.Error(0)
}

func (m *MockFileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.File), args.Error(1)
}

func (m *MockFileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFileRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.File, error) {
	args := m.Called(ctx, tenantID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.File), args.Error(1)
}

// MockStorageProvider ...
type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, path string) (string, error) {
	args := m.Called(ctx, file, header, path)
	return args.String(0), args.Error(1)
}

func (m *MockStorageProvider) Delete(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *MockStorageProvider) GetURL(ctx context.Context, path string) (string, error) {
	args := m.Called(ctx, path)
	return args.String(0), args.Error(1)
}

func (m *MockStorageProvider) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestFileUseCase_Download(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	path := "test/file.txt"

	t.Run("success", func(t *testing.T) {
		file := &domain.File{ID: id, Path: path}
		content := io.NopCloser(strings.NewReader("content"))
		
		// Create new mock instances for this run to avoid call count issues if sharing
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		repo.On("GetByID", ctx, id).Return(file, nil)
		storage.On("Get", ctx, path).Return(content, nil)

		f, r, err := uc.Download(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, file, f)
		assert.NotNil(t, r)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		repo.On("GetByID", ctx, id).Return(nil, nil)
		f, r, err := uc.Download(ctx, id)
		assert.NoError(t, err)
		assert.Nil(t, f)
		assert.Nil(t, r)
	})
}

func TestFileUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	path := "test/file.txt"

	t.Run("success", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		file := &domain.File{ID: id, Path: path}
		repo.On("GetByID", ctx, id).Return(file, nil)
		storage.On("Delete", ctx, path).Return(nil)
		repo.On("Delete", ctx, id).Return(nil)

		err := uc.Delete(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		repo.On("GetByID", ctx, id).Return(nil, nil)

		err := uc.Delete(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("storage error", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		file := &domain.File{ID: id, Path: path}
		repo.On("GetByID", ctx, id).Return(file, nil)
		storage.On("Delete", ctx, path).Return(assert.AnError)

		err := uc.Delete(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}

func TestFileUseCase_Upload(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		header := &multipart.FileHeader{Filename: "test.txt", Size: 1024}
		file := &domain.File{
			TenantID: tenantID,
			UploadedBy: userID,
			Name: "test.txt",
			MimeType: "text/plain",
			Size: 1024,
		}

		storage.On("Upload", ctx, nil, header, mock.Anything).Return("path/to/file", nil)
		storage.On("GetURL", ctx, "path/to/file").Return("http://example.com/file", nil)
		repo.On("Create", ctx, mock.MatchedBy(func(f *domain.File) bool {
			return f.Name == file.Name && f.Size == file.Size
		})).Return(nil)

		res, err := uc.Upload(ctx, nil, header, tenantID, userID, "uploads")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "test.txt", res.Name)
	})
}

func TestFileUseCase_List(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()

	t.Run("success", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		files := []*domain.File{
			{ID: uuid.New(), Name: "file1.txt"},
			{ID: uuid.New(), Name: "file2.txt"},
		}

		repo.On("List", ctx, tenantID, 10, 0).Return(files, nil)
		storage.On("GetURL", ctx, mock.Anything).Return("http://example.com/file", nil)

		res, _, err := uc.List(ctx, tenantID, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.NotEmpty(t, res[0].URL)
	})
}

func TestFileUseCase_Get(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	path := "bucket/tenant/2025/12/file.ext"

	t.Run("success", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		file := &domain.File{ID: id, Path: path}
		repo.On("GetByID", ctx, id).Return(file, nil)
		storage.On("GetURL", ctx, path).Return("http://example.com/file.ext", nil)

		res, err := uc.Get(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "http://example.com/file.ext", res.URL)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		repo.On("GetByID", ctx, id).Return(nil, nil)

		res, err := uc.Get(ctx, id)
		assert.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("url error", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		file := &domain.File{ID: id, Path: path}
		repo.On("GetByID", ctx, id).Return(file, nil)
		storage.On("GetURL", ctx, path).Return("", assert.AnError)

		res, err := uc.Get(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestFileUseCase_Upload_ErrorPaths(t *testing.T) {
	ctx := context.Background()
	tenantID := uuid.New()
	userID := uuid.New()
	header := &multipart.FileHeader{Filename: "image.png", Size: 2048}

	t.Run("repo create error triggers rollback", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		storage.On("Upload", ctx, nil, header, mock.Anything).Return("stored/path", nil)
		repo.On("Create", ctx, mock.AnythingOfType("*domain.File")).Return(assert.AnError)
		storage.On("Delete", ctx, "stored/path").Return(nil)

		res, err := uc.Upload(ctx, nil, header, tenantID, userID, "uploads")
		assert.Error(t, err)
		assert.Nil(t, res)
		storage.AssertCalled(t, "Delete", ctx, "stored/path")
	})

	t.Run("get URL error after create", func(t *testing.T) {
		repo := new(MockFileRepository)
		storage := new(MockStorageProvider)
		uc := NewFileUseCase(repo, storage)

		storage.On("Upload", ctx, nil, header, mock.Anything).Return("stored/path", nil)
		repo.On("Create", ctx, mock.AnythingOfType("*domain.File")).Return(nil)
		storage.On("GetURL", ctx, "stored/path").Return("", assert.AnError)

		res, err := uc.Upload(ctx, nil, header, tenantID, userID, "uploads")
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
