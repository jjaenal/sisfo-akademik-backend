package testutil

import "testing"

func TestStartPostgres(t *testing.T) {
	dsn, terminate := StartPostgres(t)
	if dsn == "" {
		t.Skip("container not available")
	}
	defer terminate()
	if want := "postgres://"; dsn[:len(want)] != want {
		t.Fatalf("unexpected dsn: %s", dsn)
	}
}

