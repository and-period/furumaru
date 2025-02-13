package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestExperienceReviewReactionType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		reactionType entity.ExperienceReviewReactionType
		request      int32
		expect       ExperienceReviewReactionType
	}{
		{
			name:         "like",
			reactionType: entity.ExperienceReviewReactionTypeLike,
			request:      1,
			expect:       ExperienceReviewReactionTypeLike,
		},
		{
			name:         "dislike",
			reactionType: entity.ExperienceReviewReactionTypeDislike,
			request:      2,
			expect:       ExperienceReviewReactionTypeDislike,
		},
		{
			name:         "unknown",
			reactionType: entity.ExperienceReviewReactionTypeUnknown,
			request:      0,
			expect:       ExperienceReviewReactionTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := NewExperienceReviewReactionType(tt.reactionType)
			assert.Equal(t, tt.expect, actual)

			req, _ := NewExperienceReviewReactionTypeFromRequest(tt.request)
			assert.Equal(t, actual, req)

			assert.Equal(t, tt.request, actual.Response())
			assert.Equal(t, tt.reactionType, actual.StoreEntity())
		})
	}
}

func TestExperienceReviews(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name    string
		reviews entity.ExperienceReviews
		users   map[string]*uentity.User
		expect  ExperienceReviews
	}{
		{
			name: "success",
			reviews: entity.ExperienceReviews{
				{
					ID:           "review-id",
					ExperienceID: "experience-id",
					UserID:       "user-id",
					Rate:         5,
					Title:        "title",
					Comment:      "comment",
					CreatedAt:    now,
					UpdatedAt:    now,
				},
				{
					ID:           "disabled-id",
					ExperienceID: "experience-id",
					UserID:       "disabled-id",
					Rate:         5,
					Title:        "title",
					Comment:      "comment",
					CreatedAt:    now,
					UpdatedAt:    now,
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
			expect: ExperienceReviews{
				{
					ExperienceReview: response.ExperienceReview{
						ID:           "review-id",
						ExperienceID: "experience-id",
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Rate:         5,
						Title:        "title",
						Comment:      "comment",
						PublishedAt:  now.Unix(),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewExperienceReviews(tt.reviews, tt.users)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceReviews_Response(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name    string
		reviews ExperienceReviews
		expect  []*response.ExperienceReview
	}{
		{
			name: "success",
			reviews: ExperienceReviews{
				{
					ExperienceReview: response.ExperienceReview{
						ID:           "review-id",
						ExperienceID: "experience-id",
						UserID:       "user-id",
						Username:     "username",
						AccountID:    "account-id",
						ThumbnailURL: "http://example.com/thumbnail.png",
						Rate:         5,
						Title:        "title",
						Comment:      "comment",
						PublishedAt:  now.Unix(),
					},
				},
			},
			expect: []*response.ExperienceReview{
				{
					ID:           "review-id",
					ExperienceID: "experience-id",
					UserID:       "user-id",
					Username:     "username",
					AccountID:    "account-id",
					ThumbnailURL: "http://example.com/thumbnail.png",
					Rate:         5,
					Title:        "title",
					Comment:      "comment",
					PublishedAt:  now.Unix(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.reviews.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceReviewReactions(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		reactions entity.ExperienceReviewReactions
		expect    ExperienceReviewReactions
	}{
		{
			name: "success",
			reactions: entity.ExperienceReviewReactions{
				{
					ReviewID:     "review-id",
					UserID:       "user-id",
					ReactionType: entity.ExperienceReviewReactionTypeLike,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expect: ExperienceReviewReactions{
				{
					ExperienceReviewReaction: response.ExperienceReviewReaction{
						ReviewID:     "review-id",
						ReactionType: int32(ExperienceReviewReactionTypeLike),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewExperienceReviewReactions(tt.reactions)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestExperienceReviewReactions_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		reactions ExperienceReviewReactions
		expect    []*response.ExperienceReviewReaction
	}{
		{
			name: "success",
			reactions: ExperienceReviewReactions{
				{
					ExperienceReviewReaction: response.ExperienceReviewReaction{
						ReviewID:     "review-id",
						ReactionType: int32(ExperienceReviewReactionTypeLike),
					},
				},
			},
			expect: []*response.ExperienceReviewReaction{
				{
					ReviewID:     "review-id",
					ReactionType: int32(ExperienceReviewReactionTypeLike),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.reactions.Response()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
