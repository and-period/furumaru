package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const adminRoleTable = "admin_roles"

type adminRole struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdminRole(db *mysql.Client) database.AdminRole {
	return &adminRole{
		db:  db,
		now: jst.Now,
	}
}

func (r *adminRole) List(
	ctx context.Context,
	params *database.ListAdminRolesParams,
	fields ...string,
) (entity.AdminRoles, error) {
	var roles entity.AdminRoles

	stmt := r.db.Statement(ctx, r.db.DB, adminRoleTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&roles).Error
	return roles, dbError(err)
}

func (r *adminRole) Count(
	ctx context.Context,
	params *database.ListAdminRolesParams,
) (int64, error) {
	total, err := r.db.Count(ctx, r.db.DB, &entity.AdminRole{}, nil)
	return total, dbError(err)
}

func (r *adminRole) MultiGet(
	ctx context.Context,
	roleIDs []string,
	fields ...string,
) (entity.AdminRoles, error) {
	var roles entity.AdminRoles

	stmt := r.db.Statement(ctx, r.db.DB, adminRoleTable, fields...).
		Where("id IN (?)", roleIDs)

	err := stmt.Find(&roles).Error
	return roles, dbError(err)
}

func (r *adminRole) Get(
	ctx context.Context,
	roleID string,
	fields ...string,
) (*entity.AdminRole, error) {
	var role *entity.AdminRole

	stmt := r.db.Statement(ctx, r.db.DB, adminRoleTable, fields...).
		Where("id = ?", roleID)

	if err := stmt.First(&role).Error; err != nil {
		return nil, dbError(err)
	}
	return role, nil
}
