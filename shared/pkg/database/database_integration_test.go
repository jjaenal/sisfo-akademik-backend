package database

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestConnect_WithRealPostgres(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		Env:          map[string]string{"POSTGRES_DB": "testdb", "POSTGRES_USER": "test", "POSTGRES_PASSWORD": "test"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Skipf("skip: cannot start postgres container: %v", err)
		return
	}
	defer func() { _ = pg.Terminate(ctx) }()

	host, err := pg.Host(ctx)
	if err != nil {
		t.Skipf("skip: cannot get container host: %v", err)
		return
	}
	port, err := pg.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Skipf("skip: cannot get container port: %v", err)
		return
	}

	dsn := "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	pool, err := Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	// simple query ensures pool usable
	if _, err := pool.Exec(ctx, "SELECT 1"); err != nil {
		t.Fatalf("exec failed: %v", err)
	}
}
