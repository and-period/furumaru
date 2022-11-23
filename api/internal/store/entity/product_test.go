package entity

import (
	"testing"

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
				Price:            400,
				DeliveryType:     DeliveryTypeNormal,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
			},
			expect: &Product{
				TypeID:          "type-id",
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
				Price:            400,
				DeliveryType:     DeliveryTypeNormal,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
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
				ID:        "product-id",
				Name:      "&.農園のみかん",
				MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true,"images":[{"url":"https://and-period.jp/thumbnail_240.png","size":1}]}]`)),
			},
			expect: &Product{
				ID:   "product-id",
				Name: "&.農園のみかん",
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
			},
			expect: &Product{
				ID:   "product-id",
				Name: "&.農園のみかん",
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
					ID:        "product-id",
					Name:      "&.農園のみかん",
					MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true}]`)),
				},
			},
			expect: Products{
				{
					ID:   "product-id",
					Name: "&.農園のみかん",
					Media: MultiProductMedia{
						{
							URL:         "https://and-period.jp/thumbnail.png",
							IsThumbnail: true,
						},
					},
					MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true}]`)),
				},
			},
			hasErr: false,
		},
		{
			name: "success is empty",
			products: Products{
				{
					ID:        "product-id",
					Name:      "&.農園のみかん",
					MediaJSON: datatypes.JSON(nil),
				},
			},
			expect: Products{
				{
					ID:        "product-id",
					Name:      "&.農園のみかん",
					Media:     MultiProductMedia{},
					MediaJSON: datatypes.JSON(nil),
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
