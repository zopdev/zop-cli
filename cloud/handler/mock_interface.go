// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go
//
// Generated by this command:
//
//	mockgen -source=interface.go -destination=mock_interface.go -package=handler
//

// Package handler is a generated GoMock package.
package handler

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	gofr "gofr.dev/pkg/gofr"
	list "zop.dev/cli/zop/cloud/service/list"
)

// MockAccountImporter is a mock of AccountImporter interface.
type MockAccountImporter struct {
	ctrl     *gomock.Controller
	recorder *MockAccountImporterMockRecorder
	isgomock struct{}
}

// MockAccountImporterMockRecorder is the mock recorder for MockAccountImporter.
type MockAccountImporterMockRecorder struct {
	mock *MockAccountImporter
}

// NewMockAccountImporter creates a new mock instance.
func NewMockAccountImporter(ctrl *gomock.Controller) *MockAccountImporter {
	mock := &MockAccountImporter{ctrl: ctrl}
	mock.recorder = &MockAccountImporterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountImporter) EXPECT() *MockAccountImporterMockRecorder {
	return m.recorder
}

// PostAccounts mocks base method.
func (m *MockAccountImporter) PostAccounts(ctx *gofr.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostAccounts", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostAccounts indicates an expected call of PostAccounts.
func (mr *MockAccountImporterMockRecorder) PostAccounts(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostAccounts", reflect.TypeOf((*MockAccountImporter)(nil).PostAccounts), ctx)
}

// MockAccountGetter is a mock of AccountGetter interface.
type MockAccountGetter struct {
	ctrl     *gomock.Controller
	recorder *MockAccountGetterMockRecorder
	isgomock struct{}
}

// MockAccountGetterMockRecorder is the mock recorder for MockAccountGetter.
type MockAccountGetterMockRecorder struct {
	mock *MockAccountGetter
}

// NewMockAccountGetter creates a new mock instance.
func NewMockAccountGetter(ctrl *gomock.Controller) *MockAccountGetter {
	mock := &MockAccountGetter{ctrl: ctrl}
	mock.recorder = &MockAccountGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountGetter) EXPECT() *MockAccountGetterMockRecorder {
	return m.recorder
}

// GetAccounts mocks base method.
func (m *MockAccountGetter) GetAccounts(ctx *gofr.Context) ([]*list.CloudAccountResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", ctx)
	ret0, _ := ret[0].([]*list.CloudAccountResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccounts indicates an expected call of GetAccounts.
func (mr *MockAccountGetterMockRecorder) GetAccounts(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockAccountGetter)(nil).GetAccounts), ctx)
}
