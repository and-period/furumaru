package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestProductStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ProductStatus
		expect ProductStatus
	}{
		{
			name:   "private",
			status: entity.ProductStatusPrivate,
			expect: ProductStatusUnknown,
		},
		{
			name:   "presale",
			status: entity.ProductStatusPresale,
			expect: ProductStatusPresale,
		},
		{
			name:   "for sale",
			status: entity.ProductStatusForSale,
			expect: ProductStatusForSale,
		},
		{
			name:   "out of sale",
			status: entity.ProductStatusOutOfSale,
			expect: ProductStatusOutOfSale,
		},
		{
			name:   "unknown",
			status: entity.ProductStatusUnknown,
			expect: ProductStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProductStatus(tt.status))
		})
	}
}

func TestStorageMethodType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		storageMethodType entity.StorageMethodType
		expect            StorageMethodType
	}{
		{
			name:              "success to normal",
			storageMethodType: entity.StorageMethodTypeNormal,
			expect:            StorageMethodTypeNormal,
		},
		{
			name:              "success to cook dark",
			storageMethodType: entity.StorageMethodTypeCoolDark,
			expect:            StorageMethodTypeCoolDark,
		},
		{
			name:              "success to refrigerated",
			storageMethodType: entity.StorageMethodTypeRefrigerated,
			expect:            StorageMethodTypeRefrigerated,
		},
		{
			name:              "success to frozen",
			storageMethodType: entity.StorageMethodTypeFrozen,
			expect:            StorageMethodTypeFrozen,
		},
		{
			name:              "success to unknown",
			storageMethodType: entity.StorageMethodTypeUnknown,
			expect:            StorageMethodTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewStorageMethodType(tt.storageMethodType))
		})
	}
}

func TestDeliveryType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		deliveryType entity.DeliveryType
		expect       DeliveryType
	}{
		{
			name:         "success to normal",
			deliveryType: entity.DeliveryTypeNormal,
			expect:       DeliveryTypeNormal,
		},
		{
			name:         "success to frozen",
			deliveryType: entity.DeliveryTypeFrozen,
			expect:       DeliveryTypeFrozen,
		},
		{
			name:         "success to refrigerated",
			deliveryType: entity.DeliveryTypeRefrigerated,
			expect:       DeliveryTypeRefrigerated,
		},
		{
			name:         "success to unknown",
			deliveryType: entity.DeliveryTypeUnknown,
			expect:       DeliveryTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewDeliveryType(tt.deliveryType))
		})
	}
}

func TestDeliveryType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		deliveryType DeliveryType
		expect       int32
	}{
		{
			name:         "success",
			deliveryType: DeliveryTypeNormal,
			expect:       1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.deliveryType.Response())
		})
	}
}

func TestProductWeight(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		weight     int64
		weightUnit entity.WeightUnit
		expect     float64
	}{
		{
			name:       "from 130g to 0.1kg",
			weight:     130,
			weightUnit: entity.WeightUnitGram,
			expect:     0.1,
		},
		{
			name:       "from 1300g to 1.3kg",
			weight:     1300,
			weightUnit: entity.WeightUnitGram,
			expect:     1.3,
		},
		{
			name:       "from 13500g to 13.5kg",
			weight:     13500,
			weightUnit: entity.WeightUnitGram,
			expect:     13.5,
		},
		{
			name:       "from 15kg to 15.0kg",
			weight:     15,
			weightUnit: entity.WeightUnitKilogram,
			expect:     15.0,
		},
		{
			name:       "unknown weight unit",
			weight:     1500,
			weightUnit: entity.WeightUnitUnknown,
			expect:     0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductWeight(tt.weight, tt.weightUnit)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProduct(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		product  *entity.Product
		category *Category
		rate     *ProductRate
		expect   *Product
	}{
		{
			name: "success",
			product: &entity.Product{
				ID:              "product-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Status:          entity.ProductStatusForSale,
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          1300,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
				Media: entity.MultiProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail01.png",
						IsThumbnail: true,
					},
					{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
					},
				},
				RecommendedPoints:    []string{"ポイント1", "ポイント2", "ポイント3"},
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefecture:     "滋賀県",
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              jst.Date(2022, 1, 1, 0, 0, 0, 0),
				EndAt:                jst.Date(2022, 1, 1, 0, 0, 0, 0),
				ProductRevision: entity.ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     400,
					Cost:      300,
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			category: &Category{
				Category: response.Category{
					ID:   "category-id",
					Name: "野菜",
				},
			},
			rate: &ProductRate{
				ProductRate: response.ProductRate{
					Count:   4,
					Average: 2.5,
					Detail: map[int64]int64{
						1: 2,
						2: 0,
						3: 1,
						4: 0,
						5: 1,
					},
				},
				productID: "product-id",
			},
			expect: &Product{
				Product: response.Product{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					CategoryID:      "category-id",
					ProductTypeID:   "product-type-id",
					ProductTagIDs:   []string{"product-tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Status:          int32(ProductStatusForSale),
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
					Media: []*response.ProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://and-period.jp/thumbnail02.png",
							IsThumbnail: false,
						},
					},
					Price:             400,
					RecommendedPoint1: "ポイント1",
					RecommendedPoint2: "ポイント2",
					RecommendedPoint3: "ポイント3",
					StorageMethodType: int32(StorageMethodTypeNormal),
					DeliveryType:      int32(DeliveryTypeNormal),
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					Rate: &response.ProductRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					StartAt: 1640962800,
					EndAt:   1640962800,
				},
				revisionID: 1,
			},
		},
		{
			name: "success without additional values",
			product: &entity.Product{
				ID:              "product-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
				Name:            "新鮮なじゃがいも",
				Status:          entity.ProductStatusForSale,
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          1300,
				WeightUnit:      entity.WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
				Media: entity.MultiProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail01.png",
						IsThumbnail: true,
					},
					{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
					},
				},
				RecommendedPoints:    []string{"ポイント1", "ポイント2", "ポイント3"},
				StorageMethodType:    entity.StorageMethodTypeNormal,
				DeliveryType:         entity.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefecture:     "滋賀県",
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              jst.Date(2022, 1, 1, 0, 0, 0, 0),
				EndAt:                jst.Date(2022, 1, 1, 0, 0, 0, 0),
				ProductRevision: entity.ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     400,
					Cost:      300,
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			category: nil,
			rate:     nil,
			expect: &Product{
				Product: response.Product{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					CategoryID:      "",
					ProductTypeID:   "product-type-id",
					ProductTagIDs:   []string{"product-tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Status:          int32(ProductStatusForSale),
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
					Media: []*response.ProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://and-period.jp/thumbnail02.png",
							IsThumbnail: false,
						},
					},
					Price:             400,
					RecommendedPoint1: "ポイント1",
					RecommendedPoint2: "ポイント2",
					RecommendedPoint3: "ポイント3",
					StorageMethodType: int32(StorageMethodTypeNormal),
					DeliveryType:      int32(DeliveryTypeNormal),
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					Rate: &response.ProductRate{
						Count:   0,
						Average: 0.0,
						Detail: map[int64]int64{
							1: 0,
							2: 0,
							3: 0,
							4: 0,
							5: 0,
						},
					},
					StartAt: 1640962800,
					EndAt:   1640962800,
				},
				revisionID: 1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProduct(tt.product, tt.category, tt.rate))
		})
	}
}

func TestProduct_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		product *Product
		expect  *response.Product
	}{
		{
			name: "success",
			product: &Product{
				Product: response.Product{
					ID:              "product-id",
					ProductTypeID:   "product-type-id",
					CategoryID:      "category-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Status:          int32(ProductStatusForSale),
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: []*response.ProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://and-period.jp/thumbnail02.png",
							IsThumbnail: false,
						},
					},
					Price:            400,
					DeliveryType:     int32(DeliveryTypeNormal),
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "滋賀県",
					OriginCity:       "彦根市",
					StartAt:          1640962800,
					EndAt:            1640962800,
				},
			},
			expect: &response.Product{
				ID:              "product-id",
				ProductTypeID:   "product-type-id",
				CategoryID:      "category-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Status:          int32(ProductStatusForSale),
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*response.ProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail01.png",
						IsThumbnail: true,
					},
					{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
					},
				},
				Price:            400,
				DeliveryType:     int32(DeliveryTypeNormal),
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
				StartAt:          1640962800,
				EndAt:            1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.product.Response())
		})
	}
}

func TestProducts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products entity.Products
		params   *ProductDetailsParams
		expect   Products
	}{
		{
			name: "success",
			products: entity.Products{
				{
					ID:              "product-id",
					TypeID:          "product-type-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          entity.ProductStatusForSale,
					Inventory:       100,
					Weight:          1300,
					WeightUnit:      entity.WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
					Media: entity.MultiProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://and-period.jp/thumbnail02.png",
							IsThumbnail: false,
						},
					},
					DeliveryType:         entity.DeliveryTypeNormal,
					Box60Rate:            50,
					Box80Rate:            40,
					Box100Rate:           30,
					OriginPrefecture:     "滋賀県",
					OriginPrefectureCode: 25,
					OriginCity:           "彦根市",
					StartAt:              jst.Date(2022, 1, 1, 0, 0, 0, 0),
					EndAt:                jst.Date(2022, 1, 1, 0, 0, 0, 0),
					ProductRevision: entity.ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
						CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
						UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			params: &ProductDetailsParams{
				ProductTypes: map[string]*ProductType{
					"product-type-id": {
						ProductType: response.ProductType{
							ID:         "product-type-id",
							CategoryID: "category-id",
							Name:       "じゃがいも",
						},
					},
				},
				Categories: map[string]*Category{
					"category-id": {
						Category: response.Category{
							ID:   "category-id",
							Name: "野菜",
						},
					},
				},
				ProductRates: map[string]*ProductRate{
					"product-id": {
						ProductRate: response.ProductRate{
							Count:   4,
							Average: 2.5,
							Detail: map[int64]int64{
								1: 2,
								2: 0,
								3: 1,
								4: 0,
								5: 1,
							},
						},
						productID: "product-id",
					},
				},
			},
			expect: Products{
				{
					Product: response.Product{
						ID:              "product-id",
						ProductTypeID:   "product-type-id",
						CategoryID:      "category-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Status:          int32(ProductStatusForSale),
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
						Media: []*response.ProductMedia{
							{
								URL:         "https://and-period.jp/thumbnail01.png",
								IsThumbnail: true,
							},
							{
								URL:         "https://and-period.jp/thumbnail02.png",
								IsThumbnail: false,
							},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "滋賀県",
						OriginCity:       "彦根市",
						Rate: &response.ProductRate{
							Count:   4,
							Average: 2.5,
							Detail: map[int64]int64{
								1: 2,
								2: 0,
								3: 1,
								4: 0,
								5: 1,
							},
						},
						StartAt: 1640962800,
						EndAt:   1640962800,
					},
					revisionID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProducts(tt.products, tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProducts_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   []string
	}{
		{
			name: "success",
			products: Products{
				{
					Product: response.Product{
						ID:              "product-id",
						ProductTypeID:   "product-type-id",
						CategoryID:      "category-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:        400,
						DeliveryType: int32(DeliveryTypeNormal),
						Box60Rate:    50,
						Box80Rate:    40,
						Box100Rate:   30,
						OriginCity:   "彦根市",
					},
					revisionID: 1,
				},
			},
			expect: []string{"product-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.IDs())
		})
	}
}

func TestProducts_MapByRevision(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   map[int64]*Product
	}{
		{
			name: "success",
			products: Products{
				{
					Product: response.Product{
						ID:              "product-id",
						ProductTypeID:   "product-type-id",
						CategoryID:      "category-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:        400,
						DeliveryType: int32(DeliveryTypeNormal),
						Box60Rate:    50,
						Box80Rate:    40,
						Box100Rate:   30,
						OriginCity:   "彦根市",
					},
					revisionID: 1,
				},
			},
			expect: map[int64]*Product{
				1: {
					Product: response.Product{
						ID:              "product-id",
						ProductTypeID:   "product-type-id",
						CategoryID:      "category-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:        400,
						DeliveryType: int32(DeliveryTypeNormal),
						Box60Rate:    50,
						Box80Rate:    40,
						Box100Rate:   30,
						OriginCity:   "彦根市",
					},
					revisionID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.MapByRevision())
		})
	}
}

func TestProducts_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   []*response.Product
	}{
		{
			name: "success",
			products: Products{
				{
					Product: response.Product{
						ID:              "product-id",
						ProductTypeID:   "product-type-id",
						CategoryID:      "category-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Status:          int32(ProductStatusForSale),
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "滋賀県",
						OriginCity:       "彦根市",
						StartAt:          1640962800,
						EndAt:            1640962800,
					},
					revisionID: 1,
				},
			},
			expect: []*response.Product{
				{
					ID:              "product-id",
					ProductTypeID:   "product-type-id",
					CategoryID:      "category-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Status:          int32(ProductStatusForSale),
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					ThumbnailURL:    "https://and-period.jp/thumbnail01.png",
					Media: []*response.ProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					Price:            400,
					DeliveryType:     int32(DeliveryTypeNormal),
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "滋賀県",
					OriginCity:       "彦根市",
					StartAt:          1640962800,
					EndAt:            1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.Response())
		})
	}
}

func TestProductMedia(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		media  *entity.ProductMedia
		expect *ProductMedia
	}{
		{
			name: "success",
			media: &entity.ProductMedia{
				URL:         "https://and-period.jp/thumbnail01.png",
				IsThumbnail: true,
			},
			expect: &ProductMedia{
				ProductMedia: response.ProductMedia{
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProductMedia(tt.media))
		})
	}
}

func TestProductMedia_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		media  *ProductMedia
		expect *response.ProductMedia
	}{
		{
			name: "success",
			media: &ProductMedia{
				ProductMedia: response.ProductMedia{
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
				},
			},
			expect: &response.ProductMedia{
				URL:         "https://and-period.jp/thumbnail01.png",
				IsThumbnail: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.media.Response())
		})
	}
}

func TestMultiProductMedia(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		media  entity.MultiProductMedia
		expect MultiProductMedia
	}{
		{
			name: "success",
			media: entity.MultiProductMedia{
				{
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
				},
				{
					URL:         "https://and-period.jp/thumbnail02.png",
					IsThumbnail: false,
				},
			},
			expect: MultiProductMedia{
				{
					ProductMedia: response.ProductMedia{
						URL:         "https://and-period.jp/thumbnail01.png",
						IsThumbnail: true,
					},
				},
				{
					ProductMedia: response.ProductMedia{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewMultiProductMedia(tt.media))
		})
	}
}

func TestMultiProductMedia_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		media  MultiProductMedia
		expect []*response.ProductMedia
	}{
		{
			name: "success",
			media: MultiProductMedia{
				{
					ProductMedia: response.ProductMedia{
						URL:         "https://and-period.jp/thumbnail01.png",
						IsThumbnail: true,
					},
				},
				{
					ProductMedia: response.ProductMedia{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
					},
				},
			},
			expect: []*response.ProductMedia{
				{
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
				},
				{
					URL:         "https://and-period.jp/thumbnail02.png",
					IsThumbnail: false,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.media.Response())
		})
	}
}

func TestProductRates(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews entity.AggregatedProductReviews
		expect  ProductRates
	}{
		{
			name: "success",
			reviews: entity.AggregatedProductReviews{
				{
					ProductID: "product-id",
					Count:     4,
					Average:   2.5,
					Rate1:     2,
					Rate2:     0,
					Rate3:     1,
					Rate4:     0,
					Rate5:     1,
				},
			},
			expect: ProductRates{
				{
					ProductRate: response.ProductRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					productID: "product-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductRates(tt.reviews)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductRates_MapByProductID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rates  ProductRates
		expect map[string]*ProductRate
	}{
		{
			name: "success",
			rates: ProductRates{
				{
					ProductRate: response.ProductRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					productID: "product-id",
				},
			},
			expect: map[string]*ProductRate{
				"product-id": {
					ProductRate: response.ProductRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					productID: "product-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.rates.MapByProductID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductRates_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rates  ProductRates
		expect []*response.ProductRate
	}{
		{
			name: "success",
			rates: ProductRates{
				{
					ProductRate: response.ProductRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					productID: "product-id",
				},
			},
			expect: []*response.ProductRate{
				{
					Count:   4,
					Average: 2.5,
					Detail: map[int64]int64{
						1: 2,
						2: 0,
						3: 1,
						4: 0,
						5: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.rates.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
