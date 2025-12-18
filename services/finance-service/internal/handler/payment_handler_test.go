package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPaymentHandler_Record(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUseCase := new(mocks.PaymentUseCaseMock)
		handler := NewPaymentHandler(mockUseCase)

		tenantID := uuid.New()
		invoiceID := uuid.New()
		
		req := struct {
			TenantID        uuid.UUID            `json:"tenant_id"`
			InvoiceID       uuid.UUID            `json:"invoice_id"`
			Amount          float64              `json:"amount"`
			PaymentMethod   entity.PaymentMethod `json:"payment_method"`
			ReferenceNumber string               `json:"reference_number"`
		}{
			TenantID:        tenantID,
			InvoiceID:       invoiceID,
			Amount:          100000,
			PaymentMethod:   entity.PaymentMethodTransfer,
			ReferenceNumber: "REF123",
		}

		mockUseCase.On("RecordPayment", mock.Anything, mock.MatchedBy(func(p *entity.Payment) bool {
			return p.TenantID == tenantID && p.InvoiceID == invoiceID && p.Amount == 100000
		})).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Record(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		mockUseCase := new(mocks.PaymentUseCaseMock)
		handler := NewPaymentHandler(mockUseCase)

		req := struct {
			Amount float64 `json:"amount"`
		}{
			Amount: -100,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Record(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUseCase := new(mocks.PaymentUseCaseMock)
		handler := NewPaymentHandler(mockUseCase)

		tenantID := uuid.New()
		invoiceID := uuid.New()
		
		req := struct {
			TenantID        uuid.UUID            `json:"tenant_id"`
			InvoiceID       uuid.UUID            `json:"invoice_id"`
			Amount          float64              `json:"amount"`
			PaymentMethod   entity.PaymentMethod `json:"payment_method"`
			ReferenceNumber string               `json:"reference_number"`
		}{
			TenantID:        tenantID,
			InvoiceID:       invoiceID,
			Amount:          100000,
			PaymentMethod:   entity.PaymentMethodTransfer,
			ReferenceNumber: "REF123",
		}

		mockUseCase.On("RecordPayment", mock.Anything, mock.Anything).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Record(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPaymentHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUseCase := new(mocks.PaymentUseCaseMock)
		handler := NewPaymentHandler(mockUseCase)

		paymentID := uuid.New()
		payment := &entity.Payment{
			ID:              paymentID,
			Amount:          100000,
			Status:          entity.PaymentStatusSuccess,
			CreatedAt:       time.Now(),
		}

		mockUseCase.On("GetByID", mock.Anything, paymentID).Return(payment, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: paymentID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/payments/"+paymentID.String(), nil)

		handler.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockUseCase := new(mocks.PaymentUseCaseMock)
		handler := NewPaymentHandler(mockUseCase)

		paymentID := uuid.New()

		mockUseCase.On("GetByID", mock.Anything, paymentID).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: paymentID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/payments/"+paymentID.String(), nil)

		handler.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
