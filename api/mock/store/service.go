// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	context "context"
	reflect "reflect"

	store "github.com/and-period/furumaru/api/internal/store"
	entity "github.com/and-period/furumaru/api/internal/store/entity"
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

// CreateCategory mocks base method.
func (m *MockService) CreateCategory(ctx context.Context, in *store.CreateCategoryInput) (*entity.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", ctx, in)
	ret0, _ := ret[0].(*entity.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCategory indicates an expected call of CreateCategory.
func (mr *MockServiceMockRecorder) CreateCategory(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockService)(nil).CreateCategory), ctx, in)
}

// CreateProductType mocks base method.
func (m *MockService) CreateProductType(ctx context.Context, in *store.CreateProductTypeInput) (*entity.ProductType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProductType", ctx, in)
	ret0, _ := ret[0].(*entity.ProductType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProductType indicates an expected call of CreateProductType.
func (mr *MockServiceMockRecorder) CreateProductType(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProductType", reflect.TypeOf((*MockService)(nil).CreateProductType), ctx, in)
}

// DeleteCategory mocks base method.
func (m *MockService) DeleteCategory(ctx context.Context, in *store.DeleteCategoryInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategory", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategory indicates an expected call of DeleteCategory.
func (mr *MockServiceMockRecorder) DeleteCategory(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategory", reflect.TypeOf((*MockService)(nil).DeleteCategory), ctx, in)
}

// DeleteProductType mocks base method.
func (m *MockService) DeleteProductType(ctx context.Context, in *store.DeleteProductTypeInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProductType", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProductType indicates an expected call of DeleteProductType.
func (mr *MockServiceMockRecorder) DeleteProductType(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProductType", reflect.TypeOf((*MockService)(nil).DeleteProductType), ctx, in)
}

// ListCategories mocks base method.
func (m *MockService) ListCategories(ctx context.Context, in *store.ListCategoriesInput) (entity.Categories, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCategories", ctx, in)
	ret0, _ := ret[0].(entity.Categories)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCategories indicates an expected call of ListCategories.
func (mr *MockServiceMockRecorder) ListCategories(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCategories", reflect.TypeOf((*MockService)(nil).ListCategories), ctx, in)
}

// ListProductTypes mocks base method.
func (m *MockService) ListProductTypes(ctx context.Context, in *store.ListProductTypesInput) (entity.ProductTypes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProductTypes", ctx, in)
	ret0, _ := ret[0].(entity.ProductTypes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProductTypes indicates an expected call of ListProductTypes.
func (mr *MockServiceMockRecorder) ListProductTypes(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProductTypes", reflect.TypeOf((*MockService)(nil).ListProductTypes), ctx, in)
}

// MultiGetCategories mocks base method.
func (m *MockService) MultiGetCategories(ctx context.Context, in *store.MultiGetCategoriesInput) (entity.Categories, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiGetCategories", ctx, in)
	ret0, _ := ret[0].(entity.Categories)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiGetCategories indicates an expected call of MultiGetCategories.
func (mr *MockServiceMockRecorder) MultiGetCategories(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiGetCategories", reflect.TypeOf((*MockService)(nil).MultiGetCategories), ctx, in)
}

// UpdateCategory mocks base method.
func (m *MockService) UpdateCategory(ctx context.Context, in *store.UpdateCategoryInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCategory indicates an expected call of UpdateCategory.
func (mr *MockServiceMockRecorder) UpdateCategory(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockService)(nil).UpdateCategory), ctx, in)
}

// UpdateProductType mocks base method.
func (m *MockService) UpdateProductType(ctx context.Context, in *store.UpdateProductTypeInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProductType", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProductType indicates an expected call of UpdateProductType.
func (mr *MockServiceMockRecorder) UpdateProductType(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProductType", reflect.TypeOf((*MockService)(nil).UpdateProductType), ctx, in)
}
