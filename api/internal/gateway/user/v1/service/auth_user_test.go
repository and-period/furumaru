package service

import (
	"testing"

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
					ProviderType:  entity.UserAuthProviderTypeEmail,
					Email:         "test@and-period.jp",
					PhoneNumber:   "+819012345678",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					VerifiedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			notification: &entity.UserNotification{
				UserID:    "user-id",
				Disabled:  false,
				CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &AuthUser{
				AuthUser: response.AuthUser{
					ID:                  "user-id",
					Username:            "username",
					AccountID:           "account-id",
					ThumbnailURL:        "https://and-period.jp/thumbnail.png",
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
					ProviderType:  entity.UserAuthProviderTypeEmail,
					Email:         "test@and-period.jp",
					PhoneNumber:   "+819012345678",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					VerifiedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			notification: nil,
			expect: &AuthUser{
				AuthUser: response.AuthUser{
					ID:            "user-id",
					Username:      "username",
					AccountID:     "account-id",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
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
					ID:            "user-id",
					Username:      "username",
					AccountID:     "account-id",
					ThumbnailURL:  "https://and-period.jp/thumbnail.png",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
				},
			},
			expect: &response.AuthUser{
				ID:            "user-id",
				Username:      "username",
				AccountID:     "account-id",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
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
