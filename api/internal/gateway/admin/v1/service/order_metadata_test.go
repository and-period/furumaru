package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestNewOrderMetadata(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		metadata *entity.OrderMetadata
		expect   *OrderMetadata
	}{
		{
			name: "success",
			metadata: &entity.OrderMetadata{
				OrderID:        "order-id",
				PickupAt:       now,
				PickupLocation: "東京都千代田区",
			},
			expect: &OrderMetadata{
				OrderMetadata: response.OrderMetadata{
					PickupAt:       jst.Unix(now),
					PickupLocation: "東京都千代田区",
				},
				orderID: "order-id",
			},
		},
		{
			name: "empty location",
			metadata: &entity.OrderMetadata{
				OrderID:        "order-id",
				PickupAt:       now,
				PickupLocation: "",
			},
			expect: &OrderMetadata{
				OrderMetadata: response.OrderMetadata{
					PickupAt:       jst.Unix(now),
					PickupLocation: "",
				},
				orderID: "order-id",
			},
		},
		{
			name: "zero time",
			metadata: &entity.OrderMetadata{
				OrderID:        "order-id",
				PickupAt:       time.Time{},
				PickupLocation: "東京都千代田区",
			},
			expect: &OrderMetadata{
				OrderMetadata: response.OrderMetadata{
					PickupAt:       0,
					PickupLocation: "東京都千代田区",
				},
				orderID: "order-id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderMetadata(tt.metadata)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestOrderMetadata_Response(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		metadata *OrderMetadata
		expect   *response.OrderMetadata
	}{
		{
			name: "success",
			metadata: &OrderMetadata{
				OrderMetadata: response.OrderMetadata{
					PickupAt:       jst.Unix(now),
					PickupLocation: "東京都千代田区",
				},
				orderID: "order-id",
			},
			expect: &response.OrderMetadata{
				PickupAt:       jst.Unix(now),
				PickupLocation: "東京都千代田区",
			},
		},
		{
			name: "empty values",
			metadata: &OrderMetadata{
				OrderMetadata: response.OrderMetadata{
					PickupAt:       0,
					PickupLocation: "",
				},
				orderID: "order-id",
			},
			expect: &response.OrderMetadata{
				PickupAt:       0,
				PickupLocation: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.metadata.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
