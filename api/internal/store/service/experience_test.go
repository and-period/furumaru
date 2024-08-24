package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestListExperiences(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListExperiencesInput
		expect      entity.Experiences
		expectTotal int64
		expectErr   error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListExperiencesInput{
				Name:   "",
				Limit:  20,
				Offset: 0,
			},
			expect:      []*entity.Experience{},
			expectTotal: 0,
			expectErr:   nil,
		},
		{
			name:        "invalid argument",
			setup:       func(ctx context.Context, mocks *mocks) {},
			input:       &store.ListExperiencesInput{},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListExperiences(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}))
	}
}

func TestMultiGetExperiences(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetExperiencesInput
		expect    entity.Experiences
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetExperiencesInput{
				ExperienceIDs: []string{"experience-id"},
			},
			expect:    []*entity.Experience{},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetExperiencesInput{
				ExperienceIDs: []string{""},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetExperiences(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetExperiencesByRevision(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetExperiencesByRevisionInput
		expect    entity.Experiences
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetExperiencesByRevisionInput{
				ExperienceRevisionIDs: []int64{1, 2},
			},
			expect:    []*entity.Experience{},
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.MultiGetExperiencesByRevisionInput{
				ExperienceRevisionIDs: []int64{0},
			},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetExperiencesByRevision(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetExperienceInput
		expect    *entity.Experience
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.GetExperienceInput{
				ExperienceID: "experience-id",
			},
			expect:    &entity.Experience{},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetExperienceInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.GetExperience(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestCreateExperience(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateExperienceInput
		expect    *entity.Experience
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.CreateExperienceInput{
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				PriceAdult:            1000,
				PriceJuniorHighSchool: 800,
				PriceElementarySchool: 600,
				PricePreschool:        400,
				PriceSenior:           700,
				RecommendedPoints: []string{
					"じゃがいもを収穫する楽しさを体験できます。",
					"新鮮なじゃがいもを持ち帰ることができます。",
				},
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expect:    &entity.Experience{},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateExperienceInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CreateExperience(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}))
	}
}

func TestUpdateExperience(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateExperienceInput
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.UpdateExperienceInput{
				ExperienceID:  "experience-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.UpdateExperienceMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
				},
				PriceAdult:            1000,
				PriceJuniorHighSchool: 800,
				PriceElementarySchool: 600,
				PricePreschool:        400,
				PriceSenior:           700,
				RecommendedPoints: []string{
					"じゃがいもを収穫する楽しさを体験できます。",
					"新鮮なじゃがいもを持ち帰ることができます。",
				},
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.UpdateExperienceInput{},
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateExperience(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteExperience(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteExperienceInput
		expectErr error
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.DeleteExperienceInput{
				ExperienceID: "experience-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.DeleteExperienceInput{},
			expectErr: exception.ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteExperience(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
