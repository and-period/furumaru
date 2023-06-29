package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactRead(t *testing.T) {
	assert.NotNil(t, NewContactRead(nil))
}

func TestContactRead_Get(t *testing.T) {
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

	category := testContactCategory("category-id", "お問い合わせ種別", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)

	contact := testContact("contact-id", now())
	err = db.DB.Create(&contact).Error
	require.NoError(t, err)

	cr := testContactRead("contact-read-id", "contact-id", "user-id", now())
	err = db.DB.Create(&cr).Error
	require.NoError(t, err)

	type args struct {
		contactID string
		userID    string
	}
	type want struct {
		contactRead *entity.ContactRead
		hasErr      bool
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
				contactID: "contact-id",
				userID:    "user-id",
			},
			want: want{
				contactRead: cr,
				hasErr:      false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				contactID: "other-id",
				userID:    "user-id",
			},
			want: want{
				contactRead: nil,
				hasErr:      true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			tt.setup(ctx, t, db)

			db := &contactRead{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.contactID, tt.args.userID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.contactRead, actual)
		})
	}
}

func TestContactRead_Create(t *testing.T) {
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

	category := testContactCategory("category-id", "お問い合わせ種別", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)

	contact := testContact("contact-id", now())
	err = db.DB.Create(&contact).Error
	require.NoError(t, err)

	cr := testContactRead("contact-read-id", "contact-id", "user-id", now())

	type args struct {
		contactRead *entity.ContactRead
	}
	type want struct {
		hasErr bool
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
				contactRead: cr,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				err := db.DB.Create(cr).Error
				require.NoError(t, err)
			},
			args: args{
				contactRead: cr,
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, contactReadTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &contactRead{db: db, now: now}
			err = db.Create(ctx, tt.args.contactRead)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func TestContactRead_UpdateRead(t *testing.T) {
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
	category := testContactCategory("category-id", "お問い合わせ種別", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)

	contacts := make([]*entity.Contact, 2)
	contacts[0] = testContact("contact-id01", now())
	contacts[1] = testContact("contact-id02", now())
	err = db.DB.Create(&contacts).Error
	require.NoError(t, err)

	contactReads := make([]*entity.ContactRead, 3)
	contactReads[0] = testContactRead("contact-read-id01", "contact-id01", "user-id", now())
	contactReads[1] = testContactRead("contact-read-id02", "contact-id01", "admin-id", now())
	contactReads[2] = testContactRead("contact-read-id03", "contact-id02", "", now())
	err = db.DB.Create(&contactReads).Error
	require.NoError(t, err)

	type args struct {
		params *UpdateContactReadFlagParams
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *database.Client)
		args  args
		want  want
	}{
		{
			name:  "success with user-id",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &UpdateContactReadFlagParams{
					ContactID: "contact-id01",
					UserID:    "user-id",
					Read:      true,
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "success with admin-id",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &UpdateContactReadFlagParams{
					ContactID: "contact-id01",
					UserID:    "admin-id",
					Read:      true,
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "success without userID",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &UpdateContactReadFlagParams{
					ContactID: "contact-id02",
					UserID:    "",
					Read:      true,
				},
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &UpdateContactReadFlagParams{
					ContactID: "contact-id04",
				},
			},
			want: want{
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &contactRead{db: db, now: now}
			err = db.UpdateRead(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testContactRead(id, contactID, userID string, now time.Time) *entity.ContactRead {
	return &entity.ContactRead{
		ID:        id,
		ContactID: contactID,
		UserID:    userID,
		UserType:  entity.ContactUserTypeUnknown,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
