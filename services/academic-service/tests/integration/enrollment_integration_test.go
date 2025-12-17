package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestEnrollmentIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockEnrollmentRepository(ctrl)
	mockClassRepo := mocks.NewMockClassRepository(ctrl)

	// Real UseCase with Mock Repo
	u := usecase.NewEnrollmentUseCase(mockRepo, mockClassRepo, time.Second*2)

	// Real Handler with Real UseCase
	h := handler.NewEnrollmentHandler(u)

	// Setup Router
	r := gin.New()
	enrollments := r.Group("/api/v1/enrollments")
	{
		enrollments.POST("", h.Enroll)
		enrollments.PATCH("/:id/status", h.UpdateStatus)
		enrollments.GET("/:id", h.GetByID)
	}

	t.Run("Enroll Student Success", func(t *testing.T) {
		classID := uuid.New()
		reqBody := map[string]interface{}{
			"tenant_id":  "tenant-1",
			"class_id":   classID.String(),
			"student_id": uuid.New().String(),
			"status":     "enrolled",
		}
		body, _ := json.Marshal(reqBody)

		validClass := &entity.Class{
			ID:       classID,
			TenantID: "tenant-1",
			Name:     "Class 1A",
			Capacity: 30,
		}

		// Expect Class Check and Capacity Check
		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(validClass, nil)
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return([]entity.Enrollment{}, nil)

		// Expect Enroll
		mockRepo.EXPECT().Enroll(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, e *entity.Enrollment) error {
			assert.Equal(t, "enrolled", e.Status)
			assert.Equal(t, "tenant-1", e.TenantID)
			e.ID = uuid.New()
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/enrollments", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool              `json:"success"`
			Data    entity.Enrollment `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "enrolled", resp.Data.Status)
		assert.NotEmpty(t, resp.Data.ID)
	})

	t.Run("Update Enrollment Status Success", func(t *testing.T) {
		id := uuid.New()
		status := "completed"
		reqBody := map[string]string{"status": status}
		body, _ := json.Marshal(reqBody)

		mockRepo.EXPECT().UpdateStatus(gomock.Any(), id, status).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/enrollments/"+id.String()+"/status", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Success bool `json:"success"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
	})

	t.Run("Bulk Enroll Success", func(t *testing.T) {
		classID := uuid.New()
		
		// Setup Router for this specific route (since it's under /classes)
		classes := r.Group("/api/v1/classes")
		classes.POST("/:id/students/bulk", h.BulkEnroll)

		// Create CSV content
		studentID1 := uuid.New()
		studentID2 := uuid.New()
		csvContent := studentID1.String() + "\n" + studentID2.String() + "\n"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "students.csv")
		assert.NoError(t, err)
		_, err = part.Write([]byte(csvContent))
		assert.NoError(t, err)
		err = writer.Close()
		assert.NoError(t, err)

		validClass := &entity.Class{
			ID:       classID,
			TenantID: "tenant-1",
			Name:     "Class 1A",
			Capacity: 30,
		}

		// Mock Expectations
		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(validClass, nil)
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return([]entity.Enrollment{}, nil)
		mockRepo.EXPECT().BulkEnroll(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, enrollments []*entity.Enrollment) error {
			assert.Len(t, enrollments, 2)
			assert.Equal(t, studentID1, enrollments[0].StudentID)
			assert.Equal(t, studentID2, enrollments[1].StudentID)
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/classes/"+classID.String()+"/students/bulk", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Success bool                   `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, float64(2), resp.Data["count"])
	})

	t.Run("Get Enrollment Success", func(t *testing.T) {
		id := uuid.New()
		expectedEnrollment := &entity.Enrollment{
			ID:       id,
			TenantID: "tenant-1",
			Status:   "enrolled",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expectedEnrollment, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/enrollments/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool              `json:"success"`
			Data    entity.Enrollment `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, id, resp.Data.ID)
		assert.Equal(t, "enrolled", resp.Data.Status)
	})
}
