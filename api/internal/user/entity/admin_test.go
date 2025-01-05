package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAdminType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		adminType int32
		expect    AdminType
		expectErr error
	}{
		{
			name:      "success",
			adminType: 1,
			expect:    AdminTypeAdministrator,
			expectErr: nil,
		},
		{
			name:      "invalid role",
			adminType: 0,
			expect:    AdminTypeUnknown,
			expectErr: errInvalidAdminRole,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewAdminType(tt.adminType)
			assert.ErrorIs(t, tt.expectErr, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdminType_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		adminType AdminType
		expect    error
	}{
		{
			name:      "administrator",
			adminType: AdminTypeAdministrator,
			expect:    nil,
		},
		{
			name:      "coordinator",
			adminType: AdminTypeCoordinator,
			expect:    nil,
		},
		{
			name:      "producer",
			adminType: AdminTypeProducer,
			expect:    nil,
		},
		{
			name:      "unknown",
			adminType: AdminTypeUnknown,
			expect:    errInvalidAdminRole,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.adminType.Validate()
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
		name           string
		admin          *Admin
		groups         RelatedAdminGroups
		expectGroupIDs []string
		expectStatus   AdminStatus
	}{
		{
			name: "producer",
			admin: &Admin{
				Type: AdminTypeProducer,
			},
			groups: RelatedAdminGroups{
				{
					AdminID: "admin-id",
					GroupID: "group-id",
				},
			},
			expectGroupIDs: []string{"group-id"},
			expectStatus:   AdminStatusDeactivated,
		},
		{
			name: "invited",
			admin: &Admin{
				Type:          AdminTypeCoordinator,
				FirstSignInAt: time.Time{},
			},
			groups: RelatedAdminGroups{
				{
					AdminID: "admin-id",
					GroupID: "group-id",
				},
			},
			expectGroupIDs: []string{"group-id"},
			expectStatus:   AdminStatusInvited,
		},
		{
			name: "activated",
			admin: &Admin{
				Type:          AdminTypeCoordinator,
				FirstSignInAt: now,
			},
			groups: RelatedAdminGroups{
				{
					AdminID: "admin-id",
					GroupID: "group-id",
				},
			},
			expectGroupIDs: []string{"group-id"},
			expectStatus:   AdminStatusActivated,
		},
		{
			name: "deactivated",
			admin: &Admin{
				Type:          AdminTypeCoordinator,
				FirstSignInAt: now,
				DeletedAt: gorm.DeletedAt{
					Time:  now,
					Valid: true,
				},
			},
			groups:         nil,
			expectGroupIDs: []string{},
			expectStatus:   AdminStatusDeactivated,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.admin.Fill(tt.groups)
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
					Type:      AdminTypeAdministrator,
				},
			},
			expect: map[string]*Admin{
				"admin-id": {
					ID:        "admin-id",
					CognitoID: "cognito-id",
					Type:      AdminTypeAdministrator,
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

func TestAdmins_GroupByType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		admins Admins
		expect map[AdminType]Admins
	}{
		{
			name: "success",
			admins: Admins{
				{
					ID:        "admin-id",
					CognitoID: "cognito-id",
					Type:      AdminTypeAdministrator,
				},
			},
			expect: map[AdminType]Admins{
				AdminTypeAdministrator: {
					{
						ID:        "admin-id",
						CognitoID: "cognito-id",
						Type:      AdminTypeAdministrator,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admins.GroupByType())
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
					Type:      AdminTypeAdministrator,
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
					Type:      AdminTypeAdministrator,
				},
				{
					ID:        "admin-id",
					CognitoID: "cognito-id",
					Device:    "",
					Type:      AdminTypeAdministrator,
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
