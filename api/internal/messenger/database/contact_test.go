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

	contacts := make(entity.Contacts, 3)
	contacts[0] = testContact("contact-id01", now())
	contacts[1] = testContact("contact-id02", now())
	contacts[2] = testContact("contact-id03", now())
	err = db.DB.Create(&contacts).Error
	require.NoError(t, err)

	type args struct {
		params *ListContactsParams
	}
	type want struct {
		contacts entity.Contacts
		hasErr   bool
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
				params: &ListContactsParams{
					Limit:  2,
					Offset: 1,
				},
			},
			want: want{
				contacts: contacts[1:],
				hasErr:   false,
			},
		},
		{
			name:  "success with sort",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				params: &ListContactsParams{
					Orders: []*ListContactsOrder{
						{Key: entity.ContactOrderByPriority, OrderByASC: true},
					},
				},
			},
			want: want{
				contacts: contacts,
				hasErr:   false,
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
			assert.Equal(t, tt.want.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.want.contacts, actual)
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

	contacts := make(entity.Contacts, 3)
	contacts[0] = testContact("contact-id01", now())
	contacts[1] = testContact("contact-id02", now())
	contacts[2] = testContact("contact-id03", now())
	err = db.DB.Create(&contacts).Error
	require.NoError(t, err)

	type args struct {
		params *ListContactsParams
	}
	type want struct {
		total  int64
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
				params: &ListContactsParams{},
			},
			want: want{
				total:  3,
				hasErr: false,
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
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
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

	c := testContact("contact-id", now())
	err = db.DB.Create(&c).Error
	require.NoError(t, err)

	type args struct {
		contactID string
	}
	type want struct {
		contact *entity.Contact
		hasErr  bool
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
			},
			want: want{
				contact: c,
				hasErr:  false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				contactID: "other-id",
			},
			want: want{
				contact: nil,
				hasErr:  true,
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
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
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

	c := testContact("contact-id", now())

	type args struct {
		contact *entity.Contact
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
				contact: c,
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				err = db.DB.Create(&c).Error
				require.NoError(t, err)
			},
			args: args{
				contact: c,
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

			err := delete(ctx, contactTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &contact{db: db, now: now}
			err = db.Create(ctx, tt.args.contact)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
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

	c := testContact("contact-id", now())

	type args struct {
		contactID string
		params    *UpdateContactParams
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
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				err = db.DB.Create(&c).Error
				require.NoError(t, err)
			},
			args: args{
				contactID: "contact-id",
				params: &UpdateContactParams{
					Status:   entity.ContactStatusDone,
					Priority: entity.ContactPriorityHigh,
					Note:     "対応メモです。",
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
				contactID: "contact-id",
				params: &UpdateContactParams{
					Status:   entity.ContactStatusDone,
					Priority: entity.ContactPriorityHigh,
					Note:     "対応メモです。",
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

			err := delete(ctx, contactTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &contact{db: db, now: now}
			err = db.Update(ctx, tt.args.contactID, tt.args.params)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
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

	c := testContact("contact-id", now())

	type args struct {
		contactID string
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
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {
				err = db.DB.Create(&c).Error
				require.NoError(t, err)
			},
			args: args{
				contactID: "contact-id",
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *database.Client) {},
			args: args{
				contactID: "contact-id",
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

			err := delete(ctx, contactTable)
			require.NoError(t, err)
			tt.setup(ctx, t, db)

			db := &contact{db: db, now: now}
			err = db.Delete(ctx, tt.args.contactID)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testContact(id string, now time.Time) *entity.Contact {
	return &entity.Contact{
		ID:          id,
		Title:       "お問い合わせ件名",
		Content:     "お問い合わせ内容です。",
		Username:    "あんど どっと",
		Email:       "test-user@and-period.jp",
		PhoneNumber: "+819012345678",
		Status:      entity.ContactStatusInprogress,
		Priority:    entity.ContactPriorityMiddle,
		Note:        "対応者のメモです",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
