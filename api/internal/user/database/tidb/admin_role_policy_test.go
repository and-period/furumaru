package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminRolePolicy_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	role := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&role).Error
	require.NoError(t, err)
	policies := make(entity.AdminPolicies, 3)
	policies[0] = testAdminPolicy("policy-id01", now())
	policies[1] = testAdminPolicy("policy-id02", now())
	policies[2] = testAdminPolicy("policy-id03", now())
	err = db.DB.WithContext(ctx).Table(adminPolicyTable).Create(&policies).Error
	require.NoError(t, err)

	rolePolicies := make(entity.AdminRolePolicies, 3)
	rolePolicies[0] = testAdminRolePolicy("role-id", policies[0].ID, now())
	rolePolicies[1] = testAdminRolePolicy("role-id", policies[1].ID, now())
	rolePolicies[2] = testAdminRolePolicy("role-id", policies[2].ID, now())
	err = db.DB.WithContext(ctx).Table(adminRolePolicyTable).Create(&rolePolicies).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminRolePoliciesParams
	}
	type want struct {
		rolePolicies entity.AdminRolePolicies
		err          error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListAdminRolePoliciesParams{},
			},
			want: want{
				rolePolicies: rolePolicies,
				err:          nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminRolePolicy{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.Equal(t, tt.want.rolePolicies, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminRolePolicy_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	role := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&role).Error
	require.NoError(t, err)
	policies := make(entity.AdminPolicies, 3)
	policies[0] = testAdminPolicy("policy-id01", now())
	policies[1] = testAdminPolicy("policy-id02", now())
	policies[2] = testAdminPolicy("policy-id03", now())
	err = db.DB.WithContext(ctx).Table(adminPolicyTable).Create(&policies).Error
	require.NoError(t, err)

	rolePolicies := make(entity.AdminRolePolicies, 3)
	rolePolicies[0] = testAdminRolePolicy("role-id", policies[0].ID, now())
	rolePolicies[1] = testAdminRolePolicy("role-id", policies[1].ID, now())
	rolePolicies[2] = testAdminRolePolicy("role-id", policies[2].ID, now())
	err = db.DB.WithContext(ctx).Table(adminRolePolicyTable).Create(&rolePolicies).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminRolePoliciesParams
	}
	type want struct {
		total int64
		err   error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListAdminRolePoliciesParams{},
			},
			want: want{
				total: 3,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminRolePolicy{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.Equal(t, tt.want.total, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestAdminRolePolicy_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	role := testAdminRole("role-id", now())
	err = db.DB.WithContext(ctx).Table(adminRoleTable).Create(&role).Error
	require.NoError(t, err)
	policy := testAdminPolicy("policy-id", now())
	err = db.DB.WithContext(ctx).Table(adminPolicyTable).Create(&policy).Error
	require.NoError(t, err)

	rp := testAdminRolePolicy("role-id", "policy-id", now())
	err = db.DB.WithContext(ctx).Table(adminRolePolicyTable).Create(&rp).Error
	require.NoError(t, err)

	type args struct {
		roleID   string
		policyID string
	}
	type want struct {
		rolePolicy *entity.AdminRolePolicy
		err        error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				roleID:   "role-id",
				policyID: "policy-id",
			},
			want: want{
				rolePolicy: rp,
				err:        nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				roleID:   "",
				policyID: "",
			},
			want: want{
				rolePolicy: nil,
				err:        database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminRolePolicy{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.roleID, tt.args.policyID)
			assert.Equal(t, tt.want.rolePolicy, actual)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testAdminRolePolicy(roleID, policyID string, now time.Time) *entity.AdminRolePolicy {
	return &entity.AdminRolePolicy{
		RoleID:    roleID,
		PolicyID:  policyID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
