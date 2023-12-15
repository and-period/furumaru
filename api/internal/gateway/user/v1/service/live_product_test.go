package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestLiveProduct(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name    string
		product *entity.Product
		expect  *LiveProduct
	}{
		{
			name: "success",
			product: &entity.Product{
				ID:              "product-id",
				TypeID:          "product-type-id",
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
				Media: entity.MultiProductMedia{
					{
						URL:         "https://example.com/thumbnail01.png",
						IsThumbnail: true,
						Images: common.Images{
							{URL: "https://example.com/thumbnail01_240.png", Size: common.ImageSizeSmall},
							{URL: "https://example.com/thumbnail01_675.png", Size: common.ImageSizeMedium},
							{URL: "https://example.com/thumbnail01_900.png", Size: common.ImageSizeLarge},
						},
					},
					{
						URL:         "https://example.com/thumbnail02.png",
						IsThumbnail: false,
						Images: common.Images{
							{URL: "https://example.com/thumbnail02_240.png", Size: common.ImageSizeSmall},
							{URL: "https://example.com/thumbnail02_675.png", Size: common.ImageSizeMedium},
							{URL: "https://example.com/thumbnail02_900.png", Size: common.ImageSizeLarge},
						},
					},
				},
				DeliveryType:     entity.DeliveryTypeNormal,
				Box60Rate:        50,
				Box80Rate:        40,
				Box100Rate:       30,
				OriginPrefecture: "滋賀県",
				OriginCity:       "彦根市",
				StartAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
				EndAt:            jst.Date(2022, 1, 1, 0, 0, 0, 0),
				ProductRevision: entity.ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     400,
					Cost:      300,
					CreatedAt: now,
					UpdatedAt: now,
				},
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &LiveProduct{
				LiveProduct: response.LiveProduct{
					ProductID:    "product-id",
					Name:         "新鮮なじゃがいも",
					Price:        400,
					Inventory:    100,
					ThumbnailURL: "https://example.com/thumbnail01.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLiveProduct(tt.product)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLiveProduct_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		live   *LiveProduct
		expect *response.LiveProduct
	}{
		{
			name: "success",
			live: &LiveProduct{
				LiveProduct: response.LiveProduct{
					ProductID:    "product-id",
					Name:         "新鮮なじゃがいも",
					Price:        400,
					Inventory:    100,
					ThumbnailURL: "https://example.com/thumbnail01.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
					},
				},
			},
			expect: &response.LiveProduct{
				ProductID:    "product-id",
				Name:         "新鮮なじゃがいも",
				Price:        400,
				Inventory:    100,
				ThumbnailURL: "https://example.com/thumbnail01.png",
				Thumbnails: []*response.Image{
					{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
					{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
					{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.live.Response())
		})
	}
}

func TestLiveProducts(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	tests := []struct {
		name     string
		products entity.Products
		expect   LiveProducts
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
					Status:          entity.ProductStatusForSale,
					Inventory:       100,
					Weight:          1300,
					WeightUnit:      entity.WeightUnitGram,
					Item:            1,
					ItemUnit:        "袋",
					ItemDescription: "1袋あたり100gのじゃがいも",
					Media: entity.MultiProductMedia{
						{
							URL:         "https://example.com/thumbnail01.png",
							IsThumbnail: true,
							Images: common.Images{
								{URL: "https://example.com/thumbnail01_240.png", Size: common.ImageSizeSmall},
								{URL: "https://example.com/thumbnail01_675.png", Size: common.ImageSizeMedium},
								{URL: "https://example.com/thumbnail01_900.png", Size: common.ImageSizeLarge},
							},
						},
						{
							URL:         "https://example.com/thumbnail02.png",
							IsThumbnail: false,
							Images: common.Images{
								{URL: "https://example.com/thumbnail02_240.png", Size: common.ImageSizeSmall},
								{URL: "https://example.com/thumbnail02_675.png", Size: common.ImageSizeMedium},
								{URL: "https://example.com/thumbnail02_900.png", Size: common.ImageSizeLarge},
							},
						},
					},
					DeliveryType:     entity.DeliveryTypeNormal,
					Box60Rate:        50,
					Box80Rate:        40,
					Box100Rate:       30,
					OriginPrefecture: "滋賀県",
					OriginCity:       "彦根市",
					StartAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
					EndAt:            jst.Date(2022, 1, 1, 0, 0, 0, 0),
					ProductRevision: entity.ProductRevision{
						ID:        1,
						ProductID: "product-id",
						Price:     400,
						Cost:      300,
						CreatedAt: now,
						UpdatedAt: now,
					},
					CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: LiveProducts{
				{
					LiveProduct: response.LiveProduct{
						ProductID:    "product-id",
						Name:         "新鮮なじゃがいも",
						Price:        400,
						Inventory:    100,
						ThumbnailURL: "https://example.com/thumbnail01.png",
						Thumbnails: []*response.Image{
							{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLiveProducts(tt.products)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLiveProducts_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		lives  LiveProducts
		expect []*response.LiveProduct
	}{
		{
			name: "success",
			lives: LiveProducts{
				{
					LiveProduct: response.LiveProduct{
						ProductID:    "product-id",
						Name:         "新鮮なじゃがいも",
						Price:        400,
						Inventory:    100,
						ThumbnailURL: "https://example.com/thumbnail01.png",
						Thumbnails: []*response.Image{
							{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
							{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
							{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
						},
					},
				},
			},
			expect: []*response.LiveProduct{
				{
					ProductID:    "product-id",
					Name:         "新鮮なじゃがいも",
					Price:        400,
					Inventory:    100,
					ThumbnailURL: "https://example.com/thumbnail01.png",
					Thumbnails: []*response.Image{
						{URL: "https://example.com/thumbnail01_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://example.com/thumbnail01_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://example.com/thumbnail01_900.png", Size: int32(ImageSizeLarge)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.lives.Response())
		})
	}
}
