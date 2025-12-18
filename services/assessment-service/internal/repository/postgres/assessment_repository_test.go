package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestAssessmentRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewAssessmentRepository(db)
	categoryRepo := NewGradeCategoryRepository(db)

	ctx := context.Background()

	// Need a GradeCategory first as it's a foreign key
	tenantID := uuid.New()
	category := &entity.GradeCategory{
		ID:          uuid.New(),
		TenantID:    tenantID.String(),
		Name:        "Quiz " + uuid.New().String(),
		Description: "Quiz category",
		Weight:      10,
	}
	err := categoryRepo.Create(ctx, category)
	assert.NoError(t, err)

	assessmentID := uuid.New()
	assessment := &entity.Assessment{
		ID:              assessmentID,
		TenantID:        tenantID,
		TeacherID:       uuid.New(),
		SubjectID:       uuid.New(),
		ClassID:         uuid.New(),
		GradeCategoryID: category.ID,
		Name:            "Math Quiz 1",
		Date:            time.Now().Truncate(time.Second),
		MaxScore:        100,
		Description:     "Algebra basics",
	}

	// Test Create
	err = repo.Create(ctx, assessment)
	assert.NoError(t, err)

	// Test GetByID
	fetched, err := repo.GetByID(ctx, assessmentID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, assessment.ID, fetched.ID)
	assert.Equal(t, assessment.Name, fetched.Name)
	assert.Equal(t, assessment.MaxScore, fetched.MaxScore)

	// Test Update
	assessment.Name = "Math Quiz 1 (Revised)"
	assessment.MaxScore = 90
	err = repo.Update(ctx, assessment)
	assert.NoError(t, err)

	fetchedAfterUpdate, err := repo.GetByID(ctx, assessmentID)
	assert.NoError(t, err)
	assert.Equal(t, "Math Quiz 1 (Revised)", fetchedAfterUpdate.Name)
	assert.Equal(t, 90.0, fetchedAfterUpdate.MaxScore)

	// Test GetByClassAndSubject
	list, err := repo.GetByClassAndSubject(ctx, assessment.ClassID, assessment.SubjectID)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
	found := false
	for _, a := range list {
		if a.ID == assessment.ID {
			found = true
			break
		}
	}
	assert.True(t, found)

	// Test Delete
	err = repo.Delete(ctx, assessmentID)
	assert.NoError(t, err)

	fetchedAfterDelete, err := repo.GetByID(ctx, assessmentID)
	assert.Error(t, err)
	assert.Nil(t, fetchedAfterDelete)
}
