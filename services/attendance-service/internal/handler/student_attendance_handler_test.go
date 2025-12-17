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

func TestStudentAttendanceHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	studentID := uuid.New()
	classID := uuid.New()
	semesterID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"student_id":      studentID.String(),
			"class_id":        classID.String(),
			"semester_id":     semesterID.String(),
			"attendance_date": now,
			"status":          "present",
			"notes":           "Present",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, att *entity.StudentAttendance) error {
			assert.Equal(t, studentID, att.StudentID)
			assert.Equal(t, classID, att.ClassID)
			assert.Equal(t, semesterID, att.SemesterID)
			assert.Equal(t, entity.AttendanceStatus("present"), att.Status)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/attendance", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"student_id": "invalid",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/attendance", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"student_id":      studentID.String(),
			"class_id":        classID.String(),
			"semester_id":     semesterID.String(),
			"attendance_date": now,
			"status":          "present",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/attendance", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestStudentAttendanceHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := &entity.StudentAttendance{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil) // Return nil, nil for not found (handler logic)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
