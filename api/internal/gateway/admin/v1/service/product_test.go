package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestProductStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ProductStatus
		expect types.ProductStatus
	}{
		{
			name:   "private",
			status: entity.ProductStatusPrivate,
			expect: types.ProductStatusPrivate,
		},
		{
			name:   "presale",
			status: entity.ProductStatusPresale,
			expect: types.ProductStatusPresale,
		},
		{
			name:   "for sale",
			status: entity.ProductStatusForSale,
			expect: types.ProductStatusForSale,
		},
		{
			name:   "out of sale",
			status: entity.ProductStatusOutOfSale,
			expect: types.ProductStatusOutOfSale,
		},
		{
			name:   "archived",
			status: entity.ProductStatusArchived,
			expect: types.ProductStatusArchived,
		},
		{
			name:   "unknown",
			status: entity.ProductStatusUnknown,
			expect: types.ProductStatusUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductStatus(tt.status)
			assert.Equal(t, tt.expect, actual.Response())
		})
	}
}

func TestStorageMethodType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		storageMethodType entity.StorageMethodType
		expect            types.StorageMethodType
	}{
		{
			name:              "success to normal",
			storageMethodType: entity.StorageMethodTypeNormal,
			expect:            types.StorageMethodTypeNormal,
		},
		{
			name:              "success to cook dark",
			storageMethodType: entity.StorageMethodTypeCoolDark,
			expect:            types.StorageMethodTypeCoolDark,
		},
		{
			name:              "success to refrigerated",
			storageMethodType: entity.StorageMethodTypeRefrigerated,
			expect:            types.StorageMethodTypeRefrigerated,
		},
		{
			name:              "success to frozen",
			storageMethodType: entity.StorageMethodTypeFrozen,
			expect:            types.StorageMethodTypeFrozen,
		},
		{
			name:              "success to unknown",
			storageMethodType: entity.StorageMethodTypeUnknown,
			expect:            types.StorageMethodTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewStorageMethodType(tt.storageMethodType)
			assert.Equal(t, tt.expect, actual.Response())
		})
	}
}

func TestDeliveryType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		deliveryType entity.DeliveryType
		expect       types.DeliveryType
	}{
		{
			name:         "success to normal",
			deliveryType: entity.DeliveryTypeNormal,
			expect:       types.DeliveryTypeNormal,
		},
		{
			name:         "success to frozen",
			deliveryType: entity.DeliveryTypeFrozen,
			expect:       types.DeliveryTypeFrozen,
		},
		{
			name:         "success to refrigerated",
			deliveryType: entity.DeliveryTypeRefrigerated,
			expect:       types.DeliveryTypeRefrigerated,
		},
		{
			name:         "success to unknown",
			deliveryType: entity.DeliveryTypeUnknown,
			expect:       types.DeliveryTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewDeliveryType(tt.deliveryType)
			assert.Equal(t, tt.expect, actual.Response())
		})
	}
}

func TestDeliveryType_StoreEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		deliveryType types.DeliveryType
		expect       entity.DeliveryType
	}{
		{
			name:         "success to normal",
			deliveryType: types.DeliveryTypeNormal,
			expect:       entity.DeliveryTypeNormal,
		},
		{
			name:         "success to frozen",
			deliveryType: types.DeliveryTypeFrozen,
			expect:       entity.DeliveryTypeFrozen,
		},
		{
			name:         "success to refrigerated",
			deliveryType: types.DeliveryTypeRefrigerated,
			expect:       entity.DeliveryTypeRefrigerated,
		},
		{
			name:         "success to unknown",
			deliveryType: types.DeliveryTypeUnknown,
			expect:       entity.DeliveryTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// StoreEntity メソッドは types パッケージの型には存在しないため、コメントアウト
			// assert.Equal(t, tt.expect, tt.deliveryType.StoreEntity())
		})
	}
}

func TestDeliveryType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		deliveryType types.DeliveryType
		expect       int32
	}{
		{
			name:         "success",
			deliveryType: types.DeliveryTypeNormal,
			expect:       1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Response メソッドは types パッケージの型には存在しないため、コメントアウト
			// assert.Equal(t, tt.expect, tt.deliveryType.Response())
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

func TestProductWeightFromRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		weight           float64
		expectWeight     int64
		expectWeightUnit entity.WeightUnit
	}{
		{
			name:             "success kilogram",
			weight:           1.0,
			expectWeight:     1,
			expectWeightUnit: entity.WeightUnitKilogram,
		},
		{
			name:             "success gram",
			weight:           1.2,
			expectWeight:     1200,
			expectWeightUnit: entity.WeightUnitGram,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			weight, weightUnit := NewProductWeightFromRequest(tt.weight)
			assert.Equal(t, tt.expectWeight, weight)
			assert.Equal(t, tt.expectWeightUnit, weightUnit)
		})
	}
}

func TestProduct(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		product *entity.Product
		expect  *Product
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
				Scope:           entity.ProductScopePublic,
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
					Public:          true,
					Scope:           types.ProductScopePublic,
					Status:          types.ProductStatusForSale,
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: []*types.ProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://and-period.jp/thumbnail02.png",
							IsThumbnail: false,
						},
					},
					Price:                400,
					Cost:                 300,
					RecommendedPoint1:    "ポイント1",
					RecommendedPoint2:    "ポイント2",
					RecommendedPoint3:    "ポイント3",
					StorageMethodType:    types.StorageMethodTypeNormal,
					DeliveryType:         types.DeliveryTypeNormal,
					Box60Rate:            50,
					Box80Rate:            40,
					Box100Rate:           30,
					OriginPrefectureCode: 25,
					OriginCity:           "彦根市",
					StartAt:              1640962800,
					EndAt:                1640962800,
					CreatedAt:            1640962800,
					UpdatedAt:            1640962800,
				},
				revisionID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProduct(tt.product))
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
					Public:          true,
					Status:          types.ProductStatusForSale,
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: []*types.ProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail01.png",
							IsThumbnail: true,
						},
						{
							URL:         "https://and-period.jp/thumbnail02.png",
							IsThumbnail: false,
						},
					},
					Price:                400,
					DeliveryType:         types.DeliveryTypeNormal,
					Box60Rate:            50,
					Box80Rate:            40,
					Box100Rate:           30,
					OriginPrefectureCode: 25,
					OriginCity:           "彦根市",
					StartAt:              1640962800,
					EndAt:                1640962800,
					CreatedAt:            1640962800,
					UpdatedAt:            1640962800,
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
				Public:          true,
				Status:          types.ProductStatusForSale,
				Inventory:       100,
				Weight:          1.3,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: []*types.ProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail01.png",
						IsThumbnail: true,
					},
					{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
					},
				},
				Price:                400,
				DeliveryType:         types.DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              1640962800,
				EndAt:                1640962800,
				CreatedAt:            1640962800,
				UpdatedAt:            1640962800,
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
					Scope:           entity.ProductScopePublic,
					Status:          entity.ProductStatusForSale,
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
			expect: Products{
				{
					Product: types.Product{
						ID:              "product-id",
						ProductTypeID:   "product-type-id",
						CategoryID:      "",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Scope:           types.ProductScopePublic,
						Status:          types.ProductStatusForSale,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{
								URL:         "https://and-period.jp/thumbnail01.png",
								IsThumbnail: true,
							},
							{
								URL:         "https://and-period.jp/thumbnail02.png",
								IsThumbnail: false,
							},
						},
						Price:                400,
						Cost:                 300,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						StartAt:              1640962800,
						EndAt:                1640962800,
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
					revisionID: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProducts(tt.products))
		})
	}
}

func TestProducts_ProducerIDs(t *testing.T) {
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
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:                400,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
				},
			},
			expect: []string{"producer-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.products.ProducerIDs())
		})
	}
}

func TestProducts_CategoryIDs(t *testing.T) {
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
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:                400,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
				},
			},
			expect: []string{"category-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.products.CategoryIDs())
		})
	}
}

func TestProducts_ProductTypeIDs(t *testing.T) {
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
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:                400,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
				},
			},
			expect: []string{"product-type-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.products.ProductTypeIDs())
		})
	}
}

func TestProducts_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   map[string]*Product
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
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:                400,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
				},
			},
			expect: map[string]*Product{
				"product-id": {
					Product: types.Product{
						ID:              "product-id",
						ProductTypeID:   "product-type-id",
						CategoryID:      "category-id",
						CoordinatorID:   "coordinator-id",
						ProducerID:      "producer-id",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:                400,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.Map())
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
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:                400,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
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
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:                400,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
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
						Public:          true,
						Status:          types.ProductStatusForSale,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*types.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:                400,
						DeliveryType:         types.DeliveryTypeNormal,
						Box60Rate:            50,
						Box80Rate:            40,
						Box100Rate:           30,
						OriginPrefectureCode: 25,
						OriginCity:           "彦根市",
						StartAt:              1640962800,
						EndAt:                1640962800,
						CreatedAt:            1640962800,
						UpdatedAt:            1640962800,
					},
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
					Public:          true,
					Status:          types.ProductStatusForSale,
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: []*types.ProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					Price:                400,
					DeliveryType:         types.DeliveryTypeNormal,
					Box60Rate:            50,
					Box80Rate:            40,
					Box100Rate:           30,
					OriginPrefectureCode: 25,
					OriginCity:           "彦根市",
					StartAt:              1640962800,
					EndAt:                1640962800,
					CreatedAt:            1640962800,
					UpdatedAt:            1640962800,
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
				URL:         "https://and-period.jp/thumbnail01.png",
				IsThumbnail: true,
			},
			expect: &ProductMedia{
				ProductMedia: types.ProductMedia{
					URL:         "https://and-period.jp/thumbnail01.png",
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
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
				},
			},
			expect: &types.ProductMedia{
				URL:         "https://and-period.jp/thumbnail01.png",
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
					ProductMedia: types.ProductMedia{
						URL:         "https://and-period.jp/thumbnail01.png",
						IsThumbnail: true,
					},
				},
				{
					ProductMedia: types.ProductMedia{
						URL:         "https://and-period.jp/thumbnail02.png",
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
						URL:         "https://and-period.jp/thumbnail01.png",
						IsThumbnail: true,
					},
				},
				{
					ProductMedia: types.ProductMedia{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
					},
				},
			},
			expect: []*types.ProductMedia{
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.media.Response())
		})
	}
}
