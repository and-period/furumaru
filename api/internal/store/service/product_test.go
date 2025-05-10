package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListProducts(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	params := &database.ListProductsParams{
		Name:          "みかん",
		CoordinatorID: "coordinator-id",
		ProducerID:    "",
		ProducerIDs:   []string{"producer-id"},
		EndAtGte:      jst.Date(2022, 6, 28, 18, 30, 0, 0),
		Limit:         30,
		Offset:        0,
		Orders: []*database.ListProductsOrder{
			{Key: database.ListProductsOrderByName, OrderByASC: true},
		},
	}
	products := entity.Products{
		{
			ID:              "product-id",
			TypeID:          "type-id",
			TagIDs:          []string{"tag-id"},
			CoordinatorID:   "coordinator-id",
			ProducerID:      "producer-id",
			Name:            "新鮮なじゃがいも",
			Description:     "新鮮なじゃがいもをお届けします。",
			Public:          true,
			Inventory:       100,
			Weight:          100,
			WeightUnit:      entity.WeightUnitGram,
			Item:            1,
			ItemUnit:        "袋",
			ItemDescription: "1袋あたり100gのじゃがいも",
			Media: entity.MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
			},
			ExpirationDate:    7,
			StorageMethodType: entity.StorageMethodTypeNormal,
			DeliveryType:      entity.DeliveryTypeNormal,
			Box60Rate:         50,
			Box80Rate:         40,
			Box100Rate:        30,
			OriginPrefecture:  "滋賀県",
			OriginCity:        "彦根市",
			ProductRevision: entity.ProductRevision{
				ID:        1,
				ProductID: "product-id",
				Price:     400,
				Cost:      300,
				CreatedAt: now,
				UpdatedAt: now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListProductsInput
		expect      entity.Products
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().List(gomock.Any(), params).Return(products, nil)
				mocks.db.Product.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListProductsInput{
				Name:             "みかん",
				CoordinatorID:    "coordinator-id",
				ProducerID:       "",
				ProducerIDs:      []string{"producer-id"},
				ExcludeOutOfSale: true,
				Limit:            30,
				Offset:           0,
				Orders: []*store.ListProductsOrder{
					{Key: store.ListProductsOrderByName, OrderByASC: true},
				},
			},
			expect:      products,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListProductsInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Product.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListProductsInput{
				Name:             "みかん",
				CoordinatorID:    "coordinator-id",
				ProducerID:       "",
				ProducerIDs:      []string{"producer-id"},
				ExcludeOutOfSale: true,
				Limit:            30,
				Offset:           0,
				Orders: []*store.ListProductsOrder{
					{Key: store.ListProductsOrderByName, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().List(gomock.Any(), params).Return(products, nil)
				mocks.db.Product.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListProductsInput{
				Name:             "みかん",
				CoordinatorID:    "coordinator-id",
				ProducerID:       "",
				ProducerIDs:      []string{"producer-id"},
				ExcludeOutOfSale: true,
				Limit:            30,
				Offset:           0,
				Orders: []*store.ListProductsOrder{
					{Key: store.ListProductsOrderByName, OrderByASC: true},
				},
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListProducts(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}, withNow(now)))
	}
}

func TestMultiGetProducts(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	products := entity.Products{
		{
			ID:              "product-id",
			TypeID:          "type-id",
			TagIDs:          []string{"tag-id"},
			CoordinatorID:   "coordinator-id",
			ProducerID:      "producer-id",
			Name:            "新鮮なじゃがいも",
			Description:     "新鮮なじゃがいもをお届けします。",
			Public:          true,
			Inventory:       100,
			Weight:          100,
			WeightUnit:      entity.WeightUnitGram,
			Item:            1,
			ItemUnit:        "袋",
			ItemDescription: "1袋あたり100gのじゃがいも",
			Media: entity.MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
			},
			ExpirationDate:    7,
			StorageMethodType: entity.StorageMethodTypeNormal,
			DeliveryType:      entity.DeliveryTypeNormal,
			Box60Rate:         50,
			Box80Rate:         40,
			Box100Rate:        30,
			OriginPrefecture:  "滋賀県",
			OriginCity:        "彦根市",
			ProductRevision: entity.ProductRevision{
				ID:        1,
				ProductID: "product-id",
				Price:     400,
				Cost:      300,
				CreatedAt: now,
				UpdatedAt: now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetProductsInput
		expect    entity.Products
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products, nil)
			},
			input: &store.MultiGetProductsInput{
				ProductIDs: []string{"product-id"},
			},
			expect:    products,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetProductsInput{
				ProductIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetProductsInput{
				ProductIDs: []string{"product-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetProducts(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetProductsByRevision(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 5, 2, 18, 30, 0, 0)
	products := entity.Products{
		{
			ID:              "product-id",
			TypeID:          "type-id",
			TagIDs:          []string{"tag-id"},
			CoordinatorID:   "coordinator-id",
			ProducerID:      "producer-id",
			Name:            "新鮮なじゃがいも",
			Description:     "新鮮なじゃがいもをお届けします。",
			Public:          true,
			Inventory:       100,
			Weight:          100,
			WeightUnit:      entity.WeightUnitGram,
			Item:            1,
			ItemUnit:        "袋",
			ItemDescription: "1袋あたり100gのじゃがいも",
			Media: entity.MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
			},
			ExpirationDate:    7,
			StorageMethodType: entity.StorageMethodTypeNormal,
			DeliveryType:      entity.DeliveryTypeNormal,
			Box60Rate:         50,
			Box80Rate:         40,
			Box100Rate:        30,
			OriginPrefecture:  "滋賀県",
			OriginCity:        "彦根市",
			ProductRevision: entity.ProductRevision{
				ID:        1,
				ProductID: "product-id",
				Price:     400,
				Cost:      300,
				CreatedAt: now,
				UpdatedAt: now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetProductsByRevisionInput
		expect    entity.Products
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().MultiGetByRevision(ctx, []int64{1}).Return(products, nil)
			},
			input: &store.MultiGetProductsByRevisionInput{
				ProductRevisionIDs: []int64{1},
			},
			expect:    products,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetProductsByRevisionInput{
				ProductRevisionIDs: []int64{0},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().MultiGetByRevision(ctx, []int64{1}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetProductsByRevisionInput{
				ProductRevisionIDs: []int64{1},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetProductsByRevision(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetProduct(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	product := &entity.Product{
		ID:              "product-id",
		TypeID:          "type-id",
		TagIDs:          []string{"tag-id"},
		CoordinatorID:   "coordinator-id",
		ProducerID:      "producer-id",
		Name:            "新鮮なじゃがいも",
		Description:     "新鮮なじゃがいもをお届けします。",
		Public:          true,
		Inventory:       100,
		Weight:          100,
		WeightUnit:      entity.WeightUnitGram,
		Item:            1,
		ItemUnit:        "袋",
		ItemDescription: "1袋あたり100gのじゃがいも",
		Media: entity.MultiProductMedia{
			{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
			{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
		},
		ExpirationDate:    7,
		StorageMethodType: entity.StorageMethodTypeNormal,
		DeliveryType:      entity.DeliveryTypeNormal,
		Box60Rate:         50,
		Box80Rate:         40,
		Box100Rate:        30,
		OriginPrefecture:  "滋賀県",
		OriginCity:        "彦根市",
		ProductRevision: entity.ProductRevision{
			ID:        1,
			ProductID: "product-id",
			Price:     400,
			Cost:      300,
			CreatedAt: now,
			UpdatedAt: now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetProductInput
		expect    *entity.Product
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(product, nil)
			},
			input: &store.GetProductInput{
				ProductID: "product-id",
			},
			expect:    product,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetProductInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Get(ctx, "product-id").Return(nil, assert.AnError)
			},
			input: &store.GetProductInput{
				ProductID: "product-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetProduct(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateProduct(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	shop := &entity.Shop{
		ID:            "shop-id",
		Name:          "じゃがいも農園",
		CoordinatorID: "coordinator-id",
		ProducerIDs:   []string{"producer-id"},
	}
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		AdminID: "coordinator-id",
	}
	producerIn := &user.GetProducerInput{
		ProducerID: "producer-id",
	}
	producer := &uentity.Producer{
		AdminID: "producer-id",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateProductInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, product *entity.Product) error {
						expect := &entity.Product{
							ID:              product.ID, // ignore
							ShopID:          "shop-id",
							TypeID:          "product-type-id",
							TagIDs:          []string{"product-tag-id"},
							CoordinatorID:   "coordinator-id",
							ProducerID:      "producer-id",
							Name:            "新鮮なじゃがいも",
							Description:     "新鮮なじゃがいもをお届けします。",
							Public:          true,
							Inventory:       100,
							Weight:          100,
							WeightUnit:      entity.WeightUnitGram,
							Item:            1,
							ItemUnit:        "袋",
							ItemDescription: "1袋あたり100gのじゃがいも",
							Media: entity.MultiProductMedia{
								{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
								{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
							},
							ExpirationDate:       7,
							StorageMethodType:    entity.StorageMethodTypeNormal,
							DeliveryType:         entity.DeliveryTypeNormal,
							Box60Rate:            50,
							Box80Rate:            40,
							Box100Rate:           30,
							OriginPrefecture:     "滋賀県",
							OriginPrefectureCode: 25,
							OriginCity:           "彦根市",
							StartAt:              now.AddDate(0, -1, 0),
							EndAt:                now.AddDate(0, 1, 0),
							ProductRevision: entity.ProductRevision{
								ProductID: product.ID,
								Price:     400,
								Cost:      300,
							},
						}
						assert.Equal(t, expect, product)
						return nil
					})
			},
			input: &store.CreateProductInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateProductInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid media format",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.CreateProductInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: true},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "not found coordinator or producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(nil, database.ErrNotFound)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, exception.ErrNotFound)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, exception.ErrNotFound)
			},
			input: &store.CreateProductInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(nil, assert.AnError)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
			},
			input: &store.CreateProductInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
			},
			input: &store.CreateProductInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, assert.AnError)
			},
			input: &store.CreateProductInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to new product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
			},
			input: &store.CreateProductInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: -1,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.db.Product.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateProductInput{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.CreateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateProduct(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateProduct(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateProductInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().
					Update(ctx, "product-id", gomock.Any()).
					DoAndReturn(func(ctx context.Context, productID string, params *database.UpdateProductParams) error {
						expect := &database.UpdateProductParams{
							TypeID:          "product-type-id",
							TagIDs:          []string{"product-tag-id"},
							Name:            "新鮮なじゃがいも",
							Description:     "新鮮なじゃがいもをお届けします。",
							Public:          true,
							Inventory:       100,
							Weight:          100,
							WeightUnit:      entity.WeightUnitGram,
							Item:            1,
							ItemUnit:        "袋",
							ItemDescription: "1袋あたり100gのじゃがいも",
							Media: entity.MultiProductMedia{
								{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
								{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
							},
							Price:                400,
							Cost:                 300,
							ExpirationDate:       7,
							StorageMethodType:    entity.StorageMethodTypeNormal,
							DeliveryType:         entity.DeliveryTypeNormal,
							Box60Rate:            50,
							Box80Rate:            40,
							Box100Rate:           30,
							OriginPrefectureCode: 25,
							OriginCity:           "彦根市",
							StartAt:              now.AddDate(0, -1, 0),
							EndAt:                now.AddDate(0, 1, 0),
						}
						assert.Equal(t, expect, params)
						return nil
					})
			},
			input: &store.UpdateProductInput{
				ProductID:       "product-id",
				TagIDs:          []string{"product-tag-id"},
				TypeID:          "product-type-id",
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.UpdateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateProductInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid prefecture code",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.UpdateProductInput{
				ProductID:       "product-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.UpdateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: -1,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid media format",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.UpdateProductInput{
				ProductID:       "product-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.UpdateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: true},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Product.EXPECT().Update(ctx, "product-id", gomock.Any()).Return(assert.AnError)
			},
			input: &store.UpdateProductInput{
				ProductID:       "product-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*store.UpdateProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateProduct(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteProduct(t *testing.T) {
	t.Parallel()

	videosIn := &media.ListProductVideosInput{
		ProductID: "product-id",
	}
	videos := mentity.Videos{}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteProductInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.media.EXPECT().ListProductVideos(ctx, videosIn).Return(videos, nil)
				mocks.db.Product.EXPECT().Delete(ctx, "product-id").Return(nil)
			},
			input: &store.DeleteProductInput{
				ProductID: "product-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteProductInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list product videos",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.media.EXPECT().ListProductVideos(ctx, videosIn).Return(nil, assert.AnError)
			},
			input: &store.DeleteProductInput{
				ProductID: "product-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "product has videos",
			setup: func(ctx context.Context, mocks *mocks) {
				videos := mentity.Videos{{ID: "video-id"}}
				mocks.media.EXPECT().ListProductVideos(ctx, videosIn).Return(videos, nil)
			},
			input: &store.DeleteProductInput{
				ProductID: "product-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to delete product",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.media.EXPECT().ListProductVideos(ctx, videosIn).Return(videos, nil)
				mocks.db.Product.EXPECT().Delete(ctx, "product-id").Return(assert.AnError)
			},
			input: &store.DeleteProductInput{
				ProductID: "product-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteProduct(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
