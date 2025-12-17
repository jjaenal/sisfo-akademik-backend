package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	entity "github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockTeacherAttendanceRepository is a mock of TeacherAttendanceRepository interface.
type MockTeacherAttendanceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTeacherAttendanceRepositoryMockRecorder
}

// MockTeacherAttendanceRepositoryMockRecorder is the mock recorder for MockTeacherAttendanceRepository.
type MockTeacherAttendanceRepositoryMockRecorder struct {
	mock *MockTeacherAttendanceRepository
}

// NewMockTeacherAttendanceRepository creates a new mock instance.
func NewMockTeacherAttendanceRepository(ctrl *gomock.Controller) *MockTeacherAttendanceRepository {
	mock := &MockTeacherAttendanceRepository{ctrl: ctrl}
	mock.recorder = &MockTeacherAttendanceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTeacherAttendanceRepository) EXPECT() *MockTeacherAttendanceRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTeacherAttendanceRepository) Create(ctx context.Context, attendance *entity.TeacherAttendance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, attendance)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTeacherAttendanceRepositoryMockRecorder) Create(ctx, attendance interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTeacherAttendanceRepository)(nil).Create), ctx, attendance)
}

// GetByID mocks base method.
func (m *MockTeacherAttendanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TeacherAttendance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.TeacherAttendance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockTeacherAttendanceRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTeacherAttendanceRepository)(nil).GetByID), ctx, id)
}

// GetByTeacherAndDate mocks base method.
func (m *MockTeacherAttendanceRepository) GetByTeacherAndDate(ctx context.Context, teacherID uuid.UUID, date time.Time) (*entity.TeacherAttendance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTeacherAndDate", ctx, teacherID, date)
	ret0, _ := ret[0].(*entity.TeacherAttendance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTeacherAndDate indicates an expected call of GetByTeacherAndDate.
func (mr *MockTeacherAttendanceRepositoryMockRecorder) GetByTeacherAndDate(ctx, teacherID, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTeacherAndDate", reflect.TypeOf((*MockTeacherAttendanceRepository)(nil).GetByTeacherAndDate), ctx, teacherID, date)
}

// List mocks base method.
func (m *MockTeacherAttendanceRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.TeacherAttendance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].([]*entity.TeacherAttendance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTeacherAttendanceRepositoryMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTeacherAttendanceRepository)(nil).List), ctx, filter)
}

// Update mocks base method.
func (m *MockTeacherAttendanceRepository) Update(ctx context.Context, attendance *entity.TeacherAttendance) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, attendance)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTeacherAttendanceRepositoryMockRecorder) Update(ctx, attendance interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTeacherAttendanceRepository)(nil).Update), ctx, attendance)
}
