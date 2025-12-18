package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/service"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
)

type reportCardUseCase struct {
	reportCardRepo    repository.ReportCardRepository
	gradeRepo         repository.GradeRepository
	assessmentRepo    repository.AssessmentRepository
	gradeCategoryRepo repository.GradeCategoryRepository
	fileStorage       service.FileStorage
}

func NewReportCardUseCase(
	rcRepo repository.ReportCardRepository,
	gRepo repository.GradeRepository,
	aRepo repository.AssessmentRepository,
	gcRepo repository.GradeCategoryRepository,
	fileStorage service.FileStorage,
) domainUseCase.ReportCardUseCase {
	return &reportCardUseCase{
		reportCardRepo:    rcRepo,
		gradeRepo:         gRepo,
		assessmentRepo:    aRepo,
		gradeCategoryRepo: gcRepo,
		fileStorage:       fileStorage,
	}
}

func (u *reportCardUseCase) Generate(ctx context.Context, tenantID string, studentID, classID, semesterID uuid.UUID) (*entity.ReportCard, error) {
	// 1. Check existing
	existing, err := u.reportCardRepo.GetByStudentAndSemester(ctx, studentID, semesterID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Status == entity.ReportCardStatusPublished {
		return nil, errors.New("report card already published")
	}

	// 2. Fetch Grades
	grades, err := u.gradeRepo.GetByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// 3. Aggregate Scores by Subject
	type subjectStats struct {
		TotalWeightedScore float64
		TotalWeight        float64
		SubjectName        string // Mocked for now
	}
	subjectMap := make(map[uuid.UUID]*subjectStats)

	for _, g := range grades {
		// Fetch Assessment
		assessment, err := u.assessmentRepo.GetByID(ctx, g.AssessmentID)
		if err != nil {
			continue // Skip if not found
		}
		// Filter by Class (simple check)
		if assessment.ClassID != classID {
			continue
		}

		// Fetch Category
		category, err := u.gradeCategoryRepo.GetByID(ctx, assessment.GradeCategoryID)
		if err != nil {
			continue
		}

		if _, ok := subjectMap[assessment.SubjectID]; !ok {
			subjectMap[assessment.SubjectID] = &subjectStats{
				SubjectName: "Subject " + assessment.SubjectID.String()[:8], // Placeholder
			}
		}
		stats := subjectMap[assessment.SubjectID]

		// Normalize score to 100
		normalizedScore := (g.Score / float64(assessment.MaxScore)) * 100
		stats.TotalWeightedScore += normalizedScore * category.Weight
		stats.TotalWeight += category.Weight
	}

	// 4. Build Details
	var details []entity.ReportCardDetail
	var totalPoints float64
	var totalCredits int

	for subjectID, stats := range subjectMap {
		finalScore := 0.0
		if stats.TotalWeight > 0 {
			finalScore = stats.TotalWeightedScore / stats.TotalWeight
		}

		// Grading Rule (Simple)
		gradeLetter := "E"
		points := 0.0
		if finalScore >= 90 {
			gradeLetter = "A"
			points = 4.0
		} else if finalScore >= 80 {
			gradeLetter = "B"
			points = 3.0
		} else if finalScore >= 70 {
			gradeLetter = "C"
			points = 2.0
		} else if finalScore >= 60 {
			gradeLetter = "D"
			points = 1.0
		}

		// Mock Credit
		credit := 3 // Standard credit
		
		details = append(details, entity.ReportCardDetail{
			ID:          uuid.New(),
			TenantID:    tenantID,
			ReportCardID: uuid.Nil, // Will be set on save
			SubjectID:   subjectID,
			SubjectName: stats.SubjectName,
			Credit:      credit,
			FinalScore:  finalScore,
			GradeLetter: gradeLetter,
			Comments:    "Generated",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})

		totalPoints += points * float64(credit)
		totalCredits += credit
	}

	// 5. Calculate GPA
	gpa := 0.0
	if totalCredits > 0 {
		gpa = totalPoints / float64(totalCredits)
	}

	// 6. Create or Update
	rc := &entity.ReportCard{
		ID:                uuid.New(),
		TenantID:          tenantID,
		StudentID:         studentID,
		ClassID:           classID,
		SemesterID:        semesterID,
		Status:            entity.ReportCardStatusGenerated,
		GPA:               gpa,
		TotalCredits:      totalCredits,
		Attendance:        0, // Placeholder
		AttendanceSummary: map[string]int{"present": 0, "absent": 0},
		Details:           details,
		GeneratedAt:       func() *time.Time { t := time.Now(); return &t }(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// 7. Generate PDF
	pdfBytes, err := u.generatePDF(rc)
	if err == nil {
		fileName := fmt.Sprintf("report_cards/%s/%s.pdf", tenantID, rc.ID)
		url, err := u.fileStorage.Upload(ctx, fileName, bytes.NewReader(pdfBytes))
		if err == nil {
			rc.PDFUrl = url
		}
	}

	if existing != nil {
		rc.ID = existing.ID
		rc.CreatedAt = existing.CreatedAt
		// NOTE: Updating details is complex, for now we only update header info
		// In a real implementation, we should delete old details and insert new ones
		err = u.reportCardRepo.Update(ctx, rc)
	} else {
		err = u.reportCardRepo.Create(ctx, rc)
	}
	
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func (u *reportCardUseCase) GetByStudent(ctx context.Context, studentID, semesterID uuid.UUID) (*entity.ReportCard, error) {
	return u.reportCardRepo.GetByStudentAndSemester(ctx, studentID, semesterID)
}

func (u *reportCardUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportCard, error) {
	return u.reportCardRepo.GetByID(ctx, id)
}

func (u *reportCardUseCase) Publish(ctx context.Context, id uuid.UUID) error {
	rc, err := u.reportCardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rc == nil {
		return errors.New("report card not found")
	}
	
	rc.Status = entity.ReportCardStatusPublished
	now := time.Now()
	rc.PublishedAt = &now
	
	return u.reportCardRepo.Update(ctx, rc)
}
