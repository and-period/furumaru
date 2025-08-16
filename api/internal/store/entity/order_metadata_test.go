package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderPickupMetadata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *NewOrderPickupMetadataParams
		expect *OrderMetadata
	}{
		{
			name: "success",
			params: &NewOrderPickupMetadataParams{
				OrderID:        "order-id",
				PickupAt:       jst.Date(2022, 1, 1, 10, 0, 0, 0),
				PickupLocation: "店舗A",
			},
			expect: &OrderMetadata{
				OrderID:        "order-id",
				PickupAt:       jst.Date(2022, 1, 1, 10, 0, 0, 0),
				PickupLocation: "店舗A",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderPickupMetadata(tt.params)
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
					OrderID:        "order-id01",
					PickupAt:       jst.Date(2022, 1, 1, 10, 0, 0, 0),
					PickupLocation: "店舗A",
				},
				{
					OrderID:        "order-id02",
					PickupAt:       jst.Date(2022, 1, 2, 15, 0, 0, 0),
					PickupLocation: "店舗B",
				},
			},
			expect: map[string]*OrderMetadata{
				"order-id01": {
					OrderID:        "order-id01",
					PickupAt:       jst.Date(2022, 1, 1, 10, 0, 0, 0),
					PickupLocation: "店舗A",
				},
				"order-id02": {
					OrderID:        "order-id02",
					PickupAt:       jst.Date(2022, 1, 2, 15, 0, 0, 0),
					PickupLocation: "店舗B",
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