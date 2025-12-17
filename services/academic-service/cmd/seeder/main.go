package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
)

func main() {
	log.Println("Starting academic service seeding...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	defaultTenantID := "default"

	// Repositories
	subjectRepo := postgres.NewSubjectRepository(db)
	classRepo := postgres.NewClassRepository(db)
	academicYearRepo := postgres.NewAcademicYearRepository(db)

	// Seed Academic Year
	academicYear := &entity.AcademicYear{
		ID:        uuid.New(),
		TenantID:  defaultTenantID,
		Name:      "2025/2026",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(1, 0, 0),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Printf("Seeding Academic Year: %s", academicYear.Name)
	years, err := academicYearRepo.List(ctx, defaultTenantID)
	if err != nil {
		log.Printf("Error listing academic years: %v", err)
	}
	
	yearExists := false
	for _, y := range years {
		if y.Name == academicYear.Name {
			yearExists = true
			academicYear.ID = y.ID
			log.Printf("Academic Year %s already exists", academicYear.Name)
			break
		}
	}

	if !yearExists {
		if err := academicYearRepo.Create(ctx, academicYear); err != nil {
			log.Printf("Failed to create academic year: %v", err)
		} else {
			log.Printf("Academic Year created")
		}
	}

	// Seed Subjects
	subjects := []entity.Subject{
		{
			ID:          uuid.New(),
			TenantID:    defaultTenantID,
			Code:        "MATH101",
			Name:        "Mathematics Grade 10",
			Description: "Basic Math",
			CreditUnits: 3,
			Type:        "compulsory",
		},
		{
			ID:          uuid.New(),
			TenantID:    defaultTenantID,
			Code:        "PHYS101",
			Name:        "Physics Grade 10",
			Description: "Basic Physics",
			CreditUnits: 3,
			Type:        "compulsory",
		},
		{
			ID:          uuid.New(),
			TenantID:    defaultTenantID,
			Code:        "ENG101",
			Name:        "English Grade 10",
			Description: "Basic English",
			CreditUnits: 2,
			Type:        "compulsory",
		},
	}

	existingSubjects, _, err := subjectRepo.List(ctx, defaultTenantID, 100, 0)
	if err != nil {
		log.Printf("Error listing subjects: %v", err)
	}

	for _, subj := range subjects {
		log.Printf("Seeding Subject: %s", subj.Name)
		exists := false
		for _, es := range existingSubjects {
			if es.Code == subj.Code {
				exists = true
				break
			}
		}

		if !exists {
			subj.CreatedAt = time.Now()
			subj.UpdatedAt = time.Now()
			// Need to pass pointer to Create
			s := subj
			if err := subjectRepo.Create(ctx, &s); err != nil {
				log.Printf("Failed to create subject %s: %v", subj.Name, err)
			} else {
				log.Printf("Subject %s created", subj.Name)
			}
		} else {
			log.Printf("Subject %s already exists", subj.Name)
		}
	}

	// Seed Classes
	classes := []entity.Class{
		{
			ID:             uuid.New(),
			TenantID:       defaultTenantID,
			AcademicYearID: &academicYear.ID,
			Name:           "Grade 10 Class A",
			Level:          10,
			Major:          "Science",
			Capacity:       30,
		},
		{
			ID:             uuid.New(),
			TenantID:       defaultTenantID,
			AcademicYearID: &academicYear.ID,
			Name:           "Grade 10 Class B",
			Level:          10,
			Major:          "Science",
			Capacity:       30,
		},
	}

	existingClasses, _, err := classRepo.List(ctx, defaultTenantID, 100, 0)
	if err != nil {
		log.Printf("Error listing classes: %v", err)
	}

	for _, cls := range classes {
		log.Printf("Seeding Class: %s", cls.Name)
		exists := false
		for _, ec := range existingClasses {
			if ec.Name == cls.Name {
				exists = true
				break
			}
		}

		if !exists {
			cls.CreatedAt = time.Now()
			cls.UpdatedAt = time.Now()
			c := cls
			if err := classRepo.Create(ctx, &c); err != nil {
				log.Printf("Failed to create class %s: %v", cls.Name, err)
			} else {
				log.Printf("Class %s created", cls.Name)
			}
		} else {
			log.Printf("Class %s already exists", cls.Name)
		}
	}

	log.Println("Academic service seeding completed.")
}
