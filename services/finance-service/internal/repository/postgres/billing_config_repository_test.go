package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestBillingConfigRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewBillingConfigRepository(db)

	tenantID := uuid.New()
	
	// Create
	config := &entity.BillingConfig{
		TenantID:  tenantID,
		Name:      "SPP Bulanan",
		Amount:    500000,
		Frequency: entity.BillingFrequencyMonthly,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx := context.Background()
	err := repo.Create(ctx, config)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, config.ID)

	// GetByID
	fetched, err := repo.GetByID(ctx, config.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, config.ID, fetched.ID)
	assert.Equal(t, config.Name, fetched.Name)
	assert.Equal(t, config.Amount, fetched.Amount)

	// Update
	config.Amount = 600000
	err = repo.Update(ctx, config)
	assert.NoError(t, err)

	fetchedAfterUpdate, err := repo.GetByID(ctx, config.ID)
	assert.NoError(t, err)
	assert.Equal(t, 600000.0, fetchedAfterUpdate.Amount)

	// List
	// Add another config
	config2 := &entity.BillingConfig{
		TenantID:  tenantID,
		Name:      "Uang Gedung",
		Amount:    2000000,
		Frequency: entity.BillingFrequencyOnce,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = repo.Create(ctx, config2)
	assert.NoError(t, err)

	list, err := repo.List(ctx, tenantID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}
