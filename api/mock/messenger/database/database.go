// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"

	database "github.com/and-period/furumaru/api/internal/messenger/database"
	entity "github.com/and-period/furumaru/api/internal/messenger/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockContact is a mock of Contact interface.
type MockContact struct {
	ctrl     *gomock.Controller
	recorder *MockContactMockRecorder
}

// MockContactMockRecorder is the mock recorder for MockContact.
type MockContactMockRecorder struct {
	mock *MockContact
}

// NewMockContact creates a new mock instance.
func NewMockContact(ctrl *gomock.Controller) *MockContact {
	mock := &MockContact{ctrl: ctrl}
	mock.recorder = &MockContactMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContact) EXPECT() *MockContactMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockContact) Create(ctx context.Context, contact *entity.Contact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, contact)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockContactMockRecorder) Create(ctx, contact interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockContact)(nil).Create), ctx, contact)
}

// Delete mocks base method.
func (m *MockContact) Delete(ctx context.Context, contactID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, contactID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockContactMockRecorder) Delete(ctx, contactID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockContact)(nil).Delete), ctx, contactID)
}

// Get mocks base method.
func (m *MockContact) Get(ctx context.Context, contactID string, fields ...string) (*entity.Contact, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, contactID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockContactMockRecorder) Get(ctx, contactID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, contactID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockContact)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockContact) List(ctx context.Context, params *database.ListContactsParams, fields ...string) (entity.Contacts, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.Contacts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockContactMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockContact)(nil).List), varargs...)
}

// Update mocks base method.
func (m *MockContact) Update(ctx context.Context, contactID string, params *database.UpdateContactParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, contactID, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockContactMockRecorder) Update(ctx, contactID, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockContact)(nil).Update), ctx, contactID, params)
}

// MockNotification is a mock of Notification interface.
type MockNotification struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationMockRecorder
}

// MockNotificationMockRecorder is the mock recorder for MockNotification.
type MockNotificationMockRecorder struct {
	mock *MockNotification
}

// NewMockNotification creates a new mock instance.
func NewMockNotification(ctrl *gomock.Controller) *MockNotification {
	mock := &MockNotification{ctrl: ctrl}
	mock.recorder = &MockNotificationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotification) EXPECT() *MockNotificationMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockNotification) Create(ctx context.Context, notification *entity.Notification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, notification)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockNotificationMockRecorder) Create(ctx, notification interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNotification)(nil).Create), ctx, notification)
}
