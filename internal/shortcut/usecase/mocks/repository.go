// Code generated by MockGen. DO NOT EDIT.
// Source: smallurl/internal/shortcut/usecase (interfaces: Repository)
//
// Generated by this command:
//
//	mockgen -destination=mocks/repository.go -package=mu . Repository
//

// Package mu is a generated GoMock package.
package mu

import (
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

// AddURL mocks base method.
func (m *MockRepository) AddURL(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddURL", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddURL indicates an expected call of AddURL.
func (mr *MockRepositoryMockRecorder) AddURL(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddURL", reflect.TypeOf((*MockRepository)(nil).AddURL), arg0, arg1)
}

// GetLongURL mocks base method.
func (m *MockRepository) GetLongURL(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLongURL", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLongURL indicates an expected call of GetLongURL.
func (mr *MockRepositoryMockRecorder) GetLongURL(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLongURL", reflect.TypeOf((*MockRepository)(nil).GetLongURL), arg0)
}