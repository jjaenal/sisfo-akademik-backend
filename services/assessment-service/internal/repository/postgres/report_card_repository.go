package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
)

type reportCardRepository struct {
	db *pgxpool.Pool
}

func NewReportCardRepository(db *pgxpool.Pool) repository.ReportCardRepository {
	return &reportCardRepository{db: db}
}

func (r *reportCardRepository) Create(ctx context.Context, rc *entity.ReportCard) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if rc.ID == uuid.Nil {
		rc.ID = uuid.New()
	}
	now := time.Now()
	if rc.CreatedAt.IsZero() {
		rc.CreatedAt = now
	}
	if rc.UpdatedAt.IsZero() {
		rc.UpdatedAt = now
	}

	query := `
		INSERT INTO report_cards (
			id, tenant_id, student_id, class_id, semester_id, status, 
			gpa, total_credits, attendance, comments, pdf_url, generated_at, published_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13,
			$14, $15
		)
	`
	_, err = tx.Exec(ctx, query,
		rc.ID, rc.TenantID, rc.StudentID, rc.ClassID, rc.SemesterID, rc.Status,
		rc.GPA, rc.TotalCredits, rc.Attendance, rc.Comments, rc.PDFUrl, rc.GeneratedAt, rc.PublishedAt,
		rc.CreatedAt, rc.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, d := range rc.Details {
		if d.ID == uuid.Nil {
			d.ID = uuid.New()
		}
		d.ReportCardID = rc.ID
		detailQuery := `
			INSERT INTO report_card_details (
				id, tenant_id, report_card_id, subject_id, subject_name, credit, 
				final_score, grade_letter, teacher_id, comments, 
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10,
				$11, $12
			)
		`
		_, err = tx.Exec(ctx, detailQuery,
			d.ID, rc.TenantID, d.ReportCardID, d.SubjectID, d.SubjectName, d.Credit,
			d.FinalScore, d.GradeLetter, d.TeacherID, d.Comments,
			now, now,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *reportCardRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportCard, error) {
	query := `
		SELECT 
			id, tenant_id, student_id, class_id, semester_id, status, 
			gpa, total_credits, attendance, comments, pdf_url, generated_at, published_at,
			created_at, updated_at
		FROM report_cards
		WHERE id = $1 AND deleted_at IS NULL
	`
	var rc entity.ReportCard
	err := r.db.QueryRow(ctx, query, id).Scan(
		&rc.ID, &rc.TenantID, &rc.StudentID, &rc.ClassID, &rc.SemesterID, &rc.Status,
		&rc.GPA, &rc.TotalCredits, &rc.Attendance, &rc.Comments, &rc.PDFUrl, &rc.GeneratedAt, &rc.PublishedAt,
		&rc.CreatedAt, &rc.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Fetch details
	detailsQuery := `
		SELECT 
			id, report_card_id, subject_id, subject_name, credit, 
			final_score, grade_letter, comments
		FROM report_card_details
		WHERE report_card_id = $1 AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, detailsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []entity.ReportCardDetail
	for rows.Next() {
		var d entity.ReportCardDetail
		if err := rows.Scan(
			&d.ID, &d.ReportCardID, &d.SubjectID, &d.SubjectName, &d.Credit,
			&d.FinalScore, &d.GradeLetter, &d.Comments,
		); err != nil {
			return nil, err
		}
		details = append(details, d)
	}
	rc.Details = details

	return &rc, nil
}

func (r *reportCardRepository) GetByStudentAndSemester(ctx context.Context, studentID, semesterID uuid.UUID) (*entity.ReportCard, error) {
	query := `
		SELECT 
			id, tenant_id, student_id, class_id, semester_id, status, 
			gpa, total_credits, attendance, comments, generated_at, published_at,
			created_at, updated_at
		FROM report_cards
		WHERE student_id = $1 AND semester_id = $2 AND deleted_at IS NULL
	`
	var rc entity.ReportCard
	err := r.db.QueryRow(ctx, query, studentID, semesterID).Scan(
		&rc.ID, &rc.TenantID, &rc.StudentID, &rc.ClassID, &rc.SemesterID, &rc.Status,
		&rc.GPA, &rc.TotalCredits, &rc.Attendance, &rc.Comments, &rc.GeneratedAt, &rc.PublishedAt,
		&rc.CreatedAt, &rc.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Fetch details
	detailsQuery := `
		SELECT 
			id, report_card_id, subject_id, subject_name, credit, 
			final_score, grade_letter, comments
		FROM report_card_details
		WHERE report_card_id = $1 AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, detailsQuery, rc.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []entity.ReportCardDetail
	for rows.Next() {
		var d entity.ReportCardDetail
		if err := rows.Scan(
			&d.ID, &d.ReportCardID, &d.SubjectID, &d.SubjectName, &d.Credit,
			&d.FinalScore, &d.GradeLetter, &d.Comments,
		); err != nil {
			return nil, err
		}
		details = append(details, d)
	}
	rc.Details = details

	return &rc, nil
}

func (r *reportCardRepository) Update(ctx context.Context, rc *entity.ReportCard) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// If this report card is being published, verify it exists and is generated
	rc.UpdatedAt = time.Now()

	query := `
		UPDATE report_cards SET
			status = $1, gpa = $2, total_credits = $3, attendance = $4, 
			comments = $5, generated_at = $6, published_at = $7, updated_at = $8
		WHERE id = $9
	`
	_, err = tx.Exec(ctx, query,
		rc.Status, rc.GPA, rc.TotalCredits, rc.Attendance, rc.Comments, 
		rc.GeneratedAt, rc.PublishedAt, rc.UpdatedAt, rc.ID,
	)
	if err != nil {
		return err
	}

	// Delete existing details (simplest way to handle updates for now)
	_, err = tx.Exec(ctx, "DELETE FROM report_card_details WHERE report_card_id = $1", rc.ID)
	if err != nil {
		return err
	}

	// Re-insert details
	for _, d := range rc.Details {
		if d.ID == uuid.Nil {
			d.ID = uuid.New()
		}
		d.ReportCardID = rc.ID
		detailQuery := `
			INSERT INTO report_card_details (
				id, tenant_id, report_card_id, subject_id, subject_name, credit, 
				final_score, grade_letter, teacher_id, comments, 
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10,
				$11, $12
			)
		`
		_, err = tx.Exec(ctx, detailQuery,
			d.ID, rc.TenantID, d.ReportCardID, d.SubjectID, d.SubjectName, d.Credit,
			d.FinalScore, d.GradeLetter, d.TeacherID, d.Comments,
			time.Now(), time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *reportCardRepository) List(ctx context.Context, classID, semesterID uuid.UUID) ([]*entity.ReportCard, error) {
	query := `
		SELECT 
			id, tenant_id, student_id, class_id, semester_id, status, 
			gpa, total_credits, attendance, comments, generated_at, published_at,
			created_at, updated_at
		FROM report_cards
		WHERE class_id = $1 AND semester_id = $2 AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, classID, semesterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reportCards []*entity.ReportCard
	for rows.Next() {
		var rc entity.ReportCard
		if err := rows.Scan(
			&rc.ID, &rc.TenantID, &rc.StudentID, &rc.ClassID, &rc.SemesterID, &rc.Status,
			&rc.GPA, &rc.TotalCredits, &rc.Attendance, &rc.Comments, &rc.GeneratedAt, &rc.PublishedAt,
			&rc.CreatedAt, &rc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reportCards = append(reportCards, &rc)
	}

	return reportCards, nil
}
