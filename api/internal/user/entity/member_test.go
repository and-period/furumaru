package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestMember_Name(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		member *Member
		expect string
	}{
		{
			name: "success",
			member: &Member{
				UserID:         "user-id",
				CognitoID:      "cognito-id",
				AccountID:      "account-id",
				Username:       "username",
				Lastname:       "&.",
				Firstname:      "利用者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "りようしゃ",
				ProviderType:   ProviderTypeEmail,
				Email:          "test@and-period.jp",
				PhoneNumber:    "+819012345678",
				ThumbnailURL:   "http://example.com/image.png",
				Thumbnails:     nil,
				ThumbnailsJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/media.png","size":1}]`)),
			},
			expect: "&. 利用者",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.member.Name())
		})
	}
}

func TestMember_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		member *Member
		expect *Member
		hasErr bool
	}{
		{
			name: "success",
			member: &Member{
				UserID:         "user-id",
				CognitoID:      "cognito-id",
				AccountID:      "account-id",
				Username:       "username",
				Lastname:       "&.",
				Firstname:      "利用者",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "りようしゃ",
				ProviderType:   ProviderTypeEmail,
				Email:          "test@and-period.jp",
				PhoneNumber:    "+819012345678",
				ThumbnailURL:   "http://example.com/image.png",
				Thumbnails:     nil,
				ThumbnailsJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/media.png","size":1}]`)),
			},
			expect: &Member{
				UserID:        "user-id",
				CognitoID:     "cognito-id",
				AccountID:     "account-id",
				Username:      "username",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				ProviderType:  ProviderTypeEmail,
				Email:         "test@and-period.jp",
				PhoneNumber:   "+819012345678",
				ThumbnailURL:  "http://example.com/image.png",
				Thumbnails: common.Images{
					{Size: common.ImageSizeSmall, URL: "http://example.com/media.png"},
				},
				ThumbnailsJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/media.png","size":1}]`)),
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.member.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.member)
		})
	}
}

func TestMembers_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		members Members
		expect  map[string]*Member
	}{
		{
			name: "success",
			members: Members{
				{
					UserID:    "user-id01",
					CognitoID: "cognito-id01",
					AccountID: "account-id01",
					Username:  "username",
				},
				{
					UserID:    "user-id02",
					CognitoID: "cognito-id02",
					AccountID: "account-id02",
					Username:  "username",
				},
			},
			expect: map[string]*Member{
				"user-id01": {
					UserID:    "user-id01",
					CognitoID: "cognito-id01",
					AccountID: "account-id01",
					Username:  "username",
				},
				"user-id02": {
					UserID:    "user-id02",
					CognitoID: "cognito-id02",
					AccountID: "account-id02",
					Username:  "username",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.members.Map())
		})
	}
}

func TestMembers_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		members Members
		expect  Members
		hasErr  bool
	}{
		{
			name: "success",
			members: Members{
				{
					UserID:         "user-id",
					CognitoID:      "cognito-id",
					AccountID:      "account-id",
					Username:       "username",
					Lastname:       "&.",
					Firstname:      "利用者",
					LastnameKana:   "あんどどっと",
					FirstnameKana:  "りようしゃ",
					ProviderType:   ProviderTypeEmail,
					Email:          "test@and-period.jp",
					PhoneNumber:    "+819012345678",
					ThumbnailURL:   "http://example.com/image.png",
					Thumbnails:     nil,
					ThumbnailsJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/media.png","size":1}]`)),
				},
			},
			expect: Members{
				{
					UserID:        "user-id",
					CognitoID:     "cognito-id",
					AccountID:     "account-id",
					Username:      "username",
					Lastname:      "&.",
					Firstname:     "利用者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "りようしゃ",
					ProviderType:  ProviderTypeEmail,
					Email:         "test@and-period.jp",
					PhoneNumber:   "+819012345678",
					ThumbnailURL:  "http://example.com/image.png",
					Thumbnails: common.Images{
						{Size: common.ImageSizeSmall, URL: "http://example.com/media.png"},
					},
					ThumbnailsJSON: datatypes.JSON([]byte(`[{"url":"http://example.com/media.png","size":1}]`)),
				},
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.members.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.members)
		})
	}
}
