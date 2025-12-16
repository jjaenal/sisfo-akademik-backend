package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
)

type classSubjectRepository struct {
	db *pgxpool.Pool
}

var _ repository.ClassSubjectRepository = (*classSubjectRepository)(nil)

func NewClassSubjectRepository(db *pgxpool.Pool) repository.ClassSubjectRepository {
	return &classSubjectRepository{db: db}
}

func (r *classSubjectRepository) Create(ctx context.Context, cs *entity.ClassSubject) error {
	query := `
		INSERT INTO class_subjects (
			id, tenant_id, class_id, subject_id, teacher_id,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9
		)
	`
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	now := time.Now()
	if cs.CreatedAt.IsZero() {
		cs.CreatedAt = now
	}
	if cs.UpdatedAt.IsZero() {
		cs.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		cs.ID, cs.TenantID, cs.ClassID, cs.SubjectID, cs.TeacherID,
		cs.CreatedAt, cs.UpdatedAt, cs.CreatedBy, cs.UpdatedBy,
	)
	return err
}

func (r *classSubjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ClassSubject, error) {
	query := `
		SELECT 
			id, tenant_id, class_id, subject_id, teacher_id,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM class_subjects
		WHERE id = $1 AND deleted_at IS NULL
	`
	var cs entity.ClassSubject
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cs.ID, &cs.TenantID, &cs.ClassID, &cs.SubjectID, &cs.TeacherID,
		&cs.CreatedAt, &cs.UpdatedAt, &cs.CreatedBy, &cs.UpdatedBy, &cs.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cs, nil
}

func (r *classSubjectRepository) GetByClassAndSubject(ctx context.Context, classID, subjectID uuid.UUID) (*entity.ClassSubject, error) {
	query := `
		SELECT 
			id, tenant_id, class_id, subject_id, teacher_id,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM class_subjects
		WHERE class_id = $1 AND subject_id = $2 AND deleted_at IS NULL
	`
	var cs entity.ClassSubject
	err := r.db.QueryRow(ctx, query, classID, subjectID).Scan(
		&cs.ID, &cs.TenantID, &cs.ClassID, &cs.SubjectID, &cs.TeacherID,
		&cs.CreatedAt, &cs.UpdatedAt, &cs.CreatedBy, &cs.UpdatedBy, &cs.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cs, nil
}

func (r *classSubjectRepository) ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.ClassSubject, error) {
	query := `
		SELECT 
			cs.id, cs.tenant_id, cs.class_id, cs.subject_id, cs.teacher_id,
			cs.created_at, cs.updated_at, cs.created_by, cs.updated_by, cs.deleted_at,
			s.id, s.name, s.code, s.description,
			t.id, t.name, t.nip
		FROM class_subjects cs
		JOIN subjects s ON cs.subject_id = s.id
		LEFT JOIN teachers t ON cs.teacher_id = t.id
		WHERE cs.class_id = $1 AND cs.deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.ClassSubject
	for rows.Next() {
		var cs entity.ClassSubject
		var s entity.Subject
		var sDesc *string
		var tID *uuid.UUID
		var tName *string
		var tNIP *string

		err := rows.Scan(
			&cs.ID, &cs.TenantID, &cs.ClassID, &cs.SubjectID, &cs.TeacherID,
			&cs.CreatedAt, &cs.UpdatedAt, &cs.CreatedBy, &cs.UpdatedBy, &cs.DeletedAt,
			&s.ID, &s.Name, &s.Code, &sDesc,
			&tID, &tName, &tNIP,
		)
		if err != nil {
			return nil, err
		}

		if sDesc != nil {
			s.Description = *sDesc
		}
		cs.Subject = &s
		if tID != nil {
			teacher := entity.Teacher{
				ID:   *tID,
				Name: *tName,
			}
			if tNIP != nil {
				teacher.NIP = *tNIP
			}
			cs.Teacher = &teacher
		}
		result = append(result, cs)
	}
	return result, nil
}

func (r *classSubjectRepository) ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]entity.ClassSubject, error) {
	query := `
		SELECT 
			cs.id, cs.tenant_id, cs.class_id, cs.subject_id, cs.teacher_id,
			cs.created_at, cs.updated_at, cs.created_by, cs.updated_by, cs.deleted_at,
			s.id, s.name, s.code, s.description,
			c.id, c.name, c.level, c.major
		FROM class_subjects cs
		JOIN subjects s ON cs.subject_id = s.id
		JOIN classes c ON cs.class_id = c.id
		WHERE cs.teacher_id = $1 AND cs.deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.ClassSubject
	for rows.Next() {
		var cs entity.ClassSubject
		var s entity.Subject
		var c entity.Class

		err := rows.Scan(
			&cs.ID, &cs.TenantID, &cs.ClassID, &cs.SubjectID, &cs.TeacherID,
			&cs.CreatedAt, &cs.UpdatedAt, &cs.CreatedBy, &cs.UpdatedBy, &cs.DeletedAt,
			&s.ID, &s.Name, &s.Code, &s.Description,
			&c.ID, &c.Name, &c.Level, &c.Major,
		)
		if err != nil {
			return nil, err
		}

		cs.Subject = &s
		cs.Class = &c
		result = append(result, cs)
	}
	return result, nil
}

func (r *classSubjectRepository) Update(ctx context.Context, cs *entity.ClassSubject) error {
	query := `
		UPDATE class_subjects SET
			class_id = $1, subject_id = $2, teacher_id = $3,
			updated_at = $4, updated_by = $5
		WHERE id = $6 AND deleted_at IS NULL
	`
	cs.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		cs.ClassID, cs.SubjectID, cs.TeacherID,
		cs.UpdatedAt, cs.UpdatedBy, cs.ID,
	)
	return err
}

func (r *classSubjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE class_subjects SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
