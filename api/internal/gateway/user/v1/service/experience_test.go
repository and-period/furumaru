package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
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

func TestExperience_Calc(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name           string
		experience     *Experience
		params         *CalcExperienceParams
		expectSubTotal int64
		expectDiscount int64
	}{
		{
			name: "success when part of amount",
			experience: &Experience{
				Experience: types.Experience{
					ID:               "experience-id",
					CoordinatorID:    "coordinator-id",
					ProducerID:       "producer-id",
					ExperienceTypeID: "experience-type-id",
					Title:            "じゃがいも収穫",
					Description:      "じゃがいもを収穫する体験です。",
					Status:           int32(ExperienceStatusAccepting),
					ThumbnailURL:     "http://example.com/thumbnail.png",
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
					HostPrefecture:        "滋賀県",
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
					Rate: &types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					StartAt: now.AddDate(0, 0, -1).Unix(),
					EndAt:   now.AddDate(0, 0, 1).Unix(),
				},
				revisionID: 1,
			},
			params: &CalcExperienceParams{
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 0,
				PreschoolCount:        0,
				SeniorCount:           1,
				Promotion: &Promotion{
					Promotion: types.Promotion{
						ID:           "promotion-id",
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Status:       int32(PromotionStatusEnabled),
						DiscountType: DiscountTypeAmount.Response(),
						DiscountRate: 200,
						Code:         "code0001",
						StartAt:      1640962800,
						EndAt:        1643641200,
					},
				},
			},
			expectSubTotal: 3500,
			expectDiscount: 200,
		},
		{
			name: "success when full amount",
			experience: &Experience{
				Experience: types.Experience{
					ID:               "experience-id",
					CoordinatorID:    "coordinator-id",
					ProducerID:       "producer-id",
					ExperienceTypeID: "experience-type-id",
					Title:            "じゃがいも収穫",
					Description:      "じゃがいもを収穫する体験です。",
					Status:           int32(ExperienceStatusAccepting),
					ThumbnailURL:     "http://example.com/thumbnail.png",
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
					HostPrefecture:        "滋賀県",
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
					Rate: &types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					StartAt: now.AddDate(0, 0, -1).Unix(),
					EndAt:   now.AddDate(0, 0, 1).Unix(),
				},
				revisionID: 1,
			},
			params: &CalcExperienceParams{
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 0,
				PreschoolCount:        0,
				SeniorCount:           1,
				Promotion: &Promotion{
					Promotion: types.Promotion{
						ID:           "promotion-id",
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Status:       int32(PromotionStatusEnabled),
						DiscountType: DiscountTypeAmount.Response(),
						DiscountRate: 4000,
						Code:         "code0001",
						StartAt:      1640962800,
						EndAt:        1643641200,
					},
				},
			},
			expectSubTotal: 3500,
			expectDiscount: 3500,
		},
		{
			name: "success when rate",
			experience: &Experience{
				Experience: types.Experience{
					ID:               "experience-id",
					CoordinatorID:    "coordinator-id",
					ProducerID:       "producer-id",
					ExperienceTypeID: "experience-type-id",
					Title:            "じゃがいも収穫",
					Description:      "じゃがいもを収穫する体験です。",
					Status:           int32(ExperienceStatusAccepting),
					ThumbnailURL:     "http://example.com/thumbnail.png",
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
					HostPrefecture:        "滋賀県",
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
					Rate: &types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					StartAt: now.AddDate(0, 0, -1).Unix(),
					EndAt:   now.AddDate(0, 0, 1).Unix(),
				},
				revisionID: 1,
			},
			params: &CalcExperienceParams{
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 0,
				PreschoolCount:        0,
				SeniorCount:           1,
				Promotion: &Promotion{
					Promotion: types.Promotion{
						ID:           "promotion-id",
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Status:       int32(PromotionStatusEnabled),
						DiscountType: DiscountTypeRate.Response(),
						DiscountRate: 50,
						Code:         "code0001",
						StartAt:      1640962800,
						EndAt:        1643641200,
					},
				},
			},
			expectSubTotal: 3500,
			expectDiscount: 1750,
		},
		{
			name:           "empty",
			experience:     nil,
			params:         nil,
			expectSubTotal: 0,
			expectDiscount: 0,
		},
		{
			name: "promotion none",
			experience: &Experience{
				Experience: types.Experience{
					ID:               "experience-id",
					CoordinatorID:    "coordinator-id",
					ProducerID:       "producer-id",
					ExperienceTypeID: "experience-type-id",
					Title:            "じゃがいも収穫",
					Description:      "じゃがいもを収穫する体験です。",
					Status:           int32(ExperienceStatusAccepting),
					ThumbnailURL:     "http://example.com/thumbnail.png",
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
					HostPrefecture:        "滋賀県",
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
					Rate: &types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					StartAt: now.AddDate(0, 0, -1).Unix(),
					EndAt:   now.AddDate(0, 0, 1).Unix(),
				},
				revisionID: 1,
			},
			params: &CalcExperienceParams{
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 0,
				PreschoolCount:        0,
				SeniorCount:           1,
				Promotion:             nil,
			},
			expectSubTotal: 3500,
			expectDiscount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			subtotal, discount := tt.experience.Calc(tt.params)
			assert.Equal(t, tt.expectSubTotal, subtotal)
			assert.Equal(t, tt.expectDiscount, discount)
		})
	}
}

func TestExperiences(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences entity.Experiences
		rates       map[string]*ExperienceRate
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
			rates: map[string]*ExperienceRate{
				"experience-id": {
					ExperienceRate: types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					experienceID: "experience-id",
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
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
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
						HostPrefecture:        "滋賀県",
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						HostLongitude:         136.251739,
						HostLatitude:          35.276833,
						Rate: &types.ExperienceRate{
							Count:   4,
							Average: 2.5,
							Detail: map[int64]int64{
								1: 2,
								2: 0,
								3: 1,
								4: 0,
								5: 1,
							},
						},
						StartAt: now.AddDate(0, 0, -1).Unix(),
						EndAt:   now.AddDate(0, 0, 1).Unix(),
					},
					revisionID: 1,
				},
			},
		},
		{
			name: "success withtout additional values",
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
			rates: map[string]*ExperienceRate{},
			expect: Experiences{
				{
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
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
						HostPrefecture:        "滋賀県",
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						HostLongitude:         136.251739,
						HostLatitude:          35.276833,
						Rate: &types.ExperienceRate{
							Count:   0,
							Average: 0.0,
							Detail: map[int64]int64{
								1: 0,
								2: 0,
								3: 0,
								4: 0,
								5: 0,
							},
						},
						StartAt: now.AddDate(0, 0, -1).Unix(),
						EndAt:   now.AddDate(0, 0, 1).Unix(),
					},
					revisionID: 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewExperiences(tt.experiences, tt.rates)
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
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
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
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
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
						Status:           int32(ExperienceStatusAccepting),
						ThumbnailURL:     "http://example.com/thumbnail.png",
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
			expect: []*types.Experience{
				{
					ID:               "experience-id",
					CoordinatorID:    "coordinator-id",
					ProducerID:       "producer-id",
					ExperienceTypeID: "experience-type-id",
					Title:            "じゃがいも収穫",
					Description:      "じゃがいもを収穫する体験です。",
					Status:           int32(ExperienceStatusAccepting),
					ThumbnailURL:     "http://example.com/thumbnail.png",
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

func TestExperienceRates(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		reviews entity.AggregatedExperienceReviews
		expect  ExperienceRates
	}{
		{
			name: "success",
			reviews: entity.AggregatedExperienceReviews{
				{
					ExperienceID: "experience-id",
					Count:        4,
					Average:      2.5,
					Rate1:        2,
					Rate2:        0,
					Rate3:        1,
					Rate4:        0,
					Rate5:        1,
				},
			},
			expect: ExperienceRates{
				{
					ExperienceRate: types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					experienceID: "experience-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewExperienceRates(tt.reviews)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceRates_MapByExperienceID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rates  ExperienceRates
		expect map[string]*ExperienceRate
	}{
		{
			name: "success",
			rates: ExperienceRates{
				{
					ExperienceRate: types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					experienceID: "experience-id",
				},
			},
			expect: map[string]*ExperienceRate{
				"experience-id": {
					ExperienceRate: types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					experienceID: "experience-id",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.rates.MapByExperienceID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceRates_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		rates  ExperienceRates
		expect []*types.ExperienceRate
	}{
		{
			name: "success",
			rates: ExperienceRates{
				{
					ExperienceRate: types.ExperienceRate{
						Count:   4,
						Average: 2.5,
						Detail: map[int64]int64{
							1: 2,
							2: 0,
							3: 1,
							4: 0,
							5: 1,
						},
					},
					experienceID: "experience-id",
				},
			},
			expect: []*types.ExperienceRate{
				{
					Count:   4,
					Average: 2.5,
					Detail: map[int64]int64{
						1: 2,
						2: 0,
						3: 1,
						4: 0,
						5: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.rates.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
