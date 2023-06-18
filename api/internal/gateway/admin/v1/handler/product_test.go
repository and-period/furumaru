package handler

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFilterAccessProduct(t *testing.T) {
	t.Parallel()

	listProducersIn := &user.ListProducersInput{
		CoordinatorID: "coordinator-id",
	}
	getProducersIn := &user.MultiGetProducersInput{
		ProducerIDs: []string{"producer-id"},
	}
	producers := uentity.Producers{
		{
			Admin: uentity.Admin{
				ID:            "producer-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id",
			CoordinatorID: "coordinator-id",
			Username:      "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	productIn := &store.GetProductInput{
		ProductID: "product-id",
	}
	product := &sentity.Product{
		ID:              "product-id",
		TypeID:          "product-type-id",
		ProducerID:      "producer-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Inventory:       100,
		Weight:          1300,
		WeightUnit:      entity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		Media: entity.MultiProductMedia{
			{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
			{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
		},
		Price:            400,
		DeliveryType:     entity.DeliveryTypeNormal,
		Box60Rate:        50,
		Box80Rate:        40,
		Box100Rate:       30,
		OriginPrefecture: "滋賀県",
		OriginCity:       "彦根市",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: []string{"category-id"},
	}
	categories := sentity.Categories{
		{
			ID:        "category-id",
			Name:      "野菜",
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:         "product-type-id",
			Name:       "じゃがいも",
			CategoryID: "category-id",
			IconURL:    "https://and-period.jp/icon.png",
			CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		expect  int
	}{
		{
			name:    "administrator success",
			setup:   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			expect:  http.StatusOK,
		},
		{
			name: "coordinator success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), listProducersIn).Return(producers, int64(1), nil)
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(product, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), getProducersIn).Return(producers, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect:  http.StatusOK,
		},
		{
			name: "coordinator forbidden",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				producers := uentity.Producers{}
				mocks.user.EXPECT().ListProducers(gomock.Any(), listProducersIn).Return(producers, int64(1), nil)
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(product, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), getProducersIn).Return(producers, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect:  http.StatusForbidden,
		},
		{
			name: "coordinator failed to get producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), listProducersIn).Return(nil, int64(0), assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect:  http.StatusInternalServerError,
		},
		{
			name: "coordinator failed to get product",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), listProducersIn).Return(producers, int64(1), nil)
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(nil, assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect:  http.StatusInternalServerError,
		},
		{
			name: "producer success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(product, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), getProducersIn).Return(producers, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleProducer), withAdminID("producer-id")},
			expect:  http.StatusOK,
		},
		{
			name: "producer forbidden",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(product, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), getProducersIn).Return(producers, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleProducer)},
			expect:  http.StatusForbidden,
		},
		{
			name: "producer failed to get product",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(nil, assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleProducer), withAdminID("coordinator-id")},
			expect:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const route, path = "/products/:productId", "/products/product-id"
			testMiddleware(t, tt.setup, route, path, tt.expect, func(h *handler) gin.HandlerFunc {
				return h.filterAccessProduct
			}, tt.options...)
		})
	}
}

func TestListProducts(t *testing.T) {
	t.Parallel()

	productsIn := &store.ListProductsInput{
		Name:       "じゃがいも",
		ProducerID: "producer-id",
		Limit:      20,
		Offset:     0,
		Orders:     []*store.ListProductsOrder{},
	}
	products := sentity.Products{
		{
			ID:              "product-id",
			TypeID:          "product-type-id",
			ProducerID:      "producer-id",
			Name:            "新鮮なじゃがいも",
			Description:     "新鮮なじゃがいもをお届けします。",
			Public:          true,
			Inventory:       100,
			Weight:          1300,
			WeightUnit:      entity.WeightUnitGram,
			Item:            1,
			ItemUnit:        "袋",
			ItemDescription: "1袋あたり100gのじゃがいも",
			Media: entity.MultiProductMedia{
				{
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
					Images: common.Images{
						{URL: "https://and-period.jp/thumbnail01_240.png", Size: common.ImageSizeSmall},
						{URL: "https://and-period.jp/thumbnail01_675.png", Size: common.ImageSizeMedium},
						{URL: "https://and-period.jp/thumbnail01_900.png", Size: common.ImageSizeLarge},
					},
				},
				{
					URL:         "https://and-period.jp/thumbnail02.png",
					IsThumbnail: false,
					Images:      common.Images{},
				},
			},
			Price:            400,
			DeliveryType:     entity.DeliveryTypeNormal,
			Box60Rate:        50,
			Box80Rate:        40,
			Box100Rate:       30,
			OriginPrefecture: "滋賀県",
			OriginCity:       "彦根市",
			CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	psIn := &user.ListProducersInput{
		CoordinatorID: "coordinator-id",
	}
	producersIn := &user.MultiGetProducersInput{
		ProducerIDs: []string{"producer-id"},
	}
	producers := uentity.Producers{
		{
			Admin: uentity.Admin{
				ID:            "producer-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id",
			CoordinatorID: "coordinator-id",
			Username:      "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: []string{"category-id"},
	}
	categories := sentity.Categories{
		{
			ID:        "category-id",
			Name:      "野菜",
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	productTypes := sentity.ProductTypes{
		{
			ID:         "product-type-id",
			Name:       "じゃがいも",
			CategoryID: "category-id",
			IconURL:    "https://and-period.jp/icon.png",
			CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		query   string
		expect  *testResponse
	}{
		{
			name: "success administrator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListProducts(gomock.Any(), productsIn).Return(products, int64(1), nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			query:   "?name=じゃがいも&producerId=producer-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductsResponse{
					Products: []*response.Product{
						{
							ID:              "product-id",
							TypeID:          "product-type-id",
							TypeName:        "じゃがいも",
							TypeIconURL:     "https://and-period.jp/icon.png",
							CategoryID:      "category-id",
							CategoryName:    "野菜",
							ProducerID:      "producer-id",
							StoreName:       "&.農園",
							Name:            "新鮮なじゃがいも",
							Description:     "新鮮なじゃがいもをお届けします。",
							Public:          true,
							Inventory:       100,
							Weight:          1.3,
							ItemUnit:        "袋",
							ItemDescription: "1袋あたり100gのじゃがいも",
							Media: []*response.ProductMedia{
								{
									URL:         "https://and-period.jp/thumbnail01.png",
									IsThumbnail: true,
									Images: []*response.Image{
										{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(service.ImageSizeSmall)},
										{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(service.ImageSizeMedium)},
										{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(service.ImageSizeLarge)},
									},
								},
								{
									URL:         "https://and-period.jp/thumbnail02.png",
									IsThumbnail: false,
									Images:      []*response.Image{},
								},
							},
							Price:            400,
							DeliveryType:     1,
							Box60Rate:        50,
							Box80Rate:        40,
							Box100Rate:       30,
							OriginPrefecture: "滋賀県",
							OriginCity:       "彦根市",
							CreatedAt:        1640962800,
							UpdatedAt:        1640962800,
						},
					},
					Total: 1,
				},
			},
		},
		{
			name: "success coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				productsIn := &store.ListProductsInput{
					ProducerIDs: []string{"producer-id"},
					Limit:       20,
					Offset:      0,
					Orders:      []*store.ListProductsOrder{},
				}
				mocks.user.EXPECT().ListProducers(gomock.Any(), psIn).Return(producers, int64(1), nil)
				mocks.store.EXPECT().ListProducts(gomock.Any(), productsIn).Return(sentity.Products{}, int64(0), nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductsResponse{
					Products: []*response.Product{},
					Total:    0,
				},
			},
		},
		{
			name:  "invalid limit",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid offset",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid orders",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?orders=name,public,inventory,price,originPrefecture,originCity,createdAt,updatedAt,other",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to list producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), psIn).Return(nil, int64(0), assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			query:   "",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to list products",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListProducts(gomock.Any(), productsIn).Return(nil, int64(0), assert.AnError)
			},
			query: "?name=じゃがいも&coordinatorId=coordinator-id&producerId=producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListProducts(gomock.Any(), productsIn).Return(products, int64(1), nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(nil, assert.AnError)
			},
			query: "?name=じゃがいも&coordinatorId=coordinator-id&producerId=producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get categories",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListProducts(gomock.Any(), productsIn).Return(products, int64(1), nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
			},
			query: "?name=じゃがいも&coordinatorId=coordinator-id&producerId=producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get product types",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListProducts(gomock.Any(), productsIn).Return(products, int64(1), nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
			},
			query: "?name=じゃがいも&coordinatorId=coordinator-id&producerId=producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/products%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path, tt.options...)
		})
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	productIn := &store.GetProductInput{
		ProductID: "product-id",
	}
	producersIn := &user.MultiGetProducersInput{
		ProducerIDs: []string{"producer-id"},
	}
	categoriesIn := &store.MultiGetCategoriesInput{
		CategoryIDs: []string{"category-id"},
	}
	productTypesIn := &store.MultiGetProductTypesInput{
		ProductTypeIDs: []string{"product-type-id"},
	}
	product := &sentity.Product{
		ID:              "product-id",
		TypeID:          "product-type-id",
		ProducerID:      "producer-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Inventory:       100,
		Weight:          1300,
		WeightUnit:      entity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		Media: entity.MultiProductMedia{
			{
				URL:         "https://and-period.jp/thumbnail01.png",
				IsThumbnail: true,
				Images: common.Images{
					{URL: "https://and-period.jp/thumbnail01_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/thumbnail01_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/thumbnail01_900.png", Size: common.ImageSizeLarge},
				},
			},
			{
				URL:         "https://and-period.jp/thumbnail02.png",
				IsThumbnail: false,
				Images:      common.Images{},
			},
		},
		Price:            400,
		DeliveryType:     entity.DeliveryTypeNormal,
		Box60Rate:        50,
		Box80Rate:        40,
		Box100Rate:       30,
		OriginPrefecture: "滋賀県",
		OriginCity:       "彦根市",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	producers := uentity.Producers{
		{
			Admin: uentity.Admin{
				ID:            "producer-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id",
			CoordinatorID: "coordinator-id",
			Username:      "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
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
			ID:         "product-type-id",
			Name:       "じゃがいも",
			IconURL:    "https://and-period.jp/icon.png",
			CategoryID: "category-id",
			CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		productID string
		expect    *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(product, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
			},
			productID: "product-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductResponse{
					Product: &response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						StoreName:       "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{
								URL:         "https://and-period.jp/thumbnail01.png",
								IsThumbnail: true,
								Images: []*response.Image{
									{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(service.ImageSizeSmall)},
									{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(service.ImageSizeMedium)},
									{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(service.ImageSizeLarge)},
								},
							},
							{
								URL:         "https://and-period.jp/thumbnail02.png",
								IsThumbnail: false,
								Images:      []*response.Image{},
							},
						},
						Price:            400,
						DeliveryType:     1,
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "滋賀県",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
		{
			name: "failed to get product",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(nil, assert.AnError)
			},
			productID: "product-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(product, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(categories, nil)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(nil, assert.AnError)
			},
			productID: "product-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get product type",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(product, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
			},
			productID: "product-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get category",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetProduct(gomock.Any(), productIn).Return(product, nil)
				mocks.store.EXPECT().MultiGetProductTypes(gomock.Any(), productTypesIn).Return(productTypes, nil)
				mocks.store.EXPECT().MultiGetCategories(gomock.Any(), categoriesIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetProducers(gomock.Any(), producersIn).Return(producers, nil)
			},
			productID: "product-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/products/%s"
			path := fmt.Sprintf(format, tt.productID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateProduct(t *testing.T) {
	t.Parallel()

	producerIn := &user.GetProducerInput{
		ProducerID: "producer-id",
	}
	categoryIn := &store.GetCategoryInput{
		CategoryID: "category-id",
	}
	productTypeIn := &store.GetProductTypeInput{
		ProductTypeID: "product-type-id",
	}
	productIn := &store.CreateProductInput{
		ProducerID:      "producer-id",
		TypeID:          "product-type-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Inventory:       100,
		Weight:          1300,
		WeightUnit:      entity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		Media: []*store.CreateProductMedia{
			{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
			{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
		},
		Price:            400,
		DeliveryType:     entity.DeliveryTypeNormal,
		Box60Rate:        50,
		Box80Rate:        40,
		Box100Rate:       30,
		OriginPrefecture: "滋賀県",
		OriginCity:       "彦根市",
	}
	producer := &uentity.Producer{
		Admin: uentity.Admin{
			ID:            "producer-id",
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-producer@and-period.jp",
		},
		AdminID:      "producer-id",
		Username:     "&.農園",
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
		HeaderURL:    "https://and-period.jp/header.png",
		PhoneNumber:  "+819012345678",
		PostalCode:   "1000014",
		Prefecture:   "東京都",
		City:         "千代田区",
		CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	producersIn := &user.ListProducersInput{
		CoordinatorID: "coordinator-id",
	}
	producers := uentity.Producers{producer}
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
	product := &sentity.Product{
		ID:              "product-id",
		TypeID:          "product-type-id",
		ProducerID:      "producer-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Inventory:       100,
		Weight:          1300,
		WeightUnit:      entity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		Media: entity.MultiProductMedia{
			{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
			{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
		},
		Price:            400,
		DeliveryType:     entity.DeliveryTypeNormal,
		Box60Rate:        50,
		Box80Rate:        40,
		Box100Rate:       30,
		OriginPrefecture: "滋賀県",
		OriginCity:       "彦根市",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		req     *request.CreateProductRequest
		expect  *testResponse
	}{
		{
			name: "success administrator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.store.EXPECT().GetProductType(gomock.Any(), productTypeIn).Return(productType, nil)
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(category, nil)
				mocks.media.EXPECT().UploadProductMedia(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, in *media.UploadFileInput) (string, error) {
						return in.URL, nil
					}).Times(2)
				mocks.store.EXPECT().CreateProduct(gomock.Any(), productIn).Return(product, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductResponse{
					Product: &response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						StoreName:       "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{
								URL:         "https://and-period.jp/thumbnail01.png",
								IsThumbnail: true,
								Images:      []*response.Image{},
							},
							{
								URL:         "https://and-period.jp/thumbnail02.png",
								IsThumbnail: false,
								Images:      []*response.Image{},
							},
						},
						Price:            400,
						DeliveryType:     1,
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "滋賀県",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
		{
			name: "success coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), producersIn).Return(producers, int64(1), nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.store.EXPECT().GetProductType(gomock.Any(), productTypeIn).Return(productType, nil)
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(category, nil)
				mocks.media.EXPECT().UploadProductMedia(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, in *media.UploadFileInput) (string, error) {
						return in.URL, nil
					}).Times(2)
				mocks.store.EXPECT().CreateProduct(gomock.Any(), productIn).Return(product, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProductResponse{
					Product: &response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						StoreName:       "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{
								URL:         "https://and-period.jp/thumbnail01.png",
								IsThumbnail: true,
								Images:      []*response.Image{},
							},
							{
								URL:         "https://and-period.jp/thumbnail02.png",
								IsThumbnail: false,
								Images:      []*response.Image{},
							},
						},
						Price:            400,
						DeliveryType:     1,
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "滋賀県",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
		{
			name: "failed to list producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), producersIn).Return(nil, int64(0), assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to not contain producer id",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), producersIn).Return(uentity.Producers{}, int64(0), nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusForbidden,
			},
		},
		{
			name: "not found dependencies",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, exception.ErrNotFound)
				mocks.store.EXPECT().GetProductType(gomock.Any(), productTypeIn).Return(productType, nil)
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(category, nil)
			},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().GetProductType(gomock.Any(), productTypeIn).Return(productType, nil)
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(category, nil)
			},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get product type",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.store.EXPECT().GetProductType(gomock.Any(), productTypeIn).Return(nil, assert.AnError)
			},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get category",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.store.EXPECT().GetProductType(gomock.Any(), productTypeIn).Return(productType, nil)
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(nil, assert.AnError)
			},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to upload product media",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.store.EXPECT().GetProductType(gomock.Any(), productTypeIn).Return(productType, nil)
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(category, nil)
				mocks.media.EXPECT().UploadProductMedia(gomock.Any(), gomock.Any()).Return("", assert.AnError).AnyTimes()
			},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to create product",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.store.EXPECT().GetProductType(gomock.Any(), productTypeIn).Return(productType, nil)
				mocks.store.EXPECT().GetCategory(gomock.Any(), categoryIn).Return(category, nil)
				mocks.media.EXPECT().UploadProductMedia(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, in *media.UploadFileInput) (string, error) {
						return in.URL, nil
					}).Times(2)
				mocks.store.EXPECT().CreateProduct(gomock.Any(), productIn).Return(nil, assert.AnError)
			},
			req: &request.CreateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/products"
			testPost(t, tt.setup, tt.expect, path, tt.req, tt.options...)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	t.Parallel()

	in := &store.UpdateProductInput{
		ProductID:       "product-id",
		ProducerID:      "producer-id",
		TypeID:          "product-type-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Inventory:       100,
		Weight:          1300,
		WeightUnit:      entity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		Media: []*store.UpdateProductMedia{
			{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
			{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
		},
		Price:            400,
		DeliveryType:     entity.DeliveryTypeNormal,
		Box60Rate:        50,
		Box80Rate:        40,
		Box100Rate:       30,
		OriginPrefecture: "滋賀県",
		OriginCity:       "彦根市",
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		productID string
		req       *request.UpdateProductRequest
		expect    *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProductMedia(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, in *media.UploadFileInput) (string, error) {
						return in.URL, nil
					}).Times(2)
				mocks.store.EXPECT().UpdateProduct(gomock.Any(), in).Return(nil)
			},
			productID: "product-id",
			req: &request.UpdateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.UpdateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update product",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProductMedia(gomock.Any(), gomock.Any()).Return("", assert.AnError).AnyTimes()
			},
			productID: "product-id",
			req: &request.UpdateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.UpdateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to update product",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().UploadProductMedia(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, in *media.UploadFileInput) (string, error) {
						return in.URL, nil
					}).Times(2)
				mocks.store.EXPECT().UpdateProduct(gomock.Any(), in).Return(assert.AnError)
			},
			productID: "product-id",
			req: &request.UpdateProductRequest{
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*request.UpdateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:            400,
				DeliveryType:     1,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/products/%s"
			path := fmt.Sprintf(format, tt.productID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	t.Parallel()
	in := &store.DeleteProductInput{
		ProductID: "product-id",
	}
	tests := []struct {
		name      string
		setup     func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		productID string
		expect    *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeleteProduct(gomock.Any(), in).Return(nil)
			},
			productID: "product-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete product",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeleteProduct(gomock.Any(), in).Return(assert.AnError)
			},
			productID: "product-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/products/%s"
			path := fmt.Sprintf(format, tt.productID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
