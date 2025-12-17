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
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTemplateIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockTemplateRepository(ctrl)

	// Real UseCase with Mock Repo
	u := usecase.NewTemplateUseCase(mockRepo)

	// Real Handler with Real UseCase
	h := handler.NewTemplateHandler(u)

	// Setup Router
	r := gin.New()
	templates := r.Group("/api/v1/templates")
	{
		templates.POST("", h.Create)
		templates.GET("", h.List)
		templates.GET("/:id", h.GetByID)
		templates.PUT("/:id", h.Update)
		templates.DELETE("/:id", h.Delete)
	}

	t.Run("Create Template Success", func(t *testing.T) {
		tenantID := uuid.New()
		reqBody := map[string]interface{}{
			"tenant_id":  tenantID.String(),
			"name":       "Default Template",
			"is_default": true,
			"config": map[string]string{
				"header_text":     "School Header",
				"logo_url":        "http://logo.url",
				"primary_color":   "#000000",
				"secondary_color": "#FFFFFF",
				"footer_text":     "School Footer",
			},
		}
		body, _ := json.Marshal(reqBody)

		mockRepo.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return([]*entity.ReportCardTemplate{}, nil)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tmpl *entity.ReportCardTemplate) error {
			assert.Equal(t, tenantID, tmpl.TenantID)
			assert.Equal(t, "Default Template", tmpl.Name)
			assert.True(t, tmpl.IsDefault)
			assert.Equal(t, "School Header", tmpl.Config.HeaderText)
			tmpl.ID = uuid.New()
			tmpl.CreatedAt = time.Now()
			tmpl.UpdatedAt = time.Now()
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/templates", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool                      `json:"success"`
			Data    entity.ReportCardTemplate `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "Default Template", resp.Data.Name)
	})

	t.Run("List Templates Success", func(t *testing.T) {
		tenantID := uuid.New()
		tmplList := []*entity.ReportCardTemplate{
			{
				ID:       uuid.New(),
				TenantID: tenantID,
				Name:     "Template 1",
			},
			{
				ID:       uuid.New(),
				TenantID: tenantID,
				Name:     "Template 2",
			},
		}

		mockRepo.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(tmplList, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/templates?tenant_id="+tenantID.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool                         `json:"success"`
			Data    []*entity.ReportCardTemplate `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Len(t, resp.Data, 2)
	})

	t.Run("Get Template By ID Success", func(t *testing.T) {
		id := uuid.New()
		tmpl := &entity.ReportCardTemplate{
			ID:   id,
			Name: "My Template",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(tmpl, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/templates/"+id.String(), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool                      `json:"success"`
			Data    entity.ReportCardTemplate `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "My Template", resp.Data.Name)
	})
}
