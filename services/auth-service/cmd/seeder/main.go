package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
)

func main() {
	log.Println("Starting database seeding...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	rolesRepo := repository.NewRolesRepo(db)
	ctx := context.Background()

	// Default Tenant ID
	defaultTenantID := "default"

	// Define roles to seed
	roles := []struct {
		Name     string
		IsSystem bool
	}{
		{"admin", true},
		{"student", true},
		{"teacher", true},
		{"parent", true},
		{"staff", true},
	}

	for _, role := range roles {
		log.Printf("Checking role: %s", role.Name)
		existing, err := rolesRepo.FindRoleByName(ctx, defaultTenantID, role.Name)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Role not found, create it
				log.Printf("Role %s not found, creating...", role.Name)
				_, err = rolesRepo.CreateRole(ctx, defaultTenantID, role.Name, role.IsSystem)
				if err != nil {
					log.Printf("Failed to create role %s: %v", role.Name, err)
				} else {
					log.Printf("Role %s created successfully", role.Name)
				}
				continue
			}
			// Other error
			log.Printf("Error checking role %s: %v", role.Name, err)
			continue
		}

		if existing != nil {
			log.Printf("Role %s already exists", role.Name)
		}
	}

	log.Println("Seeding completed.")
}
