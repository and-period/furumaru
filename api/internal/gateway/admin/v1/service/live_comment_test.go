package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/stretchr/testify/assert"
)

func TestLiveComments(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		comments mentity.BroadcastComments
		users    map[string]*User
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
			users: map[string]*User{
				"user-id": {
					User: response.User{
						ID:            "user-id",
						Status:        int32(UserStatusVerified),
						Registered:    true,
						AccountID:     "account-id",
						Username:      "username",
						Lastname:      "&.",
						Firstname:     "利用者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "りようしゃ",
						Email:         "test@example.com",
						PhoneNumber:   "+819012345678",
						ThumbnailURL:  "http://example.com/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "http://example.com/thumbnail_small.png", Size: int32(ImageSizeSmall)},
							{URL: "http://example.com/thumbnail_medium.png", Size: int32(ImageSizeMedium)},
							{URL: "http://example.com/thumbnail_large.png", Size: int32(ImageSizeLarge)},
						},
						CreatedAt: now.Unix(),
						UpdatedAt: now.Unix(),
					},
				},
			},
			expect: LiveComments{
				{
					LiveComment: response.LiveComment{
						ID:           "comment-id",
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
						Disabled:    false,
						PublishedAt: now.Unix(),
					},
				},
				{
					LiveComment: response.LiveComment{
						ID:           "disabled-id",
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						Comment:      "こんにちは",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Thumbnails: []*response.Image{
							{URL: "http://example.com/thumbnail_small.png", Size: int32(ImageSizeSmall)},
							{URL: "http://example.com/thumbnail_medium.png", Size: int32(ImageSizeMedium)},
							{URL: "http://example.com/thumbnail_large.png", Size: int32(ImageSizeLarge)},
						},
						Disabled:    true,
						PublishedAt: now.Unix(),
					},
				},
				{
					LiveComment: response.LiveComment{
						ID:           "unknown-id",
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
						ID:           "comment-id",
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Comment:      "こんにちは",
						Disabled:     false,
						PublishedAt:  now.Unix(),
					},
				},
			},
			expect: []*response.LiveComment{
				{
					ID:           "comment-id",
					UserID:       "user-id",
					Username:     "username",
					AccountID:    "account-id",
					ThumbnailURL: "http://example.com/thumbnail.png",
					Comment:      "こんにちは",
					Disabled:     false,
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
