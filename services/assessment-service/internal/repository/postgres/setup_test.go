package postgres

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

	// Drop all tables to ensure clean state
	tables := []string{"report_card_details", "report_cards", "grades", "assessments", "grade_categories", "report_card_templates"}
	for _, table := range tables {
		_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS "+table+" CASCADE;")
	}

	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Locate migrations folder relative to this file
	// This file is in internal/repository/postgres
	// Migrations are in migrations/
	// So we need to go up 3 levels: ../../../migrations
	migrationsDir := "../../../migrations"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("failed to read migrations dir: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" && len(file.Name()) > 7 && file.Name()[len(file.Name())-7:] == ".up.sql" {
			content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				t.Fatalf("failed to read migration file %s: %v", file.Name(), err)
			}
			_, err = db.Exec(ctx, string(content))
			if err != nil {
				t.Logf("migration %s warning: %v", file.Name(), err)
			}
		}
	}
}
