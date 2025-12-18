package email

import (
	"crypto/tls"

	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/service"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"gopkg.in/gomail.v2"
)

type smtpEmailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewSMTPEmailService(cfg config.Config) service.EmailService {
	dialer := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)
	// In development, we might want to skip verification if using a local mock SMTP or self-signed cert
	if cfg.Env == "development" {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402
	}

	return &smtpEmailService{
		dialer: dialer,
		from:   cfg.SMTPFromEmail,
	}
}

func (s *smtpEmailService) SendEmail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}
