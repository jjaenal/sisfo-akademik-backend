package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestTemplateRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewTemplateRepository(db)

	ctx := context.Background()
	tenantID := uuid.New().String()

	templateID := uuid.New()
	template := &entity.ReportCardTemplate{
		ID:        templateID,
		TenantID:  tenantID,
		Name:      "Standard Template",
		IsDefault: true,
		Config: entity.TemplateConfig{
			HeaderText:     "School Header",
			PrimaryColor:   "#000000",
		},
		CreatedAt: time.Now().Truncate(time.Second),
		UpdatedAt: time.Now().Truncate(time.Second),
	}

	// Test Create
	err := repo.Create(ctx, template)
	assert.NoError(t, err)

	// Test GetByID
	fetched, err := repo.GetByID(ctx, templateID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, template.ID, fetched.ID)
	assert.Equal(t, template.Name, fetched.Name)
	assert.Equal(t, template.Config.HeaderText, fetched.Config.HeaderText)

	// Test GetByTenantID
	list, err := repo.GetByTenantID(ctx, tenantID)
	assert.NoError(t, err)
	assert.NotEmpty(t, list)
	assert.Equal(t, 1, len(list))

	// Test GetDefault
	def, err := repo.GetDefault(ctx, tenantID)
	assert.NoError(t, err)
	assert.NotNil(t, def)
	assert.Equal(t, template.ID, def.ID)

	// Test Update
	template.Name = "Updated Template"
	template.Config.HeaderText = "New Header"
	err = repo.Update(ctx, template)
	assert.NoError(t, err)

	fetchedUpdated, err := repo.GetByID(ctx, templateID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Template", fetchedUpdated.Name)
	assert.Equal(t, "New Header", fetchedUpdated.Config.HeaderText)

	// Test Update IsDefault logic
	template2 := &entity.ReportCardTemplate{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      "Second Template",
		IsDefault: false,
		Config:    entity.TemplateConfig{},
		CreatedAt: time.Now().Truncate(time.Second),
		UpdatedAt: time.Now().Truncate(time.Second),
	}
	err = repo.Create(ctx, template2)
	assert.NoError(t, err)

	// Set template2 as default, should unset template
	template2.IsDefault = true
	err = repo.Update(ctx, template2)
	assert.NoError(t, err)

	// Check template1 is no longer default
	fetched1, err := repo.GetByID(ctx, templateID)
	assert.NoError(t, err)
	assert.False(t, fetched1.IsDefault)
	
	// Check template2 is default
	fetched2, err := repo.GetByID(ctx, template2.ID)
	assert.NoError(t, err)
	assert.True(t, fetched2.IsDefault)

	// Test Delete
	err = repo.Delete(ctx, templateID)
	assert.NoError(t, err)

	fetchedDeleted, err := repo.GetByID(ctx, templateID)
	assert.NoError(t, err)
	assert.Nil(t, fetchedDeleted)
}
