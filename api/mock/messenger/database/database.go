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

// Count mocks base method.
func (m *MockContact) Count(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockContactMockRecorder) Count(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockContact)(nil).Count), ctx)
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

// MockContactCategory is a mock of ContactCategory interface.
type MockContactCategory struct {
	ctrl     *gomock.Controller
	recorder *MockContactCategoryMockRecorder
}

// MockContactCategoryMockRecorder is the mock recorder for MockContactCategory.
type MockContactCategoryMockRecorder struct {
	mock *MockContactCategory
}

// NewMockContactCategory creates a new mock instance.
func NewMockContactCategory(ctrl *gomock.Controller) *MockContactCategory {
	mock := &MockContactCategory{ctrl: ctrl}
	mock.recorder = &MockContactCategoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContactCategory) EXPECT() *MockContactCategoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockContactCategory) Create(ctx context.Context, category *entity.ContactCategory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, category)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockContactCategoryMockRecorder) Create(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockContactCategory)(nil).Create), ctx, category)
}

// Get mocks base method.
func (m *MockContactCategory) Get(ctx context.Context, categoryID string, fields ...string) (*entity.ContactCategory, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, categoryID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.ContactCategory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockContactCategoryMockRecorder) Get(ctx, categoryID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, categoryID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockContactCategory)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockContactCategory) List(ctx context.Context, params *database.ListContactCategoriesParams, fields ...string) (entity.ContactCategories, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.ContactCategories)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockContactCategoryMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockContactCategory)(nil).List), varargs...)
}

// MockContactRead is a mock of ContactRead interface.
type MockContactRead struct {
	ctrl     *gomock.Controller
	recorder *MockContactReadMockRecorder
}

// MockContactReadMockRecorder is the mock recorder for MockContactRead.
type MockContactReadMockRecorder struct {
	mock *MockContactRead
}

// NewMockContactRead creates a new mock instance.
func NewMockContactRead(ctrl *gomock.Controller) *MockContactRead {
	mock := &MockContactRead{ctrl: ctrl}
	mock.recorder = &MockContactReadMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContactRead) EXPECT() *MockContactReadMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockContactRead) Create(ctx context.Context, contactRead *entity.ContactRead) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, contactRead)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockContactReadMockRecorder) Create(ctx, contactRead interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockContactRead)(nil).Create), ctx, contactRead)
}

// MockThread is a mock of Thread interface.
type MockThread struct {
	ctrl     *gomock.Controller
	recorder *MockThreadMockRecorder
}

// MockThreadMockRecorder is the mock recorder for MockThread.
type MockThreadMockRecorder struct {
	mock *MockThread
}

// NewMockThread creates a new mock instance.
func NewMockThread(ctrl *gomock.Controller) *MockThread {
	mock := &MockThread{ctrl: ctrl}
	mock.recorder = &MockThreadMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockThread) EXPECT() *MockThreadMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockThread) Count(ctx context.Context, params *database.ListThreadsByContactIDParams) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, params)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockThreadMockRecorder) Count(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockThread)(nil).Count), ctx, params)
}

// Create mocks base method.
func (m *MockThread) Create(ctx context.Context, thread *entity.Thread) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, thread)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockThreadMockRecorder) Create(ctx, thread interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockThread)(nil).Create), ctx, thread)
}

// Delete mocks base method.
func (m *MockThread) Delete(ctx context.Context, threadID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, threadID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockThreadMockRecorder) Delete(ctx, threadID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockThread)(nil).Delete), ctx, threadID)
}

// Get mocks base method.
func (m *MockThread) Get(ctx context.Context, threadID string, fields ...string) (*entity.Thread, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, threadID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.Thread)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockThreadMockRecorder) Get(ctx, threadID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, threadID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockThread)(nil).Get), varargs...)
}

// ListByContactID mocks base method.
func (m *MockThread) ListByContactID(ctx context.Context, params *database.ListThreadsByContactIDParams, fields ...string) (entity.Threads, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListByContactID", varargs...)
	ret0, _ := ret[0].(entity.Threads)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByContactID indicates an expected call of ListByContactID.
func (mr *MockThreadMockRecorder) ListByContactID(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByContactID", reflect.TypeOf((*MockThread)(nil).ListByContactID), varargs...)
}

// Update mocks base method.
func (m *MockThread) Update(ctx context.Context, threadID string, params *database.UpdateThreadParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, threadID, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockThreadMockRecorder) Update(ctx, threadID, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockThread)(nil).Update), ctx, threadID, params)
}

// MockMessage is a mock of Message interface.
type MockMessage struct {
	ctrl     *gomock.Controller
	recorder *MockMessageMockRecorder
}

// MockMessageMockRecorder is the mock recorder for MockMessage.
type MockMessageMockRecorder struct {
	mock *MockMessage
}

// NewMockMessage creates a new mock instance.
func NewMockMessage(ctrl *gomock.Controller) *MockMessage {
	mock := &MockMessage{ctrl: ctrl}
	mock.recorder = &MockMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessage) EXPECT() *MockMessageMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockMessage) Count(ctx context.Context, params *database.ListMessagesParams) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, params)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockMessageMockRecorder) Count(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockMessage)(nil).Count), ctx, params)
}

// Get mocks base method.
func (m *MockMessage) Get(ctx context.Context, messageID string, fields ...string) (*entity.Message, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, messageID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMessageMockRecorder) Get(ctx, messageID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, messageID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMessage)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockMessage) List(ctx context.Context, params *database.ListMessagesParams, fields ...string) (entity.Messages, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.Messages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockMessageMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockMessage)(nil).List), varargs...)
}

// MultiCreate mocks base method.
func (m *MockMessage) MultiCreate(ctx context.Context, messages entity.Messages) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiCreate", ctx, messages)
	ret0, _ := ret[0].(error)
	return ret0
}

// MultiCreate indicates an expected call of MultiCreate.
func (mr *MockMessageMockRecorder) MultiCreate(ctx, messages interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiCreate", reflect.TypeOf((*MockMessage)(nil).MultiCreate), ctx, messages)
}

// UpdateRead mocks base method.
func (m *MockMessage) UpdateRead(ctx context.Context, messageID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRead", ctx, messageID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRead indicates an expected call of UpdateRead.
func (mr *MockMessageMockRecorder) UpdateRead(ctx, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRead", reflect.TypeOf((*MockMessage)(nil).UpdateRead), ctx, messageID)
}

// MockMessageTemplate is a mock of MessageTemplate interface.
type MockMessageTemplate struct {
	ctrl     *gomock.Controller
	recorder *MockMessageTemplateMockRecorder
}

// MockMessageTemplateMockRecorder is the mock recorder for MockMessageTemplate.
type MockMessageTemplateMockRecorder struct {
	mock *MockMessageTemplate
}

// NewMockMessageTemplate creates a new mock instance.
func NewMockMessageTemplate(ctrl *gomock.Controller) *MockMessageTemplate {
	mock := &MockMessageTemplate{ctrl: ctrl}
	mock.recorder = &MockMessageTemplateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageTemplate) EXPECT() *MockMessageTemplateMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockMessageTemplate) Get(ctx context.Context, messageID string, fields ...string) (*entity.MessageTemplate, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, messageID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.MessageTemplate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMessageTemplateMockRecorder) Get(ctx, messageID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, messageID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMessageTemplate)(nil).Get), varargs...)
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

// Count mocks base method.
func (m *MockNotification) Count(ctx context.Context, params *database.ListNotificationsParams) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, params)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockNotificationMockRecorder) Count(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockNotification)(nil).Count), ctx, params)
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

// Delete mocks base method.
func (m *MockNotification) Delete(ctx context.Context, notificationID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, notificationID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockNotificationMockRecorder) Delete(ctx, notificationID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNotification)(nil).Delete), ctx, notificationID)
}

// Get mocks base method.
func (m *MockNotification) Get(ctx context.Context, notificationID string, fields ...string) (*entity.Notification, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, notificationID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockNotificationMockRecorder) Get(ctx, notificationID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, notificationID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockNotification)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockNotification) List(ctx context.Context, params *database.ListNotificationsParams, fields ...string) (entity.Notifications, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.Notifications)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockNotificationMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNotification)(nil).List), varargs...)
}

// Update mocks base method.
func (m *MockNotification) Update(ctx context.Context, notificationID string, params *database.UpdateNotificationParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, notificationID, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockNotificationMockRecorder) Update(ctx, notificationID, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNotification)(nil).Update), ctx, notificationID, params)
}

// MockPushTemplate is a mock of PushTemplate interface.
type MockPushTemplate struct {
	ctrl     *gomock.Controller
	recorder *MockPushTemplateMockRecorder
}

// MockPushTemplateMockRecorder is the mock recorder for MockPushTemplate.
type MockPushTemplateMockRecorder struct {
	mock *MockPushTemplate
}

// NewMockPushTemplate creates a new mock instance.
func NewMockPushTemplate(ctrl *gomock.Controller) *MockPushTemplate {
	mock := &MockPushTemplate{ctrl: ctrl}
	mock.recorder = &MockPushTemplateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPushTemplate) EXPECT() *MockPushTemplateMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPushTemplate) Get(ctx context.Context, pushID string, fields ...string) (*entity.PushTemplate, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, pushID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.PushTemplate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPushTemplateMockRecorder) Get(ctx, pushID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, pushID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPushTemplate)(nil).Get), varargs...)
}

// MockReceivedQueue is a mock of ReceivedQueue interface.
type MockReceivedQueue struct {
	ctrl     *gomock.Controller
	recorder *MockReceivedQueueMockRecorder
}

// MockReceivedQueueMockRecorder is the mock recorder for MockReceivedQueue.
type MockReceivedQueueMockRecorder struct {
	mock *MockReceivedQueue
}

// NewMockReceivedQueue creates a new mock instance.
func NewMockReceivedQueue(ctrl *gomock.Controller) *MockReceivedQueue {
	mock := &MockReceivedQueue{ctrl: ctrl}
	mock.recorder = &MockReceivedQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceivedQueue) EXPECT() *MockReceivedQueueMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockReceivedQueue) Create(ctx context.Context, queue *entity.ReceivedQueue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, queue)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockReceivedQueueMockRecorder) Create(ctx, queue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockReceivedQueue)(nil).Create), ctx, queue)
}

// Get mocks base method.
func (m *MockReceivedQueue) Get(ctx context.Context, queueID string, fields ...string) (*entity.ReceivedQueue, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, queueID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.ReceivedQueue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockReceivedQueueMockRecorder) Get(ctx, queueID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, queueID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockReceivedQueue)(nil).Get), varargs...)
}

// UpdateDone mocks base method.
func (m *MockReceivedQueue) UpdateDone(ctx context.Context, queueID string, done bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDone", ctx, queueID, done)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDone indicates an expected call of UpdateDone.
func (mr *MockReceivedQueueMockRecorder) UpdateDone(ctx, queueID, done interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDone", reflect.TypeOf((*MockReceivedQueue)(nil).UpdateDone), ctx, queueID, done)
}

// MockReportTemplate is a mock of ReportTemplate interface.
type MockReportTemplate struct {
	ctrl     *gomock.Controller
	recorder *MockReportTemplateMockRecorder
}

// MockReportTemplateMockRecorder is the mock recorder for MockReportTemplate.
type MockReportTemplateMockRecorder struct {
	mock *MockReportTemplate
}

// NewMockReportTemplate creates a new mock instance.
func NewMockReportTemplate(ctrl *gomock.Controller) *MockReportTemplate {
	mock := &MockReportTemplate{ctrl: ctrl}
	mock.recorder = &MockReportTemplateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReportTemplate) EXPECT() *MockReportTemplateMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockReportTemplate) Get(ctx context.Context, reportID string, fields ...string) (*entity.ReportTemplate, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, reportID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.ReportTemplate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockReportTemplateMockRecorder) Get(ctx, reportID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, reportID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockReportTemplate)(nil).Get), varargs...)
}

// MockSchedule is a mock of Schedule interface.
type MockSchedule struct {
	ctrl     *gomock.Controller
	recorder *MockScheduleMockRecorder
}

// MockScheduleMockRecorder is the mock recorder for MockSchedule.
type MockScheduleMockRecorder struct {
	mock *MockSchedule
}

// NewMockSchedule creates a new mock instance.
func NewMockSchedule(ctrl *gomock.Controller) *MockSchedule {
	mock := &MockSchedule{ctrl: ctrl}
	mock.recorder = &MockScheduleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSchedule) EXPECT() *MockScheduleMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockSchedule) List(ctx context.Context, params *database.ListSchedulesParams, fields ...string) (entity.Schedules, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.Schedules)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockScheduleMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSchedule)(nil).List), varargs...)
}

// UpdateCancel mocks base method.
func (m *MockSchedule) UpdateCancel(ctx context.Context, messageType entity.ScheduleType, messageID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCancel", ctx, messageType, messageID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCancel indicates an expected call of UpdateCancel.
func (mr *MockScheduleMockRecorder) UpdateCancel(ctx, messageType, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCancel", reflect.TypeOf((*MockSchedule)(nil).UpdateCancel), ctx, messageType, messageID)
}

// UpdateDone mocks base method.
func (m *MockSchedule) UpdateDone(ctx context.Context, messageType entity.ScheduleType, messageID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDone", ctx, messageType, messageID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDone indicates an expected call of UpdateDone.
func (mr *MockScheduleMockRecorder) UpdateDone(ctx, messageType, messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDone", reflect.TypeOf((*MockSchedule)(nil).UpdateDone), ctx, messageType, messageID)
}

// UpsertProcessing mocks base method.
func (m *MockSchedule) UpsertProcessing(ctx context.Context, schedule *entity.Schedule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertProcessing", ctx, schedule)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProcessing indicates an expected call of UpsertProcessing.
func (mr *MockScheduleMockRecorder) UpsertProcessing(ctx, schedule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProcessing", reflect.TypeOf((*MockSchedule)(nil).UpsertProcessing), ctx, schedule)
}
