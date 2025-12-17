package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/service"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
)

type httpWhatsAppService struct {
	apiURL string
	apiKey string
	client *http.Client
}

func NewHTTPWhatsAppService(cfg config.Config) service.WhatsAppService {
	return &httpWhatsAppService{
		apiURL: cfg.WhatsAppAPIURL,
		apiKey: cfg.WhatsAppAPIKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *httpWhatsAppService) SendWhatsApp(to, message string) error {
	if s.apiURL == "" {
		// Just log or return nil if not configured to avoid errors in dev
		return fmt.Errorf("whatsapp api url not configured")
	}

	payload := map[string]string{
		"phone":   to,
		"message": message,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", s.apiURL+"/send", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("whatsapp api returned status: %d", resp.StatusCode)
	}

	return nil
}
