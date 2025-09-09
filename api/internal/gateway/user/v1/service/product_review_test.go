package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestProductReviewReactionType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		reactionType entity.ProductReviewReactionType
		request      int32
		expect       ProductReviewReactionType
	}{
		{
			name:         "like",
			reactionType: entity.ProductReviewReactionTypeLike,
			request:      1,
			expect:       ProductReviewReactionTypeLike,
		},
		{
			name:         "dislike",
			reactionType: entity.ProductReviewReactionTypeDislike,
			request:      2,
			expect:       ProductReviewReactionTypeDislike,
		},
		{
			name:         "unknown",
			reactionType: entity.ProductReviewReactionTypeUnknown,
			request:      0,
			expect:       ProductReviewReactionTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := NewProductReviewReactionType(tt.reactionType)
			assert.Equal(t, tt.expect, actual)

			req, _ := NewProductReviewReactionTypeFromRequest(tt.request)
			assert.Equal(t, actual, req)

			assert.Equal(t, tt.request, actual.Response())
			assert.Equal(t, tt.reactionType, actual.StoreEntity())
		})
	}
}

func TestProductReviews(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name    string
		reviews entity.ProductReviews
		users   map[string]*uentity.User
		expect  ProductReviews
	}{
		{
			name: "success",
			reviews: entity.ProductReviews{
				{
					ID:        "review-id",
					ProductID: "product-id",
					UserID:    "user-id",
					Rate:      5,
					Title:     "title",
					Comment:   "comment",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "disabled-id",
					ProductID: "product-id",
					UserID:    "disabled-id",
					Rate:      5,
					Title:     "title",
					Comment:   "comment",
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
			expect: ProductReviews{
				{
					ProductReview: types.ProductReview{
						ID:           "review-id",
						ProductID:    "product-id",
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
			actual := NewProductReviews(tt.reviews, tt.users)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductReviews_Response(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name    string
		reviews ProductReviews
		expect  []*types.ProductReview
	}{
		{
			name: "success",
			reviews: ProductReviews{
				{
					ProductReview: types.ProductReview{
						ID:           "review-id",
						ProductID:    "product-id",
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
			expect: []*types.ProductReview{
				{
					ID:           "review-id",
					ProductID:    "product-id",
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

func TestProductReviewReactions(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name      string
		reactions entity.ProductReviewReactions
		expect    ProductReviewReactions
	}{
		{
			name: "success",
			reactions: entity.ProductReviewReactions{
				{
					ReviewID:     "review-id",
					UserID:       "user-id",
					ReactionType: entity.ProductReviewReactionTypeLike,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			expect: ProductReviewReactions{
				{
					ProductReviewReaction: types.ProductReviewReaction{
						ReviewID:     "review-id",
						ReactionType: int32(ProductReviewReactionTypeLike),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProductReviewReactions(tt.reactions)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProductReviewReactions_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		reactions ProductReviewReactions
		expect    []*types.ProductReviewReaction
	}{
		{
			name: "success",
			reactions: ProductReviewReactions{
				{
					ProductReviewReaction: types.ProductReviewReaction{
						ReviewID:     "review-id",
						ReactionType: int32(ProductReviewReactionTypeLike),
					},
				},
			},
			expect: []*types.ProductReviewReaction{
				{
					ReviewID:     "review-id",
					ReactionType: int32(ProductReviewReactionTypeLike),
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
