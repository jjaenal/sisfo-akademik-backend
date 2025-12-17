package mocks

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	entity "github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockGradeRepository is a mock of GradeRepository interface.
type MockGradeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGradeRepositoryMockRecorder
}

// MockGradeRepositoryMockRecorder is the mock recorder for MockGradeRepository.
type MockGradeRepositoryMockRecorder struct {
	mock *MockGradeRepository
}

// NewMockGradeRepository creates a new mock instance.
func NewMockGradeRepository(ctrl *gomock.Controller) *MockGradeRepository {
	mock := &MockGradeRepository{ctrl: ctrl}
	mock.recorder = &MockGradeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGradeRepository) EXPECT() *MockGradeRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockGradeRepository) Create(ctx context.Context, grade *entity.Grade) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, grade)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockGradeRepositoryMockRecorder) Create(ctx, grade interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockGradeRepository)(nil).Create), ctx, grade)
}

// Delete mocks base method.
func (m *MockGradeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockGradeRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockGradeRepository)(nil).Delete), ctx, id)
}

// GetByID mocks base method.
func (m *MockGradeRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Grade, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.Grade)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockGradeRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockGradeRepository)(nil).GetByID), ctx, id)
}

// GetByStudentAndAssessment mocks base method.
func (m *MockGradeRepository) GetByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (*entity.Grade, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByStudentAndAssessment", ctx, studentID, assessmentID)
	ret0, _ := ret[0].(*entity.Grade)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByStudentAndAssessment indicates an expected call of GetByStudentAndAssessment.
func (mr *MockGradeRepositoryMockRecorder) GetByStudentAndAssessment(ctx, studentID, assessmentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByStudentAndAssessment", reflect.TypeOf((*MockGradeRepository)(nil).GetByStudentAndAssessment), ctx, studentID, assessmentID)
}

// List mocks base method.
func (m *MockGradeRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.Grade, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].([]*entity.Grade)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockGradeRepositoryMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockGradeRepository)(nil).List), ctx, filter)
}

// Update mocks base method.
func (m *MockGradeRepository) Update(ctx context.Context, grade *entity.Grade) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, grade)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockGradeRepositoryMockRecorder) Update(ctx, grade interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGradeRepository)(nil).Update), ctx, grade)
}
