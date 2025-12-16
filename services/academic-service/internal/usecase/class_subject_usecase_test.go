package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClassSubjectUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClassSubjectRepository(ctrl)
	u := usecase.NewClassSubjectUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	id := uuid.New()
	classID := uuid.New()
	subjectID := uuid.New()
	teacherID := uuid.New()

	validCS := &entity.ClassSubject{
		ID:        id,
		TenantID:  tenantID,
		ClassID:   classID,
		SubjectID: subjectID,
		TeacherID: &teacherID,
	}

	t.Run("Create", func(t *testing.T) {
		mockRepo.EXPECT().GetByClassAndSubject(gomock.Any(), classID, subjectID).Return(nil, nil)
		mockRepo.EXPECT().Create(gomock.Any(), validCS).Return(nil)
		err := u.Create(context.Background(), validCS)
		assert.NoError(t, err)
	})

	t.Run("Create Conflict", func(t *testing.T) {
		mockRepo.EXPECT().GetByClassAndSubject(gomock.Any(), classID, subjectID).Return(validCS, nil)
		err := u.Create(context.Background(), validCS)
		assert.Error(t, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validCS, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validCS, res)
	})

	t.Run("ListByClass", func(t *testing.T) {
		list := []entity.ClassSubject{*validCS}
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return(list, nil)
		res, err := u.ListByClass(context.Background(), classID)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
	})

	t.Run("ListByTeacher", func(t *testing.T) {
		list := []entity.ClassSubject{*validCS}
		mockRepo.EXPECT().ListByTeacher(gomock.Any(), teacherID).Return(list, nil)
		res, err := u.ListByTeacher(context.Background(), teacherID)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
	})

	t.Run("AssignTeacher", func(t *testing.T) {
		newTeacherID := uuid.New()
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validCS, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		
		err := u.AssignTeacher(context.Background(), id, newTeacherID)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})
}
