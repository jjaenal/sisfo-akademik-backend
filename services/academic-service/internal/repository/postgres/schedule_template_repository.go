package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
)

type scheduleTemplateRepository struct {
	db DBPool
}

var _ repository.ScheduleTemplateRepository = (*scheduleTemplateRepository)(nil)

func NewScheduleTemplateRepository(db DBPool) repository.ScheduleTemplateRepository {
	return &scheduleTemplateRepository{db: db}
}

func (r *scheduleTemplateRepository) Create(ctx context.Context, t *entity.ScheduleTemplate) error {
	query := `
		INSERT INTO schedule_templates (
			id, tenant_id, name, description, is_active, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		t.ID, t.TenantID, t.Name, t.Description, t.IsActive, t.CreatedAt, t.UpdatedAt, t.CreatedBy, t.UpdatedBy,
	)
	return err
}

func (r *scheduleTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleTemplate, error) {
	query := `
		SELECT id, tenant_id, name, description, is_active, created_at, updated_at, created_by, updated_by
		FROM schedule_templates
		WHERE id = $1 AND deleted_at IS NULL
	`
	var t entity.ScheduleTemplate
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.TenantID, &t.Name, &t.Description, &t.IsActive, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *scheduleTemplateRepository) List(ctx context.Context, tenantID string) ([]entity.ScheduleTemplate, error) {
	query := `
		SELECT id, tenant_id, name, description, is_active, created_at, updated_at, created_by, updated_by
		FROM schedule_templates
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []entity.ScheduleTemplate
	for rows.Next() {
		var t entity.ScheduleTemplate
		if err := rows.Scan(
			&t.ID, &t.TenantID, &t.Name, &t.Description, &t.IsActive, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *scheduleTemplateRepository) Update(ctx context.Context, t *entity.ScheduleTemplate) error {
	query := `
		UPDATE schedule_templates
		SET name = $1, description = $2, is_active = $3, updated_at = $4, updated_by = $5
		WHERE id = $6 AND deleted_at IS NULL
	`
	t.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		t.Name, t.Description, t.IsActive, t.UpdatedAt, t.UpdatedBy, t.ID,
	)
	return err
}

func (r *scheduleTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE schedule_templates
		SET deleted_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *scheduleTemplateRepository) AddItem(ctx context.Context, item *entity.ScheduleTemplateItem) error {
	query := `
		INSERT INTO schedule_template_items (
			id, template_id, subject_id, day_of_week, start_time, end_time, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		item.ID, item.TemplateID, item.SubjectID, item.DayOfWeek, item.StartTime, item.EndTime, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *scheduleTemplateRepository) RemoveItem(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedule_template_items WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *scheduleTemplateRepository) ListItems(ctx context.Context, templateID uuid.UUID) ([]entity.ScheduleTemplateItem, error) {
	query := `
		SELECT id, template_id, subject_id, day_of_week, start_time, end_time, created_at, updated_at
		FROM schedule_template_items
		WHERE template_id = $1
		ORDER BY day_of_week, start_time
	`
	rows, err := r.db.Query(ctx, query, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.ScheduleTemplateItem
	for rows.Next() {
		var i entity.ScheduleTemplateItem
		if err := rows.Scan(
			&i.ID, &i.TemplateID, &i.SubjectID, &i.DayOfWeek, &i.StartTime, &i.EndTime, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
