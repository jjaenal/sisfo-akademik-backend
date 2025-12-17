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

func TestClassIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockClassRepository(ctrl)

	// Real UseCase with Mock Repo
	u := usecase.NewClassUseCase(mockRepo, time.Second*2)

	// Real Handler with Real UseCase
	h := handler.NewClassHandler(u)

	// Setup Router
	r := gin.New()
	classes := r.Group("/api/v1/classes")
	{
		classes.POST("", h.Create)
		classes.GET("/:id", h.GetByID)
		classes.GET("", h.List)
		classes.PUT("/:id", h.Update)
		classes.DELETE("/:id", h.Delete)
	}

	t.Run("Create Class Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			"name":      "Class 1A",
			"level":     10,
			"major":     "Science",
			"capacity":  30,
		}
		body, _ := json.Marshal(reqBody)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, c *entity.Class) error {
			assert.Equal(t, "tenant-1", c.TenantID)
			assert.Equal(t, "Class 1A", c.Name)
			assert.Equal(t, 10, c.Level)
			c.ID = uuid.New()
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/classes", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Success bool         `json:"success"`
			Data    entity.Class `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.NotEmpty(t, resp.Data.ID)
		assert.Equal(t, "Class 1A", resp.Data.Name)
	})

	t.Run("Get Class By ID Success", func(t *testing.T) {
		id := uuid.New()
		c := &entity.Class{
			ID:       id,
			TenantID: "tenant-1",
			Name:     "Class 1A",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(c, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/classes/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("List Classes Success", func(t *testing.T) {
		tenantID := "tenant-1"
		classesList := []entity.Class{
			{ID: uuid.New(), Name: "Class 1A"},
			{ID: uuid.New(), Name: "Class 1B"},
		}

		mockRepo.EXPECT().List(gomock.Any(), tenantID, 10, 0).Return(classesList, 2, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/classes?tenant_id="+tenantID, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Success bool `json:"success"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, float64(2), resp.Data["total"])
	})
	
	t.Run("Update Class Success", func(t *testing.T) {
		id := uuid.New()
		reqBody := map[string]interface{}{
			"name":      "Class 1A Updated",
			"level":     10,
			"major":     "Science",
			"capacity":  35,
		}
		body, _ := json.Marshal(reqBody)
		
		existing := &entity.Class{
			ID:       id,
			TenantID: "tenant-1",
			Name:     "Class 1A",
			Level:    10,
			Major:    "Science",
			Capacity: 30,
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, c *entity.Class) error {
			assert.Equal(t, "Class 1A Updated", c.Name)
			assert.Equal(t, 35, c.Capacity)
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/classes/"+id.String(), bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Class Success", func(t *testing.T) {
		id := uuid.New()

		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/classes/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
