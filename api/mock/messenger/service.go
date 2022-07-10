// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_messenger is a generated GoMock package.
package mock_messenger

import (
	context "context"
	reflect "reflect"

	messenger "github.com/and-period/furumaru/api/internal/messenger"
	entity "github.com/and-period/furumaru/api/internal/messenger/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateContact mocks base method.
func (m *MockService) CreateContact(ctx context.Context, in *messenger.CreateContactInput) (*entity.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContact", ctx, in)
	ret0, _ := ret[0].(*entity.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContact indicates an expected call of CreateContact.
func (mr *MockServiceMockRecorder) CreateContact(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContact", reflect.TypeOf((*MockService)(nil).CreateContact), ctx, in)
}

// CreateNotification mocks base method.
func (m *MockService) CreateNotification(ctx context.Context, in *messenger.CreateNotificationInput) (*entity.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNotification", ctx, in)
	ret0, _ := ret[0].(*entity.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNotification indicates an expected call of CreateNotification.
func (mr *MockServiceMockRecorder) CreateNotification(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotification", reflect.TypeOf((*MockService)(nil).CreateNotification), ctx, in)
}

// GetContact mocks base method.
func (m *MockService) GetContact(ctx context.Context, in *messenger.GetContactInput) (*entity.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContact", ctx, in)
	ret0, _ := ret[0].(*entity.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContact indicates an expected call of GetContact.
func (mr *MockServiceMockRecorder) GetContact(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContact", reflect.TypeOf((*MockService)(nil).GetContact), ctx, in)
}

// ListContacts mocks base method.
func (m *MockService) ListContacts(ctx context.Context, in *messenger.ListContactsInput) (entity.Contacts, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListContacts", ctx, in)
	ret0, _ := ret[0].(entity.Contacts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListContacts indicates an expected call of ListContacts.
func (mr *MockServiceMockRecorder) ListContacts(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListContacts", reflect.TypeOf((*MockService)(nil).ListContacts), ctx, in)
}

// NotifyRegisterAdmin mocks base method.
func (m *MockService) NotifyRegisterAdmin(ctx context.Context, in *messenger.NotifyRegisterAdminInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NotifyRegisterAdmin", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// NotifyRegisterAdmin indicates an expected call of NotifyRegisterAdmin.
func (mr *MockServiceMockRecorder) NotifyRegisterAdmin(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyRegisterAdmin", reflect.TypeOf((*MockService)(nil).NotifyRegisterAdmin), ctx, in)
}

// UpdateContact mocks base method.
func (m *MockService) UpdateContact(ctx context.Context, in *messenger.UpdateContactInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContact", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContact indicates an expected call of UpdateContact.
func (mr *MockServiceMockRecorder) UpdateContact(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContact", reflect.TypeOf((*MockService)(nil).UpdateContact), ctx, in)
}
