package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListProductTypes(t *testing.T) {
	t.Parallel()

	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: []string{"category-id"},
	}
	typesIn := &store.ListProductTypesInput{
		Name:       "いも",
		CategoryID: "category-id",
		Limit:      20,
		Offset:     0,
		Orders:     []*store.ListProductTypesOrder{},
	}
	categories := sentity.Categories{
		{
			ID:        "category-id",
			Name:      "野菜",
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:         "product-type-id01",
			Name:       "じゃがいも",
			IconURL:    "https://and-period.jp/icon.png",
			CategoryID: "category-id",
			CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			ID:         "product-type-id02",
			Name:       "さつまいも",
			IconURL:    "https://and-period.jp/icon.png",
			CategoryID: "category-id",
			CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		categoryID string
		query      string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListProductTypes(gomock.Any(), typesIn).Return(productTypes, int64(2), nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
			},
			categoryID: "category-id",
			query:      "?name=いも",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductTypesResponse{
					ProductTypes: []*response.ProductType{
						{
							ID:           "product-type-id01",
							Name:         "じゃがいも",
							IconURL:      "https://and-period.jp/icon.png",
							CategoryID:   "category-id",
							CategoryName: "野菜",
							CreatedAt:    1640962800,
							UpdatedAt:    1640962800,
						},
						{
							ID:           "product-type-id02",
							Name:         "さつまいも",
							IconURL:      "https://and-period.jp/icon.png",
							CategoryID:   "category-id",
							CategoryName: "野菜",
							CreatedAt:    1640962800,
							UpdatedAt:    1640962800,
						},
					},
					Total: 2,
				},
			},
		},
		{
			name: "success empty",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				productTypes := sentity.ProductTypes{}
				mocks.store.EXPECT().ListProductTypes(gomock.Any(), typesIn).Return(productTypes, int64(0), nil)
			},
			categoryID: "category-id",
			query:      "?name=いも",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductTypesResponse{
					ProductTypes: []*response.ProductType{},
					Total:        0,
				},
			},
		},
		{
			name: "success without category id",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				typesIn := &store.ListProductTypesInput{
					Name:   "いも",
					Limit:  20,
					Offset: 0,
					Orders: []*store.ListProductTypesOrder{},
				}
				mocks.store.EXPECT().ListProductTypes(gomock.Any(), typesIn).Return(productTypes, int64(2), nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
			},
			categoryID: "-",
			query:      "?name=いも",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductTypesResponse{
					ProductTypes: []*response.ProductType{
						{
							ID:           "product-type-id01",
							Name:         "じゃがいも",
							IconURL:      "https://and-period.jp/icon.png",
							CategoryID:   "category-id",
							CategoryName: "野菜",
							CreatedAt:    1640962800,
							UpdatedAt:    1640962800,
						},
						{
							ID:           "product-type-id02",
							Name:         "さつまいも",
							IconURL:      "https://and-period.jp/icon.png",
							CategoryID:   "category-id",
							CategoryName: "野菜",
							CreatedAt:    1640962800,
							UpdatedAt:    1640962800,
						},
					},
					Total: 2,
				},
			},
		},
		{
			name:       "invalid limit",
			setup:      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			categoryID: "category-id",
			query:      "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:       "invalid offset",
			setup:      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			categoryID: "category-id",
			query:      "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid orders",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?orders=name,other",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get product types",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListProductTypes(gomock.Any(), typesIn).Return(nil, int64(0), assert.AnError)
			},
			categoryID: "category-id",
			query:      "?name=いも",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get categories",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListProductTypes(gomock.Any(), typesIn).Return(productTypes, int64(0), nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(nil, assert.AnError)
			},
			categoryID: "category-id",
			query:      "?name=いも",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/categories/%s/product-types%s"
			path := fmt.Sprintf(format, tt.categoryID, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateProductType(t *testing.T) {
	t.Parallel()

	categoryIn := &store.GetCategoryInput{
		CategoryID: "category-id",
	}
	typeIn := &store.CreateProductTypeInput{
		Name:       "じゃがいも",
		IconURL:    "https://and-period.jp/icon.png",
		CategoryID: "category-id",
	}
	category := &sentity.Category{
		ID:        "category-id",
		Name:      "野菜",
		CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	productType := &sentity.ProductType{
		ID:         "product-type-id",
		Name:       "じゃがいも",
		IconURL:    "https://and-period.jp/icon.png",
		CategoryID: "category-id",
		CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		categoryID string
		req        *request.CreateProductTypeRequest
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(category, nil)
				mocks.store.EXPECT().CreateProductType(gomock.Any(), typeIn).Return(productType, nil)
			},
			categoryID: "category-id",
			req: &request.CreateProductTypeRequest{
				Name:    "じゃがいも",
				IconURL: "https://and-period.jp/icon.png",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductTypeResponse{
					ProductType: &response.ProductType{
						ID:           "product-type-id",
						Name:         "じゃがいも",
						IconURL:      "https://and-period.jp/icon.png",
						CategoryID:   "category-id",
						CategoryName: "野菜",
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
		{
			name: "failed to get category",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(nil, assert.AnError)
			},
			categoryID: "category-id",
			req: &request.CreateProductTypeRequest{
				Name:    "じゃがいも",
				IconURL: "https://and-period.jp/icon.png",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to create product type",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(category, nil)
				mocks.store.EXPECT().CreateProductType(gomock.Any(), typeIn).Return(nil, assert.AnError)
			},
			categoryID: "category-id",
			req: &request.CreateProductTypeRequest{
				Name:    "じゃがいも",
				IconURL: "https://and-period.jp/icon.png",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/categories/%s/product-types"
			path := fmt.Sprintf(format, tt.categoryID)
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateProductType(t *testing.T) {
	t.Parallel()

	in := &store.UpdateProductTypeInput{
		ProductTypeID: "product-type-id",
		Name:          "じゃがいも",
		IconURL:       "https://and-period.jp/icon.png",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		categoryID    string
		productTypeID string
		req           *request.UpdateProductTypeRequest
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdateProductType(gomock.Any(), in).Return(nil)
			},
			categoryID:    "category-id",
			productTypeID: "product-type-id",
			req: &request.UpdateProductTypeRequest{
				Name:    "じゃがいも",
				IconURL: "https://and-period.jp/icon.png",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update product type",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdateProductType(gomock.Any(), in).Return(assert.AnError)
			},
			categoryID:    "category-id",
			productTypeID: "product-type-id",
			req: &request.UpdateProductTypeRequest{
				Name:    "じゃがいも",
				IconURL: "https://and-period.jp/icon.png",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/categories/%s/product-types/%s"
			path := fmt.Sprintf(format, tt.categoryID, tt.productTypeID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestDeleteProductType(t *testing.T) {
	t.Parallel()

	in := &store.DeleteProductTypeInput{
		ProductTypeID: "product-type-id",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		categoryID    string
		productTypeID string
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeleteProductType(gomock.Any(), in).Return(nil)
			},
			categoryID:    "category-id",
			productTypeID: "product-type-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete product type",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeleteProductType(gomock.Any(), in).Return(assert.AnError)
			},
			categoryID:    "category-id",
			productTypeID: "product-type-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/categories/%s/product-types/%s"
			path := fmt.Sprintf(format, tt.categoryID, tt.productTypeID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
