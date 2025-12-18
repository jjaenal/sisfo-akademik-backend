package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestStudentRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewStudentRepository(db)

	tenantID := uuid.New()
	studentID := uuid.New()
	classID := uuid.New()
	
	// Create
	student := &entity.Student{
		ID:        studentID,
		TenantID:  tenantID,
		Name:      "Budi Santoso",
		Status:    entity.StudentStatusActive,
		ClassID:   &classID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.Background()
	err := repo.Create(ctx, student)
	assert.NoError(t, err)

	// GetActive
	students, err := repo.GetActive(ctx, tenantID)
	assert.NoError(t, err)
	assert.NotEmpty(t, students)
	
	found := false
	for _, s := range students {
		if s.ID == student.ID {
			found = true
			assert.Equal(t, student.Name, s.Name)
			assert.Equal(t, student.Status, s.Status)
			break
		}
	}
	assert.True(t, found)

	// Test GetActive with inactive student
	inactiveStudentID := uuid.New()
	inactiveStudent := &entity.Student{
		ID:        inactiveStudentID,
		TenantID:  tenantID,
		Name:      "Siti Aminah",
		Status:    entity.StudentStatusInactive,
		ClassID:   &classID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(ctx, inactiveStudent)
	assert.NoError(t, err)

	studentsAfterInactive, err := repo.GetActive(ctx, tenantID)
	assert.NoError(t, err)
	
	foundInactive := false
	for _, s := range studentsAfterInactive {
		if s.ID == inactiveStudentID {
			foundInactive = true
			break
		}
	}
	assert.False(t, foundInactive, "Inactive student should not be returned by GetActive")
}
