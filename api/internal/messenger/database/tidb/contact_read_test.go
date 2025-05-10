package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContactRead(t *testing.T) {
	assert.NotNil(t, NewContactRead(nil))
}

func TestContactRead_GetByContactIDAndUserID(t *testing.T) {
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
		err         error
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
				contactID: "contact-id",
				userID:    "user-id",
			},
			want: want{
				contactRead: cr,
				err:         nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				contactID: "other-id",
				userID:    "user-id",
			},
			want: want{
				contactRead: nil,
				err:         database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &contactRead{db: db, now: now}
			actual, err := db.GetByContactIDAndUserID(ctx, tt.args.contactID, tt.args.userID)
			assert.ErrorIs(t, err, tt.want.err)
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
				contactRead: cr,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err := db.DB.Create(cr).Error
				require.NoError(t, err)
			},
			args: args{
				contactRead: cr,
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, contactReadTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &contactRead{db: db, now: now}
			err = db.Create(ctx, tt.args.contactRead)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestContactRead_Update(t *testing.T) {
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
		params *database.UpdateContactReadParams
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
			name:  "success with user-id",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.UpdateContactReadParams{
					ContactID: "contact-id01",
					UserID:    "user-id",
					Read:      true,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "success with admin-id",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.UpdateContactReadParams{
					ContactID: "contact-id01",
					UserID:    "admin-id",
					Read:      true,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name:  "success without userID",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.UpdateContactReadParams{
					ContactID: "contact-id02",
					UserID:    "",
					Read:      true,
				},
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

			tt.setup(ctx, t, db)

			db := &contactRead{db: db, now: now}
			err = db.Update(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
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
