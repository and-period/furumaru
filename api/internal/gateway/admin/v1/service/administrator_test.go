package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestAdministrator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admin  *entity.Administrator
		expect *Administrator
	}{
		{
			name: "success",
			admin: &entity.Administrator{
				Admin: entity.Admin{
					ID:            "admin-id",
					Type:          entity.AdminTypeAdministrator,
					Status:        entity.AdminStatusActivated,
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin@and-period.jp",
				},
				AdminID:     "admin-id",
				PhoneNumber: "+819012345678",
				CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			expect: &Administrator{
				Administrator: response.Administrator{
					ID:            "admin-id",
					Status:        int32(AdminStatusActivated),
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAdministrator(tt.admin))
		})
	}
}

func TestAdministrator_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admin  *Administrator
		expect *response.Administrator
	}{
		{
			name: "success",
			admin: &Administrator{
				Administrator: response.Administrator{
					ID:            "admin-id",
					Status:        int32(AdminStatusActivated),
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
			expect: &response.Administrator{
				ID:            "admin-id",
				Status:        int32(AdminStatusActivated),
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin01@and-period.jp",
				PhoneNumber:   "+819012345678",
				CreatedAt:     1640962800,
				UpdatedAt:     1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admin.Response())
		})
	}
}

func TestAdministrators(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins entity.Administrators
		expect Administrators
	}{
		{
			name: "success",
			admins: entity.Administrators{
				{
					Admin: entity.Admin{
						ID:            "admin-id01",
						Type:          entity.AdminTypeAdministrator,
						Status:        entity.AdminStatusActivated,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
					},
					AdminID:     "admin-id01",
					PhoneNumber: "+819012345678",
					CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
				{
					Admin: entity.Admin{
						ID:            "admin-id02",
						Type:          entity.AdminTypeAdministrator,
						Status:        entity.AdminStatusActivated,
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin02@and-period.jp",
					},
					AdminID:     "admin-id02",
					PhoneNumber: "+819012345678",
					CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			expect: Administrators{
				{
					Administrator: response.Administrator{
						ID:            "admin-id01",
						Status:        int32(AdminStatusActivated),
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Administrator: response.Administrator{
						ID:            "admin-id02",
						Status:        int32(AdminStatusActivated),
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin02@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewAdministrators(tt.admins))
		})
	}
}

func TestAdministrators_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		admins Administrators
		expect []*response.Administrator
	}{
		{
			name: "success",
			admins: Administrators{
				{
					Administrator: response.Administrator{
						ID:            "admin-id01",
						Status:        int32(AdminStatusActivated),
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Administrator: response.Administrator{
						ID:            "admin-id02",
						Status:        int32(AdminStatusActivated),
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: []*response.Administrator{
				{
					ID:            "admin-id01",
					Status:        int32(AdminStatusActivated),
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
				{
					ID:            "admin-id02",
					Status:        int32(AdminStatusActivated),
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					PhoneNumber:   "+819012345678",
					CreatedAt:     1640962800,
					UpdatedAt:     1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.admins.Response())
		})
	}
}

func TestAdministrator_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		administrators Administrators
		expect         map[string]*Administrator
	}{
		{
			name: "success",
			administrators: Administrators{
				{
					Administrator: response.Administrator{
						ID:            "admin-id01",
						Status:        int32(AdminStatusActivated),
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				{
					Administrator: response.Administrator{
						ID:            "admin-id02",
						Status:        int32(AdminStatusActivated),
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin02@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
			expect: map[string]*Administrator{
				"admin-id01": {
					Administrator: response.Administrator{
						ID:            "admin-id01",
						Status:        int32(AdminStatusActivated),
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
				"admin-id02": {
					Administrator: response.Administrator{
						ID:            "admin-id02",
						Status:        int32(AdminStatusActivated),
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin02@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.administrators.Map())
		})
	}
}
