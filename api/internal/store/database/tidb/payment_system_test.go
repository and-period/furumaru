package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentSystem(t *testing.T) {
	assert.NotNil(t, NewPaymentSystem(nil))
}

func TestPaymentSystem_MultiGet(t *testing.T) {
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

	systems := make(entity.PaymentSystems, 2)
	systems[0] = testPaymentSystem(entity.PaymentMethodTypeCreditCard, now())
	systems[1] = testPaymentSystem(entity.PaymentMethodTypePayPay, now())
	err = db.DB.Create(&systems).Error
	require.NoError(t, err)

	type args struct {
		methodTypes []entity.PaymentMethodType
	}
	type want struct {
		systems entity.PaymentSystems
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
				methodTypes: []entity.PaymentMethodType{
					entity.PaymentMethodTypeCreditCard,
					entity.PaymentMethodTypePayPay,
				},
			},
			want: want{
				systems: systems,
				err:     nil,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &paymentSystem{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.methodTypes)
			assert.ErrorIs(t, err, tt.want.err)
			assert.ElementsMatch(t, tt.want.systems, actual)
		})
	}
}

func TestPaymentSystem_Get(t *testing.T) {
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

	system := testPaymentSystem(entity.PaymentMethodTypeCreditCard, now())
	err = db.DB.Create(&system).Error
	require.NoError(t, err)

	type args struct {
		methodType entity.PaymentMethodType
	}
	type want struct {
		system *entity.PaymentSystem
		err    error
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
				methodType: entity.PaymentMethodTypeCreditCard,
			},
			want: want{
				system: system,
				err:    nil,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &paymentSystem{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.methodType)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.system, actual)
		})
	}
}

func TestPaymentSystem_Update(t *testing.T) {
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

	type args struct {
		methodType entity.PaymentMethodType
		status     entity.PaymentSystemStatus
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
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				system := testPaymentSystem(entity.PaymentMethodTypeCreditCard, now())
				err := db.DB.Create(&system).Error
				require.NoError(t, err)
			},
			args: args{
				methodType: entity.PaymentMethodTypeCreditCard,
				status:     entity.PaymentSystemStatusOutage,
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, paymentSystemTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &paymentSystem{db: db, now: now}
			err = db.Update(ctx, tt.args.methodType, tt.args.status)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testPaymentSystem(methodType entity.PaymentMethodType, now time.Time) *entity.PaymentSystem {
	return &entity.PaymentSystem{
		MethodType: methodType,
		Status:     entity.PaymentSystemStatusInUse,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
