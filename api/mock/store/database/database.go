// Code generated by MockGen. DO NOT EDIT.
// Source: database.go

// Package mock_database is a generated GoMock package.
package mock_database

import (
	context "context"
	reflect "reflect"

	database "github.com/and-period/furumaru/api/internal/store/database"
	entity "github.com/and-period/furumaru/api/internal/store/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockCategory is a mock of Category interface.
type MockCategory struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryMockRecorder
}

// MockCategoryMockRecorder is the mock recorder for MockCategory.
type MockCategoryMockRecorder struct {
	mock *MockCategory
}

// NewMockCategory creates a new mock instance.
func NewMockCategory(ctrl *gomock.Controller) *MockCategory {
	mock := &MockCategory{ctrl: ctrl}
	mock.recorder = &MockCategoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategory) EXPECT() *MockCategoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCategory) Create(ctx context.Context, category *entity.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, category)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCategoryMockRecorder) Create(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCategory)(nil).Create), ctx, category)
}

// Delete mocks base method.
func (m *MockCategory) Delete(ctx context.Context, categoryID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, categoryID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCategoryMockRecorder) Delete(ctx, categoryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCategory)(nil).Delete), ctx, categoryID)
}

// List mocks base method.
func (m *MockCategory) List(ctx context.Context, params *database.ListCategoriesParams, fields ...string) (entity.Categories, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.Categories)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCategoryMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCategory)(nil).List), varargs...)
}

// Update mocks base method.
func (m *MockCategory) Update(ctx context.Context, categoryID, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, categoryID, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCategoryMockRecorder) Update(ctx, categoryID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCategory)(nil).Update), ctx, categoryID, name)
}

// MockProductType is a mock of ProductType interface.
type MockProductType struct {
	ctrl     *gomock.Controller
	recorder *MockProductTypeMockRecorder
}

// MockProductTypeMockRecorder is the mock recorder for MockProductType.
type MockProductTypeMockRecorder struct {
	mock *MockProductType
}

// NewMockProductType creates a new mock instance.
func NewMockProductType(ctrl *gomock.Controller) *MockProductType {
	mock := &MockProductType{ctrl: ctrl}
	mock.recorder = &MockProductTypeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductType) EXPECT() *MockProductTypeMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProductType) Create(ctx context.Context, productType *entity.ProductType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, productType)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockProductTypeMockRecorder) Create(ctx, productType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductType)(nil).Create), ctx, productType)
}

// Delete mocks base method.
func (m *MockProductType) Delete(ctx context.Context, productTypeID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, productTypeID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProductTypeMockRecorder) Delete(ctx, productTypeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProductType)(nil).Delete), ctx, productTypeID)
}

// List mocks base method.
func (m *MockProductType) List(ctx context.Context, params *database.ListProductTypesParams, fields ...string) (entity.ProductTypes, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(entity.ProductTypes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockProductTypeMockRecorder) List(ctx, params interface{}, fields ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockProductType)(nil).List), varargs...)
}

// Update mocks base method.
func (m *MockProductType) Update(ctx context.Context, productTypeID, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, productTypeID, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockProductTypeMockRecorder) Update(ctx, productTypeID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProductType)(nil).Update), ctx, productTypeID, name)
}
