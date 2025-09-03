package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderMetadata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewOrderMetadataParams
		expect *OrderMetadata
	}{
		{
			name: "success with pickup",
			params: &NewOrderMetadataParams{
				OrderID:        "order-id",
				Pickup:         true,
				PickupAt:       jst.Date(2022, 1, 1, 10, 0, 0, 0),
				PickupLocation: "店舗A",
			},
			expect: &OrderMetadata{
				OrderID:        "order-id",
				PickupAt:       jst.Date(2022, 1, 1, 10, 0, 0, 0),
				PickupLocation: "店舗A",
			},
		},
		{
			name: "success with shipping",
			params: &NewOrderMetadataParams{
				OrderID:         "order-id",
				Pickup:          false,
				ShippingAddress: &entity.Address{ID: "address-id"},
				ShippingMessage: "ご注文ありがとうございます！",
			},
			expect: &OrderMetadata{
				OrderID:         "order-id",
				ShippingMessage: "ご注文ありがとうございます！",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderMetadata(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestMultiOrderMetadata_MapByOrderID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		metadata MultiOrderMetadata
		expect   map[string]*OrderMetadata
	}{
		{
			name: "success",
			metadata: MultiOrderMetadata{
				{
					OrderID:         "order-id01",
					PickupAt:        jst.Date(2022, 1, 1, 10, 0, 0, 0),
					PickupLocation:  "店舗A",
					ShippingMessage: "メッセージA",
				},
				{
					OrderID:         "order-id02",
					PickupAt:        jst.Date(2022, 1, 2, 15, 0, 0, 0),
					PickupLocation:  "店舗B",
					ShippingMessage: "メッセージB",
				},
			},
			expect: map[string]*OrderMetadata{
				"order-id01": {
					OrderID:         "order-id01",
					PickupAt:        jst.Date(2022, 1, 1, 10, 0, 0, 0),
					PickupLocation:  "店舗A",
					ShippingMessage: "メッセージA",
				},
				"order-id02": {
					OrderID:         "order-id02",
					PickupAt:        jst.Date(2022, 1, 2, 15, 0, 0, 0),
					PickupLocation:  "店舗B",
					ShippingMessage: "メッセージB",
				},
			},
		},
		{
			name:     "empty",
			metadata: MultiOrderMetadata{},
			expect:   map[string]*OrderMetadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.metadata.MapByOrderID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}