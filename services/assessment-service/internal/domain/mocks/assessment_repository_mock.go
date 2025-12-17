package mocks

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	entity "github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockAssessmentRepository is a mock of AssessmentRepository interface.
type MockAssessmentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAssessmentRepositoryMockRecorder
}

// MockAssessmentRepositoryMockRecorder is the mock recorder for MockAssessmentRepository.
type MockAssessmentRepositoryMockRecorder struct {
	mock *MockAssessmentRepository
}

// NewMockAssessmentRepository creates a new mock instance.
func NewMockAssessmentRepository(ctrl *gomock.Controller) *MockAssessmentRepository {
	mock := &MockAssessmentRepository{ctrl: ctrl}
	mock.recorder = &MockAssessmentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAssessmentRepository) EXPECT() *MockAssessmentRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAssessmentRepository) Create(ctx context.Context, assessment *entity.Assessment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, assessment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAssessmentRepositoryMockRecorder) Create(ctx, assessment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAssessmentRepository)(nil).Create), ctx, assessment)
}

// Delete mocks base method.
func (m *MockAssessmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAssessmentRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAssessmentRepository)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockAssessmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Assessment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.Assessment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAssessmentRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAssessmentRepository)(nil).GetByID), ctx, id)
}

// List mocks base method.
func (m *MockAssessmentRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.Assessment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].([]*entity.Assessment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockAssessmentRepositoryMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAssessmentRepository)(nil).List), ctx, filter)
}

// Update mocks base method.
func (m *MockAssessmentRepository) Update(ctx context.Context, assessment *entity.Assessment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, assessment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAssessmentRepositoryMockRecorder) Update(ctx, assessment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAssessmentRepository)(nil).Update), ctx, assessment)
}
