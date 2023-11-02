package mysql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddress_List(t *testing.T) {
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
		params *database.ListAddressesParams
	}
	type want struct {
		addresses entity.Addresses
		err       error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListAddressesParams{
					UserID: "user-id",
					Limit:  20,
					Offset: 1,
				},
			},
			want: want{
				addresses: addresses[1:],
				err:       nil,
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
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.addresses, actual)
		})
	}
}

func TestAddress_Count(t *testing.T) {
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
		params *database.ListAddressesParams
	}
	type want struct {
		total int64
		err   error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListAddressesParams{
					UserID: "user-id",
					Limit:  20,
					Offset: 1,
				},
			},
			want: want{
				total: 2,
				err:   nil,
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
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

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
		err       error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				addressIDs: []string{"address-id01", "address-id02"},
			},
			want: want{
				addresses: addresses,
				err:       nil,
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
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.addresses, actual)
		})
	}
}

func TestAddress_Get(t *testing.T) {
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

	a := testAddress("address-id", "user-id", now())
	err = db.DB.Create(&a).Error
	require.NoError(t, err)

	type args struct {
		addressID string
	}
	type want struct {
		address *entity.Address
		err     error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				addressID: "address-id",
			},
			want: want{
				address: a,
				err:     nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				addressID: "",
			},
			want: want{
				address: nil,
				err:     database.ErrNotFound,
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
			actual, err := db.Get(ctx, tt.args.addressID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.address, actual)
		})
	}
}

func TestAddress_Create(t *testing.T) {
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

	a := testAddress("address-id", "user-id", now())

	type args struct {
		address *entity.Address
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success when is_default is true",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				address := testAddress("other-id", "user-id", now())
				address.IsDefault = true
				address.PhoneNumber = "+818012345678"
				err := db.DB.Create(&address).Error
				require.NoError(t, err)
			},
			args: args{
				address: func() *entity.Address {
					address := testAddress("address-id", "user-id", now())
					address.IsDefault = true
					return address
				}(),
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "success when is_default is false",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				address: a,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err := db.DB.Create(&a).Error
				require.NoError(t, err)
			},
			args: args{
				address: a,
			},
			want: want{
				err: database.ErrAlreadyExists,
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
			err := db.Create(ctx, tt.args.address)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAddress_Update(t *testing.T) {
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

	a := testAddress("address-id", "user-id", now())

	type args struct {
		addressID string
		userID    string
		params    *database.UpdateAddressParams
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success when is_default is true",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err = db.DB.Create(&a).Error
				require.NoError(t, err)
				address := testAddress("other-id", "user-id", now())
				address.IsDefault = true
				address.PhoneNumber = "+818012345678"
				err := db.DB.Create(&address).Error
				require.NoError(t, err)
			},
			args: args{
				addressID: "address-id",
				userID:    "user-id",
				params: &database.UpdateAddressParams{
					Lastname:     a.Lastname,
					Firstname:    a.Firstname,
					PostalCode:   a.PostalCode,
					Prefecture:   a.Prefecture,
					City:         a.City,
					AddressLine1: a.AddressLine1,
					AddressLine2: a.AddressLine2,
					PhoneNumber:  a.PhoneNumber,
					IsDefault:    true,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success when is_default is false",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err = db.DB.Create(&a).Error
				require.NoError(t, err)
			},
			args: args{
				addressID: "address-id",
				userID:    "user-id",
				params: &database.UpdateAddressParams{
					Lastname:     a.Lastname,
					Firstname:    a.Firstname,
					PostalCode:   a.PostalCode,
					Prefecture:   a.Prefecture,
					City:         a.City,
					AddressLine1: a.AddressLine1,
					AddressLine2: a.AddressLine2,
					PhoneNumber:  a.PhoneNumber,
					IsDefault:    true,
				},
			},
			want: want{
				err: nil,
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
			err := db.Update(ctx, tt.args.addressID, tt.args.userID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testAddress(addressID, userID string, now time.Time) *entity.Address {
	return &entity.Address{
		ID:           addressID,
		UserID:       userID,
		Hash:         fmt.Sprintf("%s:%s", userID, addressID),
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
