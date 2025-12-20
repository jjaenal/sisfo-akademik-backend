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
	tenantID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":       tenantID.String(),
			"student_id":      studentID.String(),
			"class_id":        classID.String(),
			"semester_id":     semesterID.String(),
			"attendance_date": now,
			"status":          "present",
			"notes":           "Present",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, att *entity.StudentAttendance) error {
			assert.Equal(t, tenantID.String(), att.TenantID)
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
			"tenant_id":       tenantID.String(),
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

func TestStudentAttendanceHandler_BulkCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	studentID1 := uuid.New()
	studentID2 := uuid.New()
	classID := uuid.New()
	semesterID := uuid.New()
	tenantID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		reqBody := []map[string]interface{}{
			{
				"tenant_id":       tenantID.String(),
				"student_id":      studentID1.String(),
				"class_id":        classID.String(),
				"semester_id":     semesterID.String(),
				"attendance_date": now,
				"status":          "present",
			},
			{
				"tenant_id":       tenantID.String(),
				"student_id":      studentID2.String(),
				"class_id":        classID.String(),
				"semester_id":     semesterID.String(),
				"attendance_date": now,
				"status":          "absent",
			},
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().BulkCreate(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, atts []*entity.StudentAttendance) error {
			assert.Len(t, atts, 2)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/attendance/bulk", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.BulkCreate(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		reqBody := []map[string]interface{}{
			{
				"student_id": "invalid",
			},
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/attendance/bulk", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.BulkCreate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := []map[string]interface{}{
			{
				"tenant_id":       tenantID.String(),
				"student_id":      studentID1.String(),
				"class_id":        classID.String(),
				"semester_id":     semesterID.String(),
				"attendance_date": now,
				"status":          "present",
			},
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().BulkCreate(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/attendance/bulk", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.BulkCreate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestStudentAttendanceHandler_GetDailyReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	tenantID := uuid.New()
	date := time.Now().Format("2006-01-02")
	parsedDate, _ := time.Parse("2006-01-02", date)

	t.Run("success", func(t *testing.T) {
		expected := []*entity.StudentAttendance{{TenantID: tenantID.String()}}
		mockUseCase.EXPECT().GetDailyReport(gomock.Any(), tenantID, parsedDate).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/reports/daily?tenant_id="+tenantID.String()+"&date="+date, nil)

		h.GetDailyReport(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - missing params", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/reports/daily", nil)

		h.GetDailyReport(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestStudentAttendanceHandler_GetMonthlyReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	tenantID := uuid.New()
	month := 1
	year := 2024

	t.Run("success", func(t *testing.T) {
		expected := []*entity.StudentAttendance{{TenantID: tenantID.String()}}
		mockUseCase.EXPECT().GetMonthlyReport(gomock.Any(), tenantID, month, year).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/reports/monthly?tenant_id="+tenantID.String()+"&month=1&year=2024", nil)

		h.GetMonthlyReport(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestStudentAttendanceHandler_GetClassReport(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	classID := uuid.New()
	tenantID := uuid.New()
	startDate := "2024-01-01"
	endDate := "2024-01-31"
	parsedStart, _ := time.Parse("2006-01-02", startDate)
	parsedEnd, _ := time.Parse("2006-01-02", endDate)

	t.Run("success", func(t *testing.T) {
		expected := []*entity.StudentAttendance{{ClassID: classID}}
		mockUseCase.EXPECT().GetClassReport(gomock.Any(), tenantID, classID, parsedStart, parsedEnd).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "class_id", Value: classID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/reports/class/"+classID.String()+"?tenant_id="+tenantID.String()+"&start_date="+startDate+"&end_date="+endDate, nil)

		h.GetClassReport(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestStudentAttendanceHandler_GetSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	studentID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := map[string]int{"present": 10}
		mockUseCase.EXPECT().GetSummary(gomock.Any(), studentID, uuid.Nil).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: studentID.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/students/"+studentID.String()+"/summary", nil)

		h.GetSummary(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestStudentAttendanceHandler_GetByClassAndDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	classID := uuid.New()
	date := time.Now().Format("2006-01-02")
	parsedDate, _ := time.Parse("2006-01-02", date)

	t.Run("success", func(t *testing.T) {
		expected := []*entity.StudentAttendance{{ClassID: classID}}
		mockUseCase.EXPECT().GetByClassAndDate(gomock.Any(), classID, parsedDate).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/students?class_id="+classID.String()+"&date="+date, nil)

		h.GetByClassAndDate(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - missing params", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/students", nil)

		h.GetByClassAndDate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request - invalid class id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/students?class_id=invalid&date="+date, nil)

		h.GetByClassAndDate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request - invalid date format", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/students?class_id="+classID.String()+"&date=invalid", nil)

		h.GetByClassAndDate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUseCase.EXPECT().GetByClassAndDate(gomock.Any(), classID, parsedDate).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/attendance/students?class_id="+classID.String()+"&date="+date, nil)

		h.GetByClassAndDate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestStudentAttendanceHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentAttendanceUseCase(ctrl)
	h := handler.NewStudentAttendanceHandler(mockUseCase)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"status": "absent",
			"notes":  "sick",
		}
		body, _ := json.Marshal(reqBody)

		existing := &entity.StudentAttendance{ID: id, Status: "present"}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
		mockUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, att *entity.StudentAttendance) error {
			assert.Equal(t, entity.AttendanceStatus("absent"), att.Status)
			assert.Equal(t, "sick", att.Notes)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/attendance/students/"+id.String(), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"status": "absent",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/attendance/students/"+id.String(), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.Update(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
