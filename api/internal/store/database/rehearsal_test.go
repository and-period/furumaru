package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	mock_dynamodb "github.com/and-period/furumaru/api/mock/pkg/dynamodb"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRehearsal(t *testing.T) {
	assert.NotNil(t, NewRehearsal(nil))
}

func TestRehearsal_Get(t *testing.T) {
	t.Parallel()
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}
	type mocks struct {
		dynamodb *mock_dynamodb.MockClient
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		liveID    string
		expect    *entity.Rehearsal
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, m *mocks) {
				keys := map[string]interface{}{"live_id": "live-id"}
				rehearsal := &entity.Rehearsal{LiveID: "live-id"}
				m.dynamodb.EXPECT().Get(ctx, keys, rehearsal).Return(nil)
			},
			liveID:    "live-id",
			expect:    &entity.Rehearsal{LiveID: "live-id"},
			expectErr: nil,
		},
		{
			name: "failed to get",
			setup: func(ctx context.Context, m *mocks) {
				keys := map[string]interface{}{"live_id": "live-id"}
				m.dynamodb.EXPECT().Get(ctx, keys, gomock.Any()).Return(assert.AnError)
			},
			liveID:    "live-id",
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := &mocks{
				dynamodb: mock_dynamodb.NewMockClient(ctrl),
			}
			db := &rehearsal{
				db:  m.dynamodb,
				now: now,
			}
			tt.setup(ctx, m)

			actual, err := db.Get(ctx, tt.liveID)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestRehearsal_Create(t *testing.T) {
	t.Parallel()
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}
	type mocks struct {
		dynamodb *mock_dynamodb.MockClient
	}
	r := &entity.Rehearsal{
		LiveID:    "live-id",
		CreatedAt: now(),
		UpdatedAt: now(),
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		rehearsal *entity.Rehearsal
		expect    error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, m *mocks) {
				m.dynamodb.EXPECT().Insert(ctx, r).Return(nil)
			},
			rehearsal: &entity.Rehearsal{LiveID: "live-id"},
			expect:    nil,
		},
		{
			name: "failed to insert",
			setup: func(ctx context.Context, m *mocks) {
				m.dynamodb.EXPECT().Insert(ctx, r).Return(assert.AnError)
			},
			rehearsal: &entity.Rehearsal{LiveID: "live-id"},
			expect:    exception.ErrUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := &mocks{
				dynamodb: mock_dynamodb.NewMockClient(ctrl),
			}
			db := &rehearsal{
				db:  m.dynamodb,
				now: now,
			}
			tt.setup(ctx, m)

			err := db.Create(ctx, tt.rehearsal)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}
