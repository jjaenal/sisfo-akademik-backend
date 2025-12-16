package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
)

func main() {
	port := os.Getenv("APP_HTTP_PORT")
	if port == "" {
		port = "9097"
	}
	rabbitURL := os.Getenv("APP_RABBIT_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://dev:dev@rabbitmq:5672/"
	}
	rb := rabbit.New(rabbitURL)
	if rb != nil {
		msgs, err := rb.Consume("events", "notification-service.events", []string{"auth.password_reset.requested"})
		if err != nil {
			log.Printf("rabbit consume error: %v", err)
		} else {
			go func() {
				for m := range msgs {
					log.Printf("received event: rk=%s body=%s", m.RoutingKey, string(m.Body))
					var ev struct {
						TenantID string `json:"tenant_id"`
						UserID   string `json:"user_id"`
						Email    string `json:"email"`
						Token    string `json:"token"`
						Type     string `json:"type"`
					}
					if err := json.Unmarshal(m.Body, &ev); err == nil && ev.Type == "password_reset" && ev.Email != "" && ev.Token != "" {
						sendResetEmail(ev.Email, ev.Token)
					}
				}
			}()
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"data": map[string]any{
				"service": "notification-service",
			},
			"meta": map[string]any{
				"timestamp":  time.Now().UTC(),
				"request_id": uuid.NewString(),
			},
		})
	})
	_ = http.ListenAndServe(":"+port, mux)
}

func sendResetEmail(to, token string) {
	host := os.Getenv("APP_SMTP_HOST")
	port := os.Getenv("APP_SMTP_PORT")
	user := os.Getenv("APP_SMTP_USER")
	pass := os.Getenv("APP_SMTP_PASS")
	from := os.Getenv("APP_SMTP_FROM")
	base := os.Getenv("APP_RESET_BASE_URL")
	link := token
	if base != "" {
		link = fmt.Sprintf("%s%s", base, token)
	}
	subject := "Reset Password"
	body := fmt.Sprintf("Gunakan tautan berikut untuk reset password: %s", link)
	if host == "" || port == "" || user == "" || pass == "" || from == "" {
		log.Printf("smtp not configured; would send to=%s subject=%s body=%s", to, subject, body)
		return
	}
	addr := host + ":" + port
	msg := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" + body + "\r\n")
	auth := smtp.PlainAuth("", user, pass, host)
	if err := smtp.SendMail(addr, auth, from, []string{to}, msg); err != nil {
		log.Printf("smtp send error: %v", err)
		return
	}
	log.Printf("smtp sent to=%s", to)
}
