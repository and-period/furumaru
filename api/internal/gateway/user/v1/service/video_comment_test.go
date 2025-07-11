package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	mentity "github.com/and-period/furumaru/api/internal/media/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestVideoComments(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		comments mentity.VideoComments
		users    map[string]*uentity.User
		expect   VideoComments
	}{
		{
			name: "success",
			comments: mentity.VideoComments{
				{
					ID:        "comment-id",
					VideoID:   "Video-id",
					UserID:    "user-id",
					Content:   "こんにちは",
					Disabled:  false,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "disabled-id",
					VideoID:   "Video-id",
					UserID:    "user-id",
					Content:   "こんにちは",
					Disabled:  true,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "unknown-id",
					VideoID:   "Video-id",
					UserID:    "unknown-id",
					Content:   "こんにちは",
					Disabled:  false,
					CreatedAt: now,
					UpdatedAt: now,
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
						ProviderType:  uentity.UserAuthProviderTypeEmail,
						Email:         "test@example.com",
						PhoneNumber:   "+819012345678",
						ThumbnailURL:  "http://example.com/thumbnail.png",
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
			expect: VideoComments{
				{
					VideoComment: response.VideoComment{
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Comment:      "こんにちは",
						PublishedAt:  now.Unix(),
					},
				},
				{
					VideoComment: response.VideoComment{
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
		expect   []*response.VideoComment
	}{
		{
			name: "success",
			comments: VideoComments{
				{
					VideoComment: response.VideoComment{
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Comment:      "こんにちは",
						PublishedAt:  now.Unix(),
					},
				},
			},
			expect: []*response.VideoComment{
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.comments.Response())
		})
	}
}
