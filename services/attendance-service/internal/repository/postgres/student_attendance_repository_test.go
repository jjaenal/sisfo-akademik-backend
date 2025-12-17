package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStudentAttendanceRepository_Create(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewStudentAttendanceRepository(db)

	tenantID := uuid.New().String()
	studentID := uuid.New()
	classID := uuid.New()
	semesterID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)
	attendanceDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	attendance := &entity.StudentAttendance{
		TenantID:       tenantID,
		StudentID:      studentID,
		ClassID:        classID,
		SemesterID:     semesterID,
		AttendanceDate: attendanceDate,
		Status:         entity.AttendanceStatusPresent,
		Notes:          "On time",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := repo.Create(context.Background(), attendance)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, attendance.ID)

	// Verify we can retrieve it
	stored, err := repo.GetByID(context.Background(), attendance.ID)
	require.NoError(t, err)
	require.NotNil(t, stored)

	assert.Equal(t, attendance.ID, stored.ID)
	assert.Equal(t, attendance.TenantID, stored.TenantID)
	assert.Equal(t, attendance.StudentID, stored.StudentID)
	assert.Equal(t, attendance.ClassID, stored.ClassID)
	assert.Equal(t, attendance.SemesterID, stored.SemesterID)
	assert.Equal(t, attendance.Status, stored.Status)
	assert.Equal(t, attendance.Notes, stored.Notes)
	// Check time within a small delta due to DB roundtrip precision
	assert.Equal(t, attendance.AttendanceDate, stored.AttendanceDate)
}

func TestStudentAttendanceRepository_GetByClassAndDate(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewStudentAttendanceRepository(db)

	tenantID := uuid.New().String()
	classID := uuid.New()
	semesterID := uuid.New()
	date := time.Now().Truncate(time.Second)

	// Create 3 records, 2 for the target class/date, 1 for another
	att1 := &entity.StudentAttendance{
		TenantID:       tenantID,
		StudentID:      uuid.New(),
		ClassID:        classID,
		SemesterID:     semesterID,
		AttendanceDate: date,
		Status:         entity.AttendanceStatusPresent,
	}
	att2 := &entity.StudentAttendance{
		TenantID:       tenantID,
		StudentID:      uuid.New(),
		ClassID:        classID,
		SemesterID:     semesterID,
		AttendanceDate: date,
		Status:         entity.AttendanceStatusAbsent,
	}
	att3 := &entity.StudentAttendance{ // Different class
		TenantID:       tenantID,
		StudentID:      uuid.New(),
		ClassID:        uuid.New(),
		SemesterID:     semesterID,
		AttendanceDate: date,
		Status:         entity.AttendanceStatusPresent,
	}

	require.NoError(t, repo.Create(context.Background(), att1))
	require.NoError(t, repo.Create(context.Background(), att2))
	require.NoError(t, repo.Create(context.Background(), att3))

	// Test GetByClassAndDate
	results, err := repo.GetByClassAndDate(context.Background(), classID, date)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify IDs
	ids := make(map[uuid.UUID]bool)
	for _, r := range results {
		ids[r.ID] = true
	}
	assert.True(t, ids[att1.ID])
	assert.True(t, ids[att2.ID])
	assert.False(t, ids[att3.ID])
}
