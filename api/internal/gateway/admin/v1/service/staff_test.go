package service

import (
	"testing"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/marche/api/internal/store/entity"
	uentity "github.com/and-period/marche/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

func TestStaff(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		staff  *sentity.Staff
		admin  *uentity.Admin
		expect *Staff
	}{
		{
			name: "success",
			staff: &sentity.Staff{
				StoreID: 1,
				UserID:  "kSByoE6FetnPs5Byk3a9Zx",
				Role:    sentity.StoreRoleAdministrator,
			},
			admin: &uentity.Admin{
				ID:        "kSByoE6FetnPs5Byk3a9Zx",
				Lastname:  "&.",
				Firstname: "スタッフ1",
				Email:     "test-user01@and-period.jp",
			},
			expect: &Staff{
				Staff: &response.Staff{
					ID:    "kSByoE6FetnPs5Byk3a9Zx",
					Name:  "&. スタッフ1",
					Email: "test-user01@and-period.jp",
					Role:  1,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewStaff(tt.staff, tt.admin))
		})
	}
}

func TestStaff_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		staff  *Staff
		expect *response.Staff
	}{
		{
			name: "success",
			staff: &Staff{
				Staff: &response.Staff{
					ID:    "kSByoE6FetnPs5Byk3a9Zx",
					Name:  "&. スタッフ1",
					Email: "test-user01@and-period.jp",
					Role:  1,
				},
			},
			expect: &response.Staff{
				ID:    "kSByoE6FetnPs5Byk3a9Zx",
				Name:  "&. スタッフ1",
				Email: "test-user01@and-period.jp",
				Role:  1,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.staff.Response())
		})
	}
}

func TestStaffs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		staffs sentity.Staffs
		admins map[string]*uentity.Admin
		expect Staffs
	}{
		{
			name: "success",
			staffs: sentity.Staffs{
				{
					StoreID: 1,
					UserID:  "kSByoE6FetnPs5Byk3a9Zx",
					Role:    sentity.StoreRoleAdministrator,
				},
				{
					StoreID: 1,
					UserID:  "kSByoE6FetnPs5Byk3a9Za",
					Role:    sentity.StoreRoleEditor,
				},
			},
			admins: map[string]*uentity.Admin{
				"kSByoE6FetnPs5Byk3a9Zx": {
					ID:        "kSByoE6FetnPs5Byk3a9Zx",
					Lastname:  "&.",
					Firstname: "スタッフ1",
					Email:     "test-user01@and-period.jp",
				},
				"kSByoE6FetnPs5Byk3a9Za": {
					ID:        "kSByoE6FetnPs5Byk3a9Za",
					Lastname:  "&.",
					Firstname: "スタッフ2",
					Email:     "test-user02@and-period.jp",
				},
			},
			expect: Staffs{
				{
					Staff: &response.Staff{
						ID:    "kSByoE6FetnPs5Byk3a9Zx",
						Name:  "&. スタッフ1",
						Email: "test-user01@and-period.jp",
						Role:  1,
					},
				},
				{
					Staff: &response.Staff{
						ID:    "kSByoE6FetnPs5Byk3a9Za",
						Name:  "&. スタッフ2",
						Email: "test-user02@and-period.jp",
						Role:  2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewStaffs(tt.staffs, tt.admins))
		})
	}
}

func TestStaffs_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		staffs Staffs
		expect []*response.Staff
	}{
		{
			name: "success",
			staffs: Staffs{
				{
					Staff: &response.Staff{
						ID:    "kSByoE6FetnPs5Byk3a9Zx",
						Name:  "&. スタッフ1",
						Email: "test-user01@and-period.jp",
						Role:  1,
					},
				},
				{
					Staff: &response.Staff{
						ID:    "kSByoE6FetnPs5Byk3a9Za",
						Name:  "&. スタッフ2",
						Email: "test-user02@and-period.jp",
						Role:  2,
					},
				},
			},
			expect: []*response.Staff{
				{
					ID:    "kSByoE6FetnPs5Byk3a9Zx",
					Name:  "&. スタッフ1",
					Email: "test-user01@and-period.jp",
					Role:  1,
				},
				{
					ID:    "kSByoE6FetnPs5Byk3a9Za",
					Name:  "&. スタッフ2",
					Email: "test-user02@and-period.jp",
					Role:  2,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.staffs.Response())
		})
	}
}
