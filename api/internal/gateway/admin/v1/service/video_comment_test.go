package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/stretchr/testify/assert"
)

func TestVideoComments(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		comments mentity.VideoComments
		users    map[string]*User
		expect   VideoComments
	}{
		{
			name: "success",
			comments: mentity.VideoComments{
				{
					ID:        "comment-id",
					VideoID:   "video-id",
					UserID:    "user-id",
					Content:   "こんにちは",
					Disabled:  false,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "disabled-id",
					VideoID:   "video-id",
					UserID:    "user-id",
					Content:   "こんにちは",
					Disabled:  true,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "unknown-id",
					VideoID:   "video-id",
					UserID:    "unknown-id",
					Content:   "こんにちは",
					Disabled:  false,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			users: map[string]*User{
				"user-id": {
					User: types.User{
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
						CreatedAt:     now.Unix(),
						UpdatedAt:     now.Unix(),
					},
				},
			},
			expect: VideoComments{
				{
					VideoComment: types.VideoComment{
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
				{
					VideoComment: types.VideoComment{
						ID:           "disabled-id",
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						Comment:      "こんにちは",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Disabled:     true,
						PublishedAt:  now.Unix(),
					},
				},
				{
					VideoComment: types.VideoComment{
						ID:           "unknown-id",
						UserID:       "",
						Username:     "",
						AccountID:    "",
						ThumbnailURL: "",
						Comment:      "こんにちは",
						PublishedAt:  now.Unix(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewVideoComments(tt.comments, tt.users)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestVideoComments_Response(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		comments VideoComments
		expect   []*types.VideoComment
	}{
		{
			name: "success",
			comments: VideoComments{
				{
					VideoComment: types.VideoComment{
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
			expect: []*types.VideoComment{
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.comments.Response())
		})
	}
}
