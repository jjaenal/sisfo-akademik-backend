package postgres

import (
	"context"
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
			id, student_id, class_id, semester_id, attendance_date, 
			status, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9
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
		attendance.ID, attendance.StudentID, attendance.ClassID, attendance.SemesterID, attendance.AttendanceDate,
		attendance.Status, attendance.Notes, attendance.CreatedAt, attendance.UpdatedAt,
	)
	return err
}

// GetByID retrieves a student attendance record by its ID
func (r *studentAttendanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.StudentAttendance, error) {
	query := `
		SELECT 
			id, student_id, class_id, semester_id, attendance_date, 
			status, notes, created_at, updated_at
		FROM student_attendance 
		WHERE id = $1
	`
	var attendance entity.StudentAttendance
	err := r.db.QueryRow(ctx, query, id).Scan(
		&attendance.ID, &attendance.StudentID, &attendance.ClassID, &attendance.SemesterID, &attendance.AttendanceDate,
		&attendance.Status, &attendance.Notes, &attendance.CreatedAt, &attendance.UpdatedAt,
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
			id, student_id, class_id, semester_id, attendance_date, 
			status, notes, created_at, updated_at
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
			&attendance.ID, &attendance.StudentID, &attendance.ClassID, &attendance.SemesterID, &attendance.AttendanceDate,
			&attendance.Status, &attendance.Notes, &attendance.CreatedAt, &attendance.UpdatedAt,
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
		SET 
			student_id = $2, class_id = $3, semester_id = $4, attendance_date = $5, 
			status = $6, notes = $7, updated_at = $8
		WHERE id = $1
	`
	attendance.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		attendance.ID, attendance.StudentID, attendance.ClassID, attendance.SemesterID, attendance.AttendanceDate,
		attendance.Status, attendance.Notes, attendance.UpdatedAt,
	)
	return err
}
