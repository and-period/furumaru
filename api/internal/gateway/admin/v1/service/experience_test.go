package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestExperienceStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status entity.ExperienceStatus
		expect ExperienceStatus
	}{
		{
			name:   "private",
			status: entity.ExperienceStatusPrivate,
			expect: ExperienceStatus(types.ExperienceStatusPrivate),
		},
		{
			name:   "waiting",
			status: entity.ExperienceStatusWaiting,
			expect: ExperienceStatus(types.ExperienceStatusWaiting),
		},
		{
			name:   "accepting",
			status: entity.ExperienceStatusAccepting,
			expect: ExperienceStatus(types.ExperienceStatusAccepting),
		},
		{
			name:   "sold out",
			status: entity.ExperienceStatusSoldOut,
			expect: ExperienceStatus(types.ExperienceStatusSoldOut),
		},
		{
			name:   "finished",
			status: entity.ExperienceStatusFinished,
			expect: ExperienceStatus(types.ExperienceStatusFinished),
		},
		{
			name:   "archived",
			status: entity.ExperienceStatusArchived,
			expect: ExperienceStatus(types.ExperienceStatusArchived),
		},
		{
			name:   "unknown",
			status: entity.ExperienceStatusUnknown,
			expect: ExperienceStatus(types.ExperienceStatusUnknown),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewExperienceStatus(tt.status)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceStatus_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status ExperienceStatus
		expect types.ExperienceStatus
	}{
		{
			name:   "private",
			status: ExperienceStatus(types.ExperienceStatusPrivate),
			expect: 1,
		},
		{
			name:   "waiting",
			status: ExperienceStatus(types.ExperienceStatusWaiting),
			expect: 2,
		},
		{
			name:   "accepting",
			status: ExperienceStatus(types.ExperienceStatusAccepting),
			expect: 3,
		},
		{
			name:   "sold out",
			status: ExperienceStatus(types.ExperienceStatusSoldOut),
			expect: 4,
		},
		{
			name:   "finished",
			status: ExperienceStatus(types.ExperienceStatusFinished),
			expect: 5,
		},
		{
			name:   "archived",
			status: ExperienceStatus(types.ExperienceStatusArchived),
			expect: 6,
		},
		{
			name:   "unknown",
			status: ExperienceStatus(types.ExperienceStatusUnknown),
			expect: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.status.Response())
		})
	}
}

func TestExperiences(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)

	tests := []struct {
		name        string
		experiences entity.Experiences
		expect      Experiences
	}{
		{
			name: "success",
			experiences: entity.Experiences{
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
						"じゃがいもの美味しさを再認識できます。",
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
			},
			expect: Experiences{
				{
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Public:           true,
						SoldOut:          false,
						Status:           types.ExperienceStatusAccepting,
						Media: []*types.ExperienceMedia{
							{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
						},
						PriceAdult:            1000,
						PriceJuniorHighSchool: 800,
						PriceElementarySchool: 600,
						PricePreschool:        400,
						PriceSenior:           700,
						RecommendedPoint1:     "じゃがいもを収穫する楽しさを体験できます。",
						RecommendedPoint2:     "新鮮なじゃがいもを持ち帰ることができます。",
						RecommendedPoint3:     "じゃがいもの美味しさを再認識できます。",
						PromotionVideoURL:     "http://example.com/promotion.mp4",
						Duration:              60,
						Direction:             "彦根駅から徒歩10分",
						BusinessOpenTime:      "1000",
						BusinessCloseTime:     "1800",
						HostPostalCode:        "5220061",
						HostPrefectureCode:    25,
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
						CreatedAt:             now.Unix(),
						UpdatedAt:             now.Unix(),
					},
					revisionID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewExperiences(tt.experiences)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperiences_MapByRevisionID(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)

	tests := []struct {
		name        string
		experiences Experiences
		expect      map[int64]*Experience
	}{
		{
			name: "success",
			experiences: Experiences{
				{
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Public:           true,
						SoldOut:          false,
						Status:           types.ExperienceStatusAccepting,
						Media: []*types.ExperienceMedia{
							{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
						},
						PriceAdult:            1000,
						PriceJuniorHighSchool: 800,
						PriceElementarySchool: 600,
						PricePreschool:        400,
						PriceSenior:           700,
						RecommendedPoint1:     "じゃがいもを収穫する楽しさを体験できます。",
						RecommendedPoint2:     "新鮮なじゃがいもを持ち帰ることができます。",
						RecommendedPoint3:     "じゃがいもの美味しさを再認識できます。",
						PromotionVideoURL:     "http://example.com/promotion.mp4",
						Duration:              60,
						Direction:             "彦根駅から徒歩10分",
						BusinessOpenTime:      "1000",
						BusinessCloseTime:     "1800",
						HostPostalCode:        "5220061",
						HostPrefectureCode:    25,
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
						CreatedAt:             now.Unix(),
						UpdatedAt:             now.Unix(),
					},
					revisionID: 1,
				},
			},
			expect: map[int64]*Experience{
				1: {
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Public:           true,
						SoldOut:          false,
						Status:           types.ExperienceStatusAccepting,
						Media: []*types.ExperienceMedia{
							{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
						},
						PriceAdult:            1000,
						PriceJuniorHighSchool: 800,
						PriceElementarySchool: 600,
						PricePreschool:        400,
						PriceSenior:           700,
						RecommendedPoint1:     "じゃがいもを収穫する楽しさを体験できます。",
						RecommendedPoint2:     "新鮮なじゃがいもを持ち帰ることができます。",
						RecommendedPoint3:     "じゃがいもの美味しさを再認識できます。",
						PromotionVideoURL:     "http://example.com/promotion.mp4",
						Duration:              60,
						Direction:             "彦根駅から徒歩10分",
						BusinessOpenTime:      "1000",
						BusinessCloseTime:     "1800",
						HostPostalCode:        "5220061",
						HostPrefectureCode:    25,
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
						CreatedAt:             now.Unix(),
						UpdatedAt:             now.Unix(),
					},
					revisionID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.experiences.MapByRevision())
		})
	}
}

func TestExperiences_Response(t *testing.T) {
	t.Parallel()

	now := jst.Date(2024, 8, 24, 18, 30, 0, 0)

	tests := []struct {
		name        string
		experiences Experiences
		expect      []*types.Experience
	}{
		{
			name: "success",
			experiences: Experiences{
				{
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Public:           true,
						SoldOut:          false,
						Status:           types.ExperienceStatusAccepting,
						Media: []*types.ExperienceMedia{
							{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
							{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
						},
						PriceAdult:            1000,
						PriceJuniorHighSchool: 800,
						PriceElementarySchool: 600,
						PricePreschool:        400,
						PriceSenior:           700,
						RecommendedPoint1:     "じゃがいもを収穫する楽しさを体験できます。",
						RecommendedPoint2:     "新鮮なじゃがいもを持ち帰ることができます。",
						RecommendedPoint3:     "じゃがいもの美味しさを再認識できます。",
						PromotionVideoURL:     "http://example.com/promotion.mp4",
						Duration:              60,
						Direction:             "彦根駅から徒歩10分",
						BusinessOpenTime:      "1000",
						BusinessCloseTime:     "1800",
						HostPostalCode:        "5220061",
						HostPrefectureCode:    25,
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
						CreatedAt:             now.Unix(),
						UpdatedAt:             now.Unix(),
					},
					revisionID: 1,
				},
			},
			expect: []*types.Experience{
				{
					ID:               "experience-id",
					CoordinatorID:    "coordinator-id",
					ProducerID:       "producer-id",
					ExperienceTypeID: "experience-type-id",
					Title:            "じゃがいも収穫",
					Description:      "じゃがいもを収穫する体験です。",
					Public:           true,
					SoldOut:          false,
					Status:           types.ExperienceStatusAccepting,
					Media: []*types.ExperienceMedia{
						{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
						{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
					},
					PriceAdult:            1000,
					PriceJuniorHighSchool: 800,
					PriceElementarySchool: 600,
					PricePreschool:        400,
					PriceSenior:           700,
					RecommendedPoint1:     "じゃがいもを収穫する楽しさを体験できます。",
					RecommendedPoint2:     "新鮮なじゃがいもを持ち帰ることができます。",
					RecommendedPoint3:     "じゃがいもの美味しさを再認識できます。",
					PromotionVideoURL:     "http://example.com/promotion.mp4",
					Duration:              60,
					Direction:             "彦根駅から徒歩10分",
					BusinessOpenTime:      "1000",
					BusinessCloseTime:     "1800",
					HostPostalCode:        "5220061",
					HostPrefectureCode:    25,
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					StartAt:               now.AddDate(0, 0, -1).Unix(),
					EndAt:                 now.AddDate(0, 0, 1).Unix(),
					CreatedAt:             now.Unix(),
					UpdatedAt:             now.Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.experiences.Response())
		})
	}
}

func TestMultiExperienceMedia(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		media  entity.MultiExperienceMedia
		expect MultiExperienceMedia
	}{
		{
			name: "success",
			media: entity.MultiExperienceMedia{
				{
					URL:         "http://example.com/thumbnail01.png",
					IsThumbnail: true,
				},
			},
			expect: MultiExperienceMedia{
				{
					ExperienceMedia: types.ExperienceMedia{
						URL:         "http://example.com/thumbnail01.png",
						IsThumbnail: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewMultiExperienceMedia(tt.media)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestMultiExperienceMedia_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		media  MultiExperienceMedia
		expect []*types.ExperienceMedia
	}{
		{
			name: "success",
			media: MultiExperienceMedia{
				{
					ExperienceMedia: types.ExperienceMedia{
						URL:         "http://example.com/thumbnail01.png",
						IsThumbnail: true,
					},
				},
			},
			expect: []*types.ExperienceMedia{
				{
					URL:         "http://example.com/thumbnail01.png",
					IsThumbnail: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.media.Response())
		})
	}
}
