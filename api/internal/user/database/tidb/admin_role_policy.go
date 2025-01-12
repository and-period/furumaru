package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const adminRolePolicyTable = "admin_role_policies"

type adminRolePolicy struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdminRolePolicy(db *mysql.Client) database.AdminRolePolicy {
	return &adminRolePolicy{
		db:  db,
		now: jst.Now,
	}
}

func (p *adminRolePolicy) List(
	ctx context.Context, params *database.ListAdminRolePoliciesParams, fields ...string,
) (entity.AdminRolePolicies, error) {
	var policies entity.AdminRolePolicies

	stmt := p.db.Statement(ctx, p.db.DB, adminRolePolicyTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&policies).Error
	return policies, dbError(err)
}

func (p *adminRolePolicy) Count(ctx context.Context, params *database.ListAdminRolePoliciesParams) (int64, error) {
	total, err := p.db.Count(ctx, p.db.DB, &entity.AdminRolePolicy{}, nil)
	return total, dbError(err)
}

func (p *adminRolePolicy) Get(ctx context.Context, roleID, policyID string, fields ...string) (*entity.AdminRolePolicy, error) {
	var policy *entity.AdminRolePolicy

	stmt := p.db.Statement(ctx, p.db.DB, adminRolePolicyTable, fields...).
		Where("role_id = ?", roleID).
		Where("policy_id = ?", policyID)

	if err := stmt.First(&policy).Error; err != nil {
		return nil, dbError(err)
	}
	return policy, nil
}
