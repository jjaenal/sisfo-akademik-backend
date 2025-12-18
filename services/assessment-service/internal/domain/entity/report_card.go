package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReportCardStatus string

const (
	ReportCardStatusDraft     ReportCardStatus = "DRAFT"
	ReportCardStatusGenerated ReportCardStatus = "GENERATED"
	ReportCardStatusPublished ReportCardStatus = "PUBLISHED"
)

type ReportCard struct {
	ID                uuid.UUID        `json:"id"`
	TenantID          string           `json:"tenant_id"`
	StudentID         uuid.UUID        `json:"student_id"`
	SemesterID        uuid.UUID        `json:"semester_id"`
	ClassID           uuid.UUID        `json:"class_id"`
	Status            ReportCardStatus `json:"status"`
	GPA               float64          `json:"gpa"`
	TotalCredits      int              `json:"total_credits"`
	Rank              int              `json:"rank"`
	Attendance        int              `json:"attendance"` // Renamed from TotalAttendance to match code
	AttendanceSummary map[string]int   `json:"attendance_summary"`
	Comments          string           `json:"comments"` // Renamed from TeacherComments
	PDFUrl            string           `json:"pdf_url"`
	GeneratedAt       *time.Time       `json:"generated_at,omitempty"`
	PublishedAt       *time.Time       `json:"published_at,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         *time.Time       `json:"deleted_at,omitempty"`

	// Relationships
	Details []ReportCardDetail `json:"details,omitempty"`
}

type ReportCardDetail struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     string     `json:"tenant_id"`
	ReportCardID uuid.UUID  `json:"report_card_id"`
	SubjectID    uuid.UUID  `json:"subject_id"`
	SubjectName  string     `json:"subject_name"`
	Credit       int        `json:"credit"`
	FinalScore   float64    `json:"final_score"`
	GradeLetter  string     `json:"grade_letter"`
	TeacherID    uuid.UUID  `json:"teacher_id"`
	Comments     string     `json:"comments"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
