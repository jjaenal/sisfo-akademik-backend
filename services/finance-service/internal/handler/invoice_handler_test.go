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

func TestInvoiceHandler_Generate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockInvoiceUseCase(ctrl)
	h := handler.NewInvoiceHandler(mockUseCase)

	tenantID := uuid.New()
	studentID := uuid.New()
	billingConfigID := uuid.New()

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]string{
			"tenant_id":         tenantID.String(),
			"student_id":        studentID.String(),
			"billing_config_id": billingConfigID.String(),
		}
		body, _ := json.Marshal(reqBody)

		expectedInvoice := &entity.Invoice{
			ID:        uuid.New(),
			TenantID:  tenantID,
			StudentID: studentID,
			Amount:    100000,
			Status:    entity.InvoiceStatusUnpaid,
		}

		mockUseCase.EXPECT().Generate(gomock.Any(), tenantID, studentID, billingConfigID).Return(expectedInvoice, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices/generate", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Generate(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp["success"].(bool))
	})

	t.Run("bad request", func(t *testing.T) {
		reqBody := map[string]string{
			"tenant_id": "invalid",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices/generate", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Generate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := map[string]string{
			"tenant_id":         tenantID.String(),
			"student_id":        studentID.String(),
			"billing_config_id": billingConfigID.String(),
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Generate(gomock.Any(), tenantID, studentID, billingConfigID).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/invoices/generate", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Generate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestInvoiceHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockInvoiceUseCase(ctrl)
	h := handler.NewInvoiceHandler(mockUseCase)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedInvoice := &entity.Invoice{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(expectedInvoice, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices/invalid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		h.GetByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestInvoiceHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockInvoiceUseCase(ctrl)
	h := handler.NewInvoiceHandler(mockUseCase)

	tenantID := uuid.New()
	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedInvoices := []*entity.Invoice{{TenantID: tenantID}}
		mockUseCase.EXPECT().List(gomock.Any(), tenantID, studentID, entity.InvoiceStatusUnpaid).Return(expectedInvoices, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices?tenant_id="+tenantID.String()+"&student_id="+studentID.String()+"&status=UNPAID", nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid tenant_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices?tenant_id=invalid", nil)

		h.List(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request - invalid student_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices?tenant_id="+tenantID.String()+"&student_id=invalid", nil)

		h.List(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUseCase.EXPECT().List(gomock.Any(), tenantID, uuid.Nil, entity.InvoiceStatus("")).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/invoices?tenant_id="+tenantID.String(), nil)

		h.List(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
