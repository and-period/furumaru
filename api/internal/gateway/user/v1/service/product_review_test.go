package service

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

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
						ProviderType:  uentity.ProviderTypeEmail,
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
					ProductReview: response.ProductReview{
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
		expect  []*response.ProductReview
	}{
		{
			name: "success",
			reviews: ProductReviews{
				{
					ProductReview: response.ProductReview{
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
			expect: []*response.ProductReview{
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
