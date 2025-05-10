package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListSpotTypes(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)
	params := &database.ListSpotTypesParams{
		Name:   "収穫",
		Limit:  20,
		Offset: 0,
	}
	types := entity.SpotTypes{
		{
			ID:        "spot-type-id",
			Name:      "じゃがいも収穫",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListSpotTypesInput
		expect      entity.SpotTypes
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().List(gomock.Any(), params).Return(types, nil)
				mocks.db.SpotType.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListSpotTypesInput{
				Name:   "収穫",
				Limit:  20,
				Offset: 0,
			},
			expect: []*entity.SpotType{
				{
					ID:        "spot-type-id",
					Name:      "じゃがいも収穫",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListSpotTypesInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list spot types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.SpotType.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListSpotTypesInput{
				Name:   "収穫",
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count spot types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().List(gomock.Any(), params).Return(types, nil)
				mocks.db.SpotType.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListSpotTypesInput{
				Name:   "収穫",
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListSpotTypes(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetSpotTypes(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)
	types := entity.SpotTypes{
		{
			ID:        "spot-type-id",
			Name:      "じゃがいも収穫",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetSpotTypesInput
		expect    entity.SpotTypes
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().MultiGet(gomock.Any(), []string{"spot-type-id"}).Return(types, nil)
			},
			input: &store.MultiGetSpotTypesInput{
				SpotTypeIDs: []string{"spot-type-id"},
			},
			expect:    types,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetSpotTypesInput{
				SpotTypeIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get spot types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().MultiGet(gomock.Any(), []string{"spot-type-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetSpotTypesInput{
				SpotTypeIDs: []string{"spot-type-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetSpotTypes(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetSpotType(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)
	typ := &entity.SpotType{
		ID:        "spot-type-id",
		Name:      "じゃがいも収穫",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetSpotTypeInput
		expect    *entity.SpotType
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().Get(gomock.Any(), "spot-type-id").Return(typ, nil)
			},
			input: &store.GetSpotTypeInput{
				SpotTypeID: "spot-type-id",
			},
			expect:    typ,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetSpotTypeInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get spot type",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().Get(gomock.Any(), "spot-type-id").Return(nil, assert.AnError)
			},
			input: &store.GetSpotTypeInput{
				SpotTypeID: "spot-type-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetSpotType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateSpotType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateSpotTypeInput
		expect    *entity.SpotType
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, spotType *entity.SpotType) error {
						expect := &entity.SpotType{
							ID:   spotType.ID, // ignore
							Name: "じゃがいも収穫",
						}
						assert.Equal(t, expect, spotType)
						return nil
					})
			},
			input: &store.CreateSpotTypeInput{
				Name: "じゃがいも収穫",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateSpotTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create spot type",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateSpotTypeInput{
				Name: "じゃがいも収穫",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateSpotType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateSpotType(t *testing.T) {
	t.Parallel()

	params := &database.UpdateSpotTypeParams{
		Name: "じゃがいも収穫",
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateSpotTypeInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().Update(ctx, "spot-type-id", params).Return(nil)
			},
			input: &store.UpdateSpotTypeInput{
				SpotTypeID: "spot-type-id",
				Name:       "じゃがいも収穫",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateSpotTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to update spot type",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().Update(ctx, "spot-type-id", params).Return(assert.AnError)
			},
			input: &store.UpdateSpotTypeInput{
				SpotTypeID: "spot-type-id",
				Name:       "じゃがいも収穫",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateSpotType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteSpotType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteSpotTypeInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().Delete(ctx, "spot-type-id").Return(nil)
			},
			input: &store.DeleteSpotTypeInput{
				SpotTypeID: "spot-type-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteSpotTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to delete spot type",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.SpotType.EXPECT().Delete(ctx, "spot-type-id").Return(assert.AnError)
			},
			input: &store.DeleteSpotTypeInput{
				SpotTypeID: "spot-type-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteSpotType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
