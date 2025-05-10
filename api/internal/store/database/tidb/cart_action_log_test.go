package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartActionLog(t *testing.T) {
	assert.NotNil(t, NewCartActionLog(nil))
}

func TestCartActionLog_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
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
		log *entity.CartActionLog
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
				log: testCartActionLog("session-id", entity.CartActionLogTypeAddCartItem, now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				log := testCartActionLog("session-id", entity.CartActionLogTypeAddCartItem, now())
				err := db.DB.WithContext(ctx).Table(cartActionLogTable).Create(&log).Error
				require.NoError(t, err)
			},
			args: args{
				log: testCartActionLog("session-id", entity.CartActionLogTypeAddCartItem, now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			err := delete(ctx, cartActionLogTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &cartActionLog{db: db, now: now}
			err = db.Create(ctx, tt.args.log)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testCartActionLog(sessionID string, actionType entity.CartActionLogType, now time.Time) *entity.CartActionLog {
	return &entity.CartActionLog{
		SessionID: sessionID,
		Type:      actionType,
		UserID:    "user-id",
		UserAgent: "user-agent",
		ClientIP:  "client-ip",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
