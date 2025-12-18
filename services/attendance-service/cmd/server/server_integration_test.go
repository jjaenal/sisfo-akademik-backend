package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/stretchr/testify/assert"
)

func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		url = "postgres://dev:dev@localhost:55432/devdb?sslmode=disable"
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Skip("no db available:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Skip("db connect failed:", err)
	}
	return db
}

func ensureMigrations(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	
	files := []string{
		"../../migrations/001_initial_attendance.up.sql",
		"../../migrations/002_add_semester_id.up.sql",
	}

	for _, file := range files {
		b, err := os.ReadFile(file)
		if err == nil {
			_, _ = db.Exec(ctx, string(b))
		}
	}
}

func makeCfg() config.Config {
	return config.Config{
		Env:         "test",
		ServiceName: "attendance-service",
		HTTPPort:    0,
		PostgresURL: "postgres://dev:dev@localhost:55432/devdb?sslmode=disable",
	}
}

func setupServer(t *testing.T) (*gin.Engine, *pgxpool.Pool) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	ensureMigrations(t, db)

	repo := postgres.NewStudentAttendanceRepository(db)
	uc := usecase.NewStudentAttendanceUseCase(repo, nil, 5*time.Second)
	h := handler.NewStudentAttendanceHandler(uc)

	teacherRepo := postgres.NewTeacherAttendanceRepository(db)
	teacherUC := usecase.NewTeacherAttendanceUseCase(teacherRepo, 5*time.Second)
	teacherHandler := handler.NewTeacherAttendanceHandler(teacherUC)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	attendance := v1.Group("/attendance")
	{
		attendance.POST("/students", h.Create)
		attendance.POST("/students/bulk", h.BulkCreate)
		attendance.GET("/students", h.GetByClassAndDate)
		attendance.GET("/students/:id", h.GetByID)
		attendance.GET("/students/:id/summary", h.GetSummary)
		attendance.PUT("/students/:id", h.Update)

		attendance.GET("/reports/daily", h.GetDailyReport)
		attendance.GET("/reports/monthly", h.GetMonthlyReport)
		attendance.GET("/reports/class/:class_id", h.GetClassReport)

		attendance.POST("/teachers/checkin", teacherHandler.CheckIn)
		attendance.PUT("/teachers/checkout", teacherHandler.CheckOut)
		attendance.GET("/teachers", teacherHandler.GetByTeacherAndDate)
	}

	return r, db
}

func TestIntegration_StudentAttendance(t *testing.T) {
	r, db := setupServer(t)
	defer db.Close()

	studentID := uuid.New()
	classID := uuid.New()
	tenantID := uuid.NewString()

	t.Run("Create Attendance", func(t *testing.T) {
		body := map[string]interface{}{
			"tenant_id":       tenantID,
			"student_id":      studentID.String(),
			"class_id":        classID.String(),
			"semester_id":     uuid.New().String(),
			"attendance_date": time.Now().Format(time.RFC3339),
			"status":          "present",
			"notes":           "On time",
		}
		b, _ := json.Marshal(body)
		
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/attendance/students", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Summary", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/attendance/students/"+studentID.String()+"/summary", nil)
		q := req.URL.Query()
		q.Add("semester_id", uuid.NewString()) // Optional, but good to check
		req.URL.RawQuery = q.Encode()
		
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestIntegration_Reports(t *testing.T) {
	r, db := setupServer(t)
	defer db.Close()

	studentID := uuid.New()
	classID := uuid.New()
	tenantID := uuid.NewString()
	semesterID := uuid.New()
	date := time.Now()

	// Seed data
	t.Run("Seed Data", func(t *testing.T) {
		body := map[string]interface{}{
			"tenant_id":       tenantID,
			"student_id":      studentID.String(),
			"class_id":        classID.String(),
			"semester_id":     semesterID.String(),
			"attendance_date": date.Format(time.RFC3339),
			"status":          "present",
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/attendance/students", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Daily Report", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/attendance/reports/daily", nil)
		q := req.URL.Query()
		q.Add("tenant_id", tenantID)
		q.Add("date", date.Format("2006-01-02"))
		req.URL.RawQuery = q.Encode()

		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Monthly Report", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/attendance/reports/monthly", nil)
		q := req.URL.Query()
		q.Add("tenant_id", tenantID)
		q.Add("month", strconv.Itoa(int(date.Month())))
		q.Add("year", strconv.Itoa(date.Year()))
		req.URL.RawQuery = q.Encode()

		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Get Class Report", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/attendance/reports/class/"+classID.String(), nil)
		q := req.URL.Query()
		q.Add("tenant_id", tenantID)
		q.Add("start_date", date.Add(-24*time.Hour).Format("2006-01-02"))
		q.Add("end_date", date.Add(24*time.Hour).Format("2006-01-02"))
		req.URL.RawQuery = q.Encode()

		r.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestIntegration_BulkCheckIn(t *testing.T) {
	r, db := setupServer(t)
	defer db.Close()

	classID := uuid.New()
	tenantID := uuid.NewString()
	students := []string{uuid.NewString(), uuid.NewString()}

	t.Run("Bulk Create", func(t *testing.T) {
		body := []map[string]interface{}{
			{
				"tenant_id":       tenantID,
				"class_id":        classID.String(),
				"semester_id":     uuid.New().String(),
				"student_id":      students[0],
				"attendance_date": time.Now().Format(time.RFC3339),
				"status":          "present",
				"notes":           "Bulk check-in 1",
			},
			{
				"tenant_id":       tenantID,
				"class_id":        classID.String(),
				"semester_id":     uuid.New().String(),
				"student_id":      students[1],
				"attendance_date": time.Now().Format(time.RFC3339),
				"status":          "present",
				"notes":           "Bulk check-in 2",
			},
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/attendance/students/bulk", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestIntegration_TeacherAttendance(t *testing.T) {
	r, db := setupServer(t)
	defer db.Close()

	teacherID := uuid.New()
	tenantID := uuid.NewString()
	semesterID := uuid.New()

	t.Run("Check In", func(t *testing.T) {
		body := map[string]interface{}{
			"tenant_id":       tenantID,
			"teacher_id":      teacherID.String(),
			"semester_id":     semesterID.String(),
			"attendance_date": time.Now().Format(time.RFC3339),
			"check_in_time":   time.Now().Format(time.RFC3339),
			"location_latitude":  -6.200000,
			"location_longitude": 106.816666,
			"status":          "present",
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/attendance/teachers/checkin", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Check Out", func(t *testing.T) {
		body := map[string]interface{}{
			"tenant_id":       tenantID,
			"teacher_id":      teacherID.String(),
			"date":            time.Now().Format(time.RFC3339),
			"check_out_time":  time.Now().Add(8 * time.Hour).Format(time.RFC3339),
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/attendance/teachers/checkout", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
