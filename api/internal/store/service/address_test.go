package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestMultiGetAddresses(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 12, 6, 18, 30, 0, 0)
	addresses := entity.Addresses{
		{
			ID:           "address-id",
			UserID:       "user-id",
			Hash:         "789ef22a79a364f95c66a3d3b1fda213c1316a6c7f8b6306b493d8c46d2dce75",
			IsDefault:    true,
			Lastname:     "&.",
			Firstname:    "購入者",
			PostalCode:   "1000014",
			Prefecture:   13,
			City:         "千代田区",
			AddressLine1: "永田町1-7-1",
			AddressLine2: "",
			PhoneNumber:  "+819012345678",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetAddressesInput
		expect    entity.Addresses
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().MultiGet(ctx, []string{"address-id"}).Return(addresses, nil)
			},
			input: &store.MultiGetAddressesInput{
				AddressIDs: []string{"address-id"},
			},
			expect:    addresses,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetAddressesInput{
				AddressIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get addresses",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().MultiGet(ctx, []string{"address-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetAddressesInput{
				AddressIDs: []string{"address-id"},
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetAddresses(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}
