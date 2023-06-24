package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/common"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestProduct(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewProductParams
		expect *Product
	}{
		{
			name: "success",
			params: &NewProductParams{
				TypeID:          "type-id",
				TagIDs:          []string{"tag-id"},
				ProducerID:      "producer-id",
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
				Price:             400,
				Cost:              300,
				StorageMethodType: StorageMethodTypeNormal,
				DeliveryType:      DeliveryTypeNormal,
				Box60Rate:         50,
				Box80Rate:         40,
				Box100Rate:        30,
				OriginPrefecture:  codes.PrefectureValues["shiga"],
				OriginCity:        "彦根市",
			},
			expect: &Product{
				TypeID:          "type-id",
				TagIDs:          []string{"tag-id"},
				ProducerID:      "producer-id",
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
				Price:             400,
				Cost:              300,
				StorageMethodType: StorageMethodTypeNormal,
				DeliveryType:      DeliveryTypeNormal,
				Box60Rate:         50,
				Box80Rate:         40,
				Box100Rate:        30,
				OriginPrefecture:  codes.PrefectureValues["shiga"],
				OriginCity:        "彦根市",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProduct(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.product.Validate()
			assert.Equal(t, tt.hasErr, actual != nil, actual)
		})
	}
}

func TestProduct_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		product *Product
		expect  *Product
		hasErr  bool
	}{
		{
			name: "success",
			product: &Product{
				ID:                    "product-id",
				Name:                  "&.農園のみかん",
				TagIDsJSON:            datatypes.JSON([]byte(`["tag-id01","tag-id02"]`)),
				MediaJSON:             datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true,"images":[{"url":"https://and-period.jp/thumbnail_240.png","size":1}]}]`)),
				RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
			},
			expect: &Product{
				ID:   "product-id",
				Name: "&.農園のみかん",
				TagIDs: []string{
					"tag-id01",
					"tag-id02",
				},
				TagIDsJSON: datatypes.JSON([]byte(`["tag-id01","tag-id02"]`)),
				Media: MultiProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail.png",
						IsThumbnail: true,
						Images: common.Images{
							{
								URL:  "https://and-period.jp/thumbnail_240.png",
								Size: common.ImageSizeSmall,
							},
						},
					},
				},
				MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true,"images":[{"url":"https://and-period.jp/thumbnail_240.png","size":1}]}]`)),
				RecommendedPoints: []string{
					"ポイント1",
					"ポイント2",
				},
				RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.product.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.product)
		})
	}
}

func TestProduct_FillJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		product *Product
		expect  *Product
		hasErr  bool
	}{
		{
			name: "success",
			product: &Product{
				ID:   "product-id",
				Name: "&.農園のみかん",
				TagIDs: []string{
					"tag-id01",
					"tag-id02",
				},
				Media: MultiProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail.png",
						IsThumbnail: true,
						Images: common.Images{
							{
								URL:  "https://and-period.jp/thumbnail_240.png",
								Size: common.ImageSizeSmall,
							},
						},
					},
				},
				RecommendedPoints: []string{
					"ポイント1",
					"ポイント2",
				},
			},
			expect: &Product{
				ID:   "product-id",
				Name: "&.農園のみかん",
				TagIDs: []string{
					"tag-id01",
					"tag-id02",
				},
				TagIDsJSON: datatypes.JSON([]byte(`["tag-id01","tag-id02"]`)),
				Media: MultiProductMedia{
					{
						URL:         "https://and-period.jp/thumbnail.png",
						IsThumbnail: true,
						Images: common.Images{
							{
								URL:  "https://and-period.jp/thumbnail_240.png",
								Size: common.ImageSizeSmall,
							},
						},
					},
				},
				MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true,"images":[{"url":"https://and-period.jp/thumbnail_240.png","size":1}]}]`)),
				RecommendedPoints: []string{
					"ポイント1",
					"ポイント2",
				},
				RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.product.FillJSON()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.product)
		})
	}
}

func TestProducts_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		products Products
		expect   Products
		hasErr   bool
	}{
		{
			name: "success",
			products: Products{
				{
					ID:                    "product-id",
					Name:                  "&.農園のみかん",
					TagIDsJSON:            datatypes.JSON([]byte(`["tag-id01","tag-id02"]`)),
					MediaJSON:             datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true,"images":[{"url":"https://and-period.jp/thumbnail_240.png","size":1}]}]`)),
					RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
				},
			},
			expect: Products{
				{
					ID:   "product-id",
					Name: "&.農園のみかん",
					TagIDs: []string{
						"tag-id01",
						"tag-id02",
					},
					TagIDsJSON: datatypes.JSON([]byte(`["tag-id01","tag-id02"]`)),
					Media: MultiProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail.png",
							IsThumbnail: true,
							Images: common.Images{
								{
									URL:  "https://and-period.jp/thumbnail_240.png",
									Size: common.ImageSizeSmall,
								},
							},
						},
					},
					MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true,"images":[{"url":"https://and-period.jp/thumbnail_240.png","size":1}]}]`)),
					RecommendedPoints: []string{
						"ポイント1",
						"ポイント2",
					},
					RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
				},
			},
			hasErr: false,
		},
		{
			name: "success is empty",
			products: Products{
				{
					ID:                    "product-id",
					Name:                  "&.農園のみかん",
					TagIDsJSON:            datatypes.JSON(nil),
					MediaJSON:             datatypes.JSON(nil),
					RecommendedPointsJSON: datatypes.JSON(nil),
				},
			},
			expect: Products{
				{
					ID:                    "product-id",
					Name:                  "&.農園のみかん",
					TagIDs:                []string{},
					TagIDsJSON:            datatypes.JSON(nil),
					Media:                 MultiProductMedia{},
					MediaJSON:             datatypes.JSON(nil),
					RecommendedPoints:     []string{},
					RecommendedPointsJSON: datatypes.JSON(nil),
				},
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.products.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.expect, tt.products)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductMedia(tt.url, tt.isThumbnail)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductMedia_SetImages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		images common.Images
		media  *ProductMedia
		expect *ProductMedia
	}{
		{
			name: "success",
			images: common.Images{
				{Size: common.ImageSizeSmall, URL: "http://example.com/media.png"},
			},
			media: &ProductMedia{
				URL:         "http://example.com/media.png",
				IsThumbnail: true,
			},
			expect: &ProductMedia{
				URL:         "http://example.com/media.png",
				IsThumbnail: true,
				Images: common.Images{
					{Size: common.ImageSizeSmall, URL: "http://example.com/media.png"},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.media.SetImages(tt.images)
			assert.Equal(t, tt.expect, tt.media)
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
				"https://and-period.jp/thumbnail01.png": {URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				"https://and-period.jp/thumbnail02.png": {URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				"https://and-period.jp/thumbnail03.png": {URL: "https://and-period.jp/thumbnail03.png", IsThumbnail: false},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
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
		tt := tt
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
			expect: []byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true,"images":null}]`),
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.media.Marshal()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
