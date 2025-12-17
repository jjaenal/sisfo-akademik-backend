package event

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
)

type Consumer struct {
	rabbitClient *rabbit.Client
	studentUC    usecase.StudentUseCase
}

func NewConsumer(rabbitClient *rabbit.Client, studentUC usecase.StudentUseCase) *Consumer {
	return &Consumer{
		rabbitClient: rabbitClient,
		studentUC:    studentUC,
	}
}

func (c *Consumer) Start() {
	// Ensure the queue exists and bind it
	// We use "academic-service-student-registration" queue
	msgs, err := c.rabbitClient.Consume("sisfo.events", "academic-service-student-registration", []string{"admission.student.registered"})
	if err != nil {
		log.Printf("Failed to consume student registration: %v", err)
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

	// Map event to student entity
	now := time.Now()
	name := event.FirstName + " " + event.LastName
	status := "active"
	student := &entity.Student{
		ID:            uuid.New(),
		TenantID:      event.TenantID,
		Name:          name,
		Email:         event.Email,
		Phone:         event.PhoneNumber,
		NIS:           event.RegistrationNumber, // Using Registration Number as NIS
		Status:        status,
		AdmissionDate: &now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := c.studentUC.Create(ctx, student); err != nil {
		log.Printf("Failed to create student record for %s: %v", event.Email, err)
		// Return nil to avoid retry loop if it's a validation error or duplicate
		return nil
	}

	log.Printf("Successfully created student record for %s", event.Email)
	return nil
}
