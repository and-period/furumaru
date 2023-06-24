package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddress_MultiGet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	addresses := make(entity.Addresses, 2)
	addresses[0] = testAddress("address-id01", "user-id", now())
	addresses[1] = testAddress("address-id02", "user-id", now())
	err = db.DB.Create(&addresses).Error
	require.NoError(t, err)

	type args struct {
		addressIDs []string
	}
	type want struct {
		addresses entity.Addresses
		hasErr    bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				addressIDs: []string{"address-id01", "address-id02"},
			},
			want: want{
				addresses: addresses,
				hasErr:    false,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &address{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.addressIDs)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.addresses, actual)
		})
	}
}

func testAddress(id, userID string, now time.Time) *entity.Address {
	return &entity.Address{
		ID:           id,
		UserID:       userID,
		Hash:         entity.NewAddressHash(userID, "1000014", "永田町1-7-1", id),
		IsDefault:    false,
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
	}
}
