package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/repository"
)

type studentAttendanceRepository struct {
	db *pgxpool.Pool
}

// NewStudentAttendanceRepository creates a new instance of StudentAttendanceRepository
func NewStudentAttendanceRepository(db *pgxpool.Pool) repository.StudentAttendanceRepository {
	return &studentAttendanceRepository{db: db}
}

// Create inserts a new student attendance record into the database
func (r *studentAttendanceRepository) Create(ctx context.Context, attendance *entity.StudentAttendance) error {
	query := `
		INSERT INTO student_attendance (
			id, tenant_id, student_id, class_id, semester_id, attendance_date, 
			status, notes, check_in_latitude, check_in_longitude, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, $10, $11, $12
		)
	`
	if attendance.ID == uuid.Nil {
		attendance.ID = uuid.New()
	}
	now := time.Now()
	if attendance.CreatedAt.IsZero() {
		attendance.CreatedAt = now
	}
	if attendance.UpdatedAt.IsZero() {
		attendance.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		attendance.ID, attendance.TenantID, attendance.StudentID, attendance.ClassID, attendance.SemesterID, attendance.AttendanceDate,
		attendance.Status, attendance.Notes, attendance.CheckInLatitude, attendance.CheckInLongitude, attendance.CreatedAt, attendance.UpdatedAt,
	)
	return err
}

// GetByID retrieves a student attendance record by its ID
func (r *studentAttendanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.StudentAttendance, error) {
	query := `
		SELECT 
			id, tenant_id, student_id, class_id, semester_id, attendance_date, 
			status, notes, check_in_latitude, check_in_longitude, created_at, updated_at
		FROM student_attendance 
		WHERE id = $1
	`
	var attendance entity.StudentAttendance
	err := r.db.QueryRow(ctx, query, id).Scan(
		&attendance.ID, &attendance.TenantID, &attendance.StudentID, &attendance.ClassID, &attendance.SemesterID, &attendance.AttendanceDate,
		&attendance.Status, &attendance.Notes, &attendance.CheckInLatitude, &attendance.CheckInLongitude, &attendance.CreatedAt, &attendance.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attendance, nil
}

// GetByClassAndDate retrieves student attendance records for a specific class and date
func (r *studentAttendanceRepository) GetByClassAndDate(ctx context.Context, classID uuid.UUID, date time.Time) ([]*entity.StudentAttendance, error) {
	query := `
		SELECT 
			id, tenant_id, student_id, class_id, semester_id, attendance_date, 
			status, notes, check_in_latitude, check_in_longitude, created_at, updated_at
		FROM student_attendance 
		WHERE class_id = $1 AND attendance_date = $2
	`
	rows, err := r.db.Query(ctx, query, classID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*entity.StudentAttendance
	for rows.Next() {
		var attendance entity.StudentAttendance
		err := rows.Scan(
			&attendance.ID, &attendance.TenantID, &attendance.StudentID, &attendance.ClassID, &attendance.SemesterID, &attendance.AttendanceDate,
			&attendance.Status, &attendance.Notes, &attendance.CheckInLatitude, &attendance.CheckInLongitude, &attendance.CreatedAt, &attendance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, &attendance)
	}
	return attendances, nil
}

// Update updates an existing student attendance record
func (r *studentAttendanceRepository) Update(ctx context.Context, attendance *entity.StudentAttendance) error {
	query := `
		UPDATE student_attendance 
		SET status = $1, notes = $2, updated_at = $3
		WHERE id = $4
	`
	attendance.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query, attendance.Status, attendance.Notes, attendance.UpdatedAt, attendance.ID)
	return err
}

// BulkCreate inserts multiple student attendance records into the database
func (r *studentAttendanceRepository) BulkCreate(ctx context.Context, attendances []*entity.StudentAttendance) error {
	if len(attendances) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO student_attendance (
			id, tenant_id, student_id, class_id, semester_id, attendance_date, 
			status, notes, check_in_latitude, check_in_longitude, created_at, updated_at
		) VALUES 
	`
	
	vals := []interface{}{}
	now := time.Now()
	
	for i, att := range attendances {
		if att.ID == uuid.Nil {
			att.ID = uuid.New()
		}
		if att.CreatedAt.IsZero() {
			att.CreatedAt = now
		}
		if att.UpdatedAt.IsZero() {
			att.UpdatedAt = now
		}
		
		n := i * 12
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),", 
			n+1, n+2, n+3, n+4, n+5, n+6, n+7, n+8, n+9, n+10, n+11, n+12)
		
		vals = append(vals, 
			att.ID, att.TenantID, att.StudentID, att.ClassID, att.SemesterID, att.AttendanceDate,
			att.Status, att.Notes, att.CheckInLatitude, att.CheckInLongitude, att.CreatedAt, att.UpdatedAt,
		)
	}
	
	query = strings.TrimSuffix(query, ",") // Remove trailing comma

	_, err = tx.Exec(ctx, query, vals...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetSummary retrieves attendance summary for a student
func (r *studentAttendanceRepository) GetSummary(ctx context.Context, studentID uuid.UUID, semesterID uuid.UUID) (map[string]int, error) {
	query := `
		SELECT status, COUNT(*) 
		FROM student_attendance 
		WHERE student_id = $1
	`
	args := []interface{}{studentID}
	
	if semesterID != uuid.Nil {
		query += " AND semester_id = $2"
		args = append(args, semesterID)
	}
	
	query += " GROUP BY status"
	
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	summary := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		summary[status] = count
	}
	
	return summary, nil
}
