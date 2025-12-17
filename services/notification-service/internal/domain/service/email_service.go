package service

type EmailService interface {
	SendEmail(to []string, subject, body string) error
}
