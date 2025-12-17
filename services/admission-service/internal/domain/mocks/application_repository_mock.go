package mocks

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	entity "github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockApplicationRepository is a mock of ApplicationRepository interface.
type MockApplicationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationRepositoryMockRecorder
}

// MockApplicationRepositoryMockRecorder is the mock recorder for MockApplicationRepository.
type MockApplicationRepositoryMockRecorder struct {
	mock *MockApplicationRepository
}

// NewMockApplicationRepository creates a new mock instance.
func NewMockApplicationRepository(ctrl *gomock.Controller) *MockApplicationRepository {
	mock := &MockApplicationRepository{ctrl: ctrl}
	mock.recorder = &MockApplicationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplicationRepository) EXPECT() *MockApplicationRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockApplicationRepository) Create(ctx context.Context, application *entity.Application) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, application)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockApplicationRepositoryMockRecorder) Create(ctx, application interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockApplicationRepository)(nil).Create), ctx, application)
}

// Delete mocks base method.
func (m *MockApplicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockApplicationRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockApplicationRepository)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockApplicationRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockApplicationRepository)(nil).GetByID), ctx, id)
}

// GetByRegistrationNumber mocks base method.
func (m *MockApplicationRepository) GetByRegistrationNumber(ctx context.Context, regNum string) (*entity.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRegistrationNumber", ctx, regNum)
	ret0, _ := ret[0].(*entity.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRegistrationNumber indicates an expected call of GetByRegistrationNumber.
func (mr *MockApplicationRepositoryMockRecorder) GetByRegistrationNumber(ctx, regNum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRegistrationNumber", reflect.TypeOf((*MockApplicationRepository)(nil).GetByRegistrationNumber), ctx, regNum)
}

// List mocks base method.
func (m *MockApplicationRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].([]*entity.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockApplicationRepositoryMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockApplicationRepository)(nil).List), ctx, filter)
}

// Update mocks base method.
func (m *MockApplicationRepository) Update(ctx context.Context, application *entity.Application) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, application)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockApplicationRepositoryMockRecorder) Update(ctx, application interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockApplicationRepository)(nil).Update), ctx, application)
}
