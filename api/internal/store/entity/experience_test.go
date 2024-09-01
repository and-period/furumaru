package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func TestExperience(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		params *NewExperienceParams
		expect *Experience
		hasErr bool
	}{
		{
			name: "success",
			params: &NewExperienceParams{
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫",
				Description:   "じゃがいもを収穫する体験",
				Public:        true,
				SoldOut:       false,
				Media: []*ExperienceMedia{{
					URL:         "http://example.com/thumbnail.png",
					IsThumbnail: true,
				}},
				RecommendedPoints: []string{
					"じゃがいもを収穫する",
					"じゃがいもを食べる",
					"じゃがいもを持ち帰る",
				},
				PromotionVideoURL:     "http://example.com/promotion.mp4",
				HostPostalCode:        "5220061",
				HostPrefectureCode:    25,
				HostCity:              "彦根市",
				HostAddressLine1:      "金亀町１−１",
				HostAddressLine2:      "",
				HostLongitude:         136.251739,
				HostLatitude:          35.276833,
				StartAt:               now.AddDate(0, -1, 0),
				EndAt:                 now.AddDate(0, 1, 0),
				PriceAdult:            1000,
				PriceJuniorHighSchool: 800,
				PriceElementarySchool: 600,
				PricePreschool:        400,
				PriceSenior:           200,
			},
			expect: &Experience{
				CoordinatorID: "coordinator-id",
				ProducerID:    "producer-id",
				TypeID:        "experience-type-id",
				Title:         "じゃがいも収穫",
				Description:   "じゃがいもを収穫する体験",
				Public:        true,
				SoldOut:       false,
				Status:        ExperienceStatusUnknown,
				Media: MultiExperienceMedia{{
					URL:         "http://example.com/thumbnail.png",
					IsThumbnail: true,
				}},
				RecommendedPoints: []string{
					"じゃがいもを収穫する",
					"じゃがいもを食べる",
					"じゃがいもを持ち帰る",
				},
				PromotionVideoURL:  "http://example.com/promotion.mp4",
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
				ExperienceRevision: ExperienceRevision{
					ID:                    0,
					ExperienceID:          "", // ignore
					PriceAdult:            1000,
					PriceJuniorHighSchool: 800,
					PriceElementarySchool: 600,
					PricePreschool:        400,
					PriceSenior:           200,
				},
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewExperience(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			actual.ID, actual.ExperienceRevision.ExperienceID = "", "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperience_SetStatus(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name           string
		experience     *Experience
		expectedStatus ExperienceStatus
	}{
		{
			name: "archived",
			experience: &Experience{
				DeletedAt: gorm.DeletedAt{Time: now, Valid: true},
			},
			expectedStatus: ExperienceStatusArchived,
		},
		{
			name: "private",
			experience: &Experience{
				Public:  false,
				SoldOut: false,
				StartAt: now.AddDate(0, 0, -1),
				EndAt:   now.AddDate(0, 0, 1),
			},
			expectedStatus: ExperienceStatusPrivate,
		},
		{
			name: "sold out",
			experience: &Experience{
				Public:  true,
				SoldOut: true,
				StartAt: now.AddDate(0, 0, -1),
				EndAt:   now.AddDate(0, 0, 1),
			},
			expectedStatus: ExperienceStatusSoldOut,
		},
		{
			name: "waiting",
			experience: &Experience{
				Public:  true,
				SoldOut: false,
				StartAt: now.AddDate(0, 0, 1),
				EndAt:   now.AddDate(0, 0, 2),
			},
			expectedStatus: ExperienceStatusWaiting,
		},
		{
			name: "accepting",
			experience: &Experience{
				Public:  true,
				SoldOut: false,
				StartAt: now.AddDate(0, 0, -1),
				EndAt:   now.AddDate(0, 0, 1),
			},
			expectedStatus: ExperienceStatusAccepting,
		},
		{
			name: "finished",
			experience: &Experience{
				Public:  true,
				SoldOut: false,
				StartAt: now.AddDate(0, 0, -2),
				EndAt:   now.AddDate(0, 0, -1),
			},
			expectedStatus: ExperienceStatusFinished,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.experience.SetStatus(now)
			assert.Equal(t, tt.expectedStatus, tt.experience.Status)
		})
	}
}

func TestExperience_FillJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		experience *Experience
		expect     *Experience
		hasErr     bool
	}{
		{
			name: "success",
			experience: &Experience{
				Media: MultiExperienceMedia{
					{
						URL:         "https://and-period.jp/thumbnail.png",
						IsThumbnail: true,
					},
				},
				RecommendedPoints: []string{
					"ポイント1",
					"ポイント2",
				},
				HostLongitude: 136.251739,
				HostLatitude:  35.276833,
			},
			expect: &Experience{
				Media: MultiExperienceMedia{
					{
						URL:         "https://and-period.jp/thumbnail.png",
						IsThumbnail: true,
					},
				},
				MediaJSON: datatypes.JSON([]byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true}]`)),
				RecommendedPoints: []string{
					"ポイント1",
					"ポイント2",
				},
				RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
				HostLongitude:         136.251739,
				HostLatitude:          35.276833,
				HostGeolocation: mysql.Geometry{
					X: 136.251739,
					Y: 35.276833,
				},
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.experience.FillJSON()
			assert.Equal(t, err != nil, tt.hasErr)
			assert.Equal(t, tt.experience, tt.expect)
		})
	}
}

func TestExperience_Fill(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		revisions   map[string]*ExperienceRevision
		now         time.Time
		expect      Experiences
		hasErr      bool
	}{
		{
			name: "success",
			experiences: Experiences{
				{
					ID:                    "experience-id",
					CoordinatorID:         "coordinator-id",
					ProducerID:            "producer-id",
					TypeID:                "experience-type-id",
					Title:                 "じゃがいも収穫",
					Description:           "じゃがいもを収穫する体験",
					Public:                true,
					SoldOut:               false,
					Status:                ExperienceStatusUnknown,
					ThumbnailURL:          "http://example.com/thumbnail.png",
					Media:                 []*ExperienceMedia{},
					MediaJSON:             datatypes.JSON([]byte(`[{"url":"http://example.com/thumbnail.png","isThumbnail":true}]`)),
					RecommendedPoints:     []string{},
					RecommendedPointsJSON: datatypes.JSON([]byte(`["ポイント1","ポイント2"]`)),
					PromotionVideoURL:     "http://example.com/promotion.mp4",
					HostPostalCode:        "5220061",
					HostPrefectureCode:    25,
					HostCity:              "彦根市",
					HostAddressLine1:      "金亀町１−１",
					HostAddressLine2:      "",
					HostGeolocation: mysql.Geometry{
						X: 136.251739,
						Y: 35.276833,
					},
					StartAt:   now.AddDate(0, 0, -1),
					EndAt:     now.AddDate(0, 0, 1),
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			revisions: map[string]*ExperienceRevision{},
			now:       time.Now(),
			expect: Experiences{
				{
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
					ExperienceRevision: ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:            now.AddDate(0, 0, -1),
					EndAt:              now.AddDate(0, 0, 1),
					CreatedAt:          now,
					UpdatedAt:          now,
				},
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.experiences.Fill(tt.revisions, tt.now)
			assert.Equal(t, err != nil, tt.hasErr)
			assert.Equal(t, tt.experiences, tt.expect)
		})
	}
}

func TestExperiences_IDs(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		expect      []string
	}{
		{
			name: "success",
			experiences: Experiences{
				{
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
					HostPrefecture:        "東京都",
					HostPrefectureCode:    13,
					HostCity:              "千代田区",
					ExperienceRevision:    ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:               now.AddDate(0, 0, -1),
					EndAt:                 now.AddDate(0, 0, 1),
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			expect: []string{"experience-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.IDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestExperiences_CoordinatorIDs(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		expect      []string
	}{
		{
			name: "success",
			experiences: Experiences{
				{
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
					HostPrefecture:        "東京都",
					HostPrefectureCode:    13,
					HostCity:              "千代田区",
					ExperienceRevision:    ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:               now.AddDate(0, 0, -1),
					EndAt:                 now.AddDate(0, 0, 1),
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			expect: []string{"coordinator-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.CoordinatorIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestExperiences_ProducerIDs(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		expect      []string
	}{
		{
			name: "success",
			experiences: Experiences{
				{
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
					HostPrefecture:        "東京都",
					HostPrefectureCode:    13,
					HostCity:              "千代田区",
					ExperienceRevision:    ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:               now.AddDate(0, 0, -1),
					EndAt:                 now.AddDate(0, 0, 1),
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			expect: []string{"producer-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.ProducerIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestExperiences_ExperienceTypeIDs(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		expect      []string
	}{
		{
			name: "success",
			experiences: Experiences{
				{
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
					HostPrefecture:        "東京都",
					HostPrefectureCode:    13,
					HostCity:              "千代田区",
					ExperienceRevision:    ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:               now.AddDate(0, 0, -1),
					EndAt:                 now.AddDate(0, 0, 1),
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			expect: []string{"experience-type-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.ExperienceTypeIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestExperiences_Map(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		expect      map[string]*Experience
	}{
		{
			name: "success",
			experiences: Experiences{
				{
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
					HostPrefecture:        "東京都",
					HostPrefectureCode:    13,
					HostCity:              "千代田区",
					ExperienceRevision:    ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:               now.AddDate(0, 0, -1),
					EndAt:                 now.AddDate(0, 0, 1),
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			expect: map[string]*Experience{
				"experience-id": {
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
					HostPrefecture:        "東京都",
					HostPrefectureCode:    13,
					HostCity:              "千代田区",
					ExperienceRevision:    ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:               now.AddDate(0, 0, -1),
					EndAt:                 now.AddDate(0, 0, 1),
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.Map()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperiences_FilterByPublished(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		experiences Experiences
		expect      Experiences
	}{
		{
			name: "success",
			experiences: Experiences{
				{
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
					HostPrefecture:        "東京都",
					HostPrefectureCode:    13,
					HostCity:              "千代田区",
					ExperienceRevision:    ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:               now.AddDate(0, 0, -1),
					EndAt:                 now.AddDate(0, 0, 1),
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			expect: Experiences{
				{
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
					HostPrefecture:        "東京都",
					HostPrefectureCode:    13,
					HostCity:              "千代田区",
					ExperienceRevision:    ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:               now.AddDate(0, 0, -1),
					EndAt:                 now.AddDate(0, 0, 1),
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.experiences.FilterByPublished()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestMultiExperienceMedia_Marshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		media  MultiExperienceMedia
		expect []byte
		hasErr bool
	}{
		{
			name: "success",
			media: MultiExperienceMedia{
				{
					URL:         "https://and-period.jp/thumbnail.png",
					IsThumbnail: true,
				},
			},
			expect: []byte(`[{"url":"https://and-period.jp/thumbnail.png","isThumbnail":true}]`),
			hasErr: false,
		},
		{
			name:   "success is empty",
			media:  nil,
			expect: []byte{},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.media.Marshal()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
