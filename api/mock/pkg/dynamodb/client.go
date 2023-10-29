// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_dynamodb is a generated GoMock package.
package mock_dynamodb

import (
	context "context"
	reflect "reflect"

	dynamodb "github.com/and-period/furumaru/api/pkg/dynamodb"
	gomock "go.uber.org/mock/gomock"
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

// Count mocks base method.
func (m *MockClient) Count(ctx context.Context, entity dynamodb.Entity) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, entity)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockClientMockRecorder) Count(ctx, entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockClient)(nil).Count), ctx, entity)
}

// Get mocks base method.
func (m *MockClient) Get(ctx context.Context, entity dynamodb.Entity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockClientMockRecorder) Get(ctx, entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), ctx, entity)
}

// Insert mocks base method.
func (m *MockClient) Insert(ctx context.Context, entity dynamodb.Entity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockClientMockRecorder) Insert(ctx, entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockClient)(nil).Insert), ctx, entity)
}

// MockEntity is a mock of Entity interface.
type MockEntity struct {
	ctrl     *gomock.Controller
	recorder *MockEntityMockRecorder
}

// MockEntityMockRecorder is the mock recorder for MockEntity.
type MockEntityMockRecorder struct {
	mock *MockEntity
}

// NewMockEntity creates a new mock instance.
func NewMockEntity(ctrl *gomock.Controller) *MockEntity {
	mock := &MockEntity{ctrl: ctrl}
	mock.recorder = &MockEntityMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntity) EXPECT() *MockEntityMockRecorder {
	return m.recorder
}

// PrimaryKey mocks base method.
func (m *MockEntity) PrimaryKey() map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrimaryKey")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// PrimaryKey indicates an expected call of PrimaryKey.
func (mr *MockEntityMockRecorder) PrimaryKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrimaryKey", reflect.TypeOf((*MockEntity)(nil).PrimaryKey))
}

// TableName mocks base method.
func (m *MockEntity) TableName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableName")
	ret0, _ := ret[0].(string)
	return ret0
}

// TableName indicates an expected call of TableName.
func (mr *MockEntityMockRecorder) TableName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableName", reflect.TypeOf((*MockEntity)(nil).TableName))
}
