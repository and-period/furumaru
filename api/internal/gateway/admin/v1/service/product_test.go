package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

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

func TestDeliveryType_StoreEntity(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		deliveryType DeliveryType
		expect       entity.DeliveryType
	}{
		{
			name:         "success to normal",
			deliveryType: DeliveryTypeNormal,
			expect:       entity.DeliveryTypeNormal,
		},
		{
			name:         "success to frozen",
			deliveryType: DeliveryTypeFrozen,
			expect:       entity.DeliveryTypeFrozen,
		},
		{
			name:         "success to refrigerated",
			deliveryType: DeliveryTypeRefrigerated,
			expect:       entity.DeliveryTypeRefrigerated,
		},
		{
			name:         "success to unknown",
			deliveryType: DeliveryTypeUnknown,
			expect:       entity.DeliveryTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.deliveryType.StoreEntity())
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
		tt := tt
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
				ProducerID:      "producer-id",
				TypeID:          "product-type-id",
				TagIDs:          []string{"product-tag-id"},
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
					},
					{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
					},
				},
				Price:             400,
				Cost:              300,
				RecommendedPoints: []string{"ポイント1", "ポイント2", "ポイント3"},
				StorageMethodType: entity.StorageMethodTypeNormal,
				DeliveryType:      entity.DeliveryTypeNormal,
				Box60Rate:         50,
				Box80Rate:         40,
				Box100Rate:        30,
				OriginPrefecture:  codes.PrefectureValues["shiga"],
				OriginCity:        "彦根市",
				CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Product{
				Product: response.Product{
					ID:              "product-id",
					ProducerID:      "producer-id",
					ProducerName:    "",
					CategoryID:      "",
					CategoryName:    "",
					TypeID:          "product-type-id",
					TypeName:        "",
					TypeIconURL:     "",
					TagIDs:          []string{"product-tag-id"},
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
					Price:             400,
					Cost:              300,
					RecommendedPoint1: "ポイント1",
					RecommendedPoint2: "ポイント2",
					RecommendedPoint3: "ポイント3",
					StorageMethodType: int32(StorageMethodTypeNormal),
					DeliveryType:      int32(DeliveryTypeNormal),
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "shiga",
					OriginCity:        "彦根市",
					CreatedAt:         1640962800,
					UpdatedAt:         1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewProduct(tt.product))
		})
	}
}

func TestProduct_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		product     *Product
		productType *ProductType
		producer    *Producer
		expect      *Product
	}{
		{
			name: "success",
			product: &Product{
				Product: response.Product{
					ID:              "product-id",
					TypeID:          "product-type-id",
					ProducerID:      "producer-id",
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
					DeliveryType:     int32(DeliveryTypeNormal),
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "shiga",
					OriginCity:       "彦根市",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
				},
			},
			productType: &ProductType{
				ProductType: response.ProductType{
					ID:           "product-type-id",
					Name:         "じゃがいも",
					IconURL:      "https://and-period.jp/icon.png",
					CategoryID:   "category-id",
					CategoryName: "野菜",
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
			producer: &Producer{
				Producer: response.Producer{
					ID:            "producer-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Username:      "&.農園",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					HeaderURL:     "https://and-period.jp/header.png",
					Email:         "test-producer@and-period.jp",
					PhoneNumber:   "+819012345678",
					PostalCode:    "1000014",
					Prefecture:    "東京都",
					City:          "千代田区",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
			expect: &Product{
				Product: response.Product{
					ID:              "product-id",
					TypeID:          "product-type-id",
					TypeName:        "じゃがいも",
					CategoryID:      "category-id",
					CategoryName:    "野菜",
					ProducerID:      "producer-id",
					ProducerName:    "&.農園",
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					TypeIconURL:     "https://and-period.jp/icon.png",
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
					DeliveryType:     int32(DeliveryTypeNormal),
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "shiga",
					OriginCity:       "彦根市",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.product.Fill(tt.productType, tt.producer)
			assert.Equal(t, tt.expect, tt.product)
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
					TypeID:          "product-type-id",
					TypeName:        "じゃがいも",
					TypeIconURL:     "https://and-period.jp/icon.png",
					CategoryID:      "category-id",
					CategoryName:    "野菜",
					ProducerID:      "producer-id",
					ProducerName:    "&.農園",
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
					DeliveryType:     int32(DeliveryTypeNormal),
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "shiga",
					OriginCity:       "彦根市",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
				},
			},
			expect: &response.Product{
				ID:              "product-id",
				TypeID:          "product-type-id",
				TypeName:        "じゃがいも",
				TypeIconURL:     "https://and-period.jp/icon.png",
				CategoryID:      "category-id",
				CategoryName:    "野菜",
				ProducerID:      "producer-id",
				ProducerName:    "&.農園",
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
				DeliveryType:     int32(DeliveryTypeNormal),
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "shiga",
				OriginCity:       "彦根市",
				CreatedAt:        1640962800,
				UpdatedAt:        1640962800,
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
		expect   Products
	}{
		{
			name: "success",
			products: entity.Products{
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
							Images: common.Images{
								{URL: "https://and-period.jp/thumbnail02_240.png", Size: common.ImageSizeSmall},
								{URL: "https://and-period.jp/thumbnail02_675.png", Size: common.ImageSizeMedium},
								{URL: "https://and-period.jp/thumbnail02_900.png", Size: common.ImageSizeLarge},
							},
						},
					},
					Price:            400,
					DeliveryType:     entity.DeliveryTypeNormal,
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: codes.PrefectureValues["shiga"],
					OriginCity:       "彦根市",
					CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Products{
				{
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "",
						TypeIconURL:     "",
						CategoryID:      "",
						CategoryName:    "",
						ProducerID:      "producer-id",
						ProducerName:    "",
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
									{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
									{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
									{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
								},
							},
							{
								URL:         "https://and-period.jp/thumbnail02.png",
								IsThumbnail: false,
								Images: []*response.Image{
									{URL: "https://and-period.jp/thumbnail02_240.png", Size: int32(ImageSizeSmall)},
									{URL: "https://and-period.jp/thumbnail02_675.png", Size: int32(ImageSizeMedium)},
									{URL: "https://and-period.jp/thumbnail02_900.png", Size: int32(ImageSizeLarge)},
								},
							},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
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
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						ProducerName:    "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
			expect: []string{"producer-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
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
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						ProducerName:    "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
			expect: []string{"category-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
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
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						ProducerName:    "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
			expect: []string{"product-type-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
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
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						ProducerName:    "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
			expect: map[string]*Product{
				"product-id": {
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						ProducerName:    "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.Map())
		})
	}
}

func TestProducts_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		products     Products
		productTypes map[string]*ProductType
		producers    map[string]*Producer
		expect       Products
	}{
		{
			name: "success",
			products: Products{
				{
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						CategoryID:      "category-id",
						ProducerID:      "producer-id",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
			productTypes: map[string]*ProductType{
				"product-type-id": {
					ProductType: response.ProductType{
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
			producers: map[string]*Producer{
				"producer-id": {
					Producer: response.Producer{
						ID:            "producer-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Username:      "&.農園",
						ThumbnailURL:  "https://and-period.jp/thumbnail.png",
						HeaderURL:     "https://and-period.jp/header.png",
						Email:         "test-producer@and-period.jp",
						PhoneNumber:   "+819012345678",
						PostalCode:    "1000014",
						Prefecture:    "東京都",
						City:          "千代田区",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: Products{
				{
					Product: response.Product{
						ID:              "product-id",
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						ProducerName:    "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.products.Fill(tt.productTypes, tt.producers)
			assert.Equal(t, tt.expect, tt.products)
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
						TypeID:          "product-type-id",
						TypeName:        "じゃがいも",
						TypeIconURL:     "https://and-period.jp/icon.png",
						CategoryID:      "category-id",
						CategoryName:    "野菜",
						ProducerID:      "producer-id",
						ProducerName:    "&.農園",
						Name:            "新鮮なじゃがいも",
						Description:     "新鮮なじゃがいもをお届けします。",
						Public:          true,
						Inventory:       100,
						Weight:          1.3,
						ItemUnit:        "袋",
						ItemDescription: "1袋あたり100gのじゃがいも",
						Media: []*response.ProductMedia{
							{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
							{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
						},
						Price:            400,
						DeliveryType:     int32(DeliveryTypeNormal),
						Box60Rate:        50,
						Box80Rate:        40,
						Box100Rate:       30,
						OriginPrefecture: "shiga",
						OriginCity:       "彦根市",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
			expect: []*response.Product{
				{
					ID:              "product-id",
					TypeID:          "product-type-id",
					TypeName:        "じゃがいも",
					TypeIconURL:     "https://and-period.jp/icon.png",
					CategoryID:      "category-id",
					CategoryName:    "野菜",
					ProducerID:      "producer-id",
					ProducerName:    "&.農園",
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Inventory:       100,
					Weight:          1.3,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: []*response.ProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					Price:            400,
					DeliveryType:     int32(DeliveryTypeNormal),
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "shiga",
					OriginCity:       "彦根市",
					CreatedAt:        1640962800,
					UpdatedAt:        1640962800,
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
				Images: common.Images{
					{URL: "https://and-period.jp/thumbnail01_240.png", Size: common.ImageSizeSmall},
					{URL: "https://and-period.jp/thumbnail01_675.png", Size: common.ImageSizeMedium},
					{URL: "https://and-period.jp/thumbnail01_900.png", Size: common.ImageSizeLarge},
				},
			},
			expect: &ProductMedia{
				ProductMedia: response.ProductMedia{
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
					Images: []*response.Image{
						{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
					},
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
					Images: []*response.Image{
						{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
					},
				},
			},
			expect: &response.ProductMedia{
				URL:         "https://and-period.jp/thumbnail01.png",
				IsThumbnail: true,
				Images: []*response.Image{
					{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
					{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
					{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
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
					Images: common.Images{
						{URL: "https://and-period.jp/thumbnail01_240.png", Size: common.ImageSizeSmall},
						{URL: "https://and-period.jp/thumbnail01_675.png", Size: common.ImageSizeMedium},
						{URL: "https://and-period.jp/thumbnail01_900.png", Size: common.ImageSizeLarge},
					},
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
						Images: []*response.Image{
							{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
						},
					},
				},
				{
					ProductMedia: response.ProductMedia{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
						Images:      []*response.Image{},
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
						Images: []*response.Image{
							{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
						},
					},
				},
				{
					ProductMedia: response.ProductMedia{
						URL:         "https://and-period.jp/thumbnail02.png",
						IsThumbnail: false,
						Images:      []*response.Image{},
					},
				},
			},
			expect: []*response.ProductMedia{
				{
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
					Images: []*response.Image{
						{URL: "https://and-period.jp/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
					},
				},
				{
					URL:         "https://and-period.jp/thumbnail02.png",
					IsThumbnail: false,
					Images:      []*response.Image{},
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
