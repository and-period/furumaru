package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const adminPolicyTable = "admin_policies"

type adminPolicy struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdminPolicy(db *mysql.Client) database.AdminPolicy {
	return &adminPolicy{
		db:  db,
		now: jst.Now,
	}
}

func (p *adminPolicy) List(ctx context.Context, params *database.ListAdminPoliciesParams, fields ...string) (entity.AdminPolicies, error) {
	var policies entity.AdminPolicies

	stmt := p.db.Statement(ctx, p.db.DB, adminPolicyTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&policies).Error
	return policies, dbError(err)
}

func (p *adminPolicy) Count(ctx context.Context, params *database.ListAdminPoliciesParams) (int64, error) {
	total, err := p.db.Count(ctx, p.db.DB, &entity.AdminPolicy{}, nil)
	return total, dbError(err)
}

func (p *adminPolicy) MultiGet(ctx context.Context, policyIDs []string, fields ...string) (entity.AdminPolicies, error) {
	var policies entity.AdminPolicies

	stmt := p.db.Statement(ctx, p.db.DB, adminPolicyTable, fields...).
		Where("id IN (?)", policyIDs)

	err := stmt.Find(&policies).Error
	return policies, dbError(err)
}

func (p *adminPolicy) Get(ctx context.Context, policyID string, fields ...string) (*entity.AdminPolicy, error) {
	var policy *entity.AdminPolicy

	stmt := p.db.Statement(ctx, p.db.DB, adminPolicyTable, fields...).
		Where("id = ?", policyID)

	if err := stmt.First(&policy).Error; err != nil {
		return nil, dbError(err)
	}
	return policy, nil
}
