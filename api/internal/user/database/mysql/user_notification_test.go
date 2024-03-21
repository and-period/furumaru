package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserNotification(t *testing.T) {
	assert.NotNil(t, newUserNotification(nil))
}

func TestUserNotification_Get(t *testing.T) {
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
		userID string
	}
	type want struct {
		notification *entity.UserNotification
		err          error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success when exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				user := testUser("user-id", "test@example.com", "090-1234-1234", now())
				err = db.DB.WithContext(ctx).Create(&user).Error
				require.NoError(t, err)
				notification := testUserNotification("user-id", now())
				err = db.DB.WithContext(ctx).Create(&notification).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
			},
			want: want{
				notification: testUserNotification("user-id", now()),
				err:          nil,
			},
		},
		{
			name: "success when non exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				user := testUser("user-id", "test@example.com", "090-1234-1234", now())
				err = db.DB.WithContext(ctx).Create(&user).Error
				require.NoError(t, err)
			},
			args: args{
				userID: "user-id",
			},
			want: want{
				notification: testUserNotification("user-id", now()),
				err:          nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, userNotificationTable, userTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &userNotification{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.userID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.notification, actual)
		})
	}
}

func testUserNotification(userID string, now time.Time) *entity.UserNotification {
	return &entity.UserNotification{
		UserID:        userID,
		EmailDisabled: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
