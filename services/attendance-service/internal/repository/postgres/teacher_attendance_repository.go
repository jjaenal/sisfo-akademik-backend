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

type teacherAttendanceRepository struct {
	db *pgxpool.Pool
}

func NewTeacherAttendanceRepository(db *pgxpool.Pool) repository.TeacherAttendanceRepository {
	return &teacherAttendanceRepository{db: db}
}

func (r *teacherAttendanceRepository) Create(ctx context.Context, attendance *entity.TeacherAttendance) error {
	query := `
		INSERT INTO teacher_attendance (
			id, teacher_id, semester_id, attendance_date, 
			check_in_time, check_out_time, status, notes, 
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10
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
		attendance.ID, attendance.TeacherID, attendance.SemesterID, attendance.AttendanceDate,
		attendance.CheckInTime, attendance.CheckOutTime, attendance.Status, attendance.Notes,
		attendance.CreatedAt, attendance.UpdatedAt,
	)
	return err
}

func (r *teacherAttendanceRepository) Update(ctx context.Context, attendance *entity.TeacherAttendance) error {
	query := `
		UPDATE teacher_attendance SET
			teacher_id = $2, semester_id = $3, attendance_date = $4,
			check_in_time = $5, check_out_time = $6, status = $7, notes = $8,
			updated_at = $9
		WHERE id = $1
	`
	attendance.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		attendance.ID, attendance.TeacherID, attendance.SemesterID, attendance.AttendanceDate,
		attendance.CheckInTime, attendance.CheckOutTime, attendance.Status, attendance.Notes,
		attendance.UpdatedAt,
	)
	return err
}

func (r *teacherAttendanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TeacherAttendance, error) {
	query := `
		SELECT 
			id, teacher_id, semester_id, attendance_date, 
			check_in_time, check_out_time, status, notes, 
			created_at, updated_at
		FROM teacher_attendance 
		WHERE id = $1
	`
	var attendance entity.TeacherAttendance
	err := r.db.QueryRow(ctx, query, id).Scan(
		&attendance.ID, &attendance.TeacherID, &attendance.SemesterID, &attendance.AttendanceDate,
		&attendance.CheckInTime, &attendance.CheckOutTime, &attendance.Status, &attendance.Notes,
		&attendance.CreatedAt, &attendance.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attendance, nil
}

func (r *teacherAttendanceRepository) GetByTeacherAndDate(ctx context.Context, teacherID uuid.UUID, date time.Time) (*entity.TeacherAttendance, error) {
	query := `
		SELECT 
			id, teacher_id, semester_id, attendance_date, 
			check_in_time, check_out_time, status, notes, 
			created_at, updated_at
		FROM teacher_attendance 
		WHERE teacher_id = $1 AND attendance_date = $2
	`
	var attendance entity.TeacherAttendance
	err := r.db.QueryRow(ctx, query, teacherID, date).Scan(
		&attendance.ID, &attendance.TeacherID, &attendance.SemesterID, &attendance.AttendanceDate,
		&attendance.CheckInTime, &attendance.CheckOutTime, &attendance.Status, &attendance.Notes,
		&attendance.CreatedAt, &attendance.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attendance, nil
}

func (r *teacherAttendanceRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.TeacherAttendance, error) {
	query := `
		SELECT 
			id, teacher_id, semester_id, attendance_date, 
			check_in_time, check_out_time, status, notes, 
			created_at, updated_at
		FROM teacher_attendance
	`
	
	var conditions []string
	var args []interface{}
	argCount := 1

	if teacherID, ok := filter["teacher_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("teacher_id = $%d", argCount))
		args = append(args, teacherID)
		argCount++
	}

	if semesterID, ok := filter["semester_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("semester_id = $%d", argCount))
		args = append(args, semesterID)
		argCount++
	}
	
	if startDate, ok := filter["start_date"]; ok {
		conditions = append(conditions, fmt.Sprintf("attendance_date >= $%d", argCount))
		args = append(args, startDate)
		argCount++
	}
	
	if endDate, ok := filter["end_date"]; ok {
		conditions = append(conditions, fmt.Sprintf("attendance_date <= $%d", argCount))
		args = append(args, endDate)
		argCount++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	query += " ORDER BY attendance_date DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendances []*entity.TeacherAttendance
	for rows.Next() {
		var attendance entity.TeacherAttendance
		err := rows.Scan(
			&attendance.ID, &attendance.TeacherID, &attendance.SemesterID, &attendance.AttendanceDate,
			&attendance.CheckInTime, &attendance.CheckOutTime, &attendance.Status, &attendance.Notes,
			&attendance.CreatedAt, &attendance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, &attendance)
	}
	return attendances, nil
}
