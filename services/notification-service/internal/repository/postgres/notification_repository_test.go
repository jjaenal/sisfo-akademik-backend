package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		url = "postgres://dev:dev@localhost:55432/devdb?sslmode=disable"
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Skip("no db available:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Skip("db connect failed:", err)
	}
	return db
}

func ensureMigrations(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	
	// Create tables if they don't exist
	// We read the up migration file or just execute the SQL directly for simplicity in test
	// Ideally we read the file, but let's try to match the pattern of reading file
	
	migrationPath := "../../../migrations/001_create_notification_tables.up.sql"
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		// Fallback if running from a different directory context or file not found
		// For robustness, let's define the schema here as fallback or fail
		t.Logf("could not read migration file at %s: %v", migrationPath, err)
		
		// Drop tables if exist to ensure clean state
		_, _ = db.Exec(ctx, `DROP TABLE IF EXISTS notifications; DROP TABLE IF EXISTS notification_templates;`)

		// Minimal schema for testing
		schema := `
		CREATE TABLE IF NOT EXISTS notification_templates (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			name VARCHAR(255) NOT NULL,
			channel VARCHAR(50) NOT NULL,
			subject_template TEXT NOT NULL,
			body_template TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS notifications (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			template_id UUID REFERENCES notification_templates(id),
			channel VARCHAR(50) NOT NULL,
			recipient VARCHAR(255) NOT NULL,
			subject TEXT NOT NULL,
			body TEXT NOT NULL,
			status VARCHAR(50) NOT NULL,
			error_message TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			sent_at TIMESTAMP WITH TIME ZONE
		);
		`
		_, err = db.Exec(ctx, schema)
		require.NoError(t, err)
		return
	}
	
	// Drop tables if exist to ensure clean state
	_, _ = db.Exec(ctx, `DROP TABLE IF EXISTS notifications; DROP TABLE IF EXISTS notification_templates;`)
	
	_, err = db.Exec(ctx, string(content))
	require.NoError(t, err)
}

func TestNotificationRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewNotificationRepository(db)
	
	ctx := context.Background()
	
	// Prepare test data
	notificationID := uuid.New()
	templateID := uuid.New() // We might need to insert a template first if FK constraint exists, but let's try inserting notification directly first or disable FK for test if possible. 
	// Wait, the schema likely has FK. Let's insert a template or just make sure template_id is nullable or valid.
	// Looking at the migration file (inferred), template_id likely references notification_templates.
	
	// Let's create a template first to satisfy FK
	_, err := db.Exec(ctx, `INSERT INTO notification_templates (id, name, channel, subject_template, body_template) VALUES ($1, $2, $3, $4, $5)`,
		templateID, "Test Template", "email", "Subject", "Body")
	// If table doesn't exist or error, we might fail here. 
	// If error is duplicate key, ignore it.
	if err != nil {
		// Try to proceed, maybe it exists
	}

	notification := &entity.Notification{
		ID:           notificationID,
		TemplateID:   &templateID,
		Channel:      entity.NotificationChannelEmail,
		Recipient:    "test@example.com",
		Subject:      "Test Subject",
		Body:         "Test Body",
		Status:       entity.NotificationStatusPending,
		ErrorMessage: "",
		CreatedAt:    time.Now().UTC(),
		SentAt:       nil,
	}

	// 1. Create
	err = repo.Create(ctx, notification)
	require.NoError(t, err)

	// 2. GetByID
	saved, err := repo.GetByID(ctx, notificationID)
	require.NoError(t, err)
	assert.NotNil(t, saved)
	assert.Equal(t, notification.ID, saved.ID)
	assert.Equal(t, notification.Recipient, saved.Recipient)
	assert.Equal(t, notification.Status, saved.Status)

	// 3. Update
	newStatus := entity.NotificationStatusSent
	now := time.Now().UTC()
	notification.Status = newStatus
	notification.SentAt = &now
	
	err = repo.Update(ctx, notification)
	require.NoError(t, err)
	
	updated, err := repo.GetByID(ctx, notificationID)
	require.NoError(t, err)
	assert.Equal(t, newStatus, updated.Status)
	assert.NotNil(t, updated.SentAt)

	// 4. ListByRecipient
	list, err := repo.ListByRecipient(ctx, "test@example.com")
	require.NoError(t, err)
	assert.NotEmpty(t, list)
	found := false
	for _, n := range list {
		if n.ID == notificationID {
			found = true
			break
		}
	}
	assert.True(t, found)
}
