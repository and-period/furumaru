package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestOrderExperiences(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		items       entity.OrderExperiences
		experiences map[int64]*Experience
		expect      OrderExperiences
	}{
		{
			name: "success",
			items: entity.OrderExperiences{
				{
					OrderID:               "order-id",
					ExperienceRevisionID:  1,
					AdultCount:            2,
					JuniorHighSchoolCount: 1,
					ElementarySchoolCount: 0,
					PreschoolCount:        0,
					SeniorCount:           0,
					Remarks: entity.OrderExperienceRemarks{
						Transportation: "電車",
						RequestedDate:  jst.Date(2024, 1, 2, 0, 0, 0, 0),
						RequestedTime:  jst.Date(0, 1, 1, 18, 30, 0, 0),
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			experiences: map[int64]*Experience{
				1: {
					Experience: types.Experience{
						ID:               "experience-id",
						CoordinatorID:    "coordinator-id",
						ProducerID:       "producer-id",
						ExperienceTypeID: "experience-type-id",
						Title:            "じゃがいも収穫",
						Description:      "じゃがいもを収穫する体験です。",
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
						HostCity:              "彦根市",
						HostAddressLine1:      "金亀町１−１",
						HostAddressLine2:      "",
						StartAt:               now.AddDate(0, 0, -1).Unix(),
						EndAt:                 now.AddDate(0, 0, 1).Unix(),
					},
					revisionID: 1,
				},
			},
			expect: OrderExperiences{
				{
					OrderExperience: types.OrderExperience{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						AdultPrice:            1000,
						JuniorHighSchoolCount: 1,
						JuniorHighSchoolPrice: 800,
						ElementarySchoolCount: 0,
						ElementarySchoolPrice: 600,
						PreschoolCount:        0,
						PreschoolPrice:        400,
						SeniorCount:           0,
						SeniorPrice:           700,
						Transportation:        "電車",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					orderID: "order-id",
				},
			},
		},
		{
			name: "success with nil",
			items: entity.OrderExperiences{
				{
					OrderID:               "order-id",
					ExperienceRevisionID:  1,
					AdultCount:            2,
					JuniorHighSchoolCount: 1,
					ElementarySchoolCount: 0,
					PreschoolCount:        0,
					SeniorCount:           0,
					Remarks: entity.OrderExperienceRemarks{
						Transportation: "電車",
						RequestedDate:  jst.Date(2024, 1, 2, 0, 0, 0, 0),
						RequestedTime:  jst.Date(0, 1, 1, 18, 30, 0, 0),
					},
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			experiences: map[int64]*Experience{},
			expect: OrderExperiences{
				{
					OrderExperience: types.OrderExperience{
						ExperienceID:          "",
						AdultCount:            2,
						AdultPrice:            0,
						JuniorHighSchoolCount: 1,
						JuniorHighSchoolPrice: 0,
						ElementarySchoolCount: 0,
						ElementarySchoolPrice: 0,
						PreschoolCount:        0,
						PreschoolPrice:        0,
						SeniorCount:           0,
						SeniorPrice:           0,
						Transportation:        "電車",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					orderID: "order-id",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewOrderExperiences(tt.items, tt.experiences)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestOrderExperiences_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		experiences OrderExperiences
		expect      []*types.OrderExperience
	}{
		{
			name: "success",
			experiences: OrderExperiences{
				{
					OrderExperience: types.OrderExperience{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						AdultPrice:            1000,
						JuniorHighSchoolCount: 1,
						JuniorHighSchoolPrice: 800,
						ElementarySchoolCount: 0,
						ElementarySchoolPrice: 600,
						PreschoolCount:        0,
						PreschoolPrice:        400,
						SeniorCount:           0,
						SeniorPrice:           700,
						Transportation:        "電車",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					orderID: "order-id",
				},
			},
			expect: []*types.OrderExperience{
				{
					ExperienceID:          "experience-id",
					AdultCount:            2,
					AdultPrice:            1000,
					JuniorHighSchoolCount: 1,
					JuniorHighSchoolPrice: 800,
					ElementarySchoolCount: 0,
					ElementarySchoolPrice: 600,
					PreschoolCount:        0,
					PreschoolPrice:        400,
					SeniorCount:           0,
					SeniorPrice:           700,
					Transportation:        "電車",
					RequestedDate:         "20240102",
					RequestedTime:         "1830",
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
