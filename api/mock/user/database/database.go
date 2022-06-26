// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"

	database "github.com/and-period/furumaru/api/internal/user/database"
	entity "github.com/and-period/furumaru/api/internal/user/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockAdminAuth is a mock of AdminAuth interface.
type MockAdminAuth struct {
	ctrl     *gomock.Controller
	recorder *MockAdminAuthMockRecorder
}

// MockAdminAuthMockRecorder is the mock recorder for MockAdminAuth.
type MockAdminAuthMockRecorder struct {
	mock *MockAdminAuth
}

// NewMockAdminAuth creates a new mock instance.
func NewMockAdminAuth(ctrl *gomock.Controller) *MockAdminAuth {
	mock := &MockAdminAuth{ctrl: ctrl}
	mock.recorder = &MockAdminAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdminAuth) EXPECT() *MockAdminAuthMockRecorder {
	return m.recorder
}

// GetByCognitoID mocks base method.
func (m *MockAdminAuth) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.AdminAuth, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, cognitoID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByCognitoID", varargs...)
	ret0, _ := ret[0].(*entity.AdminAuth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCognitoID indicates an expected call of GetByCognitoID.
func (mr *MockAdminAuthMockRecorder) GetByCognitoID(ctx, cognitoID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, cognitoID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCognitoID", reflect.TypeOf((*MockAdminAuth)(nil).GetByCognitoID), varargs...)
}

// MockAdministrator is a mock of Administrator interface.
type MockAdministrator struct {
	ctrl     *gomock.Controller
	recorder *MockAdministratorMockRecorder
}

// MockAdministratorMockRecorder is the mock recorder for MockAdministrator.
type MockAdministratorMockRecorder struct {
	mock *MockAdministrator
}

// NewMockAdministrator creates a new mock instance.
func NewMockAdministrator(ctrl *gomock.Controller) *MockAdministrator {
	mock := &MockAdministrator{ctrl: ctrl}
	mock.recorder = &MockAdministratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdministrator) EXPECT() *MockAdministratorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAdministrator) Create(ctx context.Context, auth *entity.AdminAuth, administrator *entity.Administrator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, auth, administrator)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAdministratorMockRecorder) Create(ctx, auth, administrator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAdministrator)(nil).Create), ctx, auth, administrator)
}

// Get mocks base method.
func (m *MockAdministrator) Get(ctx context.Context, administratorID string, fields ...string) (*entity.Administrator, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, administratorID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.Administrator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAdministratorMockRecorder) Get(ctx, administratorID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, administratorID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAdministrator)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockAdministrator) List(ctx context.Context, params *database.ListAdministratorsParams, fields ...string) (entity.Administrators, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.Administrators)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockAdministratorMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAdministrator)(nil).List), varargs...)
}

// UpdateEmail mocks base method.
func (m *MockAdministrator) UpdateEmail(ctx context.Context, administratorID, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", ctx, administratorID, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmail indicates an expected call of UpdateEmail.
func (mr *MockAdministratorMockRecorder) UpdateEmail(ctx, administratorID, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockAdministrator)(nil).UpdateEmail), ctx, administratorID, email)
}

// MockCoordinator is a mock of Coordinator interface.
type MockCoordinator struct {
	ctrl     *gomock.Controller
	recorder *MockCoordinatorMockRecorder
}

// MockCoordinatorMockRecorder is the mock recorder for MockCoordinator.
type MockCoordinatorMockRecorder struct {
	mock *MockCoordinator
}

// NewMockCoordinator creates a new mock instance.
func NewMockCoordinator(ctrl *gomock.Controller) *MockCoordinator {
	mock := &MockCoordinator{ctrl: ctrl}
	mock.recorder = &MockCoordinatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoordinator) EXPECT() *MockCoordinatorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCoordinator) Create(ctx context.Context, auth *entity.AdminAuth, coordinator *entity.Coordinator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, auth, coordinator)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCoordinatorMockRecorder) Create(ctx, auth, coordinator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCoordinator)(nil).Create), ctx, auth, coordinator)
}

// Get mocks base method.
func (m *MockCoordinator) Get(ctx context.Context, coordinatorID string, fields ...string) (*entity.Coordinator, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, coordinatorID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.Coordinator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCoordinatorMockRecorder) Get(ctx, coordinatorID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, coordinatorID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCoordinator)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockCoordinator) List(ctx context.Context, params *database.ListCoordinatorsParams, fields ...string) (entity.Coordinators, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.Coordinators)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCoordinatorMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCoordinator)(nil).List), varargs...)
}

// UpdateEmail mocks base method.
func (m *MockCoordinator) UpdateEmail(ctx context.Context, coordinatorID, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", ctx, coordinatorID, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmail indicates an expected call of UpdateEmail.
func (mr *MockCoordinatorMockRecorder) UpdateEmail(ctx, coordinatorID, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockCoordinator)(nil).UpdateEmail), ctx, coordinatorID, email)
}

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProducer) Create(ctx context.Context, auth *entity.AdminAuth, producer *entity.Producer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, auth, producer)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockProducerMockRecorder) Create(ctx, auth, producer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProducer)(nil).Create), ctx, auth, producer)
}

// Get mocks base method.
func (m *MockProducer) Get(ctx context.Context, producerID string, fields ...string) (*entity.Producer, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, producerID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.Producer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockProducerMockRecorder) Get(ctx, producerID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, producerID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockProducer)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockProducer) List(ctx context.Context, params *database.ListProducersParams, fields ...string) (entity.Producers, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.Producers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockProducerMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockProducer)(nil).List), varargs...)
}

// UpdateEmail mocks base method.
func (m *MockProducer) UpdateEmail(ctx context.Context, producerID, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", ctx, producerID, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmail indicates an expected call of UpdateEmail.
func (mr *MockProducerMockRecorder) UpdateEmail(ctx, producerID, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockProducer)(nil).UpdateEmail), ctx, producerID, email)
}

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUser) Create(ctx context.Context, user *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUserMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUser)(nil).Create), ctx, user)
}

// Delete mocks base method.
func (m *MockUser) Delete(ctx context.Context, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserMockRecorder) Delete(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUser)(nil).Delete), ctx, userID)
}

// Get mocks base method.
func (m *MockUser) Get(ctx context.Context, userID string, fields ...string) (*entity.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, userID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserMockRecorder) Get(ctx, userID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, userID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUser)(nil).Get), varargs...)
}

// GetByCognitoID mocks base method.
func (m *MockUser) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, cognitoID}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByCognitoID", varargs...)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCognitoID indicates an expected call of GetByCognitoID.
func (mr *MockUserMockRecorder) GetByCognitoID(ctx, cognitoID interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, cognitoID}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCognitoID", reflect.TypeOf((*MockUser)(nil).GetByCognitoID), varargs...)
}

// GetByEmail mocks base method.
func (m *MockUser) GetByEmail(ctx context.Context, email string, fields ...string) (*entity.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, email}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByEmail", varargs...)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail.
func (mr *MockUserMockRecorder) GetByEmail(ctx, email interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, email}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUser)(nil).GetByEmail), varargs...)
}

// UpdateAccount mocks base method.
func (m *MockUser) UpdateAccount(ctx context.Context, userID, accountID, username string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccount", ctx, userID, accountID, username)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAccount indicates an expected call of UpdateAccount.
func (mr *MockUserMockRecorder) UpdateAccount(ctx, userID, accountID, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccount", reflect.TypeOf((*MockUser)(nil).UpdateAccount), ctx, userID, accountID, username)
}

// UpdateEmail mocks base method.
func (m *MockUser) UpdateEmail(ctx context.Context, userID, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEmail", ctx, userID, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEmail indicates an expected call of UpdateEmail.
func (mr *MockUserMockRecorder) UpdateEmail(ctx, userID, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEmail", reflect.TypeOf((*MockUser)(nil).UpdateEmail), ctx, userID, email)
}

// UpdateVerified mocks base method.
func (m *MockUser) UpdateVerified(ctx context.Context, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVerified", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVerified indicates an expected call of UpdateVerified.
func (mr *MockUserMockRecorder) UpdateVerified(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVerified", reflect.TypeOf((*MockUser)(nil).UpdateVerified), ctx, userID)
}
