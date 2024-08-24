package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExperienceRevisions_ExperienceIDs(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name      string
		revisions ExperienceRevisions
		expect    []string
	}{
		{
			name: "success",
			revisions: ExperienceRevisions{
				{
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
			},
			expect: []string{"experience-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.revisions.ExperienceIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestExperienceRevisions_MapByExperienceID(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name      string
		revisions ExperienceRevisions
		expect    map[string]*ExperienceRevision
	}{
		{
			name: "success",
			revisions: ExperienceRevisions{
				{
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
			},
			expect: map[string]*ExperienceRevision{
				"experience-id": {
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.revisions.MapByExperienceID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceRevisions_Merge(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name        string
		revisions   ExperienceRevisions
		experiences map[string]*Experience
		expect      Experiences
		hasErr      bool
	}{
		{
			name: "success",
			revisions: ExperienceRevisions{
				{
					ID:                    1,
					ExperienceID:          "experience-id01",
					PriceAdult:            1000,
					PriceJuniorHighSchool: 800,
					PriceElementarySchool: 600,
					PricePreschool:        400,
					PriceSenior:           700,
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    2,
					ExperienceID:          "experience-id02",
					PriceAdult:            1000,
					PriceJuniorHighSchool: 800,
					PriceElementarySchool: 600,
					PricePreschool:        400,
					PriceSenior:           700,
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			experiences: map[string]*Experience{
				"experience-id01": {
					ID:            "experience-id01",
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
					TypeID:        "experience-type-id",
					Title:         "じゃがいも収穫",
					Description:   "じゃがいもを収穫する体験です。",
					Public:        true,
					SoldOut:       false,
					Status:        ExperienceStatusAccepting,
					ThumbnailURL:  "http://example.com/thumbnail.png",
					Media: MultiExperienceMedia{
						{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
						{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
					},
					RecommendedPoints: []string{
						"じゃがいもを収穫する楽しさを体験できます。",
						"新鮮なじゃがいもを持ち帰ることができます。",
					},
					PromotionVideoURL:  "http://example.com/promotion.mp4",
					HostPrefecture:     "滋賀県",
					HostPrefectureCode: 25,
					HostCity:           "彦根市",
					StartAt:            now.AddDate(0, 0, -1),
					EndAt:              now.AddDate(0, 0, 1),
					CreatedAt:          now,
					UpdatedAt:          now,
				},
			},
			expect: Experiences{
				{
					ID:            "experience-id01",
					CoordinatorID: "coordinator-id",
					ProducerID:    "producer-id",
					TypeID:        "experience-type-id",
					Title:         "じゃがいも収穫",
					Description:   "じゃがいもを収穫する体験です。",
					Public:        true,
					SoldOut:       false,
					Status:        ExperienceStatusAccepting,
					ThumbnailURL:  "http://example.com/thumbnail.png",
					Media: MultiExperienceMedia{
						{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
						{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
					},
					RecommendedPoints: []string{
						"じゃがいもを収穫する楽しさを体験できます。",
						"新鮮なじゃがいもを持ち帰ることができます。",
					},
					PromotionVideoURL:  "http://example.com/promotion.mp4",
					HostPrefecture:     "滋賀県",
					HostPrefectureCode: 25,
					HostCity:           "彦根市",
					ExperienceRevision: ExperienceRevision{
						ID:                    1,
						ExperienceID:          "experience-id01",
						PriceAdult:            1000,
						PriceJuniorHighSchool: 800,
						PriceElementarySchool: 600,
						PricePreschool:        400,
						PriceSenior:           700,
						CreatedAt:             now,
						UpdatedAt:             now,
					},
					StartAt:   now.AddDate(0, 0, -1),
					EndAt:     now.AddDate(0, 0, 1),
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID: "experience-id02",
					ExperienceRevision: ExperienceRevision{
						ID:                    2,
						ExperienceID:          "experience-id02",
						PriceAdult:            1000,
						PriceJuniorHighSchool: 800,
						PriceElementarySchool: 600,
						PricePreschool:        400,
						PriceSenior:           700,
						CreatedAt:             now,
						UpdatedAt:             now,
					},
				},
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := tt.revisions.Merge(tt.experiences)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}
