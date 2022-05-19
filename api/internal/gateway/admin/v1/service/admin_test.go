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
		name           string
		role           entity.AdminRole
		expect         AdminRole
		expectString   string
		expectResponse int32
	}{
		{
			name:           "administrator",
			role:           entity.AdminRoleAdministrator,
			expect:         AdminRoleAdministrator,
			expectString:   "admin",
			expectResponse: 1,
		},
		{
			name:           "producer",
			role:           entity.AdminRoleProducer,
			expect:         AdminRoleProducer,
			expectString:   "producer",
			expectResponse: 2,
		},
		{
			name:           "unknown",
			role:           entity.AdminRoleUnknown,
			expect:         AdminRoleUnknown,
			expectString:   "unknown",
			expectResponse: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminRole(tt.role)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, tt.expectString, actual.String())
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
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin01@and-period.jp",
				Role:          entity.AdminRoleAdministrator,
				ThumbnailURL:  "https://and-period.jp",
				CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Admin{
				Admin: &response.Admin{
					ID:            "admin-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					Role:          int32(AdminRoleAdministrator),
					ThumbnailURL:  "https://and-period.jp",
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
				Admin: &response.Admin{
					ID:            "admin-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					Role:          int32(AdminRoleAdministrator),
					ThumbnailURL:  "https://and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
			expect: &response.Admin{
				ID:            "admin-id",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin01@and-period.jp",
				Role:          int32(AdminRoleAdministrator),
				ThumbnailURL:  "https://and-period.jp",
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
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					Role:          entity.AdminRoleAdministrator,
					ThumbnailURL:  "https://and-period.jp",
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					ID:            "admin-id02",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "すたっふ",
					Email:         "test-admin02@and-period.jp",
					Role:          entity.AdminRoleProducer,
					ThumbnailURL:  "https://and-period.jp",
					CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Admins{
				{
					Admin: &response.Admin{
						ID:            "admin-id01",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						Role:          int32(AdminRoleAdministrator),
						ThumbnailURL:  "https://and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Admin: &response.Admin{
						ID:            "admin-id02",
						Lastname:      "&.",
						Firstname:     "スタッフ",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "すたっふ",
						Email:         "test-admin02@and-period.jp",
						Role:          int32(AdminRoleProducer),
						ThumbnailURL:  "https://and-period.jp",
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
					Admin: &response.Admin{
						ID:            "admin-id01",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						Role:          int32(AdminRoleAdministrator),
						ThumbnailURL:  "https://and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Admin: &response.Admin{
						ID:            "admin-id02",
						Lastname:      "&.",
						Firstname:     "スタッフ",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "すたっふ",
						Email:         "test-admin02@and-period.jp",
						Role:          int32(AdminRoleProducer),
						ThumbnailURL:  "https://and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: []*response.Admin{
				{
					ID:            "admin-id01",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					Role:          int32(AdminRoleAdministrator),
					ThumbnailURL:  "https://and-period.jp",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				{
					ID:            "admin-id02",
					Lastname:      "&.",
					Firstname:     "スタッフ",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "すたっふ",
					Email:         "test-admin02@and-period.jp",
					Role:          int32(AdminRoleProducer),
					ThumbnailURL:  "https://and-period.jp",
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
