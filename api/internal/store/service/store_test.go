package service

import (
	"context"
	"testing"

	"github.com/and-period/marche/api/internal/store/database"
	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/jst"
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
		input     *ListStoresInput
		expect    entity.Stores
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().List(ctx, params).Return(stores, nil)
			},
			input: &ListStoresInput{
				Limit:  20,
				Offset: 0,
			},
			expect:    stores,
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &ListStoresInput{},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get stores",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().List(ctx, params).Return(nil, errmock)
			},
			input: &ListStoresInput{
				Limit:  20,
				Offset: 0,
			},
			expect:    nil,
			expectErr: ErrInternal,
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
	store := &entity.Store{
		ID:           1,
		Name:         "&.農園",
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *GetStoreInput
		expect    *entity.Store
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Get(ctx, int64(1)).Return(store, nil)
			},
			input: &GetStoreInput{
				StoreID: 1,
			},
			expect:    store,
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &GetStoreInput{},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get store",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Get(ctx, int64(1)).Return(nil, errmock)
			},
			input: &GetStoreInput{
				StoreID: 1,
			},
			expect:    nil,
			expectErr: ErrInternal,
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

	store := &entity.Store{
		Name: "&.農園",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *CreateStoreInput
		expect    *entity.Store
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Create(ctx, store).Return(nil)
			},
			input: &CreateStoreInput{
				Name: "&.農園",
			},
			expect:    store,
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &CreateStoreInput{},
			expect:    nil,
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get store",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Create(ctx, store).Return(errmock)
			},
			input: &CreateStoreInput{
				Name: "&.農園",
			},
			expect:    nil,
			expectErr: ErrInternal,
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
		input     *UpdateStoreInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Update(ctx, int64(1), "&.農園", "https://and-period.jp/thumbnail.png").Return(nil)
			},
			input: &UpdateStoreInput{
				StoreID:      1,
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
			},
			expectErr: nil,
		},
		{
			name:      "invlid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &UpdateStoreInput{},
			expectErr: ErrInvalidArgument,
		},
		{
			name: "failed to get store",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Store.EXPECT().Update(ctx, int64(1), "&.農園", "https://and-period.jp/thumbnail.png").Return(errmock)
			},
			input: &UpdateStoreInput{
				StoreID:      1,
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
			},
			expectErr: ErrInternal,
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
