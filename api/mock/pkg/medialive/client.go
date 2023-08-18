// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_medialive is a generated GoMock package.
package mock_medialive

import (
	context "context"
	reflect "reflect"

	types "github.com/aws/aws-sdk-go-v2/service/medialive/types"
	gomock "github.com/golang/mock/gomock"
)

// MockMediaLive is a mock of MediaLive interface.
type MockMediaLive struct {
	ctrl     *gomock.Controller
	recorder *MockMediaLiveMockRecorder
}

// MockMediaLiveMockRecorder is the mock recorder for MockMediaLive.
type MockMediaLiveMockRecorder struct {
	mock *MockMediaLive
}

// NewMockMediaLive creates a new mock instance.
func NewMockMediaLive(ctrl *gomock.Controller) *MockMediaLive {
	mock := &MockMediaLive{ctrl: ctrl}
	mock.recorder = &MockMediaLiveMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMediaLive) EXPECT() *MockMediaLiveMockRecorder {
	return m.recorder
}

// CreateSchedule mocks base method.
func (m *MockMediaLive) CreateSchedule(ctx context.Context, channelID string, actions ...types.ScheduleAction) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, channelID}
	for _, a := range actions {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateSchedule", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSchedule indicates an expected call of CreateSchedule.
func (mr *MockMediaLiveMockRecorder) CreateSchedule(ctx, channelID interface{}, actions ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, channelID}, actions...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSchedule", reflect.TypeOf((*MockMediaLive)(nil).CreateSchedule), varargs...)
}

// StartChannel mocks base method.
func (m *MockMediaLive) StartChannel(ctx context.Context, channelID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartChannel", ctx, channelID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartChannel indicates an expected call of StartChannel.
func (mr *MockMediaLiveMockRecorder) StartChannel(ctx, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartChannel", reflect.TypeOf((*MockMediaLive)(nil).StartChannel), ctx, channelID)
}

// StopChannel mocks base method.
func (m *MockMediaLive) StopChannel(ctx context.Context, channelID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopChannel", ctx, channelID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopChannel indicates an expected call of StopChannel.
func (mr *MockMediaLiveMockRecorder) StopChannel(ctx, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopChannel", reflect.TypeOf((*MockMediaLive)(nil).StopChannel), ctx, channelID)
}
