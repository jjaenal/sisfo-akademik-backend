package event

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
)

type Consumer struct {
	rabbitClient *rabbit.Client
	usersUC      usecase.Users
	rolesUC      usecase.Roles
}

func NewConsumer(rabbitClient *rabbit.Client, usersUC usecase.Users, rolesUC usecase.Roles) *Consumer {
	return &Consumer{
		rabbitClient: rabbitClient,
		usersUC:      usersUC,
		rolesUC:      rolesUC,
	}
}

func (c *Consumer) Start() {
	// Ensure the queue exists and bind it
	// We use "auth-service-student-registration" queue
	msgs, err := c.rabbitClient.Consume("sisfo.events", "auth-service-student-registration", []string{"admission.student.registered"})
	if err != nil {
		log.Printf("Failed to subscribe to student registration: %v", err)
		return
	}

	go func() {
		for d := range msgs {
			if err := c.handleStudentRegistration(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()
}

type StudentRegisteredEvent struct {
	TenantID           string    `json:"tenant_id"`
	ApplicationID      string    `json:"application_id"`
	RegistrationNumber string    `json:"registration_number"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Email              string    `json:"email"`
	PhoneNumber        string    `json:"phone_number"`
	Timestamp          time.Time `json:"timestamp"`
}

func (c *Consumer) handleStudentRegistration(body []byte) error {
	var event StudentRegisteredEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Failed to unmarshal student registration event: %v", err)
		return nil // Don't retry malformed messages
	}

	ctx := context.Background()
	// Default password for new students. In production, this should be random and emailed.
	password := "Student@123!"

	user, err := c.usersUC.Register(ctx, usecase.UserRegisterInput{
		TenantID: event.TenantID,
		Email:    event.Email,
		Password: password,
	})
	if err != nil {
		// If user already exists, we might want to just proceed to role assignment or log and skip.
		log.Printf("Failed to register user for student %s: %v", event.Email, err)
		// Return nil to avoid retry loop
		return nil
	}

	// Assign 'student' role
	_, err = c.rolesUC.AssignByName(ctx, event.TenantID, user.ID, "student")
	if err != nil {
		log.Printf("Failed to assign student role to user %s: %v", user.ID, err)
		return nil
	}

	log.Printf("Successfully created user and assigned student role for %s", event.Email)
	return nil
}
