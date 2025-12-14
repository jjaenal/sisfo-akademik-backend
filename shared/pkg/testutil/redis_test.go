package testutil

import "testing"

func TestStartRedis(t *testing.T) {
	addr, terminate := StartRedis(t)
	if addr == "" {
		t.Skip("container not available")
	}
	defer terminate()
}

