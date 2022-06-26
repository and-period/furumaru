package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestNotification(t *testing.T) {
	assert.NotNil(t, NewNotification(nil))
}

func TestNotification_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 6, 26, 19, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	type args struct {
		notification *entity.Notification
	}
	type want struct {
		hasErr bool
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, m *mocks)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				notification: testNotification("notification-id", "creator", "title", datatypes.JSON([]byte(`[1,2,3]`)), true, now()),
			},
			want: want{
				hasErr: false,
			},
		},
		{
			name: "failed to duplicate entry",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {
				n := testNotification("notification-id", "creator", "title", datatypes.JSON([]byte(`[1,2,3]`)), true, now())
				err = m.db.DB.Create(&n).Error
				require.NoError(t, err)
			},
			args: args{
				notification: testNotification("notification-id", "creator", "title", datatypes.JSON([]byte(`[1,2,3]`)), true, now()),
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

			err := m.dbDelete(ctx, notificationTable)
			require.NoError(t, err)
			tt.setup(ctx, t, m)

			db := &notification{db: m.db, now: now}
			err = db.Create(ctx, tt.args.notification)
			assert.Equal(t, tt.want.hasErr, err != nil, err)
		})
	}
}

func testNotification(id, creatorName, title string, targets datatypes.JSON, public bool, now time.Time) *entity.Notification {
	return &entity.Notification{
		ID:          id,
		CreatedBy:   id,
		CreatorName: creatorName,
		UpdatedBy:   id,
		Title:       title,
		Body:        title,
		TargetsJSON: targets,
		PublishedAt: now,
		Public:      public,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
