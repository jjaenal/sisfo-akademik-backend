package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationTemplateRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewNotificationTemplateRepository(db)

	ctx := context.Background()

	// Prepare test data
	templateID := uuid.New()
	template := &entity.NotificationTemplate{
		ID:              templateID,
		Name:            "Test Template",
		Channel:         entity.NotificationChannelEmail,
		SubjectTemplate: "Test Subject",
		BodyTemplate:    "Test Body",
		IsActive:        true,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	// 1. Create
	err := repo.Create(ctx, template)
	require.NoError(t, err)

	// 2. GetByID
	saved, err := repo.GetByID(ctx, templateID)
	require.NoError(t, err)
	assert.NotNil(t, saved)
	assert.Equal(t, template.ID, saved.ID)
	assert.Equal(t, template.Name, saved.Name)

	// 3. GetByName
	savedByName, err := repo.GetByName(ctx, template.Name)
	require.NoError(t, err)
	assert.NotNil(t, savedByName)
	assert.Equal(t, template.ID, savedByName.ID)

	// 4. Update
	newName := "Updated Template"
	template.Name = newName
	template.UpdatedAt = time.Now().UTC()

	err = repo.Update(ctx, template)
	require.NoError(t, err)

	updated, err := repo.GetByID(ctx, templateID)
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)

	// 5. List
	list, err := repo.List(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, list)
	found := false
	for _, t := range list {
		if t.ID == templateID {
			found = true
			break
		}
	}
	assert.True(t, found)

	// 6. Delete
	err = repo.Delete(ctx, templateID)
	require.NoError(t, err)

	deleted, err := repo.GetByID(ctx, templateID)
	require.NoError(t, err)
	assert.Nil(t, deleted)
}
