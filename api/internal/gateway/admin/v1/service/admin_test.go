package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestAdminRole(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		role                entity.AdminRole
		expect              AdminRole
		expectString        string
		expectIsCoordinator bool
		expectResponse      int32
	}{
		{
			name:                "administrator",
			role:                entity.AdminRoleAdministrator,
			expect:              AdminRoleAdministrator,
			expectString:        "administrator",
			expectIsCoordinator: false,
			expectResponse:      1,
		},
		{
			name:                "coordinator",
			role:                entity.AdminRoleCoordinator,
			expect:              AdminRoleCoordinator,
			expectString:        "coordinator",
			expectIsCoordinator: true,
			expectResponse:      2,
		},
		{
			name:                "producer",
			role:                entity.AdminRoleProducer,
			expect:              AdminRoleProducer,
			expectString:        "producer",
			expectIsCoordinator: false,
			expectResponse:      3,
		},
		{
			name:                "unknown",
			role:                entity.AdminRoleUnknown,
			expect:              AdminRoleUnknown,
			expectString:        "unknown",
			expectIsCoordinator: false,
			expectResponse:      0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminRole(tt.role)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectString, actual.String())
			assert.Equal(t, tt.expectIsCoordinator, actual.IsCoordinator())
			assert.Equal(t, tt.expectResponse, actual.Response())
		})
	}
}

func TestAdminStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		role           entity.AdminStatus
		expect         AdminStatus
		expectResponse int32
	}{
		{
			name:           "invited",
			role:           entity.AdminStatusInvited,
			expect:         AdminStatusInvited,
			expectResponse: 1,
		},
		{
			name:           "activated",
			role:           entity.AdminStatusActivated,
			expect:         AdminStatusActivated,
			expectResponse: 2,
		},
		{
			name:           "deactivated",
			role:           entity.AdminStatusDeactivated,
			expect:         AdminStatusDeactivated,
			expectResponse: 3,
		},
		{
			name:           "unknown",
			role:           entity.AdminStatusUnknown,
			expect:         AdminStatusUnknown,
			expectResponse: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminStatus(tt.role)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectResponse, actual.Response())
		})
	}
}

func TestAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admin  *entity.Admin
		expect *Admin
	}{
		{
			name: "success",
			admin: &entity.Admin{
				ID:            "admin-id",
				Role:          entity.AdminRoleAdministrator,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin@and-period.jp",
				CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Admin{
				Admin: response.Admin{
					ID:            "admin-id",
					Role:          entity.AdminRoleAdministrator,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin@and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAdmin(tt.admin))
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
			name: "success",
			admin: &Admin{
				Admin: response.Admin{
					Lastname:  "&.",
					Firstname: "管理者",
				},
			},
			expect: "&. 管理者",
		},
		{
			name: "success only lastname",
			admin: &Admin{
				Admin: response.Admin{
					Lastname:  "&.",
					Firstname: "",
				},
			},
			expect: "&.",
		},
		{
			name: "success only firstname",
			admin: &Admin{
				Admin: response.Admin{
					Lastname:  "",
					Firstname: "管理者",
				},
			},
			expect: "管理者",
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

func TestAdmin_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admin  *Admin
		expect *response.Admin
	}{
		{
			name: "success",
			admin: &Admin{
				Admin: response.Admin{
					ID:            "admin-id",
					Role:          entity.AdminRoleAdministrator,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin@and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
			expect: &response.Admin{
				ID:            "admin-id",
				Role:          entity.AdminRoleAdministrator,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどぴりおど",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin@and-period.jp",
				CreatedAt:     1640962800,
				UpdatedAt:     1640962800,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admin.Response())
		})
	}
}

func TestAdmins(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins entity.Admins
		expect Admins
	}{
		{
			name: "success",
			admins: entity.Admins{
				{
					ID:            "admin-id01",
					Role:          entity.AdminRoleAdministrator,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					ID:            "admin-id02",
					Role:          entity.AdminRoleCoordinator,
					Lastname:      "&.",
					Firstname:     "コーディネータ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "ちゅうかいしゃ",
					Email:         "test-admin02@and-period.jp",
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Admins{
				{
					Admin: response.Admin{
						ID:            "admin-id01",
						Role:          entity.AdminRoleAdministrator,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Admin: response.Admin{
						ID:            "admin-id02",
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "コーディネータ",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "ちゅうかいしゃ",
						Email:         "test-admin02@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAdmins(tt.admins))
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
					Admin: response.Admin{
						ID:            "admin-id01",
						Role:          entity.AdminRoleAdministrator,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Admin: response.Admin{
						ID:            "admin-id02",
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "コーディネータ",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "ちゅうかいしゃ",
						Email:         "test-admin02@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: map[string]*Admin{
				"admin-id01": {
					Admin: response.Admin{
						ID:            "admin-id01",
						Role:          entity.AdminRoleAdministrator,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				"admin-id02": {
					Admin: response.Admin{
						ID:            "admin-id02",
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "コーディネータ",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "ちゅうかいしゃ",
						Email:         "test-admin02@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
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

func TestAdmins_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins Admins
		expect []*response.Admin
	}{
		{
			name: "success",
			admins: Admins{
				{
					Admin: response.Admin{
						ID:            "admin-id01",
						Role:          entity.AdminRoleAdministrator,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Admin: response.Admin{
						ID:            "admin-id02",
						Role:          entity.AdminRoleCoordinator,
						Lastname:      "&.",
						Firstname:     "コーディネータ",
						LastnameKana:  "あんどぴりおど",
						FirstnameKana: "ちゅうかいしゃ",
						Email:         "test-admin02@and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: []*response.Admin{
				{
					ID:            "admin-id01",
					Role:          entity.AdminRoleAdministrator,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				{
					ID:            "admin-id02",
					Role:          entity.AdminRoleCoordinator,
					Lastname:      "&.",
					Firstname:     "コーディネータ",
					LastnameKana:  "あんどぴりおど",
					FirstnameKana: "ちゅうかいしゃ",
					Email:         "test-admin02@and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admins.Response())
		})
	}
}
