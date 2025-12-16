package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestStudentIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockStudentRepository(ctrl)

	// Real UseCase with Mock Repo
	u := usecase.NewStudentUseCase(mockRepo, time.Second*2)

	// Real Handler with Real UseCase
	h := handler.NewStudentHandler(u)

	// Setup Router
	r := gin.New()
	students := r.Group("/api/v1/students")
	{
		students.POST("", h.Create)
		students.GET("/:id", h.GetByID)
		students.GET("", h.List)
	}

	t.Run("Create Student Flow", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			"name":      "Jane Doe",
			"status":    "active",
		}
		body, _ := json.Marshal(reqBody)

		// Expectation on Repository
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *entity.Student) error {
			assert.Equal(t, "Jane Doe", s.Name)
			assert.Equal(t, "tenant-1", s.TenantID)
			s.ID = uuid.New() // Simulate ID generation
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool           `json:"success"`
			Data    entity.Student `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "Jane Doe", resp.Data.Name)
		assert.NotEmpty(t, resp.Data.ID)
	})

	t.Run("Get Student Flow", func(t *testing.T) {
		id := uuid.New()
		expectedStudent := &entity.Student{
			ID:       id,
			TenantID: "tenant-1",
			Name:     "Jane Doe",
			Status:   "active",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expectedStudent, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/students/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Success bool           `json:"success"`
			Data    entity.Student `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, id, resp.Data.ID)
		assert.Equal(t, "Jane Doe", resp.Data.Name)
	})
}
