package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTemplateUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTemplateRepository(ctrl)
	u := usecase.NewTemplateUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		template := &entity.ReportCardTemplate{
			Name:     "Template 1",
			TenantID: "tenant-1",
		}

		mockRepo.EXPECT().GetByTenantID(gomock.Any(), template.TenantID).Return(nil, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tmpl *entity.ReportCardTemplate) error {
			assert.NotEmpty(t, tmpl.ID)
			assert.NotZero(t, tmpl.CreatedAt)
			assert.NotZero(t, tmpl.UpdatedAt)
			assert.True(t, tmpl.IsDefault) // First template should be default
			return nil
		})

		err := u.Create(context.Background(), template)
		assert.NoError(t, err)
	})

	t.Run("success not default", func(t *testing.T) {
		template := &entity.ReportCardTemplate{
			Name:     "Template 2",
			TenantID: "tenant-1",
		}
		existing := []*entity.ReportCardTemplate{{ID: uuid.New()}}

		mockRepo.EXPECT().GetByTenantID(gomock.Any(), template.TenantID).Return(existing, nil)
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		err := u.Create(context.Background(), template)
		assert.NoError(t, err)
		assert.False(t, template.IsDefault)
	})

	t.Run("empty name", func(t *testing.T) {
		template := &entity.ReportCardTemplate{}
		err := u.Create(context.Background(), template)
		assert.Error(t, err)
		assert.Equal(t, "template name is required", err.Error())
	})
}

func TestTemplateUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTemplateRepository(ctrl)
	u := usecase.NewTemplateUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		expected := &entity.ReportCardTemplate{ID: id}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestTemplateUseCase_GetByTenantID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTemplateRepository(ctrl)
	u := usecase.NewTemplateUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		tenantID := "tenant-1"
		expected := []*entity.ReportCardTemplate{{TenantID: tenantID}}
		mockRepo.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(expected, nil)

		res, err := u.GetByTenantID(context.Background(), tenantID)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestTemplateUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTemplateRepository(ctrl)
	u := usecase.NewTemplateUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		template := &entity.ReportCardTemplate{
			ID:   id,
			Name: "Updated Name",
		}
		existing := &entity.ReportCardTemplate{
			ID:   id,
			Name: "Old Name",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
		mockRepo.EXPECT().Update(gomock.Any(), existing).DoAndReturn(func(ctx context.Context, tmpl *entity.ReportCardTemplate) error {
			assert.Equal(t, "Updated Name", tmpl.Name)
			return nil
		})

		err := u.Update(context.Background(), template)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		id := uuid.New()
		template := &entity.ReportCardTemplate{ID: id}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		err := u.Update(context.Background(), template)
		assert.Error(t, err)
		assert.Equal(t, "template not found", err.Error())
	})
}

func TestTemplateUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTemplateRepository(ctrl)
	u := usecase.NewTemplateUseCase(mockRepo)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})
}
