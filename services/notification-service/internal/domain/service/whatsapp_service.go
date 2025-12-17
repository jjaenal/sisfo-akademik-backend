package service

type WhatsAppService interface {
	SendWhatsApp(to, message string) error
}
