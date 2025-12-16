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

func TestTeacherIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockTeacherRepository(ctrl)

	// Real UseCase with Mock Repo
	u := usecase.NewTeacherUseCase(mockRepo, time.Second*2)

	// Real Handler with Real UseCase
	h := handler.NewTeacherHandler(u)

	// Setup Router
	r := gin.New()
	teachers := r.Group("/api/v1/teachers")
	{
		teachers.POST("", h.Create)
		teachers.GET("", h.List)
	}

	t.Run("Create Teacher Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			"name":      "Budi Santoso",
			"nip":       "198001012000121001",
			"status":    "active",
		}
		body, _ := json.Marshal(reqBody)

		// Expect Create
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tchr *entity.Teacher) error {
			assert.Equal(t, "Budi Santoso", tchr.Name)
			assert.Equal(t, "tenant-1", tchr.TenantID)
			tchr.ID = uuid.New()
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/teachers", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool           `json:"success"`
			Data    entity.Teacher `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "Budi Santoso", resp.Data.Name)
		assert.NotEmpty(t, resp.Data.ID)
	})

	t.Run("List Teachers Success", func(t *testing.T) {
		tenantID := "tenant-1"
		teachersList := []entity.Teacher{
			{ID: uuid.New(), TenantID: tenantID, Name: "Teacher 1"},
			{ID: uuid.New(), TenantID: tenantID, Name: "Teacher 2"},
		}
		total := 2

		mockRepo.EXPECT().List(gomock.Any(), tenantID, 10, 0).Return(teachersList, total, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/teachers?tenant_id="+tenantID, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool `json:"success"`
			Data    struct {
				Teachers []entity.Teacher `json:"teachers"`
				Total    int              `json:"total"`
			} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, total, resp.Data.Total)
		assert.Len(t, resp.Data.Teachers, 2)
	})
}
