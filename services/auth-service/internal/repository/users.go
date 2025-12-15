package repository

import (
	"context"
	"errors"
	"time"

	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	TenantID     string
	Email        string
	PasswordHash string
	IsActive     bool
}

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{db: db}
}

var ErrDuplicate = errors.New("duplicate")
func ErrValidation(msg string) error { return errors.New(msg) }

func (r *UsersRepo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, password_hash, is_active
		FROM users
		WHERE id=$1 AND deleted_at IS NULL
		LIMIT 1
	`, id)
	var u User
	if err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.IsActive); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepo) FindByEmail(ctx context.Context, tenantID, email string) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, password_hash, is_active
		FROM users
		WHERE tenant_id=$1 AND email=$2 AND deleted_at IS NULL
		LIMIT 1
	`, tenantID, email)
	var u User
	if err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.IsActive); err != nil {
		return nil, err
	}
	return &u, nil
}

type CreateUserParams struct {
	TenantID string
	Email    string
	Password string
}

func (r *UsersRepo) Create(ctx context.Context, p CreateUserParams) (*User, error) {
	if p.TenantID == "" || p.Email == "" || p.Password == "" {
		return nil, ErrValidation("missing fields")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), 12)
	if err != nil {
		return nil, err
	}
	var u User
	err = r.db.QueryRow(ctx, `
		INSERT INTO users (tenant_id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, tenant_id, email, password_hash, is_active
	`, p.TenantID, p.Email, string(hash)).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.IsActive)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepo) List(ctx context.Context, tenantID string, limit, offset int) ([]User, int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, tenant_id, email, password_hash, is_active
		FROM users
		WHERE tenant_id=$1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.IsActive); err != nil {
			return nil, 0, err
		}
		out = append(out, u)
	}
	// total count
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(1) FROM users WHERE tenant_id=$1 AND deleted_at IS NULL`, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

type UpdateUserParams struct {
	Email    *string
	Password *string
	IsActive *bool
	Now      time.Time
}

func (r *UsersRepo) Update(ctx context.Context, id uuid.UUID, p UpdateUserParams) (*User, error) {
	// build dynamic set
	var setCols []string
	var args []any
	argi := 1
	if p.Email != nil {
		setCols = append(setCols, "email=$"+strconv.Itoa(argi))
		args = append(args, *p.Email)
		argi++
	}
	if p.Password != nil && *p.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*p.Password), 12)
		if err != nil {
			return nil, err
		}
		setCols = append(setCols, "password_hash=$"+strconv.Itoa(argi))
		args = append(args, string(hash))
		argi++
	}
	if p.IsActive != nil {
		setCols = append(setCols, "is_active=$"+strconv.Itoa(argi))
		args = append(args, *p.IsActive)
		argi++
	}
	setCols = append(setCols, "updated_at=$"+strconv.Itoa(argi))
	args = append(args, p.Now)
	argi++
	if len(setCols) == 0 {
		// nothing to update, return current
		return r.FindByID(ctx, id)
	}
	query := "UPDATE users SET " + joinComma(setCols) + " WHERE id=$" + strconv.Itoa(argi) + " AND deleted_at IS NULL RETURNING id, tenant_id, email, password_hash, is_active"
	args = append(args, id)
	var u User
	if err := r.db.QueryRow(ctx, query, args...).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.IsActive); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func joinComma(s []string) string {
	if len(s) == 0 {
		return ""
	}
	out := s[0]
	for i := 1; i < len(s); i++ {
		out += ", " + s[i]
	}
	return out
}
