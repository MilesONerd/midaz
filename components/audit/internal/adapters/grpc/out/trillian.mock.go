// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/LerianStudio/midaz/components/audit/internal/adapters/grpc/out (interfaces: Repository)
//
// Generated by this command:
//
//	mockgen --destination=trillian.mock.go --package=out . Repository
//

// Package out is a generated GoMock package.
package out

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateLog mocks base method.
func (m *MockRepository) CreateLog(arg0 context.Context, arg1 int64, arg2 []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLog", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLog indicates an expected call of CreateLog.
func (mr *MockRepositoryMockRecorder) CreateLog(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLog", reflect.TypeOf((*MockRepository)(nil).CreateLog), arg0, arg1, arg2)
}

// CreateTree mocks base method.
func (m *MockRepository) CreateTree(arg0 context.Context, arg1, arg2 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTree", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTree indicates an expected call of CreateTree.
func (mr *MockRepositoryMockRecorder) CreateTree(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTree", reflect.TypeOf((*MockRepository)(nil).CreateTree), arg0, arg1, arg2)
}

// GetLogByHash mocks base method.
func (m *MockRepository) GetLogByHash(arg0 context.Context, arg1 int64, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogByHash", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogByHash indicates an expected call of GetLogByHash.
func (mr *MockRepositoryMockRecorder) GetLogByHash(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogByHash", reflect.TypeOf((*MockRepository)(nil).GetLogByHash), arg0, arg1, arg2)
}
