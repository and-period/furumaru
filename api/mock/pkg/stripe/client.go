// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_stripe is a generated GoMock package.
package mock_stripe

import (
	context "context"
	reflect "reflect"

	stripe "github.com/and-period/furumaru/api/pkg/stripe"
	v73 "github.com/stripe/stripe-go/v73"
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

// AttachPayment mocks base method.
func (m *MockClient) AttachPayment(ctx context.Context, customerID, paymentID string) (*v73.PaymentMethod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachPayment", ctx, customerID, paymentID)
	ret0, _ := ret[0].(*v73.PaymentMethod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AttachPayment indicates an expected call of AttachPayment.
func (mr *MockClientMockRecorder) AttachPayment(ctx, customerID, paymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachPayment", reflect.TypeOf((*MockClient)(nil).AttachPayment), ctx, customerID, paymentID)
}

// Cancel mocks base method.
func (m *MockClient) Cancel(ctx context.Context, transactionID string, reason v73.PaymentIntentCancellationReason) (*v73.PaymentIntent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancel", ctx, transactionID, reason)
	ret0, _ := ret[0].(*v73.PaymentIntent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cancel indicates an expected call of Cancel.
func (mr *MockClientMockRecorder) Cancel(ctx, transactionID, reason interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockClient)(nil).Cancel), ctx, transactionID, reason)
}

// Capture mocks base method.
func (m *MockClient) Capture(ctx context.Context, transactionID string) (*v73.PaymentIntent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capture", ctx, transactionID)
	ret0, _ := ret[0].(*v73.PaymentIntent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Capture indicates an expected call of Capture.
func (mr *MockClientMockRecorder) Capture(ctx, transactionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capture", reflect.TypeOf((*MockClient)(nil).Capture), ctx, transactionID)
}

// CreateCustomer mocks base method.
func (m *MockClient) CreateCustomer(ctx context.Context, in *stripe.CreateCustomerParams) (*v73.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCustomer", ctx, in)
	ret0, _ := ret[0].(*v73.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCustomer indicates an expected call of CreateCustomer.
func (mr *MockClientMockRecorder) CreateCustomer(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockClient)(nil).CreateCustomer), ctx, in)
}

// DeleteCustomer mocks base method.
func (m *MockClient) DeleteCustomer(ctx context.Context, customerID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCustomer", ctx, customerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCustomer indicates an expected call of DeleteCustomer.
func (mr *MockClientMockRecorder) DeleteCustomer(ctx, customerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCustomer", reflect.TypeOf((*MockClient)(nil).DeleteCustomer), ctx, customerID)
}

// DetachPayment mocks base method.
func (m *MockClient) DetachPayment(ctx context.Context, customerID, paymentID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachPayment", ctx, customerID, paymentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DetachPayment indicates an expected call of DetachPayment.
func (mr *MockClientMockRecorder) DetachPayment(ctx, customerID, paymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachPayment", reflect.TypeOf((*MockClient)(nil).DetachPayment), ctx, customerID, paymentID)
}

// GetCard mocks base method.
func (m *MockClient) GetCard(ctx context.Context, customerID, paymentID string) (*v73.PaymentMethod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCard", ctx, customerID, paymentID)
	ret0, _ := ret[0].(*v73.PaymentMethod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCard indicates an expected call of GetCard.
func (mr *MockClientMockRecorder) GetCard(ctx, customerID, paymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCard", reflect.TypeOf((*MockClient)(nil).GetCard), ctx, customerID, paymentID)
}

// GetCustomer mocks base method.
func (m *MockClient) GetCustomer(ctx context.Context, customerID string) (*v73.Customer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCustomer", ctx, customerID)
	ret0, _ := ret[0].(*v73.Customer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCustomer indicates an expected call of GetCustomer.
func (mr *MockClientMockRecorder) GetCustomer(ctx, customerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCustomer", reflect.TypeOf((*MockClient)(nil).GetCustomer), ctx, customerID)
}

// GuestOrder mocks base method.
func (m *MockClient) GuestOrder(ctx context.Context, in *stripe.GuestOrderParams) (*v73.PaymentIntent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GuestOrder", ctx, in)
	ret0, _ := ret[0].(*v73.PaymentIntent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GuestOrder indicates an expected call of GuestOrder.
func (mr *MockClientMockRecorder) GuestOrder(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GuestOrder", reflect.TypeOf((*MockClient)(nil).GuestOrder), ctx, in)
}

// ListCards mocks base method.
func (m *MockClient) ListCards(ctx context.Context, customerID string) ([]*v73.PaymentMethod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCards", ctx, customerID)
	ret0, _ := ret[0].([]*v73.PaymentMethod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCards indicates an expected call of ListCards.
func (mr *MockClientMockRecorder) ListCards(ctx, customerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCards", reflect.TypeOf((*MockClient)(nil).ListCards), ctx, customerID)
}

// Order mocks base method.
func (m *MockClient) Order(ctx context.Context, in *stripe.OrderParams) (*v73.PaymentIntent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Order", ctx, in)
	ret0, _ := ret[0].(*v73.PaymentIntent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Order indicates an expected call of Order.
func (mr *MockClientMockRecorder) Order(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Order", reflect.TypeOf((*MockClient)(nil).Order), ctx, in)
}

// SetupCard mocks base method.
func (m *MockClient) SetupCard(ctx context.Context, customerID string) (*v73.SetupIntent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetupCard", ctx, customerID)
	ret0, _ := ret[0].(*v73.SetupIntent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetupCard indicates an expected call of SetupCard.
func (mr *MockClientMockRecorder) SetupCard(ctx, customerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetupCard", reflect.TypeOf((*MockClient)(nil).SetupCard), ctx, customerID)
}

// UpdateDefaultPayment mocks base method.
func (m *MockClient) UpdateDefaultPayment(ctx context.Context, customerID, paymentID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDefaultPayment", ctx, customerID, paymentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDefaultPayment indicates an expected call of UpdateDefaultPayment.
func (mr *MockClientMockRecorder) UpdateDefaultPayment(ctx, customerID, paymentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDefaultPayment", reflect.TypeOf((*MockClient)(nil).UpdateDefaultPayment), ctx, customerID, paymentID)
}

// MockReceiver is a mock of Receiver interface.
type MockReceiver struct {
	ctrl     *gomock.Controller
	recorder *MockReceiverMockRecorder
}

// MockReceiverMockRecorder is the mock recorder for MockReceiver.
type MockReceiverMockRecorder struct {
	mock *MockReceiver
}

// NewMockReceiver creates a new mock instance.
func NewMockReceiver(ctrl *gomock.Controller) *MockReceiver {
	mock := &MockReceiver{ctrl: ctrl}
	mock.recorder = &MockReceiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceiver) EXPECT() *MockReceiverMockRecorder {
	return m.recorder
}

// Receive mocks base method.
func (m *MockReceiver) Receive(payload []byte, signature string) (*v73.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Receive", payload, signature)
	ret0, _ := ret[0].(*v73.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Receive indicates an expected call of Receive.
func (mr *MockReceiverMockRecorder) Receive(payload, signature interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Receive", reflect.TypeOf((*MockReceiver)(nil).Receive), payload, signature)
}
