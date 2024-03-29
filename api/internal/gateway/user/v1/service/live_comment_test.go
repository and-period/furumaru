package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestLiveComments(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		comments mentity.BroadcastComments
		users    map[string]*uentity.User
		expect   LiveComments
	}{
		{
			name: "success",
			comments: mentity.BroadcastComments{
				{
					ID:          "comment-id",
					BroadcastID: "broadcast-id",
					UserID:      "user-id",
					Content:     "こんにちは",
					Disabled:    false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          "disabled-id",
					BroadcastID: "broadcast-id",
					UserID:      "user-id",
					Content:     "こんにちは",
					Disabled:    true,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					ID:          "unknown-id",
					BroadcastID: "broadcast-id",
					UserID:      "unknown-id",
					Content:     "こんにちは",
					Disabled:    false,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			users: map[string]*uentity.User{
				"user-id": {
					Member: uentity.Member{
						UserID:        "user-id",
						CognitoID:     "cognito-id",
						AccountID:     "account-id",
						Username:      "username",
						Lastname:      "&.",
						Firstname:     "利用者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "りようしゃ",
						ProviderType:  uentity.ProviderTypeEmail,
						Email:         "test@example.com",
						PhoneNumber:   "+819012345678",
						ThumbnailURL:  "http://example.com/thumbnail.png",
						Thumbnails:    common.Images{{URL: "http://example.com/thumbnail_small.png", Size: common.ImageSizeSmall}, {URL: "http://example.com/thumbnail_medium.png", Size: common.ImageSizeMedium}, {URL: "http://example.com/thumbnail_large.png", Size: common.ImageSizeLarge}},
						CreatedAt:     now,
						UpdatedAt:     now,
						VerifiedAt:    now,
					},
					Guest:      uentity.Guest{},
					ID:         "user-id",
					Status:     uentity.UserStatusVerified,
					Registered: true,
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			expect: LiveComments{
				{
					LiveComment: response.LiveComment{
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "http://example.com/thumbnail_small.png", Size: int32(ImageSizeSmall)},
							{URL: "http://example.com/thumbnail_medium.png", Size: int32(ImageSizeMedium)},
							{URL: "http://example.com/thumbnail_large.png", Size: int32(ImageSizeLarge)},
						},
						Comment:     "こんにちは",
						PublishedAt: now.Unix(),
					},
				},
				{
					LiveComment: response.LiveComment{
						UserID:       "",
						Username:     "",
						AccountID:    "",
						ThumbnailURL: "",
						Thumbnails:   []*response.Image{},
						Comment:      "こんにちは",
						PublishedAt:  now.Unix(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewLiveComments(tt.comments, tt.users)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestLiveComments_Response(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		comments LiveComments
		expect   []*response.LiveComment
	}{
		{
			name: "success",
			comments: LiveComments{
				{
					LiveComment: response.LiveComment{
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Comment:      "こんにちは",
						PublishedAt:  now.Unix(),
					},
				},
			},
			expect: []*response.LiveComment{
				{
					UserID:       "user-id",
					Username:     "username",
					AccountID:    "account-id",
					ThumbnailURL: "http://example.com/thumbnail.png",
					Comment:      "こんにちは",
					PublishedAt:  now.Unix(),
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.comments.Response())
		})
	}
}
