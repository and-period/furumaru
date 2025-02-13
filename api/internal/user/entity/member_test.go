package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				UserID:        "user-id",
				CognitoID:     "cognito-id",
				AccountID:     "account-id",
				Username:      "username",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
				ProviderType:  UserAuthProviderTypeEmail,
				Email:         "test@and-period.jp",
				PhoneNumber:   "+819012345678",
				ThumbnailURL:  "http://example.com/image.png",
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
