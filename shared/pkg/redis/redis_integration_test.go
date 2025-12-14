package redisutil

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRedis_WithRealContainer(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(60 * time.Second),
	}
	rc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Skipf("skip: cannot start redis container: %v", err)
		return
	}
	defer func() { _ = rc.Terminate(ctx) }()

	host, err := rc.Host(ctx)
	if err != nil {
		t.Skipf("skip: cannot get container host: %v", err)
		return
	}
	port, err := rc.MappedPort(ctx, "6379/tcp")
	if err != nil {
		t.Skipf("skip: cannot get container port: %v", err)
		return
	}
	addr := host + ":" + port.Port()

	client := New(addr)
	raw := client.Raw()
	if err := raw.Ping(ctx).Err(); err != nil {
		t.Fatalf("ping failed: %v", err)
	}
	// KV interface behavior on real client
	if err := Set(ctx, raw, "k", "v", time.Minute); err != nil {
		t.Fatalf("set failed: %v", err)
	}
	got, err := Get(ctx, raw, "k")
	if err != nil || got != "v" {
		t.Fatalf("get mismatch: %v %s", err, got)
	}
	if err := Del(ctx, raw, "k"); err != nil {
		t.Fatalf("del failed: %v", err)
	}
	// Counter interface behavior: Incr + Expire
	if _, err := IncrWithTTL(ctx, raw, "ctr", time.Minute); err != nil {
		t.Fatalf("incr with ttl failed: %v", err)
	}
}
