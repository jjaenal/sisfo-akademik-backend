package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReportCardStatus string

const (
	ReportCardStatusDraft     ReportCardStatus = "draft"
	ReportCardStatusGenerated ReportCardStatus = "generated"
	ReportCardStatusPublished ReportCardStatus = "published"
)

type ReportCard struct {
	ID           uuid.UUID        `json:"id"`
	TenantID     uuid.UUID        `json:"tenant_id"`
	StudentID    uuid.UUID        `json:"student_id"`
	ClassID      uuid.UUID        `json:"class_id"`
	SemesterID   uuid.UUID        `json:"semester_id"`
	Status       ReportCardStatus `json:"status"`
	GPA          float64          `json:"gpa"`
	TotalCredits int              `json:"total_credits"`
	Attendance   string           `json:"attendance"` // JSON string summary
	Comments     string           `json:"comments"`
	PDFUrl       string           `json:"pdf_url"`
	GeneratedAt  *time.Time       `json:"generated_at"`
	PublishedAt  *time.Time       `json:"published_at"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`

	Details []ReportCardDetail `json:"details,omitempty"`
}

type ReportCardDetail struct {
	ID           uuid.UUID `json:"id"`
	ReportCardID uuid.UUID `json:"report_card_id"`
	SubjectID    uuid.UUID `json:"subject_id"`
	SubjectName  string    `json:"subject_name"`
	Credit       int       `json:"credit"`
	FinalScore   float64   `json:"final_score"`
	GradeLetter  string    `json:"grade_letter"`
	Comments     string    `json:"comments"`
}

func (ReportCard) TableName() string {
	return "report_cards"
}

func (ReportCardDetail) TableName() string {
	return "report_card_details"
}
