package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
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
			expect: ExperienceStatusUnknown,
		},
		{
			name:   "waiting",
			status: entity.ExperienceStatusWaiting,
			expect: ExperienceStatusWaiting,
		},
		{
			name:   "accepting",
			status: entity.ExperienceStatusAccepting,
			expect: ExperienceStatusAccepting,
		},
		{
			name:   "sold out",
			status: entity.ExperienceStatusSoldOut,
			expect: ExperienceStatusSoldOut,
		},
		{
			name:   "finished",
			status: entity.ExperienceStatusFinished,
			expect: ExperienceStatusFinished,
		},
		{
			name:   "unknown",
			status: entity.ExperienceStatus(-1),
			expect: ExperienceStatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
		expect int32
	}{
		{
			name:   "unknown",
			status: ExperienceStatusUnknown,
			expect: 0,
		},
		{
			name:   "waiting",
			status: ExperienceStatusWaiting,
			expect: 1,
		},
		{
			name:   "accepting",
			status: ExperienceStatusAccepting,
			expect: 2,
		},
		{
			name:   "sold out",
			status: ExperienceStatusSoldOut,
			expect: 3,
		},
		{
			name:   "finished",
			status: ExperienceStatusFinished,
			expect: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.status.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperiences(t *testing.T) {
	t.Parallel()

	now := time.Now()

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
			},
			expect: Experiences{
				{
					Experience: response.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
						Media: []*response.ExperienceMedia{
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
						HostPrefecture:        "滋賀県",
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						HostLongitude:         136.251739,
						HostLatitude:          35.276833,
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
					},
					revisionID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewExperiences(tt.experiences)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperiences_MapByRevision(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		expect      map[int64]*Experience
	}{
		{
			name: "success",
			experiences: Experiences{
				{
					Experience: response.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
						Media: []*response.ExperienceMedia{
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
						HostPrefecture:        "滋賀県",
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						HostLongitude:         136.251739,
						HostLatitude:          35.276833,
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
					},
					revisionID: 1,
				},
			},
			expect: map[int64]*Experience{
				1: {
					Experience: response.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
						Media: []*response.ExperienceMedia{
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
						HostPrefecture:        "滋賀県",
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						HostLongitude:         136.251739,
						HostLatitude:          35.276833,
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
					},
					revisionID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.experiences.MapByRevision())
		})
	}
}

func TestExperiences_Response(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		expect      []*response.Experience
	}{
		{
			name: "success",
			experiences: Experiences{
				{
					Experience: response.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
						Media: []*response.ExperienceMedia{
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
						HostPrefecture:        "滋賀県",
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						HostLongitude:         136.251739,
						HostLatitude:          35.276833,
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
					},
				},
			},
			expect: []*response.Experience{
				{
					ID:               "experience-id",
					CoordinatorID:    "coordinator-id",
					ProducerID:       "producer-id",
					ExperienceTypeID: "experience-type-id",
					Title:            "じゃがいも収穫",
					Description:      "じゃがいもを収穫する体験です。",
					Status:           int32(ExperienceStatusAccepting),
					ThumbnailURL:     "http://example.com/thumbnail.png",
					Media: []*response.ExperienceMedia{
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
					HostPrefecture:        "滋賀県",
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
					StartAt:               now.AddDate(0, 0, -1).Unix(),
					EndAt:                 now.AddDate(0, 0, 1).Unix(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
