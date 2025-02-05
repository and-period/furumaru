package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAdminAuthProvider_List(t *testing.T) {
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

	admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	providers := make(entity.AdminAuthProviders, 1)
	providers[0] = testAdminAuthProvider("admin-id", entity.AdminAuthProviderTypeGoogle, now())
	err = db.DB.WithContext(ctx).Table(adminAuthProviderTable).Create(&providers).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminAuthProvidersParams
	}
	type want struct {
		providers entity.AdminAuthProviders
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
				params: &database.ListAdminAuthProvidersParams{
					AdminID: "admin-id",
				},
			},
			want: want{
				providers: providers,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &adminAuthProvider{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			require.ErrorIs(t, err, tt.want.err)
			require.ElementsMatch(t, tt.want.providers, actual)
		})
	}
}

func TestAdminAuthProvider_Get(t *testing.T) {
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

	admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	provider := testAdminAuthProvider("admin-id", entity.AdminAuthProviderTypeGoogle, now())
	err = db.DB.WithContext(ctx).Table(adminAuthProviderTable).Create(&provider).Error
	require.NoError(t, err)

	type args struct {
		adminID      string
		providerType entity.AdminAuthProviderType
	}
	type want struct {
		provider *entity.AdminAuthProvider
		err      error
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
				adminID:      "admin-id",
				providerType: entity.AdminAuthProviderTypeGoogle,
			},
			want: want{
				provider: provider,
				err:      nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				adminID:      "",
				providerType: entity.AdminAuthProviderTypeUnknown,
			},
			want: want{
				provider: nil,
				err:      database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &adminAuthProvider{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.adminID, tt.args.providerType)
			require.ErrorIs(t, err, tt.want.err)
			require.Equal(t, tt.want.provider, actual)
		})
	}
}

func TestAdminAuthProvider_Upsert(t *testing.T) {
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

	admin := testAdmin("admin-id", "cognito-id", "test-admin@and-period.jp", now())
	err = db.DB.WithContext(ctx).Table(adminTable).Create(&admin).Error
	require.NoError(t, err)

	type args struct {
		provider *entity.AdminAuthProvider
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
			name:  "success create",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				provider: testAdminAuthProvider("admin-id", entity.AdminAuthProviderTypeGoogle, now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success update",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				provider := testAdminAuthProvider("admin-id", entity.AdminAuthProviderTypeGoogle, now())
				err = db.DB.WithContext(ctx).Table(adminAuthProviderTable).Create(&provider).Error
				require.NoError(t, err)
			},
			args: args{
				provider: testAdminAuthProvider("admin-id", entity.AdminAuthProviderTypeGoogle, now()),
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

			err := delete(ctx, adminAuthProviderTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &adminAuthProvider{db: db, now: now}
			err = db.Upsert(ctx, tt.args.provider)
			require.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testAdminAuthProvider(adminID string, providerType entity.AdminAuthProviderType, now time.Time) *entity.AdminAuthProvider {
	return &entity.AdminAuthProvider{
		AdminID:      adminID,
		ProviderType: providerType,
		AccountID:    "account-id",
		Email:        "test@example.com",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
