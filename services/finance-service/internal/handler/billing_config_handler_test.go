package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBillingConfigHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockBillingConfigUseCase(ctrl)
	h := handler.NewBillingConfigHandler(mockUseCase)

	tenantID := uuid.New()

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": tenantID.String(),
			"name":      "SPP Monthly",
			"amount":    100000,
			"frequency": "MONTHLY",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, config *entity.BillingConfig) error {
			assert.Equal(t, tenantID, config.TenantID)
			assert.Equal(t, "SPP Monthly", config.Name)
			assert.Equal(t, float64(100000), config.Amount)
			assert.Equal(t, entity.BillingFrequency("MONTHLY"), config.Frequency)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing-configs", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "invalid",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing-configs", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": tenantID.String(),
			"name":      "SPP Monthly",
			"amount":    100000,
			"frequency": "MONTHLY",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/billing-configs", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestBillingConfigHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockBillingConfigUseCase(ctrl)
	h := handler.NewBillingConfigHandler(mockUseCase)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":      "SPP Updated",
			"amount":    150000,
			"frequency": "MONTHLY",
			"is_active": true,
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, config *entity.BillingConfig) error {
			assert.Equal(t, id, config.ID)
			assert.Equal(t, "SPP Updated", config.Name)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/billing-configs/"+id.String(), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/billing-configs/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		h.Update(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestBillingConfigHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockBillingConfigUseCase(ctrl)
	h := handler.NewBillingConfigHandler(mockUseCase)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := &entity.BillingConfig{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing-configs/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing-configs/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestBillingConfigHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockBillingConfigUseCase(ctrl)
	h := handler.NewBillingConfigHandler(mockUseCase)

	tenantID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := []*entity.BillingConfig{{TenantID: tenantID}}
		mockUseCase.EXPECT().List(gomock.Any(), tenantID).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing-configs?tenant_id="+tenantID.String(), nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid tenant_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/billing-configs?tenant_id=invalid", nil)

		h.List(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestBillingConfigHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockBillingConfigUseCase(ctrl)
	h := handler.NewBillingConfigHandler(mockUseCase)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockUseCase.EXPECT().Delete(gomock.Any(), id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/billing-configs/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/billing-configs/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		h.Delete(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
