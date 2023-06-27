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

func testContactRead(id, contactID, userID string, now time.Time) *entity.ContactRead {
	return &entity.ContactRead{
		ID:        id,
		ContactID: contactID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
