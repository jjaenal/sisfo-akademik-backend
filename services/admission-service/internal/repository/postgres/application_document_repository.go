package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/repository"
)

type applicationDocumentRepository struct {
	db *pgxpool.Pool
}

func NewApplicationDocumentRepository(db *pgxpool.Pool) repository.ApplicationDocumentRepository {
	return &applicationDocumentRepository{db: db}
}

func (r *applicationDocumentRepository) Create(ctx context.Context, doc *entity.ApplicationDocument) error {
	if doc.ID == uuid.Nil {
		doc.ID = uuid.New()
	}
	now := time.Now()
	if doc.UploadedAt.IsZero() {
		doc.UploadedAt = now
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = now
	}
	if doc.UpdatedAt.IsZero() {
		doc.UpdatedAt = now
	}

	query := `
		INSERT INTO application_documents (
			id, application_id, document_type, file_url, file_name, file_size,
			uploaded_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		doc.ID, doc.ApplicationID, doc.DocumentType, doc.FileURL, doc.FileName, doc.FileSize,
		doc.UploadedAt, doc.CreatedAt, doc.UpdatedAt,
	)
	return err
}

func (r *applicationDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ApplicationDocument, error) {
	query := `
		SELECT 
			id, application_id, document_type, file_url, file_name, file_size,
			uploaded_at, created_at, updated_at
		FROM application_documents
		WHERE id = $1 AND deleted_at IS NULL
	`
	var doc entity.ApplicationDocument
	err := r.db.QueryRow(ctx, query, id).Scan(
		&doc.ID, &doc.ApplicationID, &doc.DocumentType, &doc.FileURL, &doc.FileName, &doc.FileSize,
		&doc.UploadedAt, &doc.CreatedAt, &doc.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (r *applicationDocumentRepository) GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*entity.ApplicationDocument, error) {
	query := `
		SELECT 
			id, application_id, document_type, file_url, file_name, file_size,
			uploaded_at, created_at, updated_at
		FROM application_documents
		WHERE application_id = $1 AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []*entity.ApplicationDocument
	for rows.Next() {
		var doc entity.ApplicationDocument
		if err := rows.Scan(
			&doc.ID, &doc.ApplicationID, &doc.DocumentType, &doc.FileURL, &doc.FileName, &doc.FileSize,
			&doc.UploadedAt, &doc.CreatedAt, &doc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}
	return docs, nil
}

func (r *applicationDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE application_documents SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
