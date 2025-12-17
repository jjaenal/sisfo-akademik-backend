package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/repository"
)

type notificationTemplateRepository struct {
	db *pgxpool.Pool
}

func NewNotificationTemplateRepository(db *pgxpool.Pool) repository.NotificationTemplateRepository {
	return &notificationTemplateRepository{db: db}
}

func (r *notificationTemplateRepository) Create(ctx context.Context, template *entity.NotificationTemplate) error {
	query := `
		INSERT INTO notification_templates (id, name, channel, subject_template, body_template, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		template.ID,
		template.Name,
		template.Channel,
		template.SubjectTemplate,
		template.BodyTemplate,
		template.IsActive,
		template.CreatedAt,
		template.UpdatedAt,
	)
	return err
}

func (r *notificationTemplateRepository) Update(ctx context.Context, template *entity.NotificationTemplate) error {
	query := `
		UPDATE notification_templates
		SET name = $1, channel = $2, subject_template = $3, body_template = $4, is_active = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.Exec(ctx, query,
		template.Name,
		template.Channel,
		template.SubjectTemplate,
		template.BodyTemplate,
		template.IsActive,
		template.UpdatedAt,
		template.ID,
	)
	return err
}

func (r *notificationTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error) {
	query := `
		SELECT id, name, channel, subject_template, body_template, is_active, created_at, updated_at
		FROM notification_templates
		WHERE id = $1
	`
	var template entity.NotificationTemplate
	err := r.db.QueryRow(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Channel,
		&template.SubjectTemplate,
		&template.BodyTemplate,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

func (r *notificationTemplateRepository) GetByName(ctx context.Context, name string) (*entity.NotificationTemplate, error) {
	query := `
		SELECT id, name, channel, subject_template, body_template, is_active, created_at, updated_at
		FROM notification_templates
		WHERE name = $1
	`
	var template entity.NotificationTemplate
	err := r.db.QueryRow(ctx, query, name).Scan(
		&template.ID,
		&template.Name,
		&template.Channel,
		&template.SubjectTemplate,
		&template.BodyTemplate,
		&template.IsActive,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

func (r *notificationTemplateRepository) List(ctx context.Context) ([]*entity.NotificationTemplate, error) {
	query := `
		SELECT id, name, channel, subject_template, body_template, is_active, created_at, updated_at
		FROM notification_templates
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*entity.NotificationTemplate
	for rows.Next() {
		var t entity.NotificationTemplate
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Channel,
			&t.SubjectTemplate,
			&t.BodyTemplate,
			&t.IsActive,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, &t)
	}
	return templates, nil
}

func (r *notificationTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notification_templates WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
