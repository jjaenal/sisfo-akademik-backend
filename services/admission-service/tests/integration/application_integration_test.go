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
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApplicationIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repository
	mockRepo := mocks.NewMockApplicationRepository(ctrl)

	// Real UseCase with Mock Repo and nil RabbitClient
	u := usecase.NewApplicationUseCase(mockRepo, nil, time.Second*2)

	// Real Handler with Real UseCase
	h := handler.NewApplicationHandler(u)

	// Setup Router
	r := gin.New()
	applications := r.Group("/api/v1/applications")
	{
		applications.POST("", h.Submit)
		applications.GET("/status", h.GetStatus)
	}

	t.Run("Submit Application Success", func(t *testing.T) {
		tenantID := uuid.New()
		admissionPeriodID := uuid.New()
		reqBody := map[string]interface{}{
			"tenant_id":           tenantID.String(),
			"admission_period_id": admissionPeriodID.String(),
			"first_name":          "John",
			"last_name":           "Doe",
			"email":               "john.doe@example.com",
			"phone_number":        "1234567890",
			"date_of_birth":       "2010-01-01",
			"place_of_birth":      "Jakarta",
			"gender":              "Male",
			"address":             "Jalan Sudirman",
			"previous_school":     "SMP 1",
			"average_score":       85.5,
		}
		body, _ := json.Marshal(reqBody)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, app *entity.Application) error {
			assert.Equal(t, "John", app.FirstName)
			assert.Equal(t, "Doe", app.LastName)
			assert.NotEmpty(t, app.RegistrationNumber)
			assert.Equal(t, entity.ApplicationStatusSubmitted, app.Status)
			app.ID = uuid.New()
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/applications", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool               `json:"success"`
			Data    entity.Application `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "John", resp.Data.FirstName)
		assert.NotEmpty(t, resp.Data.RegistrationNumber)
	})

	t.Run("Get Application Status Success", func(t *testing.T) {
		regNum := "REG-20231010-1234"
		app := &entity.Application{
			ID:                 uuid.New(),
			RegistrationNumber: regNum,
			FirstName:          "Jane",
			LastName:           "Doe",
			Status:             entity.ApplicationStatusVerified,
		}

		mockRepo.EXPECT().GetByRegistrationNumber(gomock.Any(), regNum).Return(app, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/applications/status?registration_number="+regNum, nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool               `json:"success"`
			Data    entity.Application `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, "Jane", resp.Data.FirstName)
		assert.Equal(t, string(entity.ApplicationStatusVerified), string(resp.Data.Status))
	})
}
