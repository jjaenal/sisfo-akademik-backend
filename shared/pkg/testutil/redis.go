package testutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartRedis(t *testing.T) (addr string, terminate func()) {
	t.Helper()
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Skipf("skip: cannot start redis container: %v", err)
		return "", func() {}
	}
	terminate = func() { _ = c.Terminate(ctx) }
	host, err := c.Host(ctx)
	if err != nil {
		terminate()
		t.Fatalf("failed to get host: %v", err)
	}
	port, err := c.MappedPort(ctx, "6379/tcp")
	if err != nil {
		terminate()
		t.Fatalf("failed to get port: %v", err)
	}
	addr = fmt.Sprintf("%s:%s", host, port.Port())
	return addr, terminate
}

