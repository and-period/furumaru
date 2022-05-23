package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		role      int32
		expect    AdminRole
		expectErr error
	}{
		{
			name:      "success",
			role:      1,
			expect:    AdminRoleAdministrator,
			expectErr: nil,
		},
		{
			name:      "invalid role",
			role:      0,
			expect:    AdminRoleUnknown,
			expectErr: errInvalidAdminRole,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewAdminRole(tt.role)
			assert.ErrorIs(t, tt.expectErr, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

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
			name:   "coordinator",
			role:   AdminRoleCoordinator,
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

func TestAdministrator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewAdministratorParams
		expect *Admin
	}{
		{
			name: "success",
			params: &NewAdministratorParams{
				ID:            "admin-id",
				CognitoID:     "cognito-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expect: &Admin{
				ID:            "admin-id",
				CognitoID:     "cognito-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				Role:          AdminRoleAdministrator,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdministrator(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestProducer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewProducerParams
		expect *Admin
	}{
		{
			name: "success",
			params: &NewProducerParams{
				ID:            "admin-id",
				CognitoID:     "cognito-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &Admin{
				ID:            "admin-id",
				CognitoID:     "cognito-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				Email:         "test-admin@and-period.jp",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
				Role:          AdminRoleProducer,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewProducer(tt.params)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdmin_Name(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		admin  *Admin
		expect string
	}{
		{
			name:   "success",
			admin:  &Admin{Lastname: "&.", Firstname: "スタッフ"},
			expect: "&. スタッフ",
		},
		{
			name:   "success only lastname",
			admin:  &Admin{Lastname: "&.", Firstname: ""},
			expect: "&.",
		},
		{
			name:   "success only firstname",
			admin:  &Admin{Lastname: "", Firstname: "スタッフ"},
			expect: "スタッフ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admin.Name())
		})
	}
}

func TestAdmin_Map(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		admins Admins
		expect map[string]*Admin
	}{
		{
			name: "success",
			admins: Admins{
				{
					ID:        "admin-id",
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
			},
			expect: map[string]*Admin{
				"admin-id": {
					ID:        "admin-id",
					Lastname:  "&.",
					Firstname: "スタッフ",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admins.Map())
		})
	}
}
