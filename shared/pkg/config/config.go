package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env                 string
	ServiceName         string
	HTTPPort            int
	PostgresURL         string
	RedisAddr           string
	RabbitURL           string
	JWTAccessSecret     string
	JWTRefreshSecret    string
	JWTAccessTTL        time.Duration
	JWTRefreshTTL       time.Duration
	JWTIssuer           string
	JWTAudience         string
	CORSAllowedOrigins  []string
	RateLimitPerMinute  int
	LockoutThreshold    int
	LockoutTTL          time.Duration
	FailWindowTTL       time.Duration
	AuditRetentionDays  int
	SMTPHost            string
	SMTPPort            int
	SMTPUser            string
	SMTPPassword        string
	SMTPFromEmail       string
	SMTPFromName        string
	WhatsAppAPIURL      string
	WhatsAppAPIKey      string
	AcademicServiceURL  string
}

func Load() (Config, error) {
	v := viper.New()
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	// Also attempt to read from .env file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	_ = v.ReadInConfig()

	v.SetDefault("ENV", "development")
	v.SetDefault("SERVICE_NAME", "service")
	v.SetDefault("HTTP_PORT", 8080)
	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h")
	v.SetDefault("RATE_LIMIT_PER_MINUTE", 60)
	v.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	v.SetDefault("JWT_ISSUER", "sisfo-akademik")
	v.SetDefault("JWT_AUDIENCE", "api")
	v.SetDefault("LOCKOUT_THRESHOLD", 5)
	v.SetDefault("LOCKOUT_TTL", "15m")
	v.SetDefault("FAIL_WINDOW_TTL", "15m")
	v.SetDefault("AUDIT_RETENTION_DAYS", 90)
	v.SetDefault("SMTP_PORT", 587)
	v.SetDefault("ACADEMIC_SERVICE_URL", "http://localhost:9092")

	cfg := Config{
		Env:                v.GetString("ENV"),
		ServiceName:        v.GetString("SERVICE_NAME"),
		HTTPPort:           v.GetInt("HTTP_PORT"),
		PostgresURL:        v.GetString("POSTGRES_URL"),
		RedisAddr:          v.GetString("REDIS_ADDR"),
		RabbitURL:          v.GetString("RABBIT_URL"),
		JWTAccessSecret:    v.GetString("JWT_ACCESS_SECRET"),
		JWTRefreshSecret:   v.GetString("JWT_REFRESH_SECRET"),
		JWTAccessTTL:       mustParseDuration(v.GetString("JWT_ACCESS_TTL")),
		JWTRefreshTTL:      mustParseDuration(v.GetString("JWT_REFRESH_TTL")),
		JWTIssuer:          v.GetString("JWT_ISSUER"),
		JWTAudience:        v.GetString("JWT_AUDIENCE"),
		CORSAllowedOrigins: v.GetStringSlice("CORS_ALLOWED_ORIGINS"),
		RateLimitPerMinute: v.GetInt("RATE_LIMIT_PER_MINUTE"),
		LockoutThreshold:   v.GetInt("LOCKOUT_THRESHOLD"),
		LockoutTTL:         mustParseDuration(v.GetString("LOCKOUT_TTL")),
		FailWindowTTL:      mustParseDuration(v.GetString("FAIL_WINDOW_TTL")),
		AuditRetentionDays: v.GetInt("AUDIT_RETENTION_DAYS"),
		SMTPHost:           v.GetString("SMTP_HOST"),
		SMTPPort:           v.GetInt("SMTP_PORT"),
		SMTPUser:           v.GetString("SMTP_USER"),
		SMTPPassword:       v.GetString("SMTP_PASSWORD"),
		SMTPFromEmail:      v.GetString("SMTP_FROM_EMAIL"),
		SMTPFromName:       v.GetString("SMTP_FROM_NAME"),
		WhatsAppAPIURL:     v.GetString("WHATSAPP_API_URL"),
		WhatsAppAPIKey:     v.GetString("WHATSAPP_API_KEY"),
		AcademicServiceURL: v.GetString("ACADEMIC_SERVICE_URL"),
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func mustParseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	if d == 0 {
		d = 15 * time.Minute
	}
	return d
}

func (c Config) Validate() error {
	if c.JWTAccessSecret == "" || c.JWTRefreshSecret == "" {
		return fmt.Errorf("jwt secrets required")
	}
	if c.JWTIssuer == "" || c.JWTAudience == "" {
		return fmt.Errorf("jwt issuer/audience required")
	}
	if c.PostgresURL == "" || c.RedisAddr == "" || c.RabbitURL == "" {
		return fmt.Errorf("infrastructure urls required")
	}
	return nil
}
