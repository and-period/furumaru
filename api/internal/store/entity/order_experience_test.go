package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestOrderExperience(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		params *NewOrderExperienceParams
		expect *OrderExperience
		hasErr bool
	}{
		{
			name: "success",
			params: &NewOrderExperienceParams{
				OrderID: "order-id",
				Experience: &Experience{
					ID:            "experience-id",
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
					TypeID:        "experience-type-id",
					Title:         "じゃがいも収穫",
					Description:   "じゃがいもを収穫する体験",
					Public:        true,
					SoldOut:       false,
					Status:        ExperienceStatusAccepting,
					ThumbnailURL:  "http://example.com/thumbnail.png",
					Media: MultiExperienceMedia{
						{
							URL:         "http://example.com/thumbnail.png",
							IsThumbnail: true,
						},
					},
					MediaJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/thumbnail.png","isThumbnail":true}]`)),
					RecommendedPoints: []string{
						"ポイント1",
						"ポイント2",
					},
					RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
					PromotionVideoURL:     "http://example.com/promotion.mp4",
					Duration:              60,
					Direction:             "彦根駅から徒歩10分",
					BusinessOpenTime:      "1000",
					BusinessCloseTime:     "1800",
					HostPostalCode:        "5220061",
					HostPrefecture:        "滋賀県",
					HostPrefectureCode:    25,
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostLongitude:         136.251739,
					HostLatitude:          35.276833,
					HostGeolocation: mysql.Geometry{
						X: 136.251739,
						Y: 35.276833,
					},
					ExperienceRevision: ExperienceRevision{
						ID:                    1,
						ExperienceID:          "experience-id",
						PriceAdult:            1000,
						PriceJuniorHighSchool: 500,
						PriceElementarySchool: 300,
						PricePreschool:        0,
						PriceSenior:           200,
					},
					StartAt:   now.AddDate(0, 0, -1),
					EndAt:     now.AddDate(0, 0, 1),
					CreatedAt: now,
					UpdatedAt: now,
				},
				Transportation:        "徒歩",
				RequestedDate:         "20210101",
				RequestedTime:         "1000",
				AdultCount:            1,
				JuniorHighSchoolCount: 2,
				ElementarySchoolCount: 3,
				PreschoolCount:        4,
				SeniorCount:           5,
			},
			expect: &OrderExperience{
				OrderID:               "order-id",
				ExperienceRevisionID:  1,
				AdultCount:            1,
				JuniorHighSchoolCount: 2,
				ElementarySchoolCount: 3,
				PreschoolCount:        4,
				SeniorCount:           5,
				Remarks: OrderExperienceRemarks{
					Transportation: "徒歩",
					RequestedDate:  jst.Date(2021, 1, 1, 0, 0, 0, 0),
					RequestedTime:  jst.Date(0, 1, 1, 10, 0, 0, 0),
				},
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewOrderExperience(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestNewOrderExperienceRemarks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewOrderExperienceRemarksParams
		expect *OrderExperienceRemarks
		hasErr bool
	}{
		{
			name: "success",
			params: &NewOrderExperienceRemarksParams{
				Transportation: "徒歩",
				RequestedDate:  "20220101",
				RequestedTime:  "1000",
			},
			expect: &OrderExperienceRemarks{
				Transportation: "徒歩",
				RequestedDate:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				RequestedTime:  jst.Date(0, 1, 1, 10, 0, 0, 0),
			},
			hasErr: false,
		},
		{
			name: "empty params",
			params: &NewOrderExperienceRemarksParams{
				Transportation: "",
				RequestedDate:  "",
				RequestedTime:  "",
			},
			expect: &OrderExperienceRemarks{
				Transportation: "",
				RequestedDate:  time.Time{},
				RequestedTime:  time.Time{},
			},
			hasErr: false,
		},
		{
			name: "invalid requested date",
			params: &NewOrderExperienceRemarksParams{
				Transportation: "徒歩",
				RequestedDate:  "2022-01-01",
				RequestedTime:  "1000",
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "invalid requested time",
			params: &NewOrderExperienceRemarksParams{
				Transportation: "徒歩",
				RequestedDate:  "20220101",
				RequestedTime:  "10:00",
			},
			expect: nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewOrderExperienceRemarks(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
