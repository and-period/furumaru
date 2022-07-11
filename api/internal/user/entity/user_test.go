package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		cognitoID    string
		providerType ProviderType
		email        string
		phoneNumber  string
		expect       *User
	}{
		{
			name:         "success",
			cognitoID:    "cognito-id",
			providerType: ProviderTypeEmail,
			email:        "test-user@and-period.jp",
			phoneNumber:  "+810000000000",
			expect: &User{
				CognitoID:    "cognito-id",
				ProviderType: ProviderTypeEmail,
				Email:        "test-user@and-period.jp",
				PhoneNumber:  "+810000000000",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewUser(tt.cognitoID, tt.providerType, tt.email, tt.phoneNumber)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestUser_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		user   *User
		expect string
	}{
		{
			name: "success",
			user: &User{
				Username:     "テストユーザー",
				CognitoID:    "cognito-id",
				ProviderType: ProviderTypeEmail,
				Email:        "test-user@and-period.jp",
				PhoneNumber:  "+810000000000",
			},
			expect: "テストユーザー",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.user.Name())
		})
	}
}
