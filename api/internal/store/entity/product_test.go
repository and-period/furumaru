package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestProduct(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		params *NewProductParams
		expect *Product
		hasErr bool
	}{
		{
			name: "success",
			params: &NewProductParams{
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "type-id",
				TagIDs:          []string{"tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: MultiProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				RecommendedPoints:    []string{"おすすめポイント"},
				StorageMethodType:    StorageMethodTypeNormal,
				DeliveryType:         DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now,
				EndAt:                now.AddDate(1, 0, 0),
			},
			expect: &Product{
				ID:              "", // ignore
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "type-id",
				TagIDs:          []string{"tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Status:          0,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: MultiProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				ExpirationDate:       7,
				RecommendedPoints:    []string{"おすすめポイント"},
				StorageMethodType:    StorageMethodTypeNormal,
				DeliveryType:         DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefecture:     "滋賀県",
				OriginPrefectureCode: 25,
				OriginCity:           "彦根市",
				StartAt:              now,
				EndAt:                now.AddDate(1, 0, 0),
				ProductRevision: ProductRevision{
					Price: 400,
					Cost:  300,
				},
			},
			hasErr: false,
		},
		{
			name: "invalid prefecture",
			params: &NewProductParams{
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				TypeID:          "type-id",
				TagIDs:          []string{"tag-id"},
				Name:            "新鮮なじゃがいも",
				Description:     "新鮮なじゃがいもをお届けします。",
				Public:          true,
				Inventory:       100,
				Weight:          100,
				WeightUnit:      WeightUnitGram,
				Item:            1,
				ItemUnit:        "袋",
				ItemDescription: "1袋あたり100gのじゃがいも",
				Media: MultiProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				Price:                400,
				Cost:                 300,
				ExpirationDate:       7,
				RecommendedPoints:    []string{"おすすめポイント"},
				StorageMethodType:    StorageMethodTypeNormal,
				DeliveryType:         DeliveryTypeNormal,
				Box60Rate:            50,
				Box80Rate:            40,
				Box100Rate:           30,
				OriginPrefectureCode: -1,
				OriginCity:           "彦根市",
				StartAt:              now,
				EndAt:                now.AddDate(1, 0, 0),
			},
			expect: nil,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewProduct(tt.params)
			if tt.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			actual.ID = ""        // ignore
			actual.ProductID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProduct_ShippingType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		product *Product
		expect  ShippingType
	}{
		{
			name:    "normal",
			product: &Product{DeliveryType: DeliveryTypeNormal},
			expect:  ShippingTypeNormal,
		},
		{
			name:    "refrigerated",
			product: &Product{DeliveryType: DeliveryTypeRefrigerated},
			expect:  ShippingTypeFrozen,
		},
		{
			name:    "frozen",
			product: &Product{DeliveryType: DeliveryTypeFrozen},
			expect:  ShippingTypeFrozen,
		},
		{
			name:    "unknown",
			product: &Product{DeliveryType: DeliveryTypeUnknown},
			expect:  ShippingTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.product.ShippingType())
		})
	}
}

func TestProduct_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		product *Product
		hasErr  bool
	}{
		{
			name: "success",
			product: &Product{
				Media: MultiProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				},
				RecommendedPoints: []string{"ポイント1", "ポイント2"},
			},
			hasErr: false,
		},
		{
			name: "invalid recommended points",
			product: &Product{
				Media: MultiProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				},
				RecommendedPoints: []string{"ポイント1", "ポイント2", "ポイント3", "ポイント4"},
			},
			hasErr: true,
		},
		{
			name: "invalid media",
			product: &Product{
				Media: MultiProductMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: true},
				},
				RecommendedPoints: []string{"ポイント1", "ポイント2"},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.product.Validate()
			assert.Equal(t, tt.hasErr, actual != nil, actual)
		})
	}
}

func TestProduct_Fill(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		product  *Product
		revision *ProductRevision
		expect   *Product
	}{
		{
			name: "success",
			product: &Product{
				ID:                   "product-id",
				Name:                 "&.農園のみかん",
				OriginPrefectureCode: 25,
				Public:               true,
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
			},
			revision: &ProductRevision{
				ID:        1,
				ProductID: "product-id",
				Price:     1980,
				Cost:      880,
			},
			expect: &Product{
				ID:                   "product-id",
				Name:                 "&.農園のみかん",
				Status:               ProductStatusForSale,
				OriginPrefecture:     "滋賀県",
				OriginPrefectureCode: 25,
				Public:               true,
				StartAt:              now.AddDate(0, -1, 0),
				EndAt:                now.AddDate(0, 1, 0),
				ProductRevision: ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     1980,
					Cost:      880,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.product.Fill(tt.revision, now)
			assert.Equal(t, tt.expect, tt.product)
		})
	}
}

func TestProduct_SetStatus(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name    string
		product *Product
		expect  ProductStatus
	}{
		{
			name: "archived",
			product: &Product{
				DeletedAt: gorm.DeletedAt{Time: time.Now()},
			},
			expect: ProductStatusArchived,
		},
		{
			name: "private",
			product: &Product{
				Public: false,
			},
			expect: ProductStatusPrivate,
		},
		{
			name: "presale",
			product: &Product{
				Public:  true,
				StartAt: now.AddDate(0, 1, 0),
				EndAt:   now.AddDate(0, 1, 0),
			},
			expect: ProductStatusPresale,
		},
		{
			name: "for sale",
			product: &Product{
				Public:  true,
				StartAt: now.AddDate(0, -1, 0),
				EndAt:   now.AddDate(0, 1, 0),
			},
			expect: ProductStatusForSale,
		},
		{
			name: "out of sale",
			product: &Product{
				Public:  true,
				StartAt: now.AddDate(0, -1, 0),
				EndAt:   now.AddDate(0, -1, 0),
			},
			expect: ProductStatusOutOfSale,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.product.SetStatus(now)
			assert.Equal(t, tt.expect, tt.product.Status)
		})
	}
}

func TestProduct_WeightGram(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		product *Product
		expect  int64
	}{
		{
			name: "success gram",
			product: &Product{
				ID:         "product-id",
				Weight:     100,
				WeightUnit: WeightUnitGram,
				Box60Rate:  50,
				Box80Rate:  40,
				Box100Rate: 30,
			},
			expect: 100,
		},
		{
			name: "success kilogram",
			product: &Product{
				ID:         "product-id",
				Weight:     1,
				WeightUnit: WeightUnitKilogram,
				Box60Rate:  50,
				Box80Rate:  40,
				Box100Rate: 30,
			},
			expect: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.product.WeightGram())
		})
	}
}

func TestProducts_Fill(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		products  Products
		revisions map[string]*ProductRevision
		expect    Products
	}{
		{
			name: "success",
			products: Products{
				{
					ID:                   "product-id01",
					Name:                 "&.農園のみかん",
					Public:               false,
					OriginPrefectureCode: 25,
				},
			},
			revisions: map[string]*ProductRevision{
				"product-id01": {
					ID:        1,
					ProductID: "product-id01",
					Price:     1980,
					Cost:      880,
				},
			},
			expect: Products{
				{
					ID:                   "product-id01",
					Name:                 "&.農園のみかん",
					Public:               false,
					Status:               ProductStatusPrivate,
					OriginPrefecture:     "滋賀県",
					OriginPrefectureCode: 25,
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id01",
						Price:     1980,
						Cost:      880,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.products.Fill(tt.revisions, now)
			assert.Equal(t, tt.expect, tt.products)
		})
	}
}

func TestProducts_Box60Rate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   int64
	}{
		{
			name: "success",
			products: Products{
				{
					ID:         "product-id01",
					Weight:     100,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  40,
					Box100Rate: 30,
				},
				{
					ID:         "product-id02",
					Weight:     200,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  45,
					Box100Rate: 40,
				},
			},
			expect: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.Box60Rate())
		})
	}
}

func TestProducts_Box80Rate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   int64
	}{
		{
			name: "success",
			products: Products{
				{
					ID:         "product-id01",
					Weight:     100,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  40,
					Box100Rate: 30,
				},
				{
					ID:         "product-id02",
					Weight:     200,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  45,
					Box100Rate: 40,
				},
			},
			expect: 85,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.Box80Rate())
		})
	}
}

func TestProducts_Box100Rate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   int64
	}{
		{
			name: "success",
			products: Products{
				{
					ID:         "product-id01",
					Weight:     100,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  40,
					Box100Rate: 30,
				},
				{
					ID:         "product-id02",
					Weight:     200,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  45,
					Box100Rate: 40,
				},
			},
			expect: 70,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.Box100Rate())
		})
	}
}

func TestProducts_WeightGram(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   int64
	}{
		{
			name: "success gram",
			products: Products{
				{
					ID:         "product-id01",
					Weight:     100,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  40,
					Box100Rate: 30,
				},
				{
					ID:         "product-id02",
					Weight:     200,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  45,
					Box100Rate: 40,
				},
			},
			expect: 300,
		},
		{
			name: "success kilogram",
			products: Products{
				{
					ID:         "product-id01",
					Weight:     1,
					WeightUnit: WeightUnitKilogram,
					Box60Rate:  50,
					Box80Rate:  40,
					Box100Rate: 30,
				},
				{
					ID:         "product-id02",
					Weight:     2,
					WeightUnit: WeightUnitKilogram,
					Box60Rate:  50,
					Box80Rate:  45,
					Box100Rate: 40,
				},
			},
			expect: 3000,
		},
		{
			name: "success mix",
			products: Products{
				{
					ID:         "product-id01",
					Weight:     100,
					WeightUnit: WeightUnitGram,
					Box60Rate:  50,
					Box80Rate:  40,
					Box100Rate: 30,
				},
				{
					ID:         "product-id02",
					Weight:     2,
					WeightUnit: WeightUnitKilogram,
					Box60Rate:  50,
					Box80Rate:  45,
					Box100Rate: 40,
				},
			},
			expect: 2100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.products.WeightGram())
		})
	}
}

func TestProducts_IDs(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		products Products
		expect   []string
	}{
		{
			name: "success",
			products: Products{
				{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
					},
				},
			},
			expect: []string{"product-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.products.IDs())
		})
	}
}

func TestProducts_CoordinatorIDs(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		products Products
		expect   []string
	}{
		{
			name: "success",
			products: Products{
				{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
					},
				},
			},
			expect: []string{"coordinator-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.products.CoordinatorIDs())
		})
	}
}

func TestProducts_ProducerIDs(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		products Products
		expect   []string
	}{
		{
			name: "success",
			products: Products{
				{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
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

func TestProducts_ProductTypeIDs(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		products Products
		expect   []string
	}{
		{
			name: "success",
			products: Products{
				{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
					},
				},
			},
			expect: []string{"type-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.products.ProductTypeIDs())
		})
	}
}

func TestProducts_ProductTagIDs(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		products Products
		expect   []string
	}{
		{
			name: "success",
			products: Products{
				{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
					},
				},
			},
			expect: []string{"tag-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.products.ProductTagIDs())
		})
	}
}

func TestProducts_Map(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		products Products
		expect   map[string]*Product
	}{
		{
			name: "success",
			products: Products{
				{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
					},
				},
			},
			expect: map[string]*Product{
				"product-id": {
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
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
	now := time.Now()
	tests := []struct {
		name     string
		products Products
		expect   map[int64]*Product
	}{
		{
			name: "success",
			products: Products{
				{
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
					},
				},
			},
			expect: map[int64]*Product{
				1: {
					ID:              "product-id",
					CoordinatorID:   "coordinator-id",
					ProducerID:      "producer-id",
					TypeID:          "type-id",
					TagIDs:          []string{"tag-id"},
					Name:            "新鮮なじゃがいも",
					Description:     "新鮮なじゃがいもをお届けします。",
					Public:          true,
					Status:          0,
					Inventory:       100,
					Weight:          100,
					WeightUnit:      WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: MultiProductMedia{
						{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
						{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
					},
					ExpirationDate:    7,
					RecommendedPoints: []string{"おすすめポイント"},
					StorageMethodType: StorageMethodTypeNormal,
					DeliveryType:      DeliveryTypeNormal,
					Box60Rate:         50,
					Box80Rate:         40,
					Box100Rate:        30,
					OriginPrefecture:  "滋賀県",
					OriginCity:        "彦根市",
					StartAt:           now,
					EndAt:             now.AddDate(1, 0, 0),
					ProductRevision: ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
					},
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

func TestProducts_Filter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		products   Products
		productIDs []string
		expect     Products
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01"},
				{ID: "product-id02"},
				{ID: "product-id03"},
			},
			productIDs: []string{
				"product-id01",
				"product-id03",
			},
			expect: Products{
				{ID: "product-id01"},
				{ID: "product-id03"},
			},
		},
		{
			name: "empty",
			products: Products{
				{ID: "product-id01"},
				{ID: "product-id02"},
				{ID: "product-id03"},
			},
			productIDs: []string{},
			expect:     Products{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.products.Filter(tt.productIDs...)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProducts_FilterBySales(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   Products
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01", Status: ProductStatusPrivate},
				{ID: "product-id02", Status: ProductStatusPresale},
				{ID: "product-id03", Status: ProductStatusForSale},
				{ID: "product-id04", Status: ProductStatusOutOfSale},
				{ID: "product-id05", Status: ProductStatusArchived},
			},
			expect: Products{
				{ID: "product-id03", Status: ProductStatusForSale},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.products.FilterBySales()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProducts_FilterByPublished(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		products Products
		expect   Products
	}{
		{
			name: "success",
			products: Products{
				{ID: "product-id01", Status: ProductStatusPrivate},
				{ID: "product-id02", Status: ProductStatusPresale},
				{ID: "product-id03", Status: ProductStatusForSale},
				{ID: "product-id04", Status: ProductStatusOutOfSale},
				{ID: "product-id05", Status: ProductStatusArchived},
			},
			expect: Products{
				{ID: "product-id02", Status: ProductStatusPresale},
				{ID: "product-id03", Status: ProductStatusForSale},
				{ID: "product-id04", Status: ProductStatusOutOfSale},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.products.FilterByPublished()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductMedia(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		url         string
		isThumbnail bool
		expect      *ProductMedia
	}{
		{
			name:        "success",
			url:         "http://example.com/media.png",
			isThumbnail: true,
			expect: &ProductMedia{
				URL:         "http://example.com/media.png",
				IsThumbnail: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductMedia(tt.url, tt.isThumbnail)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestMultiProductMedia_MapByURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		media  MultiProductMedia
		expect map[string]*ProductMedia
	}{
		{
			name: "success",
			media: MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				{URL: "https://and-period.jp/thumbnail03.png", IsThumbnail: false},
			},
			expect: map[string]*ProductMedia{
				"https://and-period.jp/thumbnail01.png": {
					URL:         "https://and-period.jp/thumbnail01.png",
					IsThumbnail: true,
				},
				"https://and-period.jp/thumbnail02.png": {
					URL:         "https://and-period.jp/thumbnail02.png",
					IsThumbnail: false,
				},
				"https://and-period.jp/thumbnail03.png": {
					URL:         "https://and-period.jp/thumbnail03.png",
					IsThumbnail: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.media.MapByURL()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestMultiProductMedia_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		media  MultiProductMedia
		expect error
	}{
		{
			name: "success",
			media: MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				{URL: "https://and-period.jp/thumbnail03.png", IsThumbnail: false},
			},
			expect: nil,
		},
		{
			name:   "success is empty",
			media:  nil,
			expect: nil,
		},
		{
			name: "failed to multiple thumbnails",
			media: MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				{URL: "https://and-period.jp/thumbnail03.png", IsThumbnail: true},
			},
			expect: errOnlyOneThumbnail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.media.Validate()
			assert.ErrorIs(t, tt.expect, err)
		})
	}
}

func TestMultiProductMedia_Marshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		media  MultiProductMedia
		expect []byte
		hasErr bool
	}{
		{
			name: "success",
			media: MultiProductMedia{
				{
					URL:         "https://and-period.jp/thumbnail.png",
					IsThumbnail: true,
				},
			},
			expect: []byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true}]`),
			hasErr: false,
		},
		{
			name:   "success is empty",
			media:  nil,
			expect: []byte{},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.media.Marshal()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
