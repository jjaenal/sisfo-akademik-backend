package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
)

type scheduleRepository struct {
	db DBPool
}

var _ repository.ScheduleRepository = (*scheduleRepository)(nil)

func NewScheduleRepository(db DBPool) repository.ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(ctx context.Context, s *entity.Schedule) error {
	query := `
		INSERT INTO schedules (
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13
		)
	`
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	now := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		s.ID, s.TenantID, s.ClassID, s.SubjectID, s.TeacherID, s.DayOfWeek,
		s.StartTime, s.EndTime, s.Room, s.CreatedAt, s.UpdatedAt, s.CreatedBy, s.UpdatedBy,
	)
	return err
}

func (r *scheduleRepository) CheckConflicts(ctx context.Context, s *entity.Schedule) ([]entity.Schedule, error) {
	query := `
		SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by
		FROM schedules
		WHERE 
			tenant_id = $1 
			AND deleted_at IS NULL
			AND day_of_week = $2
			AND (
				class_id = $3 OR 
				teacher_id = $4 OR 
				(room = $5 AND room != '')
			)
			AND (
				(start_time < $6 AND end_time > $7) -- Overlap condition
			)
			AND id != $8 -- Exclude self (for updates)
	`

	// Note: s.ID might be uuid.Nil for creation, so passing it is fine (it won't match any existing ID)

	rows, err := r.db.Query(ctx, query,
		s.TenantID,
		s.DayOfWeek,
		s.ClassID,
		s.TeacherID,
		s.Room,
		s.EndTime,   // Overlap logic: Existing Start < New End
		s.StartTime, // Overlap logic: Existing End > New Start
		s.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conflicts []entity.Schedule
	for rows.Next() {
		var cs entity.Schedule
		if err := rows.Scan(
			&cs.ID, &cs.TenantID, &cs.ClassID, &cs.SubjectID, &cs.TeacherID, &cs.DayOfWeek,
			&cs.StartTime, &cs.EndTime, &cs.Room, &cs.CreatedAt, &cs.UpdatedAt, &cs.CreatedBy, &cs.UpdatedBy,
		); err != nil {
			return nil, err
		}
		conflicts = append(conflicts, cs)
	}
	return conflicts, nil
}

func (r *scheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error) {
	query := `
		SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by, deleted_at
		FROM schedules
		WHERE id = $1 AND deleted_at IS NULL
	`
	var s entity.Schedule
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.TenantID, &s.ClassID, &s.SubjectID, &s.TeacherID, &s.DayOfWeek,
		&s.StartTime, &s.EndTime, &s.Room, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *scheduleRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Schedule, int, error) {
	countQuery := `SELECT COUNT(*) FROM schedules WHERE tenant_id = $1 AND deleted_at IS NULL`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by, deleted_at
		FROM schedules
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY day_of_week ASC, start_time ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var schedules []entity.Schedule
	for rows.Next() {
		var s entity.Schedule
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.ClassID, &s.SubjectID, &s.TeacherID, &s.DayOfWeek,
			&s.StartTime, &s.EndTime, &s.Room, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		schedules = append(schedules, s)
	}
	return schedules, total, nil
}

func (r *scheduleRepository) ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.Schedule, error) {
	query := `
		SELECT 
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by, deleted_at
		FROM schedules
		WHERE class_id = $1 AND deleted_at IS NULL
		ORDER BY day_of_week ASC, start_time ASC
	`
	rows, err := r.db.Query(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []entity.Schedule
	for rows.Next() {
		var s entity.Schedule
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.ClassID, &s.SubjectID, &s.TeacherID, &s.DayOfWeek,
			&s.StartTime, &s.EndTime, &s.Room, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *scheduleRepository) Update(ctx context.Context, s *entity.Schedule) error {
	query := `
		UPDATE schedules SET
			class_id = $1, subject_id = $2, teacher_id = $3, day_of_week = $4,
			start_time = $5, end_time = $6, room = $7, updated_at = $8, updated_by = $9
		WHERE id = $10 AND deleted_at IS NULL
	`
	s.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		s.ClassID, s.SubjectID, s.TeacherID, s.DayOfWeek,
		s.StartTime, s.EndTime, s.Room, s.UpdatedAt, s.UpdatedBy, s.ID,
	)
	return err
}

func (r *scheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE schedules SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}

func (r *scheduleRepository) BulkCreate(ctx context.Context, schedules []*entity.Schedule) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		INSERT INTO schedules (
			id, tenant_id, class_id, subject_id, teacher_id, day_of_week,
			start_time, end_time, room, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13
		)
	`

	for _, s := range schedules {
		if s.ID == uuid.Nil {
			s.ID = uuid.New()
		}
		now := time.Now()
		if s.CreatedAt.IsZero() {
			s.CreatedAt = now
		}
		if s.UpdatedAt.IsZero() {
			s.UpdatedAt = now
		}

		_, err := tx.Exec(ctx, query,
			s.ID, s.TenantID, s.ClassID, s.SubjectID, s.TeacherID, s.DayOfWeek,
			s.StartTime, s.EndTime, s.Room, s.CreatedAt, s.UpdatedAt, s.CreatedBy, s.UpdatedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
