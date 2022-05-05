package database

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	assert.NotNil(t, NewStore(nil))
}

func TestStore_List(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, storeTable)
	stores := make(entity.Stores, 2)
	stores[0] = testStore(1, "&.農園", now())
	stores[1] = testStore(2, "&.水産", now())
	err = m.db.DB.Create(&stores).Error
	require.NoError(t, err)

	type args struct {
		params *ListStoresParams
	}
	type want struct {
		stores entity.Stores
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
				params: &ListStoresParams{
					Limit:  1,
					Offset: 1,
				},
			},
			want: want{
				stores: stores[1:],
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

			tt.setup(ctx, t, m)

			db := &store{db: m.db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			for i := range actual {
				fillIgnoreStoreField(actual[i], now())
			}
			assert.ElementsMatch(t, tt.want.stores, actual)
		})
	}
}

func TestStore_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, err := newMocks(ctrl)
	require.NoError(t, err)
	current := jst.Date(2022, 1, 2, 18, 30, 0, 0)
	now := func() time.Time {
		return current
	}

	_ = m.dbDelete(ctx, storeTable)
	s := testStore(1, "&.農園", now())
	err = m.db.DB.Create(&s).Error
	require.NoError(t, err)

	type args struct {
		storeID int64
	}
	type want struct {
		store  *entity.Store
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
				storeID: 1,
			},
			want: want{
				store:  s,
				hasErr: false,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, m *mocks) {},
			args: args{
				storeID: 0,
			},
			want: want{
				store:  nil,
				hasErr: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, m)

			db := &store{db: m.db, now: now}
			actual, err := db.Get(ctx, tt.args.storeID)
			if tt.want.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			fillIgnoreStoreField(actual, now())
			assert.Equal(t, tt.want.store, actual)
		})
	}
}

func testStore(id int64, name string, now time.Time) *entity.Store {
	return &entity.Store{
		ID:           id,
		Name:         name,
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func fillIgnoreStoreField(s *entity.Store, now time.Time) {
	if s == nil {
		return
	}
	s.CreatedAt = now
	s.UpdatedAt = now
}
