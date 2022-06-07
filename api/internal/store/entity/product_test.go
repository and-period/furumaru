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
