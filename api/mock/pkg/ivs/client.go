// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_ivs is a generated GoMock package.
package mock_ivs

import (
	context "context"
	reflect "reflect"

	ivs "github.com/and-period/furumaru/api/pkg/ivs"
	ivs0 "github.com/aws/aws-sdk-go-v2/service/ivs"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreateChannel mocks base method.
func (m *MockClient) CreateChannel(ctx context.Context, params *ivs.CreateChannelParams) (*ivs0.CreateChannelOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChannel", ctx, params)
	ret0, _ := ret[0].(*ivs0.CreateChannelOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChannel indicates an expected call of CreateChannel.
func (mr *MockClientMockRecorder) CreateChannel(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChannel", reflect.TypeOf((*MockClient)(nil).CreateChannel), ctx, params)
}
