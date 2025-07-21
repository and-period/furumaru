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

const adminGroupUserTable = "admin_group_users"

type adminGroupUser struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdminGroupUser(db *mysql.Client) database.AdminGroupUser {
	return &adminGroupUser{
		db:  db,
		now: jst.Now,
	}
}

func (u *adminGroupUser) List(
	ctx context.Context, params *database.ListAdminGroupUsersParams, fields ...string,
) (entity.AdminGroupUsers, error) {
	var users entity.AdminGroupUsers

	stmt := u.db.Statement(ctx, u.db.DB, adminGroupUserTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&users).Error
	return users, dbError(err)
}

func (u *adminGroupUser) Count(
	ctx context.Context,
	params *database.ListAdminGroupUsersParams,
) (int64, error) {
	total, err := u.db.Count(ctx, u.db.DB, &entity.AdminGroupUser{}, nil)
	return total, dbError(err)
}

func (u *adminGroupUser) Get(
	ctx context.Context,
	groupID, adminID string,
	fields ...string,
) (*entity.AdminGroupUser, error) {
	var user *entity.AdminGroupUser

	stmt := u.db.Statement(ctx, u.db.DB, adminGroupUserTable, fields...).
		Where("group_id = ?", groupID).
		Where("admin_id = ?", adminID)

	if err := stmt.First(&user).Error; err != nil {
		return nil, dbError(err)
	}
	return user, nil
}

func (u *adminGroupUser) Upsert(ctx context.Context, user *entity.AdminGroupUser) error {
	now := u.now()
	user.CreatedAt, user.UpdatedAt = now, now

	updates := map[string]interface{}{
		"updated_admin_id": user.UpdatedAdminID,
		"updated_at":       now,
		"expired_at":       nil,
	}
	if !user.ExpiredAt.IsZero() {
		updates["expired_at"] = user.ExpiredAt
	}
	claues := clause.OnConflict{
		Columns:   []clause.Column{{Name: "group_id"}, {Name: "admin_id"}},
		DoUpdates: clause.Assignments(updates),
	}
	err := u.db.DB.WithContext(ctx).Clauses(claues).Create(&user).Error
	return dbError(err)
}

func (u *adminGroupUser) Delete(ctx context.Context, groupID, adminID string) error {
	stmt := u.db.DB.WithContext(ctx).
		Where("group_id = ?", groupID).
		Where("admin_id = ?", adminID)
	err := stmt.Delete(&entity.AdminGroupUser{}).Error
	return dbError(err)
}
