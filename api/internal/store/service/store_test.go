package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestListStores(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	params := &database.ListStoresParams{
		Limit:  20,
		Offset: 0,
	}
	stores := entity.Stores{
		{
			ID:           1,
			Name:         "&.農園",
			ThumbnailURL: "https://and-period.jp/thumbnail.png",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           2,
			Name:         "&.水産",
			ThumbnailURL: "https://and-period.jp/thumbnail.png",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.ListStoresInput
		expect    entity.Stores
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().List(ctx, params).Return(stores, nil)
			},
			input: &store.ListStoresInput{
				Limit:  20,
				Offset: 0,
			},
			expect:    stores,
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.ListStoresInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get stores",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().List(ctx, params).Return(nil, errmock)
			},
			input: &store.ListStoresInput{
				Limit:  20,
				Offset: 0,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			actual, err := service.ListStores(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestGetStore(t *testing.T) {
	t.Parallel()

	now := jst.Now()
	s := &entity.Store{
		ID:           1,
		Name:         "&.農園",
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetStoreInput
		expect    *entity.Store
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Get(ctx, int64(1)).Return(s, nil)
			},
			input: &store.GetStoreInput{
				StoreID: 1,
			},
			expect:    s,
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetStoreInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get store",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Get(ctx, int64(1)).Return(nil, errmock)
			},
			input: &store.GetStoreInput{
				StoreID: 1,
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			actual, err := service.GetStore(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateStore(t *testing.T) {
	t.Parallel()

	s := &entity.Store{
		Name: "&.農園",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateStoreInput
		expect    *entity.Store
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Create(ctx, s).Return(nil)
			},
			input: &store.CreateStoreInput{
				Name: "&.農園",
			},
			expect:    s,
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateStoreInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get store",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Create(ctx, s).Return(errmock)
			},
			input: &store.CreateStoreInput{
				Name: "&.農園",
			},
			expect:    nil,
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			actual, err := service.CreateStore(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUpdateStore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateStoreInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Update(ctx, int64(1), "&.農園", "https://and-period.jp/thumbnail.png").Return(nil)
			},
			input: &store.UpdateStoreInput{
				StoreID:      1,
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
			},
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateStoreInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get store",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Update(ctx, int64(1), "&.農園", "https://and-period.jp/thumbnail.png").Return(errmock)
			},
			input: &store.UpdateStoreInput{
				StoreID:      1,
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
			},
			expectErr: exception.ErrUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *storeService) {
			err := service.UpdateStore(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
