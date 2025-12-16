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

func TestSchoolIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockSchoolRepository(ctrl)

	// Real UseCase with Mock Repo
	u := usecase.NewSchoolUseCase(mockRepo, time.Second*2)

	// Real Handler with Real UseCase
	h := handler.NewSchoolHandler(u)

	// Setup Router
	r := gin.New()
	schools := r.Group("/api/v1/schools")
	{
		schools.POST("", h.Create)
		schools.GET("/:id", h.GetByID)
		schools.GET("/tenant/:tenant_id", h.GetByTenantID)
	}

	t.Run("Create School Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":     "tenant-1",
			"name":          "SMA Negeri 1",
			"address":       "Jl. Merdeka No. 1",
			"phone":         "021-12345678",
			"email":         "info@sman1.sch.id",
			"website":       "https://sman1.sch.id",
			"logo_url":      "https://sman1.sch.id/logo.png",
			"accreditation": "A",
			"headmaster":    "Dr. Budi Santoso",
		}
		body, _ := json.Marshal(reqBody)

		// Expect Create
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, s *entity.School) error {
			assert.Equal(t, "SMA Negeri 1", s.Name)
			assert.Equal(t, "tenant-1", s.TenantID)
			assert.Equal(t, "A", s.Accreditation)
			s.ID = uuid.New()
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/schools", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool          `json:"success"`
			Data    entity.School `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "SMA Negeri 1", resp.Data.Name)
		assert.NotEmpty(t, resp.Data.ID)
	})

	t.Run("Get School By ID Success", func(t *testing.T) {
		id := uuid.New()
		expectedSchool := &entity.School{
			ID:       id,
			TenantID: "tenant-1",
			Name:     "SMA Negeri 1",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expectedSchool, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/schools/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool          `json:"success"`
			Data    entity.School `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, id, resp.Data.ID)
		assert.Equal(t, "SMA Negeri 1", resp.Data.Name)
	})

	t.Run("Get School By Tenant ID Success", func(t *testing.T) {
		tenantID := "tenant-1"
		expectedSchool := &entity.School{
			ID:       uuid.New(),
			TenantID: tenantID,
			Name:     "SMA Negeri 1",
		}

		mockRepo.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(expectedSchool, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/schools/tenant/"+tenantID, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool          `json:"success"`
			Data    entity.School `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, tenantID, resp.Data.TenantID)
	})
}
