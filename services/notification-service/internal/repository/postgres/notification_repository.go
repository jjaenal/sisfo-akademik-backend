package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/repository"
)

type notificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	query := `
		INSERT INTO notifications (id, template_id, channel, recipient, subject, body, status, error_message, created_at, sent_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		notification.ID,
		notification.TemplateID,
		notification.Channel,
		notification.Recipient,
		notification.Subject,
		notification.Body,
		notification.Status,
		notification.ErrorMessage,
		notification.CreatedAt,
		notification.SentAt,
	)
	return err
}

func (r *notificationRepository) Update(ctx context.Context, notification *entity.Notification) error {
	query := `
		UPDATE notifications
		SET status = $1, error_message = $2, sent_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query,
		notification.Status,
		notification.ErrorMessage,
		notification.SentAt,
		notification.ID,
	)
	return err
}

func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	query := `
		SELECT id, template_id, channel, recipient, subject, body, status, error_message, created_at, sent_at
		FROM notifications
		WHERE id = $1
	`
	var n entity.Notification
	err := r.db.QueryRow(ctx, query, id).Scan(
		&n.ID,
		&n.TemplateID,
		&n.Channel,
		&n.Recipient,
		&n.Subject,
		&n.Body,
		&n.Status,
		&n.ErrorMessage,
		&n.CreatedAt,
		&n.SentAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepository) ListByRecipient(ctx context.Context, recipient string) ([]*entity.Notification, error) {
	query := `
		SELECT id, template_id, channel, recipient, subject, body, status, error_message, created_at, sent_at
		FROM notifications
		WHERE recipient = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, recipient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*entity.Notification
	for rows.Next() {
		var n entity.Notification
		if err := rows.Scan(
			&n.ID,
			&n.TemplateID,
			&n.Channel,
			&n.Recipient,
			&n.Subject,
			&n.Body,
			&n.Status,
			&n.ErrorMessage,
			&n.CreatedAt,
			&n.SentAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}
