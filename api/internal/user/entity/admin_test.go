package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminRole_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		role   AdminRole
		expect error
	}{
		{
			name:   "administrator",
			role:   AdminRoleAdministrator,
			expect: nil,
		},
		{
			name:   "producer",
			role:   AdminRoleProducer,
			expect: nil,
		},
		{
			name:   "unknown",
			role:   AdminRoleUnknown,
			expect: errInvalidAdminRole,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.role.Validate()
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}

func TestAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		adminID       string
		cognitoID     string
		lastname      string
		firstname     string
		lastnameKana  string
		firstnameKana string
		email         string
		role          AdminRole
		expect        *Admin
	}{
		{
			name:          "success",
			adminID:       "admin-id",
			cognitoID:     "cognito-id",
			lastname:      "&.",
			firstname:     "スタッフ",
			lastnameKana:  "あんどどっと",
			firstnameKana: "すたっふ",
			email:         "test-admin@and-period.jp",
			role:          AdminRoleAdministrator,
			expect: &Admin{
				ID:            "admin-id",
				CognitoID:     "cognito-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				Role:          AdminRoleAdministrator,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdmin(
				tt.adminID, tt.cognitoID,
				tt.lastname, tt.firstname, tt.lastnameKana, tt.firstnameKana,
				tt.email, tt.role,
			)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
