package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

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
				MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true}]`)),
			},
			expect: &Product{
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
					},
				},
				MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true}]`)),
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.media.Marshal()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
