package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAddresses(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 12, 6, 18, 30, 0, 0)
	params := &database.ListAddressesParams{
		UserID: "user-id",
		Limit:  20,
		Offset: 0,
	}
	addresses := entity.Addresses{
		{
			ID:             "address-id",
			UserID:         "user-id",
			Hash:           "c1f66591133a1a70cc6b29f21ede4389efe6864bb7ade2e17f734471352df1a9",
			IsDefault:      true,
			Lastname:       "&.",
			Firstname:      "購入者",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "+819012345678",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *user.ListAddressesInput
		expect      entity.Addresses
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().List(gomock.Any(), params).Return(addresses, nil)
				mocks.db.Address.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListAddressesInput{
				UserID: "user-id",
				Limit:  20,
				Offset: 0,
			},
			expect:      addresses,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.ListAddressesInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to list addresses",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Address.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &user.ListAddressesInput{
				UserID: "user-id",
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count addresses",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().List(gomock.Any(), params).Return(addresses, nil)
				mocks.db.Address.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &user.ListAddressesInput{
				UserID: "user-id",
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListAddresses(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetAddresses(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 12, 6, 18, 30, 0, 0)
	addresses := entity.Addresses{
		{
			ID:             "address-id",
			UserID:         "user-id",
			Hash:           "c1f66591133a1a70cc6b29f21ede4389efe6864bb7ade2e17f734471352df1a9",
			IsDefault:      true,
			Lastname:       "&.",
			Firstname:      "購入者",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "+819012345678",
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.MultiGetAddressesInput
		expect    entity.Addresses
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().MultiGet(ctx, []string{"address-id"}).Return(addresses, nil)
			},
			input: &user.MultiGetAddressesInput{
				AddressIDs: []string{"address-id"},
			},
			expect:    addresses,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.MultiGetAddressesInput{
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
			input: &user.MultiGetAddressesInput{
				AddressIDs: []string{"address-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
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

func TestGetAddress(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 12, 6, 18, 30, 0, 0)
	address := func(userID string) *entity.Address {
		return &entity.Address{
			ID:             "address-id",
			UserID:         userID,
			Hash:           "c1f66591133a1a70cc6b29f21ede4389efe6864bb7ade2e17f734471352df1a9",
			IsDefault:      true,
			Lastname:       "&.",
			Firstname:      "購入者",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "+819012345678",
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.GetAddressInput
		expect    *entity.Address
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id").Return(address("user-id"), nil)
			},
			input: &user.GetAddressInput{
				AddressID: "address-id",
				UserID:    "user-id",
			},
			expect:    address("user-id"),
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.GetAddressInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id").Return(nil, assert.AnError)
			},
			input: &user.GetAddressInput{
				AddressID: "address-id",
				UserID:    "user-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
		{
			name: "belongs to another user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id").Return(address("other-id"), nil)
			},
			input: &user.GetAddressInput{
				AddressID: "address-id",
				UserID:    "user-id",
			},
			expect:    nil,
			expectErr: exception.ErrForbidden,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetAddress(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateAddress(t *testing.T) {
	t.Parallel()

	address := func(addressID string) *entity.Address {
		return &entity.Address{
			ID:             addressID,
			UserID:         "user-id",
			Hash:           "c1f66591133a1a70cc6b29f21ede4389efe6864bb7ade2e17f734471352df1a9",
			IsDefault:      true,
			Lastname:       "&.",
			Firstname:      "購入者",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "+819012345678",
		}
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.CreateAddressInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, actual *entity.Address) error {
						expect := address(actual.ID)
						assert.Equal(t, expect, actual)
						return nil
					})
			},
			input: &user.CreateAddressInput{
				UserID:         "user-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
				IsDefault:      true,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.CreateAddressInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "failed to new address",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &user.CreateAddressInput{
				UserID:         "user-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				PrefectureCode: -1,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
				IsDefault:      true,
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &user.CreateAddressInput{
				UserID:         "user-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
				IsDefault:      true,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateAddress(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateAddress(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 12, 6, 18, 30, 0, 0)
	address := func(userID string) *entity.Address {
		return &entity.Address{
			ID:             "address-id",
			UserID:         userID,
			Hash:           "c1f66591133a1a70cc6b29f21ede4389efe6864bb7ade2e17f734471352df1a9",
			IsDefault:      true,
			Lastname:       "&.",
			Firstname:      "購入者",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "+819012345678",
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}
	params := &database.UpdateAddressParams{
		Lastname:       "&.",
		Firstname:      "購入者",
		PostalCode:     "1000014",
		PrefectureCode: 13,
		City:           "千代田区",
		AddressLine1:   "永田町1-7-1",
		AddressLine2:   "",
		PhoneNumber:    "+819012345678",
		IsDefault:      true,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.UpdateAddressInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id", "user_id").Return(address("user-id"), nil)
				mocks.db.Address.EXPECT().Update(ctx, "address-id", "user-id", params).Return(nil)
			},
			input: &user.UpdateAddressInput{
				AddressID:      "address-id",
				UserID:         "user-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
				IsDefault:      true,
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.UpdateAddressInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id", "user_id").Return(nil, assert.AnError)
			},
			input: &user.UpdateAddressInput{
				AddressID:      "address-id",
				UserID:         "user-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
				IsDefault:      true,
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id", "user_id").Return(address("other-id"), nil)
			},
			input: &user.UpdateAddressInput{
				AddressID:      "address-id",
				UserID:         "user-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
				IsDefault:      true,
			},
			expectErr: exception.ErrForbidden,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id", "user_id").Return(address("user-id"), nil)
				mocks.db.Address.EXPECT().Update(ctx, "address-id", "user-id", params).Return(assert.AnError)
			},
			input: &user.UpdateAddressInput{
				AddressID:      "address-id",
				UserID:         "user-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
				IsDefault:      true,
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateAddress(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteAddress(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 12, 6, 18, 30, 0, 0)
	address := func(userID string) *entity.Address {
		return &entity.Address{
			ID:             "address-id",
			UserID:         userID,
			Hash:           "c1f66591133a1a70cc6b29f21ede4389efe6864bb7ade2e17f734471352df1a9",
			IsDefault:      true,
			Lastname:       "&.",
			Firstname:      "購入者",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "+819012345678",
			CreatedAt:      now,
			UpdatedAt:      now,
		}
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *user.DeleteAddressInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id", "user_id").Return(address("user-id"), nil)
				mocks.db.Address.EXPECT().Delete(ctx, "address-id").Return(nil)
			},
			input: &user.DeleteAddressInput{
				AddressID: "address-id",
				UserID:    "user-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &user.DeleteAddressInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id", "user_id").Return(nil, assert.AnError)
			},
			input: &user.DeleteAddressInput{
				AddressID: "address-id",
				UserID:    "user-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id", "user_id").Return(address("other-id"), nil)
			},
			input: &user.DeleteAddressInput{
				AddressID: "address-id",
				UserID:    "user-id",
			},
			expectErr: exception.ErrForbidden,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Address.EXPECT().Get(ctx, "address-id", "user_id").Return(address("user-id"), nil)
				mocks.db.Address.EXPECT().Delete(ctx, "address-id").Return(assert.AnError)
			},
			input: &user.DeleteAddressInput{
				AddressID: "address-id",
				UserID:    "user-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteAddress(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}