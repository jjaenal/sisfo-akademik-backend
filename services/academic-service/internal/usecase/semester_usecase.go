package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
)

type semesterUseCase struct {
	repo           repository.SemesterRepository
	contextTimeout time.Duration
}

var _ domainUseCase.SemesterUseCase = (*semesterUseCase)(nil)

func NewSemesterUseCase(repo repository.SemesterRepository, timeout time.Duration) domainUseCase.SemesterUseCase {
	return &semesterUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *semesterUseCase) Create(ctx context.Context, semester *entity.Semester) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := semester.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	if semester.IsActive {
		if err := u.deactivateOthers(ctx, semester.AcademicYearID, uuid.Nil); err != nil {
			return err
		}
	}

	return u.repo.Create(ctx, semester)
}

func (u *semesterUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Semester, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *semesterUseCase) List(ctx context.Context, tenantID string) ([]entity.Semester, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID)
}

func (u *semesterUseCase) ListByAcademicYear(ctx context.Context, academicYearID uuid.UUID) ([]entity.Semester, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.ListByAcademicYear(ctx, academicYearID)
}

func (u *semesterUseCase) Update(ctx context.Context, semester *entity.Semester) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := semester.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	if semester.IsActive {
		if err := u.deactivateOthers(ctx, semester.AcademicYearID, semester.ID); err != nil {
			return err
		}
	}

	return u.repo.Update(ctx, semester)
}

func (u *semesterUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}

func (u *semesterUseCase) SetActive(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	semester, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if semester == nil {
		return errors.New("semester not found")
	}

	if err := u.deactivateOthers(ctx, semester.AcademicYearID, id); err != nil {
		return err
	}

	semester.IsActive = true
	return u.repo.Update(ctx, semester)
}

func (u *semesterUseCase) deactivateOthers(ctx context.Context, academicYearID uuid.UUID, excludeID uuid.UUID) error {
	semesters, err := u.repo.ListByAcademicYear(ctx, academicYearID)
	if err != nil {
		return err
	}

	for _, s := range semesters {
		if s.IsActive && s.ID != excludeID {
			s.IsActive = false
			if err := u.repo.Update(ctx, &s); err != nil {
				return err
			}
		}
	}
	return nil
}
