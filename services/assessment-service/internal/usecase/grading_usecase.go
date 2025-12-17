package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
)

type gradingUseCase struct {
	assessmentRepo repository.AssessmentRepository
	gradeRepo      repository.GradeRepository
	contextTimeout time.Duration
}

func NewGradingUseCase(assessmentRepo repository.AssessmentRepository, gradeRepo repository.GradeRepository, timeout time.Duration) domainUseCase.GradingUseCase {
	return &gradingUseCase{
		assessmentRepo: assessmentRepo,
		gradeRepo:      gradeRepo,
		contextTimeout: timeout,
	}
}

func (u *gradingUseCase) CreateAssessment(ctx context.Context, assessment *entity.Assessment) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errs := assessment.Validate(); len(errs) > 0 {
		return errors.New("validation error")
	}

	return u.assessmentRepo.Create(ctx, assessment)
}

func (u *gradingUseCase) InputGrade(ctx context.Context, grade *entity.Grade) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errs := grade.Validate(); len(errs) > 0 {
		return errors.New("validation error")
	}

	// Check if grade already exists
	existing, _ := u.gradeRepo.GetByStudentAndAssessment(ctx, grade.StudentID, grade.AssessmentID)
	if existing != nil {
		grade.ID = existing.ID
		return u.gradeRepo.Update(ctx, grade)
	}

	return u.gradeRepo.Create(ctx, grade)
}

func (u *gradingUseCase) GetStudentGrades(ctx context.Context, studentID, classID, semesterID uuid.UUID) ([]*entity.Grade, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// TODO: Filter by class and semester needs join query or more specific repo methods
	// For now, returning all grades for the student
	filter := map[string]interface{}{
		"student_id": studentID,
	}
	return u.gradeRepo.List(ctx, filter)
}

func (u *gradingUseCase) CalculateFinalScore(ctx context.Context, studentID, subjectID, semesterID uuid.UUID) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1. Get all assessments for the subject in the semester
	// This requires a more complex query or multiple calls.
	// For MVP, we assume we can list assessments and then filter.
	// Ideally: assessmentRepo.GetBySubjectAndSemester(subjectID, semesterID)
	
	filter := map[string]interface{}{
		"subject_id": subjectID,
	}
	assessments, err := u.assessmentRepo.List(ctx, filter)
	if err != nil {
		return 0, err
	}

	if len(assessments) == 0 {
		return 0, nil
	}

	totalScore := 0.0
	totalWeight := 0.0

	for _, assessment := range assessments {
		grade, err := u.gradeRepo.GetByStudentAndAssessment(ctx, studentID, assessment.ID)
		if err != nil {
			continue // Skip missing grades or handle as 0
		}
		if grade != nil {
			// Simplified calculation: Assuming weight is on GradeCategory (which we need to fetch)
			// For now, simple average or sum
			totalScore += grade.Score
			totalWeight += 1 // Placeholder for weight
		}
	}

	if totalWeight == 0 {
		return 0, nil
	}

	return totalScore / totalWeight, nil
}
