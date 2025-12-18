package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
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

	ctx := context.Background()
	err := repo.Create(ctx, attendance)
	assert.NoError(t, err)

	// Test GetByClassAndDate
	results, err := repo.GetByClassAndDate(ctx, classID, attendanceDate)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, attendance.ID, results[0].ID)

	// Test GetByID
	fetched, err := repo.GetByID(ctx, attendance.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, attendance.ID, fetched.ID)

	// Test Update
	attendance.Status = entity.AttendanceStatusAbsent
	attendance.Notes = "Updated notes"
	err = repo.Update(ctx, attendance)
	assert.NoError(t, err)

	fetchedAfterUpdate, err := repo.GetByID(ctx, attendance.ID)
	assert.NoError(t, err)
	assert.Equal(t, entity.AttendanceStatusAbsent, fetchedAfterUpdate.Status)
	assert.Equal(t, "Updated notes", fetchedAfterUpdate.Notes)

	// Test GetSummary
	summary, err := repo.GetSummary(ctx, studentID, semesterID)
	assert.NoError(t, err)
	assert.Equal(t, 1, summary[string(entity.AttendanceStatusAbsent)])

	// Test BulkCreate
	studentID2 := uuid.New()
	attendances := []*entity.StudentAttendance{
		{
			ID:             uuid.New(),
			TenantID:       tenantID,
			StudentID:      studentID2,
			ClassID:        classID,
			SemesterID:     semesterID,
			AttendanceDate: attendanceDate,
			Status:         entity.AttendanceStatusPresent,
			Notes:          "Bulk 1",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
		{
			ID:             uuid.New(),
			TenantID:       tenantID,
			StudentID:      uuid.New(),
			ClassID:        classID,
			SemesterID:     semesterID,
			AttendanceDate: attendanceDate,
			Status:         entity.AttendanceStatusLate,
			Notes:          "Bulk 2",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}
	err = repo.BulkCreate(ctx, attendances)
	assert.NoError(t, err)

	resultsBulk, err := repo.GetByClassAndDate(ctx, classID, attendanceDate)
	assert.NoError(t, err)
	// 1 from initial create + 2 from bulk create = 3
	assert.Len(t, resultsBulk, 3)
}
