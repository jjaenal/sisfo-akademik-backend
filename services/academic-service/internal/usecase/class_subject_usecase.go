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

type classSubjectUseCase struct {
	repo           repository.ClassSubjectRepository
	contextTimeout time.Duration
}

var _ domainUseCase.ClassSubjectUseCase = (*classSubjectUseCase)(nil)

func NewClassSubjectUseCase(repo repository.ClassSubjectRepository, timeout time.Duration) domainUseCase.ClassSubjectUseCase {
	return &classSubjectUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *classSubjectUseCase) Create(ctx context.Context, cs *entity.ClassSubject) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Check if already exists
	existing, err := u.repo.GetByClassAndSubject(ctx, cs.ClassID, cs.SubjectID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("subject already assigned to this class")
	}

	return u.repo.Create(ctx, cs)
}

func (u *classSubjectUseCase) AssignTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	cs, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if cs == nil {
		return errors.New("class subject not found")
	}

	cs.TeacherID = &teacherID
	return u.repo.Update(ctx, cs)
}

func (u *classSubjectUseCase) GetByClassAndSubject(ctx context.Context, classID, subjectID uuid.UUID) (*entity.ClassSubject, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByClassAndSubject(ctx, classID, subjectID)
}

func (u *classSubjectUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.ClassSubject, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *classSubjectUseCase) ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.ClassSubject, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.ListByClass(ctx, classID)
}

func (u *classSubjectUseCase) ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]entity.ClassSubject, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.ListByTeacher(ctx, teacherID)
}

func (u *classSubjectUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
