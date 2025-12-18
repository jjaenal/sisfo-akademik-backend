package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGradeRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewGradeRepository(db)
	assessmentRepo := NewAssessmentRepository(db)
	categoryRepo := NewGradeCategoryRepository(db)

	ctx := context.Background()

	// 1. Create GradeCategory
	tenantID := "tenant-1"
	category := &entity.GradeCategory{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        "Quiz " + uuid.New().String(),
		Description: "Quiz category",
		Weight:      10,
	}
	err := categoryRepo.Create(ctx, category)
	assert.NoError(t, err)

	// 2. Create Assessment
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
	err = assessmentRepo.Create(ctx, assessment)
	assert.NoError(t, err)

	// 3. Create Grade
	gradeID := uuid.New()
	studentID := uuid.New()
	grade := &entity.Grade{
		ID:           gradeID,
		TenantID:     tenantID,
		AssessmentID: assessment.ID,
		StudentID:    studentID,
		Score:        85.5,
		Status:       entity.GradeStatusFinal,
		GradedBy:     assessment.TeacherID,
	}

	// Test Create
	err = repo.Create(ctx, grade)
	assert.NoError(t, err)

	// Test GetByStudentAndAssessment
	fetchedByStudent, err := repo.GetByStudentAndAssessment(ctx, grade.StudentID, grade.AssessmentID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedByStudent)
	assert.Equal(t, grade.ID, fetchedByStudent.ID)
	assert.Equal(t, grade.Score, fetchedByStudent.Score)

	// Test GetByStudentID
	fetchedList, err := repo.GetByStudentID(ctx, grade.StudentID)
	assert.NoError(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.Equal(t, grade.ID, fetchedList[0].ID)
}
