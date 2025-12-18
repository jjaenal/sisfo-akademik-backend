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
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/handler"
)

func TestReportCardHandler_Generate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockReportCardUseCase(ctrl)
	h := handler.NewReportCardHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		tenantID := "tenant-123"
		studentID := uuid.New()
		classID := uuid.New()
		semesterID := uuid.New()

		reqBody := map[string]string{
			"tenant_id":   tenantID,
			"student_id":  studentID.String(),
			"class_id":    classID.String(),
			"semester_id": semesterID.String(),
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Generate(gomock.Any(), tenantID, studentID, classID, semesterID).Return(&entity.ReportCard{ID: uuid.New()}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/report-cards", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Generate(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid_input", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/report-cards", bytes.NewBufferString("invalid json"))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Generate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase_error", func(t *testing.T) {
		tenantID := "tenant-123"
		studentID := uuid.New()
		classID := uuid.New()
		semesterID := uuid.New()

		reqBody := map[string]string{
			"tenant_id":   tenantID,
			"student_id":  studentID.String(),
			"class_id":    classID.String(),
			"semester_id": semesterID.String(),
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Generate(gomock.Any(), tenantID, studentID, classID, semesterID).Return(nil, errors.New("internal error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/report-cards", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Generate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestReportCardHandler_GetPDF(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockReportCardUseCase(ctrl)
	h := handler.NewReportCardHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		pdfBytes := []byte("%PDF-1.4...")

		mockUseCase.EXPECT().GetPDF(gomock.Any(), id).Return(pdfBytes, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/report-cards/"+id.String()+"/pdf", nil)

		h.GetPDF(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/pdf", w.Header().Get("Content-Type"))
		assert.Equal(t, "attachment; filename=report_card.pdf", w.Header().Get("Content-Disposition"))
		assert.Equal(t, pdfBytes, w.Body.Bytes())
	})

	t.Run("invalid_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/report-cards/invalid-uuid/pdf", nil)

		h.GetPDF(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase_error", func(t *testing.T) {
		id := uuid.New()

		mockUseCase.EXPECT().GetPDF(gomock.Any(), id).Return(nil, errors.New("internal error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/report-cards/"+id.String()+"/pdf", nil)

		h.GetPDF(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
