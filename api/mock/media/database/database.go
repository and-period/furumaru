// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"

	database "github.com/and-period/furumaru/api/internal/media/database"
	entity "github.com/and-period/furumaru/api/internal/media/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockBroadcast is a mock of Broadcast interface.
type MockBroadcast struct {
	ctrl     *gomock.Controller
	recorder *MockBroadcastMockRecorder
}

// MockBroadcastMockRecorder is the mock recorder for MockBroadcast.
type MockBroadcastMockRecorder struct {
	mock *MockBroadcast
}

// NewMockBroadcast creates a new mock instance.
func NewMockBroadcast(ctrl *gomock.Controller) *MockBroadcast {
	mock := &MockBroadcast{ctrl: ctrl}
	mock.recorder = &MockBroadcastMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBroadcast) EXPECT() *MockBroadcastMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockBroadcast) Create(ctx context.Context, broadcast *entity.Broadcast) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, broadcast)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockBroadcastMockRecorder) Create(ctx, broadcast interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBroadcast)(nil).Create), ctx, broadcast)
}

// GetByScheduleID mocks base method.
func (m *MockBroadcast) GetByScheduleID(ctx context.Context, scheduleID string, fields ...string) (*entity.Broadcast, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, scheduleID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByScheduleID", varargs...)
	ret0, _ := ret[0].(*entity.Broadcast)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByScheduleID indicates an expected call of GetByScheduleID.
func (mr *MockBroadcastMockRecorder) GetByScheduleID(ctx, scheduleID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, scheduleID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByScheduleID", reflect.TypeOf((*MockBroadcast)(nil).GetByScheduleID), varargs...)
}

// Update mocks base method.
func (m *MockBroadcast) Update(ctx context.Context, broadcastID string, params *database.UpdateBroadcastParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, broadcastID, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockBroadcastMockRecorder) Update(ctx, broadcastID, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockBroadcast)(nil).Update), ctx, broadcastID, params)
}
