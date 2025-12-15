package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
)

type UserRegisterInput struct {
	TenantID string
	Email    string
	Password string
}

type UserUpdateInput struct {
	Email    *string
	Password *string
	IsActive *bool
}

type Users interface {
	Register(ctx context.Context, in UserRegisterInput) (*repository.User, error)
	Get(ctx context.Context, id uuid.UUID) (*repository.User, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]repository.User, int, error)
	Update(ctx context.Context, id uuid.UUID, in UserUpdateInput) (*repository.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type users struct {
	repo *repository.UsersRepo
}

func NewUsers(repo *repository.UsersRepo) Users {
	return &users{repo: repo}
}

func (u *users) Register(ctx context.Context, in UserRegisterInput) (*repository.User, error) {
	email := strings.TrimSpace(in.Email)
	pwd := strings.TrimSpace(in.Password)
	if email == "" || pwd == "" {
		return nil, repository.ErrValidation("email and password required")
	}
	if len(pwd) < 8 {
		return nil, repository.ErrValidation("password must be at least 8 characters")
	}
	return u.repo.Create(ctx, repository.CreateUserParams{
		TenantID: in.TenantID,
		Email:    strings.ToLower(email),
		Password: pwd,
	})
}

func (u *users) Get(ctx context.Context, id uuid.UUID) (*repository.User, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *users) List(ctx context.Context, tenantID string, limit, offset int) ([]repository.User, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return u.repo.List(ctx, tenantID, limit, offset)
}

func (u *users) Update(ctx context.Context, id uuid.UUID, in UserUpdateInput) (*repository.User, error) {
	params := repository.UpdateUserParams{
		Email:    in.Email,
		Password: in.Password,
		IsActive: in.IsActive,
		Now:      time.Now().UTC(),
	}
	if in.Email != nil {
		e := strings.TrimSpace(*in.Email)
		params.Email = &e
	}
	if in.Password != nil {
		p := strings.TrimSpace(*in.Password)
		if p != "" && len(p) < 8 {
			return nil, repository.ErrValidation("password must be at least 8 characters")
		}
		params.Password = &p
	}
	return u.repo.Update(ctx, id, params)
}

func (u *users) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.SoftDelete(ctx, id)
}
