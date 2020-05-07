// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/remipassmoilesel/gitsearch/cli (interfaces: CliDisplay)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	index "gitlab.com/remipassmoilesel/gitsearch/index"
	reflect "reflect"
)

// MockCliDisplay is a mock of CliDisplay interface
type MockCliDisplay struct {
	ctrl     *gomock.Controller
	recorder *MockCliDisplayMockRecorder
}

// MockCliDisplayMockRecorder is the mock recorder for MockCliDisplay
type MockCliDisplayMockRecorder struct {
	mock *MockCliDisplay
}

// NewMockCliDisplay creates a new mock instance
func NewMockCliDisplay(ctrl *gomock.Controller) *MockCliDisplay {
	mock := &MockCliDisplay{ctrl: ctrl}
	mock.recorder = &MockCliDisplayMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCliDisplay) EXPECT() *MockCliDisplayMockRecorder {
	return m.recorder
}

// IndexBuild mocks base method
func (m *MockCliDisplay) IndexBuild(arg0 index.BuildOperationResult) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IndexBuild", arg0)
}

// IndexBuild indicates an expected call of IndexBuild
func (mr *MockCliDisplayMockRecorder) IndexBuild(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexBuild", reflect.TypeOf((*MockCliDisplay)(nil).IndexBuild), arg0)
}

// IndexClean mocks base method
func (m *MockCliDisplay) IndexClean(arg0 index.CleanOperationResult) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IndexClean", arg0)
}

// IndexClean indicates an expected call of IndexClean
func (mr *MockCliDisplayMockRecorder) IndexClean(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IndexClean", reflect.TypeOf((*MockCliDisplay)(nil).IndexClean), arg0)
}

// Search mocks base method
func (m *MockCliDisplay) Search(arg0 string, arg1 index.SearchResult, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Search indicates an expected call of Search
func (mr *MockCliDisplayMockRecorder) Search(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockCliDisplay)(nil).Search), arg0, arg1, arg2)
}

// ShowFile mocks base method
func (m *MockCliDisplay) ShowFile(arg0 index.IndexedFile, arg1 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ShowFile indicates an expected call of ShowFile
func (mr *MockCliDisplayMockRecorder) ShowFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowFile", reflect.TypeOf((*MockCliDisplay)(nil).ShowFile), arg0, arg1)
}

// StartServer mocks base method
func (m *MockCliDisplay) StartServer(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartServer", arg0)
}

// StartServer indicates an expected call of StartServer
func (mr *MockCliDisplayMockRecorder) StartServer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartServer", reflect.TypeOf((*MockCliDisplay)(nil).StartServer), arg0)
}
