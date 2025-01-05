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

func TestContact(t *testing.T) {
	assert.NotNil(t, NewContact(nil))
}

func TestContact_List(t *testing.T) {
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

	contacts := make(entity.Contacts, 3)
	contacts[0] = testContact("contact-id01", now())
	contacts[1] = testContact("contact-id02", now())
	contacts[2] = testContact("contact-id03", now())
	err = db.DB.Create(&contacts).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListContactsParams
	}
	type want struct {
		contacts entity.Contacts
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
				params: &database.ListContactsParams{
					Limit:  3,
					Offset: 0,
				},
			},
			want: want{
				contacts: contacts,
				err:      nil,
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

			db := &contact{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.contacts, actual)
		})
	}
}

func TestContact_Count(t *testing.T) {
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

	contacts := make(entity.Contacts, 3)
	contacts[0] = testContact("contact-id01", now())
	contacts[1] = testContact("contact-id02", now())
	contacts[2] = testContact("contact-id03", now())
	err = db.DB.Create(&contacts).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListContactsParams
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
				params: &database.ListContactsParams{
					Limit:  3,
					Offset: 0,
				},
			},
			want: want{
				total: 3,
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

			db := &contact{db: db, now: now}
			actual, err := db.Count(ctx)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestContact_Get(t *testing.T) {
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

	c := testContact("contact-id", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		contactID string
	}
	type want struct {
		contact *entity.Contact
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
				contactID: "contact-id",
			},
			want: want{
				contact: c,
				err:     nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				contactID: "other-id",
			},
			want: want{
				contact: nil,
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

			db := &contact{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.contactID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.contact, actual)
		})
	}
}

func TestContact_Create(t *testing.T) {
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

	c := testContact("contact-id", now())

	type args struct {
		contact *entity.Contact
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
				contact: c,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				err = db.DB.Create(&c).Error
				require.NoError(t, err)
			},
			args: args{
				contact: c,
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, contactTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &contact{db: db, now: now}
			err = db.Create(ctx, tt.args.contact)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestContact_Update(t *testing.T) {
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

	type args struct {
		contactID string
		params    *database.UpdateContactParams
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
				contact := testContact("contact-id", now())
				err = db.DB.Create(&contact).Error
				require.NoError(t, err)
			},
			args: args{
				contactID: "contact-id",
				params: &database.UpdateContactParams{
					Title:       "件名",
					CategoryID:  "category-id",
					Content:     "内容です。",
					Username:    "あんど ぴりおど",
					UserID:      "user-id",
					Email:       "test-user@and-period.jp",
					PhoneNumber: "+819012345678",
					Status:      entity.ContactStatusDone,
					ResponderID: "responder-id",
					Note:        "メモです",
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, contactTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &contact{db: db, now: now}
			err = db.Update(ctx, tt.args.contactID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestContact_Delete(t *testing.T) {
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

	type args struct {
		contactID string
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
				contact := testContact("contact-id", now())
				err = db.DB.Create(&contact).Error
				require.NoError(t, err)
			},
			args: args{
				contactID: "contact-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, contactTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &contact{db: db, now: now}
			err = db.Delete(ctx, tt.args.contactID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testContact(id string, now time.Time) *entity.Contact {
	return &entity.Contact{
		ID:          id,
		Title:       "お問い合わせ件名",
		CategoryID:  "category-id",
		Content:     "お問い合わせ内容です。",
		Username:    "あんど ぴりおど",
		UserID:      "user-id",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+819012345678",
		Status:      entity.ContactStatusInprogress,
		ResponderID: "responder-id",
		Note:        "対応者のメモです",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
