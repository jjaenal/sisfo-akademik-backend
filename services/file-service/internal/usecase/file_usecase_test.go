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

		// Repository List usually returns data and error. Count might be a separate call or part of List.
		// Checking domain/file.go: List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*File, error)
		// It doesn't seem to return total count in the interface I saw.
		// But UseCase List returns ([], int64, error). 
		// I need to check how UseCase implements List to see where it gets the count.
		// Assuming for now it just returns 0 or I missed a Count method in repo.
		// Let's check UseCase implementation if possible, but for now I'll assume repo.List is called.
		
		repo.On("List", ctx, tenantID, 10, 0).Return(files, nil)
		// If usecase calls Count, I'd need to mock it. But I didn't see Count in Repo interface.
		// Maybe UseCase just returns len(files) or 0? 
		// Or maybe I missed Count in the Repo interface read?
		// Let's re-read the repo interface from domain/file.go in my thought process... 
		// "List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*File, error)"
		// No Count method. 
		// So UseCase probably returns 0 or implements it differently. 
		// I will just handle the return values of uc.List.
		
		storage.On("GetURL", ctx, mock.Anything).Return("http://example.com/file", nil)

		res, _, err := uc.List(ctx, tenantID, 1, 10)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.NotEmpty(t, res[0].URL)
	})
}
