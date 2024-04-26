// Code generated by MockGen. DO NOT EDIT.
// Source: services/session.service.go
//
// Generated by this command:
//
//	mockgen -source services/session.service.go -destination mocks/session.service.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/michelm117/cycling-coach-lab/model"
	gomock "go.uber.org/mock/gomock"
)

// MockSessionServicer is a mock of SessionServicer interface.
type MockSessionServicer struct {
	ctrl     *gomock.Controller
	recorder *MockSessionServicerMockRecorder
}

// MockSessionServicerMockRecorder is the mock recorder for MockSessionServicer.
type MockSessionServicerMockRecorder struct {
	mock *MockSessionServicer
}

// NewMockSessionServicer creates a new mock instance.
func NewMockSessionServicer(ctrl *gomock.Controller) *MockSessionServicer {
	mock := &MockSessionServicer{ctrl: ctrl}
	mock.recorder = &MockSessionServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionServicer) EXPECT() *MockSessionServicerMockRecorder {
	return m.recorder
}

// AuthenticateUserBySessionID mocks base method.
func (m *MockSessionServicer) AuthenticateUserBySessionID(sessionID string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateUserBySessionID", sessionID)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthenticateUserBySessionID indicates an expected call of AuthenticateUserBySessionID.
func (mr *MockSessionServicerMockRecorder) AuthenticateUserBySessionID(sessionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateUserBySessionID", reflect.TypeOf((*MockSessionServicer)(nil).AuthenticateUserBySessionID), sessionID)
}

// DeleteAllSessions mocks base method.
func (m *MockSessionServicer) DeleteAllSessions() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllSessions")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllSessions indicates an expected call of DeleteAllSessions.
func (mr *MockSessionServicerMockRecorder) DeleteAllSessions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllSessions", reflect.TypeOf((*MockSessionServicer)(nil).DeleteAllSessions))
}

// DeleteExpiredSessions mocks base method.
func (m *MockSessionServicer) DeleteExpiredSessions() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExpiredSessions")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExpiredSessions indicates an expected call of DeleteExpiredSessions.
func (mr *MockSessionServicerMockRecorder) DeleteExpiredSessions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpiredSessions", reflect.TypeOf((*MockSessionServicer)(nil).DeleteExpiredSessions))
}

// DeleteSession mocks base method.
func (m *MockSessionServicer) DeleteSession(sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSession", sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSession indicates an expected call of DeleteSession.
func (mr *MockSessionServicerMockRecorder) DeleteSession(sessionID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSession", reflect.TypeOf((*MockSessionServicer)(nil).DeleteSession), sessionID)
}

// GetByUUID mocks base method.
func (m *MockSessionServicer) GetByUUID(uuid string) (*model.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUUID", uuid)
	ret0, _ := ret[0].(*model.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUUID indicates an expected call of GetByUUID.
func (mr *MockSessionServicerMockRecorder) GetByUUID(uuid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUUID", reflect.TypeOf((*MockSessionServicer)(nil).GetByUUID), uuid)
}

// GetByUserID mocks base method.
func (m *MockSessionServicer) GetByUserID(userID int) ([]model.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", userID)
	ret0, _ := ret[0].([]model.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserID indicates an expected call of GetByUserID.
func (mr *MockSessionServicerMockRecorder) GetByUserID(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockSessionServicer)(nil).GetByUserID), userID)
}

// SaveSession mocks base method.
func (m *MockSessionServicer) SaveSession(userID int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSession", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveSession indicates an expected call of SaveSession.
func (mr *MockSessionServicerMockRecorder) SaveSession(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSession", reflect.TypeOf((*MockSessionServicer)(nil).SaveSession), userID)
}

// ScheduleSessionCleanup mocks base method.
func (m *MockSessionServicer) ScheduleSessionCleanup() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ScheduleSessionCleanup")
}

// ScheduleSessionCleanup indicates an expected call of ScheduleSessionCleanup.
func (mr *MockSessionServicerMockRecorder) ScheduleSessionCleanup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleSessionCleanup", reflect.TypeOf((*MockSessionServicer)(nil).ScheduleSessionCleanup))
}
