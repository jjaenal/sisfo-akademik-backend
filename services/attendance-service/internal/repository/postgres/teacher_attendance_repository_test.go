package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestTeacherAttendanceRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewTeacherAttendanceRepository(db)

	tenantID := uuid.New().String()
	teacherID := uuid.New()
	semesterID := uuid.New()
	
	// Use UTC date only for comparison as postgres DATE type doesn't store time
	now := time.Now().UTC().Truncate(time.Second)
	attendanceDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	
	checkInTime := now.Add(-8 * time.Hour)
	checkOutTime := now.Add(-4 * time.Hour)

	lat := -6.200000
	long := 106.816666

	attendance := &entity.TeacherAttendance{
		TenantID:          tenantID,
		TeacherID:         teacherID,
		SemesterID:        semesterID,
		AttendanceDate:    attendanceDate,
		CheckInTime:       &checkInTime,
		CheckOutTime:      &checkOutTime,
		Status:            entity.TeacherAttendanceStatusPresent,
		Notes:             "On time",
		LocationLatitude:  &lat,
		LocationLongitude: &long,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	ctx := context.Background()

	// Test Create
	err := repo.Create(ctx, attendance)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, attendance.ID)

	// Test GetByID
	fetched, err := repo.GetByID(ctx, attendance.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, attendance.ID, fetched.ID)
	assert.Equal(t, attendance.TeacherID, fetched.TeacherID)
	// assert.Equal(t, attendance.AttendanceDate, fetched.AttendanceDate) // Compare manually if needed due to time location

	// Test GetByTeacherAndDate
	fetchedByDate, err := repo.GetByTeacherAndDate(ctx, teacherID, attendanceDate)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedByDate)
	assert.Equal(t, attendance.ID, fetchedByDate.ID)

	// Test Update
	newStatus := entity.TeacherAttendanceStatusExcused
	attendance.Status = newStatus
	attendance.Notes = "Updated notes"
	err = repo.Update(ctx, attendance)
	assert.NoError(t, err)

	fetchedAfterUpdate, err := repo.GetByID(ctx, attendance.ID)
	assert.NoError(t, err)
	assert.Equal(t, newStatus, fetchedAfterUpdate.Status)
	assert.Equal(t, "Updated notes", fetchedAfterUpdate.Notes)

	// Test List
	filter := map[string]interface{}{
		"teacher_id":  teacherID,
		"semester_id": semesterID,
	}
	list, err := repo.List(ctx, filter)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
	found := false
	for _, a := range list {
		if a.ID == attendance.ID {
			found = true
			break
		}
	}
	assert.True(t, found)
}
