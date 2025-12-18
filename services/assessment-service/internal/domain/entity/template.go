package entity

import (
	"time"

	"github.com/google/uuid"
)

type TemplateConfig struct {
	HeaderText     string `json:"header_text"`
	LogoURL        string `json:"logo_url"`
	PrimaryColor   string `json:"primary_color"`   // Hex code
	SecondaryColor string `json:"secondary_color"` // Hex code
	FooterText     string `json:"footer_text"`
}

type ReportCardTemplate struct {
	ID        uuid.UUID      `json:"id"`
	TenantID  string         `json:"tenant_id"`
	Name      string         `json:"name"`
	Config    TemplateConfig `json:"config"`
	IsDefault bool           `json:"is_default"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
