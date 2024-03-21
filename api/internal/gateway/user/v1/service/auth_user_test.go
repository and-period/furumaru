package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestAuthUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		user         *entity.User
		notification *entity.UserNotification
		expect       *AuthUser
	}{
		{
			name: "success with notification",
			user: &entity.User{
				ID:         "user-id",
				Registered: true,
				CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				Member: entity.Member{
					UserID:        "user-id",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					ProviderType:  entity.ProviderTypeEmail,
					Email:         "test@and-period.jp",
					PhoneNumber:   "+819012345678",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Thumbnails: common.Images{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: common.ImageSizeSmall},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: common.ImageSizeMedium},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: common.ImageSizeLarge},
					},
					VerifiedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			notification: &entity.UserNotification{
				UserID:        "user-id",
				EmailDisabled: false,
				CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &AuthUser{
				AuthUser: response.AuthUser{
					ID:           "user-id",
					Username:     "username",
					AccountID:    "account-id",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					Lastname:            "&.",
					Firstname:           "利用者",
					LastnameKana:        "あんどどっと",
					FirstnameKana:       "りようしゃ",
					Email:               "test@and-period.jp",
					NotificationEnabled: true,
				},
			},
		},
		{
			name: "success without notification",
			user: &entity.User{
				ID:         "user-id",
				Registered: true,
				CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				Member: entity.Member{
					UserID:        "user-id",
					Username:      "username",
					AccountID:     "account-id",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					ProviderType:  entity.ProviderTypeEmail,
					Email:         "test@and-period.jp",
					PhoneNumber:   "+819012345678",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Thumbnails: common.Images{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: common.ImageSizeSmall},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: common.ImageSizeMedium},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: common.ImageSizeLarge},
					},
					VerifiedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			notification: nil,
			expect: &AuthUser{
				AuthUser: response.AuthUser{
					ID:           "user-id",
					Username:     "username",
					AccountID:    "account-id",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					Email:         "test@and-period.jp",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAuthUser(tt.user, tt.notification))
		})
	}
}

func TestAuthUser_Response(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		auth   *AuthUser
		expect *response.AuthUser
	}{
		{
			name: "success",
			auth: &AuthUser{
				AuthUser: response.AuthUser{
					ID:           "user-id",
					Username:     "username",
					AccountID:    "account-id",
					ThumbnailURL: "https://and-period.jp/thumbnail.png",
					Thumbnails: []*response.Image{
						{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
						{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
						{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
					},
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
				},
			},
			expect: &response.AuthUser{
				ID:           "user-id",
				Username:     "username",
				AccountID:    "account-id",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
				Thumbnails: []*response.Image{
					{URL: "https://and-period.jp/thumbnail_240.png", Size: int32(ImageSizeSmall)},
					{URL: "https://and-period.jp/thumbnail_675.png", Size: int32(ImageSizeMedium)},
					{URL: "https://and-period.jp/thumbnail_900.png", Size: int32(ImageSizeLarge)},
				},
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.auth.Response())
		})
	}
}
