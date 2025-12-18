package handler_test

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
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTeacherAttendanceHandler_CheckIn(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTeacherAttendanceUseCase(ctrl)
	h := handler.NewTeacherAttendanceHandler(mockUseCase)

	tenantID := uuid.New().String()
	teacherID := uuid.New()
	semesterID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":       tenantID,
			"teacher_id":      teacherID.String(),
			"semester_id":     semesterID.String(),
			"attendance_date": now,
			"check_in_time":   now,
			"status":          "present",
			"notes":           "On time",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().CheckIn(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, att *entity.TeacherAttendance) error {
			assert.Equal(t, tenantID, att.TenantID)
			assert.Equal(t, teacherID, att.TeacherID)
			assert.Equal(t, semesterID, att.SemesterID)
			assert.Equal(t, entity.TeacherAttendanceStatus("present"), att.Status)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/teacher/check-in", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CheckIn(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"teacher_id": "invalid",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/teacher/check-in", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CheckIn(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":       tenantID,
			"teacher_id":      teacherID.String(),
			"semester_id":     semesterID.String(),
			"attendance_date": now,
			"check_in_time":   now,
			"status":          "present",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().CheckIn(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/teacher/check-in", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CheckIn(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTeacherAttendanceHandler_CheckOut(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTeacherAttendanceUseCase(ctrl)
	h := handler.NewTeacherAttendanceHandler(mockUseCase)

	teacherID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"teacher_id":     teacherID.String(),
			"date":           now,
			"check_out_time": now,
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().CheckOut(gomock.Any(), teacherID, gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/teacher/check-out", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CheckOut(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"teacher_id": "invalid",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/teacher/check-out", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CheckOut(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"teacher_id":     teacherID.String(),
			"date":           now,
			"check_out_time": now,
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().CheckOut(gomock.Any(), teacherID, gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/teacher/check-out", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.CheckOut(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
