package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm/clause"
)

const adminGroupRoleTable = "admin_group_roles"

type adminGroupRole struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdminGroupRole(db *mysql.Client) database.AdminGroupRole {
	return &adminGroupRole{
		db:  db,
		now: jst.Now,
	}
}

func (r *adminGroupRole) List(
	ctx context.Context, params *database.ListAdminGroupRolesParams, fields ...string,
) (entity.AdminGroupRoles, error) {
	var roles entity.AdminGroupRoles

	stmt := r.db.Statement(ctx, r.db.DB, adminGroupRoleTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&roles).Error
	return roles, dbError(err)
}

func (r *adminGroupRole) Count(ctx context.Context, params *database.ListAdminGroupRolesParams) (int64, error) {
	total, err := r.db.Count(ctx, r.db.DB, &entity.AdminGroupRole{}, nil)
	return total, dbError(err)
}

func (r *adminGroupRole) Get(ctx context.Context, groupID, roleID string, fields ...string) (*entity.AdminGroupRole, error) {
	var role *entity.AdminGroupRole

	stmt := r.db.Statement(ctx, r.db.DB, adminGroupRoleTable, fields...).
		Where("group_id = ?", groupID).
		Where("role_id = ?", roleID)

	if err := stmt.First(&role).Error; err != nil {
		return nil, dbError(err)
	}
	return role, nil
}

func (r *adminGroupRole) Upsert(ctx context.Context, role *entity.AdminGroupRole) error {
	now := r.now()
	role.CreatedAt, role.UpdatedAt = now, now

	updates := map[string]interface{}{
		"updated_admin_id": role.UpdatedAdminID,
		"updated_at":       now,
		"deleted_at":       nil,
	}
	clauses := clause.OnConflict{
		Columns:   []clause.Column{{Name: "group_id"}, {Name: "role_id"}},
		DoUpdates: clause.Assignments(updates),
	}
	err := r.db.DB.WithContext(ctx).Clauses(clauses).Create(&role).Error
	return dbError(err)
}

func (r *adminGroupRole) Delete(ctx context.Context, groupID, roleID string) error {
	now := r.now()

	updates := map[string]interface{}{
		"updated_at": now,
		"deleted_at": now,
	}
	stmt := r.db.DB.WithContext(ctx).
		Table(adminGroupRoleTable).
		Where("group_id = ?", groupID).
		Where("role_id = ?", roleID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}
