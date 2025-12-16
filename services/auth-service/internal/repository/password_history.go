package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasswordHistory struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	PasswordHash string
	CreatedAt    time.Time
}

type PasswordHistoryRepo struct {
	db *pgxpool.Pool
}

func NewPasswordHistoryRepo(db *pgxpool.Pool) *PasswordHistoryRepo {
	return &PasswordHistoryRepo{db: db}
}

func (r *PasswordHistoryRepo) Add(ctx context.Context, userID uuid.UUID, passwordHash string, now time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO password_history (user_id, password_hash, created_at)
		VALUES ($1, $2, $3)
	`, userID, passwordHash, now)
	return err
}

func (r *PasswordHistoryRepo) Recent(ctx context.Context, userID uuid.UUID, limit int) ([]PasswordHistory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, password_hash, created_at
		FROM password_history
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PasswordHistory
	for rows.Next() {
		var ph PasswordHistory
		if err := rows.Scan(&ph.ID, &ph.UserID, &ph.PasswordHash, &ph.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, ph)
	}
	return out, nil
}

func (r *PasswordHistoryRepo) Prune(ctx context.Context, userID uuid.UUID, keep int) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM password_history
		WHERE user_id=$1
		AND id IN (
			SELECT id FROM password_history
			WHERE user_id=$1
			ORDER BY created_at DESC
			OFFSET $2
		)
	`, userID, keep)
	return err
}
