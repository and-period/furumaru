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

const adminGroupTable = "admin_groups"

type adminGroup struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdminGroup(db *mysql.Client) database.AdminGroup {
	return &adminGroup{
		db:  db,
		now: jst.Now,
	}
}

func (g *adminGroup) List(
	ctx context.Context,
	params *database.ListAdminGroupsParams,
	fields ...string,
) (entity.AdminGroups, error) {
	var groups entity.AdminGroups

	stmt := g.db.Statement(ctx, g.db.DB, adminGroupTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&groups).Error
	return groups, dbError(err)
}

func (g *adminGroup) Count(
	ctx context.Context,
	params *database.ListAdminGroupsParams,
) (int64, error) {
	total, err := g.db.Count(ctx, g.db.DB, &entity.AdminGroup{}, nil)
	return total, dbError(err)
}

func (g *adminGroup) MultiGet(
	ctx context.Context,
	groupIDs []string,
	fields ...string,
) (entity.AdminGroups, error) {
	var groups entity.AdminGroups

	stmt := g.db.Statement(ctx, g.db.DB, adminGroupTable, fields...).
		Where("id IN (?)", groupIDs)

	err := stmt.Find(&groups).Error
	return groups, dbError(err)
}

func (g *adminGroup) Get(
	ctx context.Context,
	groupID string,
	fields ...string,
) (*entity.AdminGroup, error) {
	var group *entity.AdminGroup

	stmt := g.db.Statement(ctx, g.db.DB, adminGroupTable, fields...).
		Where("id = ?", groupID)

	if err := stmt.First(&group).Error; err != nil {
		return nil, dbError(err)
	}
	return group, nil
}

func (g *adminGroup) Upsert(ctx context.Context, group *entity.AdminGroup) error {
	now := g.now()
	group.CreatedAt, group.UpdatedAt = now, now

	updates := map[string]interface{}{
		"name":             group.Name,
		"description":      group.Description,
		"updated_admin_id": group.UpdatedAdminID,
		"updated_at":       now,
		"deleted_at":       nil,
	}
	clauses := clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(updates),
	}
	err := g.db.DB.WithContext(ctx).Clauses(clauses).Create(&group).Error
	return dbError(err)
}

func (g *adminGroup) Delete(ctx context.Context, groupID string) error {
	stmt := g.db.DB.WithContext(ctx).
		Where("id = ?", groupID)
	err := stmt.Delete(&entity.AdminGroup{}).Error
	return dbError(err)
}
