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

	// Get all grades for student
	grades, err := u.gradeRepo.GetByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Filter by class and semester
	// This requires fetching assessment for each grade, which is N+1.
	// Ideally we should have a repository method that joins tables.
	// For now, we will do it in memory but it might be slow.
	// Optimization: Get all assessments for class and semester first, then filter grades.
	
	// However, GetByClassAndSemester is not available on AssessmentRepo yet?
	// We have GetByClassAndSubject.
	
	// Let's implement a simple filter loop.
	var filteredGrades []*entity.Grade
	for _, grade := range grades {
		assessment, err := u.assessmentRepo.GetByID(ctx, grade.AssessmentID)
		if err != nil {
			continue 
		}
		if assessment.ClassID == classID && assessment.SemesterID == semesterID {
			grade.Assessment = assessment
			filteredGrades = append(filteredGrades, grade)
		}
	}

	return filteredGrades, nil
}

func (u *gradingUseCase) CalculateFinalScore(ctx context.Context, studentID, classID, subjectID, semesterID uuid.UUID) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// 1. Get all assessments for the subject in the semester (filtered by class)
	assessments, err := u.assessmentRepo.GetByClassAndSubject(ctx, classID, subjectID)
	if err != nil {
		return 0, err
	}

	if len(assessments) == 0 {
		return 0, nil
	}

	totalScore := 0.0
	totalWeight := 0.0

	for _, assessment := range assessments {
		// Filter by semester if needed (GetByClassAndSubject returns all for class/subject, need to filter by semester?)
		// assessmentRepo.GetByClassAndSubject does NOT filter by semester in query currently.
		// Let's check repository implementation.
		// Repository query: SELECT ... FROM assessments WHERE class_id = $1 AND subject_id = $2 AND deleted_at IS NULL
		// So we should filter by semester here.
		if assessment.SemesterID != semesterID {
			continue
		}

		grade, err := u.gradeRepo.GetByStudentAndAssessment(ctx, studentID, assessment.ID)
		if err != nil {
			continue // Skip missing grades or handle as 0
		}
		
		// TODO: Fetch GradeCategory to get weight. 
		// For now, assuming equal weight or just summing up scores / count (simple average)
		// Or better, assume max_score is the weight base.
		
		if grade != nil {
			// Basic calculation: (Score / MaxScore) * 100
			// If we want weighted average, we need categories.
			// Let's implement simple average of percentages for now.
			if assessment.MaxScore > 0 {
				percentage := (grade.Score / assessment.MaxScore) * 100
				totalScore += percentage
				totalWeight += 1
			}
		}
	}

	if totalWeight == 0 {
		return 0, nil
	}

	return totalScore / totalWeight, nil
}
