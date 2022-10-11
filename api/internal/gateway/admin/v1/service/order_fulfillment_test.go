package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestFulfillmentStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.FulfillmentStatus
		expect FulfillmentStatus
	}{
		{
			name:   "unfulfilled",
			status: entity.FulfillmentStatusUnfulfilled,
			expect: FulfillmentStatusUnfulfilled,
		},
		{
			name:   "fulfilled",
			status: entity.FulfillmentStatusFulfilled,
			expect: FulfillmentStatusFulfilled,
		},
		{
			name:   "unknown",
			status: entity.FulfillmentStatusUnknown,
			expect: FulfillmentStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewFulfillmentStatus(tt.status))
		})
	}
}

func TestFulfillmentStatus_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status FulfillmentStatus
		expect int32
	}{
		{
			name:   "success",
			status: FulfillmentStatusFulfilled,
			expect: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestShippingCarrier(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ShippingCarrier
		expect ShippingCarrier
	}{
		{
			name:   "yamato",
			status: entity.ShippingCarrierYamato,
			expect: ShippingCarrierYamato,
		},
		{
			name:   "sagawa",
			status: entity.ShippingCarrierSagawa,
			expect: ShippingCarrierSagawa,
		},
		{
			name:   "unknown",
			status: entity.ShippingCarrierUnknown,
			expect: ShippingCarrierUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewShippingCarrier(tt.status))
		})
	}
}

func TestShippingCarrier_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ShippingCarrier
		expect int32
	}{
		{
			name:   "success",
			status: ShippingCarrierYamato,
			expect: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestShippingSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ShippingSize
		expect ShippingSize
	}{
		{
			name:   "size 60",
			status: entity.ShippingSize60,
			expect: ShippingSize60,
		},
		{
			name:   "size 80",
			status: entity.ShippingSize80,
			expect: ShippingSize80,
		},
		{
			name:   "size 100",
			status: entity.ShippingSize100,
			expect: ShippingSize100,
		},
		{
			name:   "unknown",
			status: entity.ShippingSizeUnknown,
			expect: ShippingSizeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewShippingSize(tt.status))
		})
	}
}

func TestShippingSize_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ShippingSize
		expect int32
	}{
		{
			name:   "success",
			status: ShippingSize60,
			expect: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestOrderFulfillment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		fulfillment *entity.OrderFulfillment
		status      entity.FulfillmentStatus
		expect      *OrderFulfillment
	}{
		{
			name: "success",
			fulfillment: &entity.OrderFulfillment{
				ID:              "fulfillment-id",
				OrderID:         "order-id",
				ShippingID:      "shipping-id",
				TrackingNumber:  "",
				ShippingCarrier: entity.ShippingCarrierUnknown,
				ShippingMethod:  entity.DeliveryTypeNormal,
				BoxSize:         entity.ShippingSize60,
				BoxCount:        1,
				WeightTotal:     1000,
				Lastname:        "&.",
				Firstname:       "スタッフ",
				PostalCode:      "1000014",
				Prefecture:      "東京都",
				City:            "千代田区",
				AddressLine1:    "永田町1-7-1",
				AddressLine2:    "",
				PhoneNumber:     "+819012345678",
				CreatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			status: entity.FulfillmentStatusUnfulfilled,
			expect: &OrderFulfillment{
				OrderFulfillment: response.OrderFulfillment{
					TrackingNumber:  "",
					Status:          FulfillmentStatusUnfulfilled.Response(),
					ShippingCarrier: ShippingCarrierUnknown.Response(),
					ShippingMethod:  DeliveryTypeNormal.Response(),
					BoxSize:         ShippingSize60.Response(),
					BoxCount:        1,
					WeightTotal:     1.0,
					Lastname:        "&.",
					Firstname:       "スタッフ",
					PostalCode:      "1000014",
					Prefecture:      "東京都",
					City:            "千代田区",
					AddressLine1:    "永田町1-7-1",
					AddressLine2:    "",
					PhoneNumber:     "+819012345678",
				},
				id:         "fulfillment-id",
				orderID:    "order-id",
				shippingID: "shipping-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderFulfillment(tt.fulfillment, tt.status))
		})
	}
}

func TestOrderFulfillment_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		fulfillment *OrderFulfillment
		expect      *response.OrderFulfillment
	}{
		{
			name: "success",
			fulfillment: &OrderFulfillment{
				OrderFulfillment: response.OrderFulfillment{
					TrackingNumber:  "",
					Status:          FulfillmentStatusUnfulfilled.Response(),
					ShippingCarrier: ShippingCarrierUnknown.Response(),
					ShippingMethod:  DeliveryTypeNormal.Response(),
					BoxSize:         ShippingSize60.Response(),
					BoxCount:        1,
					WeightTotal:     1.0,
					Lastname:        "&.",
					Firstname:       "スタッフ",
					PostalCode:      "1000014",
					Prefecture:      "東京都",
					City:            "千代田区",
					AddressLine1:    "永田町1-7-1",
					AddressLine2:    "",
					PhoneNumber:     "+819012345678",
				},
				id:         "fulfillment-id",
				orderID:    "order-id",
				shippingID: "shipping-id",
			},
			expect: &response.OrderFulfillment{
				TrackingNumber:  "",
				Status:          FulfillmentStatusUnfulfilled.Response(),
				ShippingCarrier: ShippingCarrierUnknown.Response(),
				ShippingMethod:  DeliveryTypeNormal.Response(),
				BoxSize:         ShippingSize60.Response(),
				BoxCount:        1,
				WeightTotal:     1.0,
				Lastname:        "&.",
				Firstname:       "スタッフ",
				PostalCode:      "1000014",
				Prefecture:      "東京都",
				City:            "千代田区",
				AddressLine1:    "永田町1-7-1",
				AddressLine2:    "",
				PhoneNumber:     "+819012345678",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.fulfillment.Response())
		})
	}
}
