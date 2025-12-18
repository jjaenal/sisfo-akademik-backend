package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestReportCardRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewReportCardRepository(db)

	ctx := context.Background()

	// 1. Create ReportCard
	rcID := uuid.New()
	tenantID := "tenant-1"
	generatedAt := time.Now().Truncate(time.Second)
	
	rc := &entity.ReportCard{
		ID:           rcID,
		TenantID:     tenantID,
		StudentID:    uuid.New(),
		ClassID:      uuid.New(),
		SemesterID:   uuid.New(),
		Status:       entity.ReportCardStatusDraft,
		GPA:          3.5,
		TotalCredits: 20,
		Attendance:   10,
		Comments:     "Excellent performance",
		PDFUrl:       "http://example.com/report.pdf",
		GeneratedAt:  &generatedAt,
		Details: []entity.ReportCardDetail{
			{
				ID:          uuid.New(),
				SubjectID:   uuid.New(),
				SubjectName: "Mathematics",
				Credit:      4,
				FinalScore:  85.5,
				GradeLetter: "A",
				TeacherID:   uuid.New(),
				Comments:    "Great job",
			},
		},
	}

	// Test Create
	err := repo.Create(ctx, rc)
	assert.NoError(t, err)

	// Test GetByID
	fetched, err := repo.GetByID(ctx, rcID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, rc.ID, fetched.ID)
	assert.Equal(t, rc.TenantID, fetched.TenantID)
	assert.Equal(t, rc.GPA, fetched.GPA)
	assert.Equal(t, rc.PDFUrl, fetched.PDFUrl)
	assert.Len(t, fetched.Details, 1)
	assert.Equal(t, rc.Details[0].SubjectName, fetched.Details[0].SubjectName)
}
