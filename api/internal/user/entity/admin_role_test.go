package entity

import (
	"bytes"
	"encoding/csv"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminRole_Write(t *testing.T) {
	t.Parallel()

	type fields struct {
		policies AdminPolicies
		roles    AdminRolePolicies
		groups   AdminGroupRoles
	}
	tests := []struct {
		name   string
		fields fields
		expect string
	}{
		{
			name: "success",
			fields: fields{
				policies: AdminPolicies{
					{
						ID:          "spot_list",
						Name:        "スポット 一覧取得",
						Description: "スポット 一覧取得の権限",
						Priority:    1,
						Path:        "/v1/spots",
						Method:      "GET",
						Action:      AdminPolicyActionAllow,
					},
					{
						ID:          "spot_get",
						Name:        "スポット 詳細取得",
						Description: "スポット 詳細取得の権限",
						Priority:    2,
						Path:        "/v1/spots/*",
						Method:      "GET",
						Action:      AdminPolicyActionAllow,
					},
					{
						ID:          "spot_create",
						Name:        "スポット 登録",
						Description: "スポット 登録の権限",
						Priority:    3,
						Path:        "/v1/spots",
						Method:      "POST",
						Action:      AdminPolicyActionAllow,
					},
					{
						ID:          "spot_update",
						Name:        "スポット 更新",
						Description: "スポット 更新の権限",
						Priority:    4,
						Path:        "/v1/spots/*",
						Method:      "PATCH",
						Action:      AdminPolicyActionAllow,
					},
					{
						ID:          "spot_delete",
						Name:        "スポット 削除",
						Description: "スポット 削除の権限",
						Priority:    5,
						Path:        "/v1/spots/*",
						Method:      "DELETE",
						Action:      AdminPolicyActionAllow,
					},
				},
				roles: AdminRolePolicies{
					{
						RoleID:   "spot_editor",
						PolicyID: "spot_list",
					},
					{
						RoleID:   "spot_editor",
						PolicyID: "spot_get",
					},
					{
						RoleID:   "spot_editor",
						PolicyID: "spot_create",
					},
					{
						RoleID:   "spot_editor",
						PolicyID: "spot_update",
					},
					{
						RoleID:   "spot_viewer",
						PolicyID: "spot_list",
					},
					{
						RoleID:   "spot_viewer",
						PolicyID: "spot_get",
					},
				},
				groups: AdminGroupRoles{
					{
						GroupID: "administrator",
						RoleID:  "spot_editor",
					},
					{
						GroupID: "coordinator",
						RoleID:  "spot_viewer",
					},
				},
			},
			expect: "./testdata/admin-role.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			writer := csv.NewWriter(buf)

			expect, err := os.ReadFile(tt.expect)
			require.NoError(t, err)

			err = tt.fields.policies.Write(writer)
			require.NoError(t, err)

			err = tt.fields.roles.Write(writer)
			require.NoError(t, err)

			err = tt.fields.groups.Write(writer)
			require.NoError(t, err)

			writer.Flush()

			assert.Equal(t, string(expect), buf.String())
		})
	}
}
