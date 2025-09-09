package service

import (
	"net/url"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
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
				ThumbnailURL:    "https://example.com/thumbnail01.png",
				Media: entity.MultiProductMedia{
					{
						URL:         "https://example.com/thumbnail01.png",
						IsThumbnail: true,
					},
					{
						URL:         "https://example.com/thumbnail02.png",
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
				Category: types.Category{
					ID:   "category-id",
					Name: "野菜",
				},
			},
			rate: &ProductRate{
				ProductRate: types.ProductRate{
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
				Product: types.Product{
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
					ThumbnailURL:    "https://example.com/thumbnail01.png",
					Media: []*types.ProductMedia{
						{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://example.com/thumbnail02.png",
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
					Rate: &types.ProductRate{
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
				cost:       300,
				status:     ProductStatusForSale,
				media: MultiProductMedia{
					{
						ProductMedia: types.ProductMedia{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
						},
					},
					{
						ProductMedia: types.ProductMedia{
							URL:         "https://example.com/thumbnail02.png",
							IsThumbnail: false,
						},
					},
				},
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
				ThumbnailURL:    "https://example.com/thumbnail01.png",
				Media: entity.MultiProductMedia{
					{
						URL:         "https://example.com/thumbnail01.png",
						IsThumbnail: true,
					},
					{
						URL:         "https://example.com/thumbnail02.png",
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
				Product: types.Product{
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
					ThumbnailURL:    "https://example.com/thumbnail01.png",
					Media: []*types.ProductMedia{
						{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://example.com/thumbnail02.png",
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
					Rate: &types.ProductRate{
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
				cost:       300,
				status:     ProductStatusForSale,
				media: MultiProductMedia{
					{
						ProductMedia: types.ProductMedia{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
						},
					},
					{
						ProductMedia: types.ProductMedia{
							URL:         "https://example.com/thumbnail02.png",
							IsThumbnail: false,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProduct(tt.product, tt.category, tt.rate))
		})
	}
}

func TestProduct_MediaURLs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		product *Product
		expect  []string
	}{
		{
			name: "success",
			product: &Product{
				media: MultiProductMedia{
					{
						ProductMedia: types.ProductMedia{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
						},
					},
					{
						ProductMedia: types.ProductMedia{
							URL:         "https://example.com/thumbnail02.png",
							IsThumbnail: false,
						},
					},
				},
			},
			expect: []string{
				"https://example.com/thumbnail01.png",
				"https://example.com/thumbnail02.png",
			},
		},
		{
			name: "empty media",
			product: &Product{
				media: MultiProductMedia{},
			},
			expect: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.product.MediaURLs())
		})
	}
}

func TestProduct_MerchantCenterItemCondition(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		product *Product
		expect  string
	}{
		{
			name: "presale with inventory",
			product: &Product{
				Product: types.Product{
					Inventory: 100,
				},
				status: ProductStatusPresale,
			},
			expect: "preorder",
		},
		{
			name: "presale without inventory",
			product: &Product{
				Product: types.Product{
					Inventory: 0,
				},
				status: ProductStatusPresale,
			},
			expect: "out_of_stock",
		},
		{
			name: "for sale with inventory",
			product: &Product{
				Product: types.Product{
					Inventory: 100,
				},
				status: ProductStatusForSale,
			},
			expect: "in_stock",
		},
		{
			name: "for sale without inventory",
			product: &Product{
				Product: types.Product{
					Inventory: 0,
				},
				status: ProductStatusForSale,
			},
			expect: "out_of_stock",
		},
		{
			name: "out of sale",
			product: &Product{
				Product: types.Product{
					Inventory: 100,
				},
				status: ProductStatusOutOfSale,
			},
			expect: "out_of_stock",
		},
		{
			name: "unknown status",
			product: &Product{
				Product: types.Product{
					Inventory: 100,
				},
				status: ProductStatusUnknown,
			},
			expect: "new",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.product.MerchantCenterItemCondition())
		})
	}
}

func TestProduct_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		product *Product
		expect  *types.Product
	}{
		{
			name: "success",
			product: &Product{
				Product: types.Product{
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
					Media: []*types.ProductMedia{
						{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://example.com/thumbnail02.png",
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
			expect: &types.Product{
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
				Media: []*types.ProductMedia{
					{
						URL:         "https://example.com/thumbnail01.png",
						IsThumbnail: true,
					},
					{
						URL:         "https://example.com/thumbnail02.png",
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
					ThumbnailURL:    "https://example.com/thumbnail01.png",
					Media: entity.MultiProductMedia{
						{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://example.com/thumbnail02.png",
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
						ProductType: types.ProductType{
							ID:         "product-type-id",
							CategoryID: "category-id",
							Name:       "じゃがいも",
						},
					},
				},
				Categories: map[string]*Category{
					"category-id": {
						Category: types.Category{
							ID:   "category-id",
							Name: "野菜",
						},
					},
				},
				ProductRates: map[string]*ProductRate{
					"product-id": {
						ProductRate: types.ProductRate{
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
					Product: types.Product{
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
						ThumbnailURL:    "https://example.com/thumbnail01.png",
						Media: []*types.ProductMedia{
							{
								URL:         "https://example.com/thumbnail01.png",
								IsThumbnail: true,
							},
							{
								URL:         "https://example.com/thumbnail02.png",
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
						Rate: &types.ProductRate{
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
					cost:       300,
					status:     ProductStatusForSale,
					media: MultiProductMedia{
						{
							ProductMedia: types.ProductMedia{
								URL:         "https://example.com/thumbnail01.png",
								IsThumbnail: true,
							},
						},
						{
							ProductMedia: types.ProductMedia{
								URL:         "https://example.com/thumbnail02.png",
								IsThumbnail: false,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
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
					Product: types.Product{
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
						Media: []*types.ProductMedia{
							{URL: "https://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "https://example.com/thumbnail02.png", IsThumbnail: false},
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
					Product: types.Product{
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
						Media: []*types.ProductMedia{
							{URL: "https://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "https://example.com/thumbnail02.png", IsThumbnail: false},
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
					Product: types.Product{
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
						Media: []*types.ProductMedia{
							{URL: "https://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "https://example.com/thumbnail02.png", IsThumbnail: false},
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
		expect   []*types.Product
	}{
		{
			name: "success",
			products: Products{
				{
					Product: types.Product{
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
						ThumbnailURL:    "https://example.com/thumbnail01.png",
						Media: []*types.ProductMedia{
							{URL: "https://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "https://example.com/thumbnail02.png", IsThumbnail: false},
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
			expect: []*types.Product{
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
					ThumbnailURL:    "https://example.com/thumbnail01.png",
					Media: []*types.ProductMedia{
						{URL: "https://example.com/thumbnail01.png", IsThumbnail: true},
						{URL: "https://example.com/thumbnail02.png", IsThumbnail: false},
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
				URL:         "https://example.com/thumbnail01.png",
				IsThumbnail: true,
			},
			expect: &ProductMedia{
				ProductMedia: types.ProductMedia{
					URL:         "https://example.com/thumbnail01.png",
					IsThumbnail: true,
				},
			},
		},
	}
	for _, tt := range tests {
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
		expect *types.ProductMedia
	}{
		{
			name: "success",
			media: &ProductMedia{
				ProductMedia: types.ProductMedia{
					URL:         "https://example.com/thumbnail01.png",
					IsThumbnail: true,
				},
			},
			expect: &types.ProductMedia{
				URL:         "https://example.com/thumbnail01.png",
				IsThumbnail: true,
			},
		},
	}
	for _, tt := range tests {
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
					URL:         "https://example.com/thumbnail01.png",
					IsThumbnail: true,
				},
				{
					URL:         "https://example.com/thumbnail02.png",
					IsThumbnail: false,
				},
			},
			expect: MultiProductMedia{
				{
					ProductMedia: types.ProductMedia{
						URL:         "https://example.com/thumbnail01.png",
						IsThumbnail: true,
					},
				},
				{
					ProductMedia: types.ProductMedia{
						URL:         "https://example.com/thumbnail02.png",
						IsThumbnail: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewMultiProductMedia(tt.media))
		})
	}
}

func TestMultiProductMedia_URLs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		media  MultiProductMedia
		expect []string
	}{
		{
			name: "success",
			media: MultiProductMedia{
				{
					ProductMedia: types.ProductMedia{
						URL:         "https://example.com/thumbnail01.png",
						IsThumbnail: true,
					},
				},
				{
					ProductMedia: types.ProductMedia{
						URL:         "https://example.com/thumbnail02.png",
						IsThumbnail: false,
					},
				},
			},
			expect: []string{
				"https://example.com/thumbnail01.png",
				"https://example.com/thumbnail02.png",
			},
		},
		{
			name:   "empty media",
			media:  MultiProductMedia{},
			expect: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.media.URLs())
		})
	}
}

func TestMultiProductMedia_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		media  MultiProductMedia
		expect []*types.ProductMedia
	}{
		{
			name: "success",
			media: MultiProductMedia{
				{
					ProductMedia: types.ProductMedia{
						URL:         "https://example.com/thumbnail01.png",
						IsThumbnail: true,
					},
				},
				{
					ProductMedia: types.ProductMedia{
						URL:         "https://example.com/thumbnail02.png",
						IsThumbnail: false,
					},
				},
			},
			expect: []*types.ProductMedia{
				{
					URL:         "https://example.com/thumbnail01.png",
					IsThumbnail: true,
				},
				{
					URL:         "https://example.com/thumbnail02.png",
					IsThumbnail: false,
				},
			},
		},
	}
	for _, tt := range tests {
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
					ProductRate: types.ProductRate{
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
					ProductRate: types.ProductRate{
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
					ProductRate: types.ProductRate{
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
		expect []*types.ProductRate
	}{
		{
			name: "success",
			rates: ProductRates{
				{
					ProductRate: types.ProductRate{
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
			expect: []*types.ProductRate{
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

func TestNewMerchantCenterItem(t *testing.T) {
	t.Parallel()
	now := jst.Date(2025, 1, 15, 18, 30, 0, 0)
	webURL := func() *url.URL {
		u, _ := url.Parse("https://example.com")
		return u
	}
	tests := []struct {
		name   string
		params *NewMerchantCenterItemParams
		expect *MerchantCenterItem
	}{
		{
			name: "success with all fields",
			params: &NewMerchantCenterItemParams{
				Product: &Product{
					Product: types.Product{
						ID:                "product-id",
						Name:              "新鮮なじゃがいも",
						Description:       "新鮮なじゃがいもをお届けします。",
						ThumbnailURL:      "https://example.com/thumbnail01.png",
						Price:             400,
						Inventory:         100,
						Weight:            1.3,
						CategoryID:        "category-id",
						ProductTypeID:     "product-type-id",
						CoordinatorID:     "coordinator-id",
						RecommendedPoint1: "ポイント1",
						StartAt:           jst.Date(2025, 1, 1, 0, 0, 0, 0).Unix(),
						EndAt:             jst.Date(2025, 2, 1, 0, 0, 0, 0).Unix(),
					},
					cost:   300,
					status: ProductStatusForSale,
					media: MultiProductMedia{
						{
							ProductMedia: types.ProductMedia{
								URL:         "https://example.com/thumbnail01.png",
								IsThumbnail: true,
							},
						},
						{
							ProductMedia: types.ProductMedia{
								URL:         "https://example.com/thumbnail02.png",
								IsThumbnail: false,
							},
						},
					},
				},
				Coordinator: &Coordinator{
					Coordinator: types.Coordinator{
						Username: "コーディネータ名",
					},
				},
				ProductType: &ProductType{
					ProductType: types.ProductType{
						Name: "じゃがいも",
					},
				},
				Category: &Category{
					Category: types.Category{
						Name: "野菜",
					},
				},
				Now:    now,
				WebURL: webURL,
			},
			expect: &MerchantCenterItem{
				MerchantCenterItem: types.MerchantCenterItem{
					ID:                   "product-id",
					Title:                "新鮮なじゃがいも",
					Description:          "新鮮なじゃがいもをお届けします。",
					Link:                 "https://example.com/items/product-id",
					ImageLink:            "https://example.com/thumbnail01.png",
					AdditionalImageLinks: []string{"https://example.com/thumbnail01.png", "https://example.com/thumbnail02.png"},
					Availability:         "in_stock",
					AvailabilityDate:     "2025-01-01T00:00:00+09:00",
					CostOfGoodsSold:      "300 JPY",
					ExpirationDate:       "2025-02-01T00:00:00+09:00",
					Price:                "400 JPY",
					UnitPricingMeasure:   "ct",
					ProductType:          "野菜 > じゃがいも",
					Brand:                "コーディネータ名",
					Condition:            "new",
					ProductWeight:        "1 kg",
					ProductHighlight:     "ポイント1",
					ShippingWeight:       "1 kg",
					ShipsFromCountry:     "JP",
				},
			},
		},
		{
			name: "success with empty description",
			params: &NewMerchantCenterItemParams{
				Product: &Product{
					Product: types.Product{
						ID:                "product-id",
						Name:              "新鮮なじゃがいも",
						Description:       "",
						ThumbnailURL:      "https://example.com/thumbnail01.png",
						Price:             400,
						Inventory:         100,
						Weight:            1.3,
						CategoryID:        "category-id",
						ProductTypeID:     "product-type-id",
						CoordinatorID:     "coordinator-id",
						RecommendedPoint1: "ポイント1",
						StartAt:           jst.Date(2025, 1, 1, 0, 0, 0, 0).Unix(),
						EndAt:             jst.Date(2025, 2, 1, 0, 0, 0, 0).Unix(),
					},
					cost:   300,
					status: ProductStatusForSale,
					media: MultiProductMedia{
						{
							ProductMedia: types.ProductMedia{
								URL:         "https://example.com/thumbnail01.png",
								IsThumbnail: true,
							},
						},
						{
							ProductMedia: types.ProductMedia{
								URL:         "https://example.com/thumbnail02.png",
								IsThumbnail: false,
							},
						},
					},
				},
				Coordinator: &Coordinator{
					Coordinator: types.Coordinator{
						Username: "コーディネータ名",
					},
				},
				ProductType: &ProductType{
					ProductType: types.ProductType{
						Name: "じゃがいも",
					},
				},
				Category: &Category{
					Category: types.Category{
						Name: "野菜",
					},
				},
				Now:    now,
				WebURL: webURL,
			},
			expect: &MerchantCenterItem{
				MerchantCenterItem: types.MerchantCenterItem{
					ID:                   "product-id",
					Title:                "新鮮なじゃがいも",
					Description:          "新鮮なじゃがいも",
					Link:                 "https://example.com/items/product-id",
					ImageLink:            "https://example.com/thumbnail01.png",
					AdditionalImageLinks: []string{"https://example.com/thumbnail01.png", "https://example.com/thumbnail02.png"},
					Availability:         "in_stock",
					AvailabilityDate:     "2025-01-01T00:00:00+09:00",
					CostOfGoodsSold:      "300 JPY",
					ExpirationDate:       "2025-02-01T00:00:00+09:00",
					Price:                "400 JPY",
					UnitPricingMeasure:   "ct",
					ProductType:          "野菜 > じゃがいも",
					Brand:                "コーディネータ名",
					Condition:            "new",
					ProductWeight:        "1 kg",
					ProductHighlight:     "ポイント1",
					ShippingWeight:       "1 kg",
					ShipsFromCountry:     "JP",
				},
			},
		},
		{
			name: "success with out of stock",
			params: &NewMerchantCenterItemParams{
				Product: &Product{
					Product: types.Product{
						ID:            "product-id",
						Name:          "新鮮なじゃがいも",
						Description:   "新鮮なじゃがいもをお届けします。",
						ThumbnailURL:  "https://example.com/thumbnail01.png",
						Price:         400,
						Inventory:     0,
						Weight:        1.3,
						CategoryID:    "category-id",
						ProductTypeID: "product-type-id",
						CoordinatorID: "coordinator-id",
					},
					cost:   300,
					status: ProductStatusForSale,
					media:  MultiProductMedia{},
				},
				Coordinator: &Coordinator{
					Coordinator: types.Coordinator{
						Username: "コーディネータ名",
					},
				},
				ProductType: &ProductType{
					ProductType: types.ProductType{
						Name: "じゃがいも",
					},
				},
				Category: &Category{
					Category: types.Category{
						Name: "野菜",
					},
				},
				Now:    now,
				WebURL: webURL,
			},
			expect: &MerchantCenterItem{
				MerchantCenterItem: types.MerchantCenterItem{
					ID:                   "product-id",
					Title:                "新鮮なじゃがいも",
					Description:          "新鮮なじゃがいもをお届けします。",
					Link:                 "https://example.com/items/product-id",
					ImageLink:            "https://example.com/thumbnail01.png",
					AdditionalImageLinks: []string{},
					Availability:         "out_of_stock",
					AvailabilityDate:     "0001-01-01T09:00:00+09:00",
					CostOfGoodsSold:      "300 JPY",
					Price:                "400 JPY",
					UnitPricingMeasure:   "ct",
					ProductType:          "野菜 > じゃがいも",
					Brand:                "コーディネータ名",
					Condition:            "new",
					ProductWeight:        "1 kg",
					ProductHighlight:     "",
					ShippingWeight:       "1 kg",
					ShipsFromCountry:     "JP",
				},
			},
		},
		{
			name: "success with nil values",
			params: &NewMerchantCenterItemParams{
				Product: &Product{
					Product: types.Product{
						ID:            "product-id",
						Name:          "新鮮なじゃがいも",
						Description:   "新鮮なじゃがいもをお届けします。",
						ThumbnailURL:  "https://example.com/thumbnail01.png",
						Price:         400,
						Inventory:     100,
						Weight:        1.3,
						CategoryID:    "category-id",
						ProductTypeID: "product-type-id",
						CoordinatorID: "coordinator-id",
					},
					cost:   0,
					status: ProductStatusUnknown,
					media:  MultiProductMedia{},
				},
				Coordinator: nil,
				ProductType: nil,
				Category:    nil,
				Now:         now,
				WebURL:      webURL,
			},
			expect: &MerchantCenterItem{
				MerchantCenterItem: types.MerchantCenterItem{
					ID:                   "product-id",
					Title:                "新鮮なじゃがいも",
					Description:          "新鮮なじゃがいもをお届けします。",
					Link:                 "https://example.com/items/product-id",
					ImageLink:            "https://example.com/thumbnail01.png",
					AdditionalImageLinks: []string{},
					Availability:         "new",
					AvailabilityDate:     "0001-01-01T09:00:00+09:00",
					CostOfGoodsSold:      "0 JPY",
					Price:                "400 JPY",
					UnitPricingMeasure:   "ct",
					ProductType:          "",
					Brand:                "",
					Condition:            "new",
					ProductWeight:        "1 kg",
					ProductHighlight:     "",
					ShippingWeight:       "1 kg",
					ShipsFromCountry:     "JP",
				},
			},
		},
		{
			name: "success with very long description",
			params: &NewMerchantCenterItemParams{
				Product: &Product{
					Product: types.Product{
						ID:            "product-id",
						Name:          "新鮮なじゃがいも",
						Description:   string(make([]byte, 5010)),
						ThumbnailURL:  "https://example.com/thumbnail01.png",
						Price:         400,
						Inventory:     100,
						Weight:        1.3,
						StartAt:       now.Add(-time.Hour * 24 * 30).Unix(), // 30日前
						EndAt:         now.Add(time.Hour * 24 * 14).Unix(),  // 14日後
						CategoryID:    "category-id",
						ProductTypeID: "product-type-id",
						CoordinatorID: "coordinator-id",
					},
					cost:   300,
					status: ProductStatusForSale,
					media: MultiProductMedia{
						{
							ProductMedia: types.ProductMedia{
								URL:         "https://example.com/thumbnail01.png",
								IsThumbnail: true,
							},
						},
					},
				},
				Coordinator: &Coordinator{
					Coordinator: types.Coordinator{
						Username: "コーディネータ名",
					},
				},
				ProductType: &ProductType{
					ProductType: types.ProductType{
						Name: "じゃがいも",
					},
				},
				Category: &Category{
					Category: types.Category{
						Name: "野菜",
					},
				},
				Now:    now,
				WebURL: webURL,
			},
			expect: &MerchantCenterItem{
				MerchantCenterItem: types.MerchantCenterItem{
					ID:                   "product-id",
					Title:                "新鮮なじゃがいも",
					Description:          string(make([]byte, 4997)) + "...",
					Link:                 "https://example.com/items/product-id",
					ImageLink:            "https://example.com/thumbnail01.png",
					AdditionalImageLinks: []string{"https://example.com/thumbnail01.png"},
					Availability:         "in_stock",
					AvailabilityDate:     "2024-12-16T18:30:00+09:00",
					CostOfGoodsSold:      "300 JPY",
					ExpirationDate:       "2025-01-29T18:30:00+09:00",
					Price:                "400 JPY",
					UnitPricingMeasure:   "ct",
					ProductType:          "野菜 > じゃがいも",
					Brand:                "コーディネータ名",
					Condition:            "new",
					ProductWeight:        "1 kg",
					ProductHighlight:     "",
					ShippingWeight:       "1 kg",
					ShipsFromCountry:     "JP",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewMerchantCenterItem(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestNewMerchantCenterItems(t *testing.T) {
	t.Parallel()
	now := jst.Date(2025, 1, 15, 18, 30, 0, 0)
	webURL := func() *url.URL {
		u, _ := url.Parse("https://example.com")
		return u
	}
	tests := []struct {
		name   string
		params *NewMerchantCenterItemsParams
		expect MerchantCenterItems
	}{
		{
			name: "success",
			params: &NewMerchantCenterItemsParams{
				Products: Products{
					{
						Product: types.Product{
							ID:            "product-id-1",
							Name:          "新鮮なじゃがいも",
							Description:   "新鮮なじゃがいもをお届けします。",
							ThumbnailURL:  "https://example.com/thumbnail01.png",
							Price:         400,
							Inventory:     100,
							Weight:        1.3,
							CategoryID:    "category-id",
							ProductTypeID: "product-type-id",
							CoordinatorID: "coordinator-id",
						},
						cost:   300,
						status: ProductStatusForSale,
						media: MultiProductMedia{
							{
								ProductMedia: types.ProductMedia{
									URL:         "https://example.com/thumbnail01.png",
									IsThumbnail: true,
								},
							},
						},
					},
					{
						Product: types.Product{
							ID:            "product-id-2",
							Name:          "新鮮な人参",
							Description:   "新鮮な人参をお届けします。",
							ThumbnailURL:  "https://example.com/thumbnail02.png",
							Price:         300,
							Inventory:     0,
							Weight:        0.5,
							CategoryID:    "category-id",
							ProductTypeID: "product-type-id-2",
							CoordinatorID: "coordinator-id-2",
						},
						cost:   150,
						status: ProductStatusForSale,
						media: MultiProductMedia{
							{
								ProductMedia: types.ProductMedia{
									URL:         "https://example.com/thumbnail02.png",
									IsThumbnail: true,
								},
							},
							{
								ProductMedia: types.ProductMedia{
									URL:         "https://example.com/thumbnail03.png",
									IsThumbnail: false,
								},
							},
						},
					},
				},
				Coordinators: map[string]*Coordinator{
					"coordinator-id": {
						Coordinator: types.Coordinator{
							Username: "コーディネータ名1",
						},
					},
					"coordinator-id-2": {
						Coordinator: types.Coordinator{
							Username: "コーディネータ名2",
						},
					},
				},
				Details: &ProductDetailsParams{
					ProductTypes: map[string]*ProductType{
						"product-type-id": {
							ProductType: types.ProductType{
								Name: "じゃがいも",
							},
						},
						"product-type-id-2": {
							ProductType: types.ProductType{
								Name: "人参",
							},
						},
					},
					Categories: map[string]*Category{
						"category-id": {
							Category: types.Category{
								Name: "野菜",
							},
						},
					},
				},
				Now:    now,
				WebURL: webURL,
			},
			expect: MerchantCenterItems{
				{
					MerchantCenterItem: types.MerchantCenterItem{
						ID:                   "product-id-1",
						Title:                "新鮮なじゃがいも",
						Description:          "新鮮なじゃがいもをお届けします。",
						Link:                 "https://example.com/items/product-id-1",
						ImageLink:            "https://example.com/thumbnail01.png",
						AdditionalImageLinks: []string{"https://example.com/thumbnail01.png"},
						Availability:         "in_stock",
						AvailabilityDate:     "0001-01-01T09:00:00+09:00",
						CostOfGoodsSold:      "300 JPY",
						Price:                "400 JPY",
						UnitPricingMeasure:   "ct",
						ProductType:          "野菜 > じゃがいも",
						Brand:                "コーディネータ名1",
						Condition:            "new",
						ProductWeight:        "1 kg",
						ProductHighlight:     "",
						ShippingWeight:       "1 kg",
						ShipsFromCountry:     "JP",
					},
				},
				{
					MerchantCenterItem: types.MerchantCenterItem{
						ID:                   "product-id-2",
						Title:                "新鮮な人参",
						Description:          "新鮮な人参をお届けします。",
						Link:                 "https://example.com/items/product-id-2",
						ImageLink:            "https://example.com/thumbnail02.png",
						AdditionalImageLinks: []string{"https://example.com/thumbnail02.png", "https://example.com/thumbnail03.png"},
						Availability:         "out_of_stock",
						AvailabilityDate:     "0001-01-01T09:00:00+09:00",
						CostOfGoodsSold:      "150 JPY",
						Price:                "300 JPY",
						UnitPricingMeasure:   "ct",
						ProductType:          "野菜 > 人参",
						Brand:                "コーディネータ名2",
						Condition:            "new",
						ProductWeight:        "0 kg",
						ProductHighlight:     "",
						ShippingWeight:       "0 kg",
						ShipsFromCountry:     "JP",
					},
				},
			},
		},
		{
			name: "success with empty products",
			params: &NewMerchantCenterItemsParams{
				Products:     Products{},
				Coordinators: map[string]*Coordinator{},
				Details: &ProductDetailsParams{
					ProductTypes: map[string]*ProductType{},
					Categories:   map[string]*Category{},
				},
				WebURL: func() *url.URL {
					u, _ := url.Parse("https://example.com")
					return u
				},
			},
			expect: MerchantCenterItems{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewMerchantCenterItems(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestMerchantCenterItem_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		item   *MerchantCenterItem
		expect *types.MerchantCenterItem
	}{
		{
			name: "success",
			item: &MerchantCenterItem{
				MerchantCenterItem: types.MerchantCenterItem{
					ID:                    "product-id",
					Title:                 "新鮮なじゃがいも",
					Description:           "新鮮なじゃがいもをお届けします。",
					Link:                  "https://example.com/products/product-id",
					ImageLink:             "https://example.com/thumbnail01.png",
					Condition:             "new",
					Availability:          "in_stock",
					Price:                 "400 JPY",
					Brand:                 "コーディネータ名",
					GoogleProductCategory: "野菜",
					ProductType:           "じゃがいも",
					ShippingWeight:        "1 g",
				},
			},
			expect: &types.MerchantCenterItem{
				ID:                    "product-id",
				Title:                 "新鮮なじゃがいも",
				Description:           "新鮮なじゃがいもをお届けします。",
				Link:                  "https://example.com/products/product-id",
				ImageLink:             "https://example.com/thumbnail01.png",
				Condition:             "new",
				Availability:          "in_stock",
				Price:                 "400 JPY",
				Brand:                 "コーディネータ名",
				GoogleProductCategory: "野菜",
				ProductType:           "じゃがいも",
				ShippingWeight:        "1 g",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.item.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestMerchantCenterItems_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		items  MerchantCenterItems
		expect []*types.MerchantCenterItem
	}{
		{
			name: "success",
			items: MerchantCenterItems{
				{
					MerchantCenterItem: types.MerchantCenterItem{
						ID:                    "product-id-1",
						Title:                 "新鮮なじゃがいも",
						Description:           "新鮮なじゃがいもをお届けします。",
						Link:                  "https://example.com/products/product-id-1",
						ImageLink:             "https://example.com/thumbnail01.png",
						Condition:             "new",
						Availability:          "in_stock",
						Price:                 "400 JPY",
						Brand:                 "コーディネータ名1",
						GoogleProductCategory: "野菜",
						ProductType:           "じゃがいも",
						ShippingWeight:        "1 g",
					},
				},
				{
					MerchantCenterItem: types.MerchantCenterItem{
						ID:                    "product-id-2",
						Title:                 "新鮮な人参",
						Description:           "新鮮な人参をお届けします。",
						Link:                  "https://example.com/products/product-id-2",
						ImageLink:             "https://example.com/thumbnail02.png",
						Condition:             "new",
						Availability:          "out_of_stock",
						Price:                 "300 JPY",
						Brand:                 "コーディネータ名2",
						GoogleProductCategory: "野菜",
						ProductType:           "人参",
						ShippingWeight:        "0 g",
					},
				},
			},
			expect: []*types.MerchantCenterItem{
				{
					ID:                    "product-id-1",
					Title:                 "新鮮なじゃがいも",
					Description:           "新鮮なじゃがいもをお届けします。",
					Link:                  "https://example.com/products/product-id-1",
					ImageLink:             "https://example.com/thumbnail01.png",
					Condition:             "new",
					Availability:          "in_stock",
					Price:                 "400 JPY",
					Brand:                 "コーディネータ名1",
					GoogleProductCategory: "野菜",
					ProductType:           "じゃがいも",
					ShippingWeight:        "1 g",
				},
				{
					ID:                    "product-id-2",
					Title:                 "新鮮な人参",
					Description:           "新鮮な人参をお届けします。",
					Link:                  "https://example.com/products/product-id-2",
					ImageLink:             "https://example.com/thumbnail02.png",
					Condition:             "new",
					Availability:          "out_of_stock",
					Price:                 "300 JPY",
					Brand:                 "コーディネータ名2",
					GoogleProductCategory: "野菜",
					ProductType:           "人参",
					ShippingWeight:        "0 g",
				},
			},
		},
		{
			name:   "success with empty items",
			items:  MerchantCenterItems{},
			expect: []*types.MerchantCenterItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.items.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
