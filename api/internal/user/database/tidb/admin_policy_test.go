package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminPolicy_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	policies := make(entity.AdminPolicies, 3)
	policies[0] = testAdminPolicy("policy-id01", now().Add(-2*time.Hour))
	policies[1] = testAdminPolicy("policy-id02", now().Add(-time.Hour))
	policies[2] = testAdminPolicy("policy-id03", now())
	err = db.DB.WithContext(ctx).Table(adminPolicyTable).Create(&policies).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminPoliciesParams
	}
	type want struct {
		policies entity.AdminPolicies
		err      error
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
				params: &database.ListAdminPoliciesParams{},
			},
			want: want{
				policies: policies,
				err:      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminPolicy{db: db, now: now}
			actual, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.policies, actual)
		})
	}
}

func TestAdminPolicy_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	policies := make(entity.AdminPolicies, 3)
	policies[0] = testAdminPolicy("policy-id01", now().Add(-2*time.Hour))
	policies[1] = testAdminPolicy("policy-id02", now().Add(-time.Hour))
	policies[2] = testAdminPolicy("policy-id03", now())
	err = db.DB.WithContext(ctx).Table(adminPolicyTable).Create(&policies).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListAdminPoliciesParams
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
				params: &database.ListAdminPoliciesParams{},
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

			db := &adminPolicy{db: db, now: now}
			actual, err := db.Count(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.total, actual)
		})
	}
}

func TestAdminPolicy_MultiGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	policies := make(entity.AdminPolicies, 3)
	policies[0] = testAdminPolicy("policy-id01", now().Add(-2*time.Hour))
	policies[1] = testAdminPolicy("policy-id02", now().Add(-time.Hour))
	policies[2] = testAdminPolicy("policy-id03", now())
	err = db.DB.WithContext(ctx).Table(adminPolicyTable).Create(&policies).Error
	require.NoError(t, err)

	type args struct {
		policyIDs []string
	}
	type want struct {
		policies entity.AdminPolicies
		err      error
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
				policyIDs: []string{"policy-id01", "policy-id02", "policy-id03"},
			},
			want: want{
				policies: policies,
				err:      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminPolicy{db: db, now: now}
			actual, err := db.MultiGet(ctx, tt.args.policyIDs)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.policies, actual)
		})
	}
}

func TestAdminPolicy_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	ctx := t.Context()
	err := deleteAll(ctx)
	require.NoError(t, err)

	p := testAdminPolicy("policy-id", now())
	err = db.DB.WithContext(ctx).Table(adminPolicyTable).Create(&p).Error
	require.NoError(t, err)

	type args struct {
		policyID string
	}
	type want struct {
		policy *entity.AdminPolicy
		err    error
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
				policyID: "policy-id",
			},
			want: want{
				policy: p,
				err:    nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				policyID: "",
			},
			want: want{
				policy: nil,
				err:    database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &adminPolicy{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.policyID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.policy, actual)
		})
	}
}

func testAdminPolicy(policyID string, now time.Time) *entity.AdminPolicy {
	return &entity.AdminPolicy{
		ID:          policyID,
		Name:        "スポット一覧取得",
		Description: "スポット一覧取得の許可",
		Priority:    1,
		Path:        "/v1/spots",
		Method:      "GET",
		Action:      entity.AdminPolicyActionAllow,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
