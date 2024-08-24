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

func TestListExperienceTypes(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)
	params := &database.ListExperienceTypesParams{
		Name:   "収穫",
		Limit:  20,
		Offset: 0,
	}
	types := entity.ExperienceTypes{
		{
			ID:        "experience-type-id",
			Name:      "じゃがいも収穫",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListExperienceTypesInput
		expect      entity.ExperienceTypes
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceType.EXPECT().List(gomock.Any(), params).Return(types, nil)
				mocks.db.ExperienceType.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListExperienceTypesInput{
				Name:   "収穫",
				Limit:  20,
				Offset: 0,
			},
			expect: []*entity.ExperienceType{
				{
					ID:        "experience-type-id",
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
			input:       &store.ListExperienceTypesInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list experience types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceType.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.ExperienceType.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListExperienceTypesInput{
				Name:   "収穫",
				Limit:  20,
				Offset: 0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count experience types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceType.EXPECT().List(gomock.Any(), params).Return(types, nil)
				mocks.db.ExperienceType.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListExperienceTypesInput{
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
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListExperienceTypes(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetExperienceTypes(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)
	types := entity.ExperienceTypes{
		{
			ID:        "experience-type-id",
			Name:      "じゃがいも収穫",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetExperienceTypesInput
		expect    entity.ExperienceTypes
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceType.EXPECT().MultiGet(gomock.Any(), []string{"experience-type-id"}).Return(types, nil)
			},
			input: &store.MultiGetExperienceTypesInput{
				ExperienceTypeIDs: []string{"experience-type-id"},
			},
			expect:    types,
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetExperienceTypesInput{
				ExperienceTypeIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to multi get experience types",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceType.EXPECT().MultiGet(gomock.Any(), []string{"experience-type-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetExperienceTypesInput{
				ExperienceTypeIDs: []string{"experience-type-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetExperienceTypes(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetExperienceType(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)
	typ := &entity.ExperienceType{
		ID:        "experience-type-id",
		Name:      "じゃがいも収穫",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetExperienceTypeInput
		expect    *entity.ExperienceType
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceType.EXPECT().Get(gomock.Any(), "experience-type-id").Return(typ, nil)
			},
			input: &store.GetExperienceTypeInput{
				ExperienceTypeID: "experience-type-id",
			},
			expect:    typ,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetExperienceTypeInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get experience type",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ExperienceType.EXPECT().Get(gomock.Any(), "experience-type-id").Return(nil, assert.AnError)
			},
			input: &store.GetExperienceTypeInput{
				ExperienceTypeID: "experience-type-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetExperienceType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateExperienceType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateExperienceTypeInput
		expect    *entity.ExperienceType
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.CreateExperienceTypeInput{
				Name: "じゃがいも収穫",
			},
			expect:    &entity.ExperienceType{},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateExperienceTypeInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CreateExperienceType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUpdateExperienceType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateExperienceTypeInput
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.UpdateExperienceTypeInput{
				ExperienceTypeID: "experience-type-id",
				Name:             "じゃがいも収穫",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateExperienceTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateExperienceType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteExperienceType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteExperienceTypeInput
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.DeleteExperienceTypeInput{
				ExperienceTypeID: "experience-type-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteExperienceTypeInput{},
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteExperienceType(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
