package service

import (
	"context"
	"testing"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/media"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/geolocation"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListExperiences(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	params := &database.ListExperiencesParams{
		Name:           "収穫",
		CoordinatorID:  "coordinator-id",
		ProducerID:     "producer-id",
		OnlyPublished:  true,
		ExcludeDeleted: true,
		EndAtGte:       now,
		Limit:          20,
		Offset:         0,
	}
	experiences := entity.Experiences{
		{
			ID:            "experience-id",
			CoordinatorID: "coordinator-id",
			ProducerID:    "producer-id",
			TypeID:        "experience-type-id",
			Title:         "じゃがいも収穫",
			Description:   "じゃがいもを収穫する体験です。",
			Public:        true,
			SoldOut:       false,
			Status:        entity.ExperienceStatusAccepting,
			ThumbnailURL:  "http://example.com/thumbnail.png",
			Media: []*entity.ExperienceMedia{
				{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
				{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
			},
			RecommendedPoints: []string{
				"じゃがいもを収穫する楽しさを体験できます。",
				"新鮮なじゃがいもを持ち帰ることができます。",
			},
			PromotionVideoURL:  "http://example.com/promotion.mp4",
			Duration:           60,
			Direction:          "彦根駅から徒歩10分",
			BusinessOpenTime:   "1000",
			BusinessCloseTime:  "1800",
			HostPostalCode:     "5220061",
			HostPrefecture:     "滋賀県",
			HostPrefectureCode: 25,
			HostCity:           "彦根市",
			HostAddressLine1:   "金亀町１−１",
			HostAddressLine2:   "",
			HostLongitude:      136.251739,
			HostLatitude:       35.276833,
			StartAt:            now.AddDate(0, 0, -1),
			EndAt:              now.AddDate(0, 0, 1),
			ExperienceRevision: entity.ExperienceRevision{
				ID:                    1,
				ExperienceID:          "experience-id",
				PriceAdult:            1000,
				PriceJuniorHighSchool: 800,
				PriceElementarySchool: 600,
				PricePreschool:        400,
				PriceSenior:           700,
				CreatedAt:             now,
				UpdatedAt:             now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListExperiencesInput
		expect      entity.Experiences
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().List(gomock.Any(), params).Return(experiences, nil)
				mocks.db.Experience.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListExperiencesInput{
				Name:            "収穫",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				OnlyPublished:   true,
				ExcludeFinished: true,
				ExcludeDeleted:  true,
				Limit:           20,
				Offset:          0,
				NoLimit:         false,
			},
			expect:      experiences,
			expectTotal: 1,
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
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().List(gomock.Any(), params).Return(nil, assert.AnError)
				mocks.db.Experience.EXPECT().Count(gomock.Any(), params).Return(int64(1), nil)
			},
			input: &store.ListExperiencesInput{
				Name:            "収穫",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				OnlyPublished:   true,
				ExcludeFinished: true,
				ExcludeDeleted:  true,
				Limit:           20,
				Offset:          0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
		{
			name: "failed to count",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().List(gomock.Any(), params).Return(experiences, nil)
				mocks.db.Experience.EXPECT().Count(gomock.Any(), params).Return(int64(0), assert.AnError)
			},
			input: &store.ListExperiencesInput{
				Name:            "収穫",
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				OnlyPublished:   true,
				ExcludeFinished: true,
				ExcludeDeleted:  true,
				Limit:           20,
				Offset:          0,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, total, err := service.ListExperiences(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
			assert.Equal(t, tt.expectTotal, total)
		}, withNow(now)))
	}
}

func TestListExperiencesByGeolocation(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	params := &database.ListExperiencesByGeolocationParams{
		CoordinatorID:  "coordinator-id",
		ProducerID:     "producer-id",
		Longitude:      136.251739,
		Latitude:       35.276833,
		Radius:         1000,
		OnlyPublished:  true,
		ExcludeDeleted: true,
		EndAtGte:       now,
	}
	experiences := entity.Experiences{
		{
			ID:            "experience-id",
			CoordinatorID: "coordinator-id",
			ProducerID:    "producer-id",
			TypeID:        "experience-type-id",
			Title:         "じゃがいも収穫",
			Description:   "じゃがいもを収穫する体験です。",
			Public:        true,
			SoldOut:       false,
			Status:        entity.ExperienceStatusAccepting,
			ThumbnailURL:  "http://example.com/thumbnail.png",
			Media: []*entity.ExperienceMedia{
				{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
				{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
			},
			RecommendedPoints: []string{
				"じゃがいもを収穫する楽しさを体験できます。",
				"新鮮なじゃがいもを持ち帰ることができます。",
			},
			PromotionVideoURL:  "http://example.com/promotion.mp4",
			Duration:           60,
			Direction:          "彦根駅から徒歩10分",
			BusinessOpenTime:   "1000",
			BusinessCloseTime:  "1800",
			HostPostalCode:     "5220061",
			HostPrefecture:     "滋賀県",
			HostPrefectureCode: 25,
			HostCity:           "彦根市",
			HostAddressLine1:   "金亀町１−１",
			HostAddressLine2:   "",
			HostLongitude:      136.251739,
			HostLatitude:       35.276833,
			StartAt:            now.AddDate(0, 0, -1),
			EndAt:              now.AddDate(0, 0, 1),
			ExperienceRevision: entity.ExperienceRevision{
				ID:                    1,
				ExperienceID:          "experience-id",
				PriceAdult:            1000,
				PriceJuniorHighSchool: 800,
				PriceElementarySchool: 600,
				PricePreschool:        400,
				PriceSenior:           700,
				CreatedAt:             now,
				UpdatedAt:             now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name        string
		setup       func(ctx context.Context, mocks *mocks)
		input       *store.ListExperiencesByGeolocationInput
		expect      entity.Experiences
		expectTotal int64
		expectErr   error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().ListByGeolocation(ctx, params).Return(experiences, nil)
			},
			input: &store.ListExperiencesByGeolocationInput{
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				Longitude:       136.251739,
				Latitude:        35.276833,
				Radius:          1000,
				OnlyPublished:   true,
				ExcludeFinished: true,
				ExcludeDeleted:  true,
			},
			expect:      experiences,
			expectTotal: 1,
			expectErr:   nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.ListExperiencesByGeolocationInput{
				Longitude: 139.81083,
				Latitude:  35.71014,
				Radius:    -1,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInvalidArgument,
		},
		{
			name: "failed to list",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().ListByGeolocation(ctx, params).Return(nil, assert.AnError)
			},
			input: &store.ListExperiencesByGeolocationInput{
				CoordinatorID:   "coordinator-id",
				ProducerID:      "producer-id",
				Longitude:       136.251739,
				Latitude:        35.276833,
				Radius:          1000,
				OnlyPublished:   true,
				ExcludeFinished: true,
				ExcludeDeleted:  true,
			},
			expect:      nil,
			expectTotal: 0,
			expectErr:   exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.ListExperiencesByGeolocation(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestMultiGetExperiences(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	experiences := entity.Experiences{
		{
			ID:            "experience-id",
			CoordinatorID: "coordinator-id",
			ProducerID:    "producer-id",
			TypeID:        "experience-type-id",
			Title:         "じゃがいも収穫",
			Description:   "じゃがいもを収穫する体験です。",
			Public:        true,
			SoldOut:       false,
			Status:        entity.ExperienceStatusAccepting,
			ThumbnailURL:  "http://example.com/thumbnail.png",
			Media: []*entity.ExperienceMedia{
				{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
				{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
			},
			RecommendedPoints: []string{
				"じゃがいもを収穫する楽しさを体験できます。",
				"新鮮なじゃがいもを持ち帰ることができます。",
			},
			PromotionVideoURL:  "http://example.com/promotion.mp4",
			Duration:           60,
			Direction:          "彦根駅から徒歩10分",
			BusinessOpenTime:   "1000",
			BusinessCloseTime:  "1800",
			HostPostalCode:     "5220061",
			HostPrefecture:     "滋賀県",
			HostPrefectureCode: 25,
			HostCity:           "彦根市",
			HostAddressLine1:   "金亀町１−１",
			HostAddressLine2:   "",
			HostLongitude:      136.251739,
			HostLatitude:       35.276833,
			StartAt:            now.AddDate(0, 0, -1),
			EndAt:              now.AddDate(0, 0, 1),
			ExperienceRevision: entity.ExperienceRevision{
				ID:                    1,
				ExperienceID:          "experience-id",
				PriceAdult:            1000,
				PriceJuniorHighSchool: 800,
				PriceElementarySchool: 600,
				PricePreschool:        400,
				PriceSenior:           700,
				CreatedAt:             now,
				UpdatedAt:             now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetExperiencesInput
		expect    entity.Experiences
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().MultiGet(ctx, []string{"experience-id"}).Return(experiences, nil)
			},
			input: &store.MultiGetExperiencesInput{
				ExperienceIDs: []string{"experience-id"},
			},
			expect:    experiences,
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
		{
			name: "failed to multi get",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().MultiGet(ctx, []string{"experience-id"}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetExperiencesInput{
				ExperienceIDs: []string{"experience-id"},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetExperiences(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestMultiGetExperiencesByRevision(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	experiences := entity.Experiences{
		{
			ID:            "experience-id",
			CoordinatorID: "coordinator-id",
			ProducerID:    "producer-id",
			TypeID:        "experience-type-id",
			Title:         "じゃがいも収穫",
			Description:   "じゃがいもを収穫する体験です。",
			Public:        true,
			SoldOut:       false,
			Status:        entity.ExperienceStatusAccepting,
			ThumbnailURL:  "http://example.com/thumbnail.png",
			Media: []*entity.ExperienceMedia{
				{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
				{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
			},
			RecommendedPoints: []string{
				"じゃがいもを収穫する楽しさを体験できます。",
				"新鮮なじゃがいもを持ち帰ることができます。",
			},
			PromotionVideoURL:  "http://example.com/promotion.mp4",
			Duration:           60,
			Direction:          "彦根駅から徒歩10分",
			BusinessOpenTime:   "1000",
			BusinessCloseTime:  "1800",
			HostPostalCode:     "5220061",
			HostPrefecture:     "滋賀県",
			HostPrefectureCode: 25,
			HostCity:           "彦根市",
			HostAddressLine1:   "金亀町１−１",
			HostAddressLine2:   "",
			HostLongitude:      136.251739,
			HostLatitude:       35.276833,
			StartAt:            now.AddDate(0, 0, -1),
			EndAt:              now.AddDate(0, 0, 1),
			ExperienceRevision: entity.ExperienceRevision{
				ID:                    1,
				ExperienceID:          "experience-id",
				PriceAdult:            1000,
				PriceJuniorHighSchool: 800,
				PriceElementarySchool: 600,
				PricePreschool:        400,
				PriceSenior:           700,
				CreatedAt:             now,
				UpdatedAt:             now,
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.MultiGetExperiencesByRevisionInput
		expect    entity.Experiences
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().MultiGetByRevision(ctx, []int64{1, 2}).Return(experiences, nil)
			},
			input: &store.MultiGetExperiencesByRevisionInput{
				ExperienceRevisionIDs: []int64{1, 2},
			},
			expect:    experiences,
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
		{
			name: "failed to multi get by revision",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().MultiGetByRevision(ctx, []int64{1, 2}).Return(nil, assert.AnError)
			},
			input: &store.MultiGetExperiencesByRevisionInput{
				ExperienceRevisionIDs: []int64{1, 2},
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.MultiGetExperiencesByRevision(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.ElementsMatch(t, tt.expect, actual)
		}))
	}
}

func TestGetExperience(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	experience := &entity.Experience{
		ID:            "experience-id",
		CoordinatorID: "coordinator-id",
		ProducerID:    "producer-id",
		TypeID:        "experience-type-id",
		Title:         "じゃがいも収穫",
		Description:   "じゃがいもを収穫する体験です。",
		Public:        true,
		SoldOut:       false,
		Status:        entity.ExperienceStatusAccepting,
		ThumbnailURL:  "http://example.com/thumbnail.png",
		Media: []*entity.ExperienceMedia{
			{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
			{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
		},
		RecommendedPoints: []string{
			"じゃがいもを収穫する楽しさを体験できます。",
			"新鮮なじゃがいもを持ち帰ることができます。",
		},
		PromotionVideoURL:  "http://example.com/promotion.mp4",
		Duration:           60,
		Direction:          "彦根駅から徒歩10分",
		BusinessOpenTime:   "1000",
		BusinessCloseTime:  "1800",
		HostPostalCode:     "5220061",
		HostPrefecture:     "滋賀県",
		HostPrefectureCode: 25,
		HostCity:           "彦根市",
		HostAddressLine1:   "金亀町１−１",
		HostAddressLine2:   "",
		HostLongitude:      136.251739,
		HostLatitude:       35.276833,
		StartAt:            now.AddDate(0, 0, -1),
		EndAt:              now.AddDate(0, 0, 1),
		ExperienceRevision: entity.ExperienceRevision{
			ID:                    1,
			ExperienceID:          "experience-id",
			PriceAdult:            1000,
			PriceJuniorHighSchool: 800,
			PriceElementarySchool: 600,
			PricePreschool:        400,
			PriceSenior:           700,
			CreatedAt:             now,
			UpdatedAt:             now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.GetExperienceInput
		expect    *entity.Experience
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().Get(ctx, "experience-id").Return(experience, nil)
			},
			input: &store.GetExperienceInput{
				ExperienceID: "experience-id",
			},
			expect:    experience,
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.GetExperienceInput{},
			expect:    nil,
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Experience.EXPECT().Get(ctx, "experience-id").Return(nil, assert.AnError)
			},
			input: &store.GetExperienceInput{
				ExperienceID: "experience-id",
			},
			expect:    nil,
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

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
	shop := &entity.Shop{
		ID:            "shop-id",
		Name:          "じゃがいも農園",
		CoordinatorID: "coordinator-id",
		ProducerIDs:   []string{"producer-id"},
	}
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		AdminID: "coordinator-id",
	}
	producerIn := &user.GetProducerInput{
		ProducerID: "producer-id",
	}
	producer := &uentity.Producer{
		AdminID: "producer-id",
	}
	locationIn := &geolocation.GetGeolocationInput{
		Address: &geolocation.Address{
			PostalCode:   "5220061",
			Prefecture:   "滋賀県",
			City:         "彦根市",
			AddressLine1: "金亀町１−１",
			AddressLine2: "",
		},
	}
	location := &geolocation.GetGeolocationOutput{
		Longitude: 136.251739,
		Latitude:  35.276833,
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CreateExperienceInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.geolocation.EXPECT().GetGeolocation(ctx, locationIn).Return(location, nil)
				mocks.db.Experience.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, experience *entity.Experience) error {
						expect := &entity.Experience{
							ID:            experience.ID, // ignore
							ShopID:        "shop-id",
							CoordinatorID: "coordinator-id",
							ProducerID:    "producer-id",
							TypeID:        "experience-type-id",
							Title:         "じゃがいも収穫体験",
							Description:   "じゃがいもを収穫する体験です。",
							Public:        true,
							SoldOut:       false,
							Status:        entity.ExperienceStatusUnknown,
							ThumbnailURL:  "",
							Media: entity.MultiExperienceMedia{
								{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
								{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
							},
							RecommendedPoints: []string{
								"じゃがいもを収穫する楽しさを体験できます。",
								"新鮮なじゃがいもを持ち帰ることができます。",
							},
							PromotionVideoURL:  "http://example.com/promotion.mp4",
							Duration:           60,
							Direction:          "彦根駅から徒歩10分",
							BusinessOpenTime:   "1000",
							BusinessCloseTime:  "1800",
							HostPostalCode:     "5220061",
							HostPrefecture:     "滋賀県",
							HostPrefectureCode: 25,
							HostCity:           "彦根市",
							HostAddressLine1:   "金亀町１−１",
							HostAddressLine2:   "",
							HostLongitude:      136.251739,
							HostLatitude:       35.276833,
							StartAt:            now.AddDate(0, -1, 0),
							EndAt:              now.AddDate(0, 1, 0),
							ExperienceRevision: entity.ExperienceRevision{
								ID:                    0,
								ExperienceID:          experience.ID, // ignore
								PriceAdult:            1000,
								PriceJuniorHighSchool: 800,
								PriceElementarySchool: 600,
								PricePreschool:        400,
								PriceSenior:           700,
							},
						}
						assert.Equal(t, expect, experience)
						return nil
					})
			},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CreateExperienceInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid experience media",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: true},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get shop",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(nil, database.ErrNotFound)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
			},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, exception.ErrNotFound)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
			},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get producer",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(nil, assert.AnError)
			},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "invalid prefecture validation",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
			},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: -1,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get geolocation",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.geolocation.EXPECT().GetGeolocation(ctx, locationIn).Return(nil, assert.AnError)
			},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "invalid experience validation",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				location := &geolocation.GetGeolocationOutput{
					Longitude: 270.0,
					Latitude:  90.0,
				}
				mocks.geolocation.EXPECT().GetGeolocation(ctx, locationIn).Return(location, nil)
			},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create experience",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Shop.EXPECT().Get(gomock.Any(), "shop-id").Return(shop, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().GetProducer(gomock.Any(), producerIn).Return(producer, nil)
				mocks.geolocation.EXPECT().GetGeolocation(ctx, locationIn).Return(location, nil)
				mocks.db.Experience.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &store.CreateExperienceInput{
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫体験",
				Description:   "じゃがいもを収穫する体験です。",
				Public:        true,
				SoldOut:       false,
				Media: []*store.CreateExperienceMedia{
					{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
					{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			_, err := service.CreateExperience(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestUpdateExperience(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 28, 18, 30, 0, 0)
	locationIn := &geolocation.GetGeolocationInput{
		Address: &geolocation.Address{
			PostalCode:   "5220061",
			Prefecture:   "滋賀県",
			City:         "彦根市",
			AddressLine1: "金亀町１−１",
			AddressLine2: "",
		},
	}
	location := &geolocation.GetGeolocationOutput{
		Longitude: 136.251739,
		Latitude:  35.276833,
	}
	params := &database.UpdateExperienceParams{
		TypeID:      "experience-type-id",
		Title:       "じゃがいも収穫体験",
		Description: "じゃがいもを収穫する体験です。",
		Public:      true,
		SoldOut:     false,
		Media: entity.MultiExperienceMedia{
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
		PromotionVideoURL:  "http://example.com/promotion.mp4",
		Duration:           60,
		Direction:          "彦根駅から徒歩10分",
		BusinessOpenTime:   "1000",
		BusinessCloseTime:  "1800",
		HostPostalCode:     "5220061",
		HostPrefectureCode: 25,
		HostCity:           "彦根市",
		HostAddressLine1:   "金亀町１−１",
		HostAddressLine2:   "",
		HostLongitude:      136.251739,
		HostLatitude:       35.276833,
		StartAt:            now.AddDate(0, -1, 0),
		EndAt:              now.AddDate(0, 1, 0),
	}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.UpdateExperienceInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.geolocation.EXPECT().GetGeolocation(ctx, locationIn).Return(location, nil)
				mocks.db.Experience.EXPECT().Update(ctx, "experience-id", params).Return(nil)
			},
			input: &store.UpdateExperienceInput{
				ExperienceID: "experience-id",
				TypeID:       "experience-type-id",
				Title:        "じゃがいも収穫体験",
				Description:  "じゃがいもを収穫する体験です。",
				Public:       true,
				SoldOut:      false,
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
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
		{
			name:  "invalid prefecture code",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.UpdateExperienceInput{
				ExperienceID: "experience-id",
				TypeID:       "experience-type-id",
				Title:        "じゃがいも収穫体験",
				Description:  "じゃがいもを収穫する体験です。",
				Public:       true,
				SoldOut:      false,
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				HostPrefectureCode: -1,
				HostCity:           "彦根市",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid experience media",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.UpdateExperienceInput{
				ExperienceID: "experience-id",
				TypeID:       "experience-type-id",
				Title:        "じゃがいも収穫体験",
				Description:  "じゃがいもを収穫する体験です。",
				Public:       true,
				SoldOut:      false,
				Media: []*store.UpdateExperienceMedia{
					{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
					{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: true},
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid prefecture validation",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.UpdateExperienceInput{
				ExperienceID: "experience-id",
				TypeID:       "experience-type-id",
				Title:        "じゃがいも収穫体験",
				Description:  "じゃがいもを収穫する体験です。",
				Public:       true,
				SoldOut:      false,
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: -1,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get geolocation",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.geolocation.EXPECT().GetGeolocation(ctx, locationIn).Return(nil, assert.AnError)
			},
			input: &store.UpdateExperienceInput{
				ExperienceID: "experience-id",
				TypeID:       "experience-type-id",
				Title:        "じゃがいも収穫体験",
				Description:  "じゃがいもを収穫する体験です。",
				Public:       true,
				SoldOut:      false,
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to update experience",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.geolocation.EXPECT().GetGeolocation(ctx, locationIn).Return(location, nil)
				mocks.db.Experience.EXPECT().Update(ctx, "experience-id", params).Return(assert.AnError)
			},
			input: &store.UpdateExperienceInput{
				ExperienceID: "experience-id",
				TypeID:       "experience-type-id",
				Title:        "じゃがいも収穫体験",
				Description:  "じゃがいもを収穫する体験です。",
				Public:       true,
				SoldOut:      false,
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
				PromotionVideoURL:  "http://example.com/promotion.mp4",
				Duration:           60,
				Direction:          "彦根駅から徒歩10分",
				BusinessOpenTime:   "1000",
				BusinessCloseTime:  "1800",
				HostPostalCode:     "5220061",
				HostPrefectureCode: 25,
				HostCity:           "彦根市",
				HostAddressLine1:   "金亀町１−１",
				HostAddressLine2:   "",
				StartAt:            now.AddDate(0, -1, 0),
				EndAt:              now.AddDate(0, 1, 0),
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.UpdateExperience(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestDeleteExperience(t *testing.T) {
	t.Parallel()

	videosIn := &media.ListExperienceVideosInput{
		ExperienceID: "experience-id",
	}
	videos := mentity.Videos{}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.DeleteExperienceInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.media.EXPECT().ListExperienceVideos(ctx, videosIn).Return(videos, nil)
				mocks.db.Experience.EXPECT().Delete(ctx, "experience-id").Return(nil)
			},
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
		{
			name: "failed to list experience videos",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.media.EXPECT().ListExperienceVideos(ctx, videosIn).Return(nil, assert.AnError)
			},
			input: &store.DeleteExperienceInput{
				ExperienceID: "experience-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "experience has videos",
			setup: func(ctx context.Context, mocks *mocks) {
				videos := mentity.Videos{{ID: "video-id"}}
				mocks.media.EXPECT().ListExperienceVideos(ctx, videosIn).Return(videos, nil)
			},
			input: &store.DeleteExperienceInput{
				ExperienceID: "experience-id",
			},
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to delete experience",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.media.EXPECT().ListExperienceVideos(ctx, videosIn).Return(videos, nil)
				mocks.db.Experience.EXPECT().Delete(ctx, "experience-id").Return(assert.AnError)
			},
			input: &store.DeleteExperienceInput{
				ExperienceID: "experience-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.DeleteExperience(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}
