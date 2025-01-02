package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAdminRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		role      int32
		expect    LegacyAdminRole
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
		role   LegacyAdminRole
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

func TestAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		params *NewAdminParams
		expect *Admin
	}{
		{
			name: "success",
			params: &NewAdminParams{
				CognitoID:     "cognito-id",
				Type:          AdminTypeAdministrator,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
			expect: &Admin{
				CognitoID:     "cognito-id",
				Role:          AdminRoleAdministrator,
				Type:          AdminTypeAdministrator,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdmin(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdmin_Name(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 6, 25, 18, 30, 0, 0)
	tests := []struct {
		name   string
		admin  *Admin
		expect string
	}{
		{
			name: "success",
			admin: &Admin{
				ID:            "admin-id",
				Role:          AdminRoleAdministrator,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "すたっふ",
				Email:         "test-admin@and-period.jp",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			expect: "&. スタッフ",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.admin.Name()
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdmin_Fill(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name         string
		admin        *Admin
		expectStatus AdminStatus
	}{
		{
			name: "producer",
			admin: &Admin{
				Role: AdminRoleProducer,
			},
			expectStatus: AdminStatusDeactivated,
		},
		{
			name: "invited",
			admin: &Admin{
				Role:          AdminRoleCoordinator,
				FirstSignInAt: time.Time{},
			},
			expectStatus: AdminStatusInvited,
		},
		{
			name: "activated",
			admin: &Admin{
				Role:          AdminRoleCoordinator,
				FirstSignInAt: now,
			},
			expectStatus: AdminStatusActivated,
		},
		{
			name: "deactivated",
			admin: &Admin{
				Role:          AdminRoleCoordinator,
				FirstSignInAt: now,
				DeletedAt: gorm.DeletedAt{
					Time:  now,
					Valid: true,
				},
			},
			expectStatus: AdminStatusDeactivated,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.admin.Fill()
			assert.Equal(t, tt.expectStatus, tt.admin.Status)
		})
	}
}

func TestAdmins_Map(t *testing.T) {
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
					CognitoID: "cognito-id",
					Role:      AdminRoleAdministrator,
				},
			},
			expect: map[string]*Admin{
				"admin-id": {
					ID:        "admin-id",
					CognitoID: "cognito-id",
					Role:      AdminRoleAdministrator,
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

func TestAdmins_GroupByRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		admins Admins
		expect map[LegacyAdminRole]Admins
	}{
		{
			name: "success",
			admins: Admins{
				{
					ID:        "admin-id",
					CognitoID: "cognito-id",
					Role:      AdminRoleAdministrator,
				},
			},
			expect: map[LegacyAdminRole]Admins{
				AdminRoleAdministrator: {
					{
						ID:        "admin-id",
						CognitoID: "cognito-id",
						Role:      AdminRoleAdministrator,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admins.GroupByRole())
		})
	}
}

func TestAdmins_IDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		s      Admins
		expect []string
	}{
		{
			name: "success",
			s: Admins{
				{
					ID:        "admin-id",
					CognitoID: "cognito-id",
					Role:      AdminRoleAdministrator,
				},
			},
			expect: []string{"admin-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.s.IDs())
		})
	}
}

func TestAdmins_Devices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		s      Admins
		expect []string
	}{
		{
			name: "success",
			s: Admins{
				{
					ID:        "admin-id",
					CognitoID: "cognito-id",
					Device:    "instance-id",
					Role:      AdminRoleAdministrator,
				},
				{
					ID:        "admin-id",
					CognitoID: "cognito-id",
					Device:    "",
					Role:      AdminRoleAdministrator,
				},
			},
			expect: []string{"instance-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.s.Devices())
		})
	}
}
