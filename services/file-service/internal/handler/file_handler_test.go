package handler

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/domain"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/middleware"
)

type mockUC struct {
	uploadCalled   bool
	getCalled      bool
	listCalled     bool
	deleteCalled   bool
	downloadCalled bool

	fileResp   *domain.File
	filesResp  []*domain.File
	totalCount int64
	reader     io.ReadCloser

	uploadErr   error
	getErr      error
	downloadErr error
	deleteErr   error
	listErr     error
}

func (m *mockUC) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, tenantID, uploadedBy uuid.UUID, bucket string) (*domain.File, error) {
	m.uploadCalled = true
	return m.fileResp, m.uploadErr
}
func (m *mockUC) Get(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	m.getCalled = true
	return m.fileResp, m.getErr
}
func (m *mockUC) Download(ctx context.Context, id uuid.UUID) (*domain.File, io.ReadCloser, error) {
	m.downloadCalled = true
	return m.fileResp, m.reader, m.downloadErr
}
func (m *mockUC) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled = true
	return m.deleteErr
}
func (m *mockUC) List(ctx context.Context, tenantID uuid.UUID, page, limit int) ([]*domain.File, int64, error) {
	m.listCalled = true
	return m.filesResp, m.totalCount, m.listErr
}

func setupRouter(h *FileHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h.RegisterRoutes(r)
	return r
}

func withClaims(req *http.Request, claims jwtutil.Claims) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.ClaimsKey, claims)
	return req.WithContext(ctx)
}

func TestUpload_Success(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	mock := &mockUC{
		fileResp: &domain.File{
			ID:           uuid.New(),
			TenantID:     tenantID,
			OriginalName: "test.txt",
			MimeType:     "text/plain",
			Size:         4,
			Path:         "uploads/test.txt",
			Bucket:       "uploads",
			UploadedBy:   userID,
		},
	}
	h := NewFileHandler(mock)
	r := setupRouter(h)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = io.Copy(part, strings.NewReader("data"))
	_ = writer.WriteField("bucket", "uploads")
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withClaims(req, jwtutil.Claims{TenantID: tenantID.String(), UserID: userID})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d want 200", rr.Code)
	}
	if !mock.uploadCalled {
		t.Fatalf("usecase Upload not called")
	}
}

func TestUpload_MissingClaims(t *testing.T) {
	mock := &mockUC{fileResp: &domain.File{ID: uuid.New()}}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = io.Copy(part, strings.NewReader("data"))
	_ = writer.Close()
	req := httptest.NewRequest("POST", "/api/v1/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d want 401", rr.Code)
	}
}

func TestUpload_InvalidTenantID(t *testing.T) {
	mock := &mockUC{fileResp: &domain.File{ID: uuid.New()}}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = io.Copy(part, strings.NewReader("data"))
	_ = writer.Close()
	req := httptest.NewRequest("POST", "/api/v1/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withClaims(req, jwtutil.Claims{TenantID: "not-a-uuid", UserID: uuid.New()})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d want 401", rr.Code)
	}
}

func TestUpload_MissingFile(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	mock := &mockUC{}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	req := httptest.NewRequest("POST", "/api/v1/files/upload", nil)
	req = withClaims(req, jwtutil.Claims{TenantID: tenantID.String(), UserID: userID})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("code=%d want 400", rr.Code)
	}
}

func TestUpload_InternalError(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	mock := &mockUC{uploadErr: io.EOF}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, _ = io.Copy(part, strings.NewReader("data"))
	_ = writer.Close()
	req := httptest.NewRequest("POST", "/api/v1/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = withClaims(req, jwtutil.Claims{TenantID: tenantID.String(), UserID: userID})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("code=%d want 500", rr.Code)
	}
}

func TestGet_Success(t *testing.T) {
	mock := &mockUC{
		fileResp: &domain.File{
			ID:           uuid.New(),
			OriginalName: "a.pdf",
			MimeType:     "application/pdf",
			Size:         10,
		},
	}
	h := NewFileHandler(mock)
	r := setupRouter(h)

	id := uuid.New()
	req := httptest.NewRequest("GET", "/api/v1/files/"+id.String(), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d want 200", rr.Code)
	}
	if !mock.getCalled {
		t.Fatalf("usecase Get not called")
	}
}

func TestGet_InvalidID(t *testing.T) {
	mock := &mockUC{}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/files/invalid-id", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("code=%d want 400", rr.Code)
	}
}

func TestGet_InternalError(t *testing.T) {
	mock := &mockUC{getErr: io.EOF}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	id := uuid.New()
	req := httptest.NewRequest("GET", "/api/v1/files/"+id.String(), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("code=%d want 500", rr.Code)
	}
}

func TestList_Success(t *testing.T) {
	tenantID := uuid.New()
	mock := &mockUC{
		filesResp: []*domain.File{
			{ID: uuid.New(), OriginalName: "x.jpg"},
			{ID: uuid.New(), OriginalName: "y.png"},
		},
		totalCount: 2,
	}
	h := NewFileHandler(mock)
	r := setupRouter(h)

	req := httptest.NewRequest("GET", "/api/v1/files?page=1&limit=10", nil)
	req = withClaims(req, jwtutil.Claims{TenantID: tenantID.String(), UserID: uuid.New()})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d want 200", rr.Code)
	}
	if !mock.listCalled {
		t.Fatalf("usecase List not called")
	}
}

func TestList_MissingClaims(t *testing.T) {
	mock := &mockUC{}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/files?page=1&limit=10", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d want 401", rr.Code)
	}
}

func TestList_InvalidTenantID(t *testing.T) {
	mock := &mockUC{}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/files?page=1&limit=10", nil)
	req = withClaims(req, jwtutil.Claims{TenantID: "bad-uuid", UserID: uuid.New()})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d want 401", rr.Code)
	}
}

func TestList_InternalError(t *testing.T) {
	tenantID := uuid.New()
	mock := &mockUC{listErr: io.EOF}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/files?page=1&limit=10", nil)
	req = withClaims(req, jwtutil.Claims{TenantID: tenantID.String(), UserID: uuid.New()})
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("code=%d want 500", rr.Code)
	}
}

func TestDelete_Success(t *testing.T) {
	mock := &mockUC{}
	h := NewFileHandler(mock)
	r := setupRouter(h)

	id := uuid.New()
	req := httptest.NewRequest("DELETE", "/api/v1/files/"+id.String(), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d want 200", rr.Code)
	}
	if !mock.deleteCalled {
		t.Fatalf("usecase Delete not called")
	}
}

func TestDelete_InvalidID(t *testing.T) {
	mock := &mockUC{}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	req := httptest.NewRequest("DELETE", "/api/v1/files/invalid-id", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("code=%d want 400", rr.Code)
	}
}

func TestDelete_InternalError(t *testing.T) {
	mock := &mockUC{deleteErr: io.EOF}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	id := uuid.New()
	req := httptest.NewRequest("DELETE", "/api/v1/files/"+id.String(), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("code=%d want 500", rr.Code)
	}
}

func TestDownload_Success(t *testing.T) {
	content := bytes.NewBufferString("data")
	mock := &mockUC{
		fileResp: &domain.File{
			ID:           uuid.New(),
			OriginalName: "dl.txt",
			MimeType:     "text/plain",
			Size:         int64(content.Len()),
		},
		reader: io.NopCloser(content),
	}
	h := NewFileHandler(mock)
	r := setupRouter(h)

	id := uuid.New()
	req := httptest.NewRequest("GET", "/api/v1/files/"+id.String()+"/download", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("code=%d want 200", rr.Code)
	}
	if !mock.downloadCalled {
		t.Fatalf("usecase Download not called")
	}
	if rr.Header().Get("Content-Disposition") == "" {
		t.Fatalf("content disposition should be set")
	}
}

func TestDownload_NotFound(t *testing.T) {
	mock := &mockUC{
		fileResp: nil,
		reader:   nil,
	}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/files/"+uuid.New().String()+"/download", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("code=%d want 404", rr.Code)
	}
}

func TestDownload_InvalidID(t *testing.T) {
	mock := &mockUC{}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	req := httptest.NewRequest("GET", "/api/v1/files/invalid-id/download", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("code=%d want 400", rr.Code)
	}
}

func TestDownload_InternalError(t *testing.T) {
	mock := &mockUC{downloadErr: io.EOF, fileResp: &domain.File{ID: uuid.New(), OriginalName: "x", MimeType: "text/plain", Size: 1}}
	h := NewFileHandler(mock)
	r := setupRouter(h)
	id := uuid.New()
	req := httptest.NewRequest("GET", "/api/v1/files/"+id.String()+"/download", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("code=%d want 500", rr.Code)
	}
}
