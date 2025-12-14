package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	if err := os.Setenv("APP_JWT_ACCESS_SECRET", "a"); err != nil {
		t.Fatalf("setenv error: %v", err)
	}
	if err := os.Setenv("APP_JWT_REFRESH_SECRET", "b"); err != nil {
		t.Fatalf("setenv error: %v", err)
	}
	if err := os.Setenv("APP_POSTGRES_URL", "postgres://u:p@host/db"); err != nil {
		t.Fatalf("setenv error: %v", err)
	}
	if err := os.Setenv("APP_REDIS_ADDR", "localhost:6379"); err != nil {
		t.Fatalf("setenv error: %v", err)
	}
	if err := os.Setenv("APP_RABBIT_URL", "amqp://u:p@host/"); err != nil {
		t.Fatalf("setenv error: %v", err)
	}
	defer func() {
		_ = os.Unsetenv("APP_JWT_ACCESS_SECRET")
		_ = os.Unsetenv("APP_JWT_REFRESH_SECRET")
		_ = os.Unsetenv("APP_POSTGRES_URL")
		_ = os.Unsetenv("APP_REDIS_ADDR")
		_ = os.Unsetenv("APP_RABBIT_URL")
	}()
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.ServiceName == "" || cfg.HTTPPort == 0 {
		t.Fatalf("defaults not applied")
	}
}

func TestLoadMissingSecrets(t *testing.T) {
	_ = os.Unsetenv("APP_JWT_ACCESS_SECRET")
	_ = os.Unsetenv("APP_JWT_REFRESH_SECRET")
	_ = os.Unsetenv("APP_POSTGRES_URL")
	_ = os.Unsetenv("APP_REDIS_ADDR")
	_ = os.Unsetenv("APP_RABBIT_URL")
	_, err := Load()
	if err == nil {
		t.Fatalf("expected validation error for missing secrets/config")
	}
}

func TestMustParseDuration(t *testing.T) {
	if d := mustParseDuration("1h"); d != time.Hour {
		t.Fatalf("mustParseDuration 1h mismatch")
	}
	if d := mustParseDuration("invalid"); d != 15*time.Minute {
		t.Fatalf("mustParseDuration invalid should default to 15m")
	}
	if d := mustParseDuration("0s"); d != 15*time.Minute {
		t.Fatalf("mustParseDuration 0s should default to 15m")
	}
}

func TestValidateMissingInfra(t *testing.T) {
	cfg := Config{
		JWTAccessSecret:  "a",
		JWTRefreshSecret: "b",
		// Missing PostgresURL, RedisAddr, RabbitURL
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected infra validation error")
	}
}
