package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGradeCategoryRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewGradeCategoryRepository(db)

	categoryID := uuid.New()
	tenantID := "tenant-1"
	category := &entity.GradeCategory{
		ID:          categoryID,
		TenantID:    tenantID,
		Name:        "Test Category " + uuid.New().String(),
		Description: "Test Description",
		Weight:      20,
	}

	ctx := context.Background()

	// Test Create
	err := repo.Create(ctx, category)
	assert.NoError(t, err)

	// Test GetByID
	fetched, err := repo.GetByID(ctx, categoryID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, category.TenantID, fetched.TenantID)
	assert.Equal(t, category.Name, fetched.Name)
	assert.Equal(t, category.Weight, fetched.Weight)

	// Test Update
	category.Name = "Updated Category"
	category.Weight = 30
	err = repo.Update(ctx, category)
	assert.NoError(t, err)

	fetchedAfterUpdate, err := repo.GetByID(ctx, categoryID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Category", fetchedAfterUpdate.Name)
	assert.Equal(t, 30.0, fetchedAfterUpdate.Weight)

	// Test GetByTenantID
	list, err := repo.GetByTenantID(ctx, tenantID)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
	found := false
	for _, c := range list {
		if c.ID == category.ID {
			found = true
			break
		}
	}
	assert.True(t, found)

	// Test Delete
	err = repo.Delete(ctx, categoryID)
	assert.NoError(t, err)

	fetchedAfterDelete, err := repo.GetByID(ctx, categoryID)
	assert.Error(t, err)
	assert.Nil(t, fetchedAfterDelete)
}
