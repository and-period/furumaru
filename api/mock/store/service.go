// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_store is a generated GoMock package.
package mock_store

import (
	context "context"
	reflect "reflect"

	store "github.com/and-period/furumaru/api/internal/store"
	entity "github.com/and-period/furumaru/api/internal/store/entity"
	gomock "go.uber.org/mock/gomock"
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

// AddCartItem mocks base method.
func (m *MockService) AddCartItem(ctx context.Context, in *store.AddCartItemInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCartItem", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCartItem indicates an expected call of AddCartItem.
func (mr *MockServiceMockRecorder) AddCartItem(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCartItem", reflect.TypeOf((*MockService)(nil).AddCartItem), ctx, in)
}

// AggregateOrders mocks base method.
func (m *MockService) AggregateOrders(ctx context.Context, in *store.AggregateOrdersInput) (entity.AggregatedOrders, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AggregateOrders", ctx, in)
	ret0, _ := ret[0].(entity.AggregatedOrders)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AggregateOrders indicates an expected call of AggregateOrders.
func (mr *MockServiceMockRecorder) AggregateOrders(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AggregateOrders", reflect.TypeOf((*MockService)(nil).AggregateOrders), ctx, in)
}

// ApproveSchedule mocks base method.
func (m *MockService) ApproveSchedule(ctx context.Context, in *store.ApproveScheduleInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApproveSchedule", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApproveSchedule indicates an expected call of ApproveSchedule.
func (mr *MockServiceMockRecorder) ApproveSchedule(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApproveSchedule", reflect.TypeOf((*MockService)(nil).ApproveSchedule), ctx, in)
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

// CreateLive mocks base method.
func (m *MockService) CreateLive(ctx context.Context, in *store.CreateLiveInput) (*entity.Live, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLive", ctx, in)
	ret0, _ := ret[0].(*entity.Live)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLive indicates an expected call of CreateLive.
func (mr *MockServiceMockRecorder) CreateLive(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLive", reflect.TypeOf((*MockService)(nil).CreateLive), ctx, in)
}

// CreateProduct mocks base method.
func (m *MockService) CreateProduct(ctx context.Context, in *store.CreateProductInput) (*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", ctx, in)
	ret0, _ := ret[0].(*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockServiceMockRecorder) CreateProduct(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockService)(nil).CreateProduct), ctx, in)
}

// CreateProductTag mocks base method.
func (m *MockService) CreateProductTag(ctx context.Context, in *store.CreateProductTagInput) (*entity.ProductTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProductTag", ctx, in)
	ret0, _ := ret[0].(*entity.ProductTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProductTag indicates an expected call of CreateProductTag.
func (mr *MockServiceMockRecorder) CreateProductTag(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProductTag", reflect.TypeOf((*MockService)(nil).CreateProductTag), ctx, in)
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

// CreatePromotion mocks base method.
func (m *MockService) CreatePromotion(ctx context.Context, in *store.CreatePromotionInput) (*entity.Promotion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePromotion", ctx, in)
	ret0, _ := ret[0].(*entity.Promotion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePromotion indicates an expected call of CreatePromotion.
func (mr *MockServiceMockRecorder) CreatePromotion(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePromotion", reflect.TypeOf((*MockService)(nil).CreatePromotion), ctx, in)
}

// CreateSchedule mocks base method.
func (m *MockService) CreateSchedule(ctx context.Context, in *store.CreateScheduleInput) (*entity.Schedule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSchedule", ctx, in)
	ret0, _ := ret[0].(*entity.Schedule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSchedule indicates an expected call of CreateSchedule.
func (mr *MockServiceMockRecorder) CreateSchedule(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSchedule", reflect.TypeOf((*MockService)(nil).CreateSchedule), ctx, in)
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

// DeleteLive mocks base method.
func (m *MockService) DeleteLive(ctx context.Context, in *store.DeleteLiveInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLive", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLive indicates an expected call of DeleteLive.
func (mr *MockServiceMockRecorder) DeleteLive(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLive", reflect.TypeOf((*MockService)(nil).DeleteLive), ctx, in)
}

// DeleteProduct mocks base method.
func (m *MockService) DeleteProduct(ctx context.Context, in *store.DeleteProductInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct.
func (mr *MockServiceMockRecorder) DeleteProduct(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockService)(nil).DeleteProduct), ctx, in)
}

// DeleteProductTag mocks base method.
func (m *MockService) DeleteProductTag(ctx context.Context, in *store.DeleteProductTagInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProductTag", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProductTag indicates an expected call of DeleteProductTag.
func (mr *MockServiceMockRecorder) DeleteProductTag(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProductTag", reflect.TypeOf((*MockService)(nil).DeleteProductTag), ctx, in)
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

// DeletePromotion mocks base method.
func (m *MockService) DeletePromotion(ctx context.Context, in *store.DeletePromotionInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePromotion", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePromotion indicates an expected call of DeletePromotion.
func (mr *MockServiceMockRecorder) DeletePromotion(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePromotion", reflect.TypeOf((*MockService)(nil).DeletePromotion), ctx, in)
}

// GetCart mocks base method.
func (m *MockService) GetCart(ctx context.Context, in *store.GetCartInput) (*entity.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCart", ctx, in)
	ret0, _ := ret[0].(*entity.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCart indicates an expected call of GetCart.
func (mr *MockServiceMockRecorder) GetCart(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCart", reflect.TypeOf((*MockService)(nil).GetCart), ctx, in)
}

// GetCategory mocks base method.
func (m *MockService) GetCategory(ctx context.Context, in *store.GetCategoryInput) (*entity.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategory", ctx, in)
	ret0, _ := ret[0].(*entity.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategory indicates an expected call of GetCategory.
func (mr *MockServiceMockRecorder) GetCategory(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategory", reflect.TypeOf((*MockService)(nil).GetCategory), ctx, in)
}

// GetDefaultShipping mocks base method.
func (m *MockService) GetDefaultShipping(ctx context.Context, in *store.GetDefaultShippingInput) (*entity.Shipping, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultShipping", ctx, in)
	ret0, _ := ret[0].(*entity.Shipping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDefaultShipping indicates an expected call of GetDefaultShipping.
func (mr *MockServiceMockRecorder) GetDefaultShipping(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultShipping", reflect.TypeOf((*MockService)(nil).GetDefaultShipping), ctx, in)
}

// GetLive mocks base method.
func (m *MockService) GetLive(ctx context.Context, in *store.GetLiveInput) (*entity.Live, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLive", ctx, in)
	ret0, _ := ret[0].(*entity.Live)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLive indicates an expected call of GetLive.
func (mr *MockServiceMockRecorder) GetLive(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLive", reflect.TypeOf((*MockService)(nil).GetLive), ctx, in)
}

// GetOrder mocks base method.
func (m *MockService) GetOrder(ctx context.Context, in *store.GetOrderInput) (*entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", ctx, in)
	ret0, _ := ret[0].(*entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrder indicates an expected call of GetOrder.
func (mr *MockServiceMockRecorder) GetOrder(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockService)(nil).GetOrder), ctx, in)
}

// GetProduct mocks base method.
func (m *MockService) GetProduct(ctx context.Context, in *store.GetProductInput) (*entity.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", ctx, in)
	ret0, _ := ret[0].(*entity.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct.
func (mr *MockServiceMockRecorder) GetProduct(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockService)(nil).GetProduct), ctx, in)
}

// GetProductTag mocks base method.
func (m *MockService) GetProductTag(ctx context.Context, in *store.GetProductTagInput) (*entity.ProductTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductTag", ctx, in)
	ret0, _ := ret[0].(*entity.ProductTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductTag indicates an expected call of GetProductTag.
func (mr *MockServiceMockRecorder) GetProductTag(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductTag", reflect.TypeOf((*MockService)(nil).GetProductTag), ctx, in)
}

// GetProductType mocks base method.
func (m *MockService) GetProductType(ctx context.Context, in *store.GetProductTypeInput) (*entity.ProductType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductType", ctx, in)
	ret0, _ := ret[0].(*entity.ProductType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductType indicates an expected call of GetProductType.
func (mr *MockServiceMockRecorder) GetProductType(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductType", reflect.TypeOf((*MockService)(nil).GetProductType), ctx, in)
}

// GetPromotion mocks base method.
func (m *MockService) GetPromotion(ctx context.Context, in *store.GetPromotionInput) (*entity.Promotion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPromotion", ctx, in)
	ret0, _ := ret[0].(*entity.Promotion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPromotion indicates an expected call of GetPromotion.
func (mr *MockServiceMockRecorder) GetPromotion(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPromotion", reflect.TypeOf((*MockService)(nil).GetPromotion), ctx, in)
}

// GetSchedule mocks base method.
func (m *MockService) GetSchedule(ctx context.Context, in *store.GetScheduleInput) (*entity.Schedule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchedule", ctx, in)
	ret0, _ := ret[0].(*entity.Schedule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchedule indicates an expected call of GetSchedule.
func (mr *MockServiceMockRecorder) GetSchedule(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchedule", reflect.TypeOf((*MockService)(nil).GetSchedule), ctx, in)
}

// GetShippingByCoordinatorID mocks base method.
func (m *MockService) GetShippingByCoordinatorID(ctx context.Context, in *store.GetShippingByCoordinatorIDInput) (*entity.Shipping, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShippingByCoordinatorID", ctx, in)
	ret0, _ := ret[0].(*entity.Shipping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShippingByCoordinatorID indicates an expected call of GetShippingByCoordinatorID.
func (mr *MockServiceMockRecorder) GetShippingByCoordinatorID(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShippingByCoordinatorID", reflect.TypeOf((*MockService)(nil).GetShippingByCoordinatorID), ctx, in)
}

// ListCategories mocks base method.
func (m *MockService) ListCategories(ctx context.Context, in *store.ListCategoriesInput) (entity.Categories, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCategories", ctx, in)
	ret0, _ := ret[0].(entity.Categories)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListCategories indicates an expected call of ListCategories.
func (mr *MockServiceMockRecorder) ListCategories(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCategories", reflect.TypeOf((*MockService)(nil).ListCategories), ctx, in)
}

// ListLives mocks base method.
func (m *MockService) ListLives(ctx context.Context, in *store.ListLivesInput) (entity.Lives, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListLives", ctx, in)
	ret0, _ := ret[0].(entity.Lives)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListLives indicates an expected call of ListLives.
func (mr *MockServiceMockRecorder) ListLives(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListLives", reflect.TypeOf((*MockService)(nil).ListLives), ctx, in)
}

// ListOrders mocks base method.
func (m *MockService) ListOrders(ctx context.Context, in *store.ListOrdersInput) (entity.Orders, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListOrders", ctx, in)
	ret0, _ := ret[0].(entity.Orders)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListOrders indicates an expected call of ListOrders.
func (mr *MockServiceMockRecorder) ListOrders(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListOrders", reflect.TypeOf((*MockService)(nil).ListOrders), ctx, in)
}

// ListProductTags mocks base method.
func (m *MockService) ListProductTags(ctx context.Context, in *store.ListProductTagsInput) (entity.ProductTags, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProductTags", ctx, in)
	ret0, _ := ret[0].(entity.ProductTags)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListProductTags indicates an expected call of ListProductTags.
func (mr *MockServiceMockRecorder) ListProductTags(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProductTags", reflect.TypeOf((*MockService)(nil).ListProductTags), ctx, in)
}

// ListProductTypes mocks base method.
func (m *MockService) ListProductTypes(ctx context.Context, in *store.ListProductTypesInput) (entity.ProductTypes, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProductTypes", ctx, in)
	ret0, _ := ret[0].(entity.ProductTypes)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListProductTypes indicates an expected call of ListProductTypes.
func (mr *MockServiceMockRecorder) ListProductTypes(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProductTypes", reflect.TypeOf((*MockService)(nil).ListProductTypes), ctx, in)
}

// ListProducts mocks base method.
func (m *MockService) ListProducts(ctx context.Context, in *store.ListProductsInput) (entity.Products, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProducts", ctx, in)
	ret0, _ := ret[0].(entity.Products)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListProducts indicates an expected call of ListProducts.
func (mr *MockServiceMockRecorder) ListProducts(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProducts", reflect.TypeOf((*MockService)(nil).ListProducts), ctx, in)
}

// ListPromotions mocks base method.
func (m *MockService) ListPromotions(ctx context.Context, in *store.ListPromotionsInput) (entity.Promotions, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPromotions", ctx, in)
	ret0, _ := ret[0].(entity.Promotions)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListPromotions indicates an expected call of ListPromotions.
func (mr *MockServiceMockRecorder) ListPromotions(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPromotions", reflect.TypeOf((*MockService)(nil).ListPromotions), ctx, in)
}

// ListSchedules mocks base method.
func (m *MockService) ListSchedules(ctx context.Context, in *store.ListSchedulesInput) (entity.Schedules, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSchedules", ctx, in)
	ret0, _ := ret[0].(entity.Schedules)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListSchedules indicates an expected call of ListSchedules.
func (mr *MockServiceMockRecorder) ListSchedules(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSchedules", reflect.TypeOf((*MockService)(nil).ListSchedules), ctx, in)
}

// ListShippingsByCoordinatorIDs mocks base method.
func (m *MockService) ListShippingsByCoordinatorIDs(ctx context.Context, in *store.ListShippingsByCoordinatorIDsInput) (entity.Shippings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListShippingsByCoordinatorIDs", ctx, in)
	ret0, _ := ret[0].(entity.Shippings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListShippingsByCoordinatorIDs indicates an expected call of ListShippingsByCoordinatorIDs.
func (mr *MockServiceMockRecorder) ListShippingsByCoordinatorIDs(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListShippingsByCoordinatorIDs", reflect.TypeOf((*MockService)(nil).ListShippingsByCoordinatorIDs), ctx, in)
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

// MultiGetProductTags mocks base method.
func (m *MockService) MultiGetProductTags(ctx context.Context, in *store.MultiGetProductTagsInput) (entity.ProductTags, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiGetProductTags", ctx, in)
	ret0, _ := ret[0].(entity.ProductTags)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiGetProductTags indicates an expected call of MultiGetProductTags.
func (mr *MockServiceMockRecorder) MultiGetProductTags(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiGetProductTags", reflect.TypeOf((*MockService)(nil).MultiGetProductTags), ctx, in)
}

// MultiGetProductTypes mocks base method.
func (m *MockService) MultiGetProductTypes(ctx context.Context, in *store.MultiGetProductTypesInput) (entity.ProductTypes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiGetProductTypes", ctx, in)
	ret0, _ := ret[0].(entity.ProductTypes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiGetProductTypes indicates an expected call of MultiGetProductTypes.
func (mr *MockServiceMockRecorder) MultiGetProductTypes(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiGetProductTypes", reflect.TypeOf((*MockService)(nil).MultiGetProductTypes), ctx, in)
}

// MultiGetProducts mocks base method.
func (m *MockService) MultiGetProducts(ctx context.Context, in *store.MultiGetProductsInput) (entity.Products, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiGetProducts", ctx, in)
	ret0, _ := ret[0].(entity.Products)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiGetProducts indicates an expected call of MultiGetProducts.
func (mr *MockServiceMockRecorder) MultiGetProducts(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiGetProducts", reflect.TypeOf((*MockService)(nil).MultiGetProducts), ctx, in)
}

// MultiGetProductsByRevision mocks base method.
func (m *MockService) MultiGetProductsByRevision(ctx context.Context, in *store.MultiGetProductsByRevisionInput) (entity.Products, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiGetProductsByRevision", ctx, in)
	ret0, _ := ret[0].(entity.Products)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiGetProductsByRevision indicates an expected call of MultiGetProductsByRevision.
func (mr *MockServiceMockRecorder) MultiGetProductsByRevision(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiGetProductsByRevision", reflect.TypeOf((*MockService)(nil).MultiGetProductsByRevision), ctx, in)
}

// MultiGetPromotions mocks base method.
func (m *MockService) MultiGetPromotions(ctx context.Context, in *store.MultiGetPromotionsInput) (entity.Promotions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiGetPromotions", ctx, in)
	ret0, _ := ret[0].(entity.Promotions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiGetPromotions indicates an expected call of MultiGetPromotions.
func (mr *MockServiceMockRecorder) MultiGetPromotions(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiGetPromotions", reflect.TypeOf((*MockService)(nil).MultiGetPromotions), ctx, in)
}

// MultiGetSchedules mocks base method.
func (m *MockService) MultiGetSchedules(ctx context.Context, in *store.MultiGetSchedulesInput) (entity.Schedules, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiGetSchedules", ctx, in)
	ret0, _ := ret[0].(entity.Schedules)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiGetSchedules indicates an expected call of MultiGetSchedules.
func (mr *MockServiceMockRecorder) MultiGetSchedules(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiGetSchedules", reflect.TypeOf((*MockService)(nil).MultiGetSchedules), ctx, in)
}

// MultiGetShippingsByRevision mocks base method.
func (m *MockService) MultiGetShippingsByRevision(ctx context.Context, in *store.MultiGetShippingsByRevisionInput) (entity.Shippings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MultiGetShippingsByRevision", ctx, in)
	ret0, _ := ret[0].(entity.Shippings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MultiGetShippingsByRevision indicates an expected call of MultiGetShippingsByRevision.
func (mr *MockServiceMockRecorder) MultiGetShippingsByRevision(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MultiGetShippingsByRevision", reflect.TypeOf((*MockService)(nil).MultiGetShippingsByRevision), ctx, in)
}

// RemoveCartItem mocks base method.
func (m *MockService) RemoveCartItem(ctx context.Context, in *store.RemoveCartItemInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCartItem", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCartItem indicates an expected call of RemoveCartItem.
func (mr *MockServiceMockRecorder) RemoveCartItem(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCartItem", reflect.TypeOf((*MockService)(nil).RemoveCartItem), ctx, in)
}

// SearchPostalCode mocks base method.
func (m *MockService) SearchPostalCode(ctx context.Context, in *store.SearchPostalCodeInput) (*entity.PostalCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchPostalCode", ctx, in)
	ret0, _ := ret[0].(*entity.PostalCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchPostalCode indicates an expected call of SearchPostalCode.
func (mr *MockServiceMockRecorder) SearchPostalCode(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchPostalCode", reflect.TypeOf((*MockService)(nil).SearchPostalCode), ctx, in)
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

// UpdateDefaultShipping mocks base method.
func (m *MockService) UpdateDefaultShipping(ctx context.Context, in *store.UpdateDefaultShippingInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDefaultShipping", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDefaultShipping indicates an expected call of UpdateDefaultShipping.
func (mr *MockServiceMockRecorder) UpdateDefaultShipping(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDefaultShipping", reflect.TypeOf((*MockService)(nil).UpdateDefaultShipping), ctx, in)
}

// UpdateLive mocks base method.
func (m *MockService) UpdateLive(ctx context.Context, in *store.UpdateLiveInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLive", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLive indicates an expected call of UpdateLive.
func (mr *MockServiceMockRecorder) UpdateLive(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLive", reflect.TypeOf((*MockService)(nil).UpdateLive), ctx, in)
}

// UpdateProduct mocks base method.
func (m *MockService) UpdateProduct(ctx context.Context, in *store.UpdateProductInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockServiceMockRecorder) UpdateProduct(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockService)(nil).UpdateProduct), ctx, in)
}

// UpdateProductMedia mocks base method.
func (m *MockService) UpdateProductMedia(ctx context.Context, in *store.UpdateProductMediaInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProductMedia", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProductMedia indicates an expected call of UpdateProductMedia.
func (mr *MockServiceMockRecorder) UpdateProductMedia(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProductMedia", reflect.TypeOf((*MockService)(nil).UpdateProductMedia), ctx, in)
}

// UpdateProductTag mocks base method.
func (m *MockService) UpdateProductTag(ctx context.Context, in *store.UpdateProductTagInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProductTag", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProductTag indicates an expected call of UpdateProductTag.
func (mr *MockServiceMockRecorder) UpdateProductTag(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProductTag", reflect.TypeOf((*MockService)(nil).UpdateProductTag), ctx, in)
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

// UpdateProductTypeIcons mocks base method.
func (m *MockService) UpdateProductTypeIcons(ctx context.Context, in *store.UpdateProductTypeIconsInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProductTypeIcons", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProductTypeIcons indicates an expected call of UpdateProductTypeIcons.
func (mr *MockServiceMockRecorder) UpdateProductTypeIcons(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProductTypeIcons", reflect.TypeOf((*MockService)(nil).UpdateProductTypeIcons), ctx, in)
}

// UpdatePromotion mocks base method.
func (m *MockService) UpdatePromotion(ctx context.Context, in *store.UpdatePromotionInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePromotion", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePromotion indicates an expected call of UpdatePromotion.
func (mr *MockServiceMockRecorder) UpdatePromotion(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePromotion", reflect.TypeOf((*MockService)(nil).UpdatePromotion), ctx, in)
}

// UpdateSchedule mocks base method.
func (m *MockService) UpdateSchedule(ctx context.Context, in *store.UpdateScheduleInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSchedule", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSchedule indicates an expected call of UpdateSchedule.
func (mr *MockServiceMockRecorder) UpdateSchedule(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSchedule", reflect.TypeOf((*MockService)(nil).UpdateSchedule), ctx, in)
}

// UpdateScheduleThumbnails mocks base method.
func (m *MockService) UpdateScheduleThumbnails(ctx context.Context, in *store.UpdateScheduleThumbnailsInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateScheduleThumbnails", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateScheduleThumbnails indicates an expected call of UpdateScheduleThumbnails.
func (mr *MockServiceMockRecorder) UpdateScheduleThumbnails(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateScheduleThumbnails", reflect.TypeOf((*MockService)(nil).UpdateScheduleThumbnails), ctx, in)
}

// UpsertShipping mocks base method.
func (m *MockService) UpsertShipping(ctx context.Context, in *store.UpsertShippingInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertShipping", ctx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertShipping indicates an expected call of UpsertShipping.
func (mr *MockServiceMockRecorder) UpsertShipping(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertShipping", reflect.TypeOf((*MockService)(nil).UpsertShipping), ctx, in)
}
