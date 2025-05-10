package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdministrator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewAdministratorParams
		expect *Administrator
	}{
		{
			name: "success",
			params: &NewAdministratorParams{
				Admin: &Admin{
					ID:            "admin-id",
					CognitoID:     "cognito-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-admin@and-period.jp",
				},
				PhoneNumber: "+819012345678",
			},
			expect: &Administrator{
				AdminID:     "admin-id",
				PhoneNumber: "+819012345678",
				Admin: Admin{
					CognitoID:     "cognito-id",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "すたっふ",
					Email:         "test-admin@and-period.jp",
				},
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdministrator(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdministrators_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		administrators Administrators
		expect         []string
	}{
		{
			name: "success",
			administrators: Administrators{
				{AdminID: "administrator-id01"},
				{AdminID: "administrator-id02"},
			},
			expect: []string{
				"administrator-id01",
				"administrator-id02",
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.administrators.IDs())
		})
	}
}

func TestAdministrators_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		administrators Administrators
		admins         map[string]*Admin
		groups         map[string]AdminGroupUsers
		expect         Administrators
	}{
		{
			name: "success",
			administrators: Administrators{
				{
					AdminID: "admin-id01",
				},
				{
					AdminID: "admin-id02",
				},
			},
			admins: map[string]*Admin{
				"admin-id01": {
					ID:        "admin-id01",
					CognitoID: "cognito-id",
					Type:      AdminTypeAdministrator,
				},
			},
			groups: map[string]AdminGroupUsers{
				"admin-id01": {
					{
						GroupID: "group-id",
						AdminID: "admin-id01",
					},
				},
			},
			expect: Administrators{
				{
					AdminID: "admin-id01",
					Admin: Admin{
						ID:        "admin-id01",
						CognitoID: "cognito-id",
						Type:      AdminTypeAdministrator,
						Status:    AdminStatusInvited,
						GroupIDs:  []string{"group-id"},
					},
				},
				{
					AdminID: "admin-id02",
					Admin: Admin{
						ID:       "admin-id02",
						Type:     AdminTypeAdministrator,
						Status:   AdminStatusInvited,
						GroupIDs: []string{},
					},
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.administrators.Fill(tt.admins, tt.groups)
			assert.Equal(t, tt.expect, tt.administrators)
		})
	}
}
