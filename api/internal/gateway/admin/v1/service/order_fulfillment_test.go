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

func TestShippingType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status entity.ShippingType
		expect ShippingType
	}{
		{
			name:   "normal",
			status: entity.ShippingTypeNormal,
			expect: ShippingTypeNormal,
		},
		{
			name:   "frozen",
			status: entity.ShippingTypeFrozen,
			expect: ShippingTypeFrozen,
		},
		{
			name:   "unknown",
			status: entity.ShippingTypeUnknown,
			expect: ShippingTypeUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewShippingType(tt.status))
		})
	}
}

func TestShippingType_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status ShippingType
		expect int32
	}{
		{
			name:   "success",
			status: ShippingTypeNormal,
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
		address     *Address
		expect      *OrderFulfillment
	}{
		{
			name: "success",
			fulfillment: &entity.OrderFulfillment{
				ID:                "fulfillment-id",
				OrderID:           "order-id",
				AddressRevisionID: 1,
				TrackingNumber:    "",
				Status:            entity.FulfillmentStatusFulfilled,
				ShippingCarrier:   entity.ShippingCarrierUnknown,
				ShippingType:      entity.ShippingTypeNormal,
				BoxNumber:         1,
				BoxSize:           entity.ShippingSize60,
				BoxRate:           80,
				CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				ShippedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			address: &Address{
				Address: response.Address{
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-5678",
				},
				revisionID: 1,
			},
			expect: &OrderFulfillment{
				OrderFulfillment: response.OrderFulfillment{
					FulfillmentID:   "fulfillment-id",
					TrackingNumber:  "",
					Status:          FulfillmentStatusFulfilled.Response(),
					ShippingCarrier: ShippingCarrierUnknown.Response(),
					ShippingType:    ShippingTypeNormal.Response(),
					BoxNumber:       1,
					BoxSize:         ShippingSize60.Response(),
					BoxRate:         80,
					ShippedAt:       1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-5678",
					},
				},
				orderID: "order-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderFulfillment(tt.fulfillment, tt.address))
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
					FulfillmentID:   "fulfillment-id",
					TrackingNumber:  "",
					Status:          FulfillmentStatusFulfilled.Response(),
					ShippingCarrier: ShippingCarrierUnknown.Response(),
					ShippingType:    ShippingTypeNormal.Response(),
					BoxNumber:       1,
					BoxSize:         ShippingSize60.Response(),
					BoxRate:         80,
					ShippedAt:       1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-5678",
					},
				},
				orderID: "order-id",
			},
			expect: &response.OrderFulfillment{
				FulfillmentID:   "fulfillment-id",
				TrackingNumber:  "",
				Status:          FulfillmentStatusFulfilled.Response(),
				ShippingCarrier: ShippingCarrierUnknown.Response(),
				ShippingType:    ShippingTypeNormal.Response(),
				BoxNumber:       1,
				BoxSize:         ShippingSize60.Response(),
				BoxRate:         80,
				ShippedAt:       1640962800,
				Address: &response.Address{
					Lastname:       "&.",
					Firstname:      "購入者",
					PostalCode:     "1000014",
					PrefectureCode: 13,
					City:           "千代田区",
					AddressLine1:   "永田町1-7-1",
					AddressLine2:   "",
					PhoneNumber:    "090-1234-5678",
				},
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

func TestOrderFulfillments(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fulfillments entity.OrderFulfillments
		addresses    map[int64]*Address
		expect       OrderFulfillments
	}{
		{
			name: "success",
			fulfillments: entity.OrderFulfillments{
				{
					ID:                "fulfillment-id",
					OrderID:           "order-id",
					AddressRevisionID: 1,
					TrackingNumber:    "",
					Status:            entity.FulfillmentStatusFulfilled,
					ShippingCarrier:   entity.ShippingCarrierUnknown,
					ShippingType:      entity.ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           entity.ShippingSize60,
					BoxRate:           80,
					CreatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
					ShippedAt:         jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			addresses: map[int64]*Address{
				1: {
					Address: response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-5678",
					},
					revisionID: 1,
				},
			},
			expect: OrderFulfillments{
				{
					OrderFulfillment: response.OrderFulfillment{
						FulfillmentID:   "fulfillment-id",
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingType:    ShippingTypeNormal.Response(),
						BoxNumber:       1,
						BoxSize:         ShippingSize60.Response(),
						BoxRate:         80,
						ShippedAt:       1640962800,
						Address: &response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "090-1234-5678",
						},
					},
					orderID: "order-id",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewOrderFulfillments(tt.fulfillments, tt.addresses))
		})
	}
}

func TestOrderFulfillments_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		fulfillments OrderFulfillments
		expect       []*response.OrderFulfillment
	}{
		{
			name: "success",
			fulfillments: OrderFulfillments{
				{
					OrderFulfillment: response.OrderFulfillment{
						FulfillmentID:   "fulfillment-id",
						TrackingNumber:  "",
						Status:          FulfillmentStatusFulfilled.Response(),
						ShippingCarrier: ShippingCarrierUnknown.Response(),
						ShippingType:    ShippingTypeNormal.Response(),
						BoxNumber:       1,
						BoxSize:         ShippingSize60.Response(),
						BoxRate:         80,
						ShippedAt:       1640962800,
						Address: &response.Address{
							Lastname:       "&.",
							Firstname:      "購入者",
							PostalCode:     "1000014",
							PrefectureCode: 13,
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "090-1234-5678",
						},
					},
					orderID: "order-id",
				},
			},
			expect: []*response.OrderFulfillment{
				{
					FulfillmentID:   "fulfillment-id",
					TrackingNumber:  "",
					Status:          FulfillmentStatusFulfilled.Response(),
					ShippingCarrier: ShippingCarrierUnknown.Response(),
					ShippingType:    ShippingTypeNormal.Response(),
					BoxNumber:       1,
					BoxSize:         ShippingSize60.Response(),
					BoxRate:         80,
					ShippedAt:       1640962800,
					Address: &response.Address{
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-5678",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.fulfillments.Response())
		})
	}
}
