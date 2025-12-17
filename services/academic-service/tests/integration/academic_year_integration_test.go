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

func TestAcademicYearIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockAcademicYearRepository(ctrl)

	// Real UseCase with Mock Repo
	u := usecase.NewAcademicYearUseCase(mockRepo, time.Second*2)

	// Real Handler with Real UseCase
	h := handler.NewAcademicYearHandler(u)

	// Setup Router
	r := gin.New()
	academicYears := r.Group("/api/v1/academic-years")
	{
		academicYears.POST("", h.Create)
		academicYears.GET("/:id", h.GetByID)
		academicYears.GET("/tenant/:tenant_id", h.List)
		academicYears.PUT("/:id", h.Update)
	}

	t.Run("Create Academic Year Success", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(1, 0, 0)
		reqBody := map[string]interface{}{
			"tenant_id":  "tenant-1",
			"name":       "2023/2024",
			"start_date": startDate,
			"end_date":   endDate,
			"is_active":  true,
		}
		body, _ := json.Marshal(reqBody)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, ay *entity.AcademicYear) error {
			assert.Equal(t, "tenant-1", ay.TenantID)
			assert.Equal(t, "2023/2024", ay.Name)
			assert.True(t, ay.IsActive)
			ay.ID = uuid.New()
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/academic-years", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Success bool                `json:"success"`
			Data    entity.AcademicYear `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.NotEmpty(t, resp.Data.ID)
		assert.Equal(t, "2023/2024", resp.Data.Name)
	})

	t.Run("Get Academic Year By ID Success", func(t *testing.T) {
		id := uuid.New()
		ay := &entity.AcademicYear{
			ID:       id,
			TenantID: "tenant-1",
			Name:     "2023/2024",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(ay, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/academic-years/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("List Academic Years Success", func(t *testing.T) {
		tenantID := "tenant-1"
		ays := []entity.AcademicYear{
			{ID: uuid.New(), Name: "2023/2024"},
			{ID: uuid.New(), Name: "2024/2025"},
		}

		mockRepo.EXPECT().List(gomock.Any(), tenantID).Return(ays, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/academic-years/tenant/"+tenantID, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Success bool                  `json:"success"`
			Data    []entity.AcademicYear `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Data, 2)
	})
}
