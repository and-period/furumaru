package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

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
