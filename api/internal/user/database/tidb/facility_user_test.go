package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFacilityUser_GetByExternalID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admin := testAdmin("producer-id", "cognito-id", "test-admin01@and-period.jp", now())
	err = db.DB.Create(&admin).Error
	require.NoError(t, err)
	p := testProducer("producer-id", "", now())
	p.Admin = *admin
	err = db.DB.Create(&p).Error
	require.NoError(t, err)

	u := testFacilityUser("user-id", "producer-id", "test-user@and-period.jp", now())
	err = db.DB.Create(&u).Error
	require.NoError(t, err)
	err = db.DB.Create(&u.FacilityUser).Error
	require.NoError(t, err)

	type args struct {
		providerType entity.UserAuthProviderType
		externalID   string
		producerID   string
	}
	type want struct {
		facilityUser *entity.FacilityUser
		err          error
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
				providerType: entity.UserAuthProviderTypeLINE,
				externalID:   "user-id",
				producerID:   "producer-id",
			},
			want: want{
				facilityUser: &u.FacilityUser,
				err:          nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args:  args{},
			want: want{
				facilityUser: nil,
				err:          database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			tt.setup(ctx, t, db)

			db := &facilityUser{db: db, now: now}
			actual, err := db.GetByExternalID(ctx, tt.args.providerType, tt.args.externalID, tt.args.producerID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.facilityUser, actual)
		})
	}
}

func TestFacilityUser_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admin := testAdmin("producer-id", "cognito-id", "test-admin01@and-period.jp", now())
	err = db.DB.Create(&admin).Error
	require.NoError(t, err)
	p := testProducer("producer-id", "", now())
	p.Admin = *admin
	err = db.DB.Create(&p).Error
	require.NoError(t, err)

	type args struct {
		user *entity.User
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
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				user: testFacilityUser("user-id", "producer-id", "test-user@and-period.jp", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "duplicate user entity",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				u := testFacilityUser("user-id", "producer-id", "test-user@and-period.jp", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.FacilityUser).Error
				require.NoError(t, err)
			},
			args: args{
				user: testFacilityUser("user-id", "producer-id", "test-user@and-period.jp", now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, facilityUserTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &facilityUser{db: db, now: now}
			err = db.Create(ctx, tt.args.user)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestFacilityUser_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admin := testAdmin("producer-id", "cognito-id", "test-admin01@and-period.jp", now())
	err = db.DB.Create(&admin).Error
	require.NoError(t, err)
	p := testProducer("producer-id", "", now())
	p.Admin = *admin
	err = db.DB.Create(&p).Error
	require.NoError(t, err)

	u := testFacilityUser("user-id", "producer-id", "test-user@and-period.jp", now())
	err = db.DB.Create(&u).Error
	require.NoError(t, err)
	err = db.DB.Create(&u.FacilityUser).Error
	require.NoError(t, err)

	type args struct {
		userID string
		params *database.UpdateFacilityUserParams
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
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				userID: "user-id",
				params: &database.UpdateFacilityUserParams{
					Lastname:      "田中",
					Firstname:     "太郎",
					LastnameKana:  "たなか",
					FirstnameKana: "たろう",
					PhoneNumber:   "090-1234-5678",
					LastCheckInAt: now().AddDate(0, 0, 2),
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found user",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				userID: "not-found-user-id",
				params: &database.UpdateFacilityUserParams{
					Lastname:      "田中",
					Firstname:     "太郎",
					LastnameKana:  "たなか",
					FirstnameKana: "たろう",
					PhoneNumber:   "090-1234-5678",
					LastCheckInAt: now().AddDate(0, 0, 2),
				},
			},
			want: want{
				err: nil, // GORMのUpdatesは対象レコードが0件でもエラーを返さない
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()

			tt.setup(ctx, t, db)

			db := &facilityUser{db: db, now: now}
			err = db.Update(ctx, tt.args.userID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)

			// 更新が成功した場合、データが正しく更新されたか確認
			if tt.want.err == nil && tt.args.userID == "user-id" {
				var updated entity.FacilityUser
				err := db.db.DB.WithContext(ctx).
					Table(facilityUserTable).
					Where("user_id = ?", tt.args.userID).
					First(&updated).Error
				require.NoError(t, err)
				assert.Equal(t, tt.args.params.Lastname, updated.Lastname)
				assert.Equal(t, tt.args.params.Firstname, updated.Firstname)
				assert.Equal(t, tt.args.params.LastnameKana, updated.LastnameKana)
				assert.Equal(t, tt.args.params.FirstnameKana, updated.FirstnameKana)
				assert.Equal(t, tt.args.params.PhoneNumber, updated.PhoneNumber)
				assert.Equal(t, tt.args.params.LastCheckInAt, updated.LastCheckInAt)
			}
		})
	}
}

func TestFacilityUser_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	admin := testAdmin("producer-id", "cognito-id", "test-admin01@and-period.jp", now())
	err = db.DB.Create(&admin).Error
	require.NoError(t, err)
	p := testProducer("producer-id", "", now())
	p.Admin = *admin
	err = db.DB.Create(&p).Error
	require.NoError(t, err)

	type args struct {
		userID string
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
				u := testFacilityUser("user-id", "producer-id", "test-user@and-period.jp", now())
				err = db.DB.Create(&u).Error
				require.NoError(t, err)
				err = db.DB.Create(&u.FacilityUser).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "not found user",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				userID: "not-found-user-id",
			},
			want: want{
				err: nil, // GORMのUpdatesは対象レコードが0件でもエラーを返さない
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, facilityUserTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &facilityUser{db: db, now: now}
			err = db.Delete(ctx, tt.args.userID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testFacility(id, producerID, email string, now time.Time) *entity.FacilityUser {
	return &entity.FacilityUser{
		UserID:        id,
		ProducerID:    producerID,
		ProviderType:  entity.UserAuthProviderTypeLINE,
		ExternalID:    id,
		Email:         email,
		LastCheckInAt: now.AddDate(0, 0, 1),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
