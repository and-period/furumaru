package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		userID       string
		cognitoID    string
		providerType ProviderType
		email        string
		phoneNumber  string
		expect       *User
	}{
		{
			name:         "success",
			userID:       "user-id",
			cognitoID:    "cognito-id",
			providerType: ProviderTypeEmail,
			email:        "test-user@and-period.jp",
			phoneNumber:  "+810000000000",
			expect: &User{
				ID:           "user-id",
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
			actual := NewUser(tt.userID, tt.cognitoID, tt.providerType, tt.email, tt.phoneNumber)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
