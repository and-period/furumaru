package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
