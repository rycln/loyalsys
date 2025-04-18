// Code generated by MockGen. DO NOT EDIT.
// Source: register.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/rycln/loyalsys/internal/models"
)

// MockregServicer is a mock of regServicer interface.
type MockregServicer struct {
	ctrl     *gomock.Controller
	recorder *MockregServicerMockRecorder
}

// MockregServicerMockRecorder is the mock recorder for MockregServicer.
type MockregServicerMockRecorder struct {
	mock *MockregServicer
}

// NewMockregServicer creates a new mock instance.
func NewMockregServicer(ctrl *gomock.Controller) *MockregServicer {
	mock := &MockregServicer{ctrl: ctrl}
	mock.recorder = &MockregServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockregServicer) EXPECT() *MockregServicerMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockregServicer) CreateUser(arg0 context.Context, arg1 *models.User) (models.UserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(models.UserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockregServicerMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockregServicer)(nil).CreateUser), arg0, arg1)
}
