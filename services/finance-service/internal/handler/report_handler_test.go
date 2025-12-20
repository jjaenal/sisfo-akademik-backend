package handler_test

import (
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

func TestReportHandler_GetDailyRevenue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockReportUseCase(ctrl)
	h := handler.NewReportHandler(mockUseCase)

	tenantID := uuid.New()
	startDate := "2023-01-01"
	endDate := "2023-01-31"

	t.Run("success", func(t *testing.T) {
		reports := []*entity.DailyRevenue{
			{Date: "2023-01-01", Amount: 1000},
		}

		mockUseCase.EXPECT().GetDailyRevenue(gomock.Any(), tenantID, gomock.Any(), gomock.Any()).Return(reports, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/revenue/daily?tenant_id="+tenantID.String()+"&start_date="+startDate+"&end_date="+endDate, nil)

		h.GetDailyRevenue(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - missing params", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/revenue/daily?tenant_id="+tenantID.String(), nil)

		h.GetDailyRevenue(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("bad request - invalid date", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/revenue/daily?tenant_id="+tenantID.String()+"&start_date=invalid&end_date="+endDate, nil)

		h.GetDailyRevenue(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUseCase.EXPECT().GetDailyRevenue(gomock.Any(), tenantID, gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/revenue/daily?tenant_id="+tenantID.String()+"&start_date="+startDate+"&end_date="+endDate, nil)

		h.GetDailyRevenue(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestReportHandler_GetMonthlyRevenue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockReportUseCase(ctrl)
	h := handler.NewReportHandler(mockUseCase)

	tenantID := uuid.New()
	year := "2023"

	t.Run("success", func(t *testing.T) {
		reports := []*entity.MonthlyRevenue{
			{Month: "January", Amount: 1000},
		}

		mockUseCase.EXPECT().GetMonthlyRevenue(gomock.Any(), tenantID, 2023).Return(reports, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/revenue/monthly?tenant_id="+tenantID.String()+"&year="+year, nil)

		h.GetMonthlyRevenue(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	
	t.Run("bad request - missing year", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/revenue/monthly?tenant_id="+tenantID.String(), nil)

		h.GetMonthlyRevenue(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request - invalid year", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/revenue/monthly?tenant_id="+tenantID.String()+"&year=invalid", nil)

		h.GetMonthlyRevenue(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestReportHandler_GetOutstandingInvoices(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockReportUseCase(ctrl)
	h := handler.NewReportHandler(mockUseCase)

	tenantID := uuid.New()

	t.Run("success", func(t *testing.T) {
		reports := []*entity.OutstandingReport{
			{InvoiceNumber: "INV-001", Amount: 1000},
		}

		mockUseCase.EXPECT().GetOutstandingInvoices(gomock.Any(), tenantID).Return(reports, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/outstanding?tenant_id="+tenantID.String(), nil)

		h.GetOutstandingInvoices(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUseCase.EXPECT().GetOutstandingInvoices(gomock.Any(), tenantID).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/outstanding?tenant_id="+tenantID.String(), nil)

		h.GetOutstandingInvoices(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestReportHandler_GetStudentHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockReportUseCase(ctrl)
	h := handler.NewReportHandler(mockUseCase)

	tenantID := uuid.New()
	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		history := &entity.StudentHistory{
			Invoices: []*entity.Invoice{},
			Payments: []*entity.Payment{},
		}

		mockUseCase.EXPECT().GetStudentHistory(gomock.Any(), tenantID, studentID).Return(history, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "student_id", Value: studentID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/student/"+studentID.String()+"/history?tenant_id="+tenantID.String(), nil)

		h.GetStudentHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	
	t.Run("bad request - invalid student id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "student_id", Value: "invalid"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/student/invalid/history?tenant_id="+tenantID.String(), nil)

		h.GetStudentHistory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUseCase.EXPECT().GetStudentHistory(gomock.Any(), tenantID, studentID).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "student_id", Value: studentID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/finance/reports/student/"+studentID.String()+"/history?tenant_id="+tenantID.String(), nil)

		h.GetStudentHistory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
