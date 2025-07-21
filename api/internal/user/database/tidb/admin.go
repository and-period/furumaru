package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const adminTable = "admins"

type admin struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdmin(db *mysql.Client) database.Admin {
	return &admin{
		db:  db,
		now: jst.Now,
	}
}

func (a *admin) MultiGet(
	ctx context.Context, adminIDs []string, fields ...string,
) (entity.Admins, error) {
	var admins entity.Admins

	stmt := a.db.Statement(ctx, a.db.DB, adminTable, fields...).Where("id IN (?)", adminIDs)

	if err := stmt.Find(&admins).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, admins...); err != nil {
		return nil, dbError(err)
	}
	return admins, nil
}

func (a *admin) Get(
	ctx context.Context, adminID string, fields ...string,
) (*entity.Admin, error) {
	admin, err := a.get(ctx, a.db.DB, adminID, fields...)
	return admin, dbError(err)
}

func (a *admin) GetByCognitoID(
	ctx context.Context, cognitoID string, fields ...string,
) (*entity.Admin, error) {
	var admin *entity.Admin

	stmt := a.db.Statement(ctx, a.db.DB, adminTable, fields...).
		Where("cognito_id = ?", cognitoID)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, admin); err != nil {
		return nil, dbError(err)
	}
	return admin, nil
}

func (a *admin) GetByEmail(
	ctx context.Context,
	email string,
	fields ...string,
) (*entity.Admin, error) {
	var admin *entity.Admin

	stmt := a.db.Statement(ctx, a.db.DB, adminTable, fields...).
		Where("email = ?", email)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, admin); err != nil {
		return nil, dbError(err)
	}
	return admin, nil
}

func (a *admin) UpdateEmail(ctx context.Context, adminID, email string) error {
	params := map[string]interface{}{
		"email":      email,
		"updated_at": a.now(),
	}
	stmt := a.db.DB.WithContext(ctx).
		Table(adminTable).
		Where("id = ?", adminID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (a *admin) UpdateDevice(ctx context.Context, adminID, device string) error {
	params := map[string]interface{}{
		"device":     device,
		"updated_at": a.now(),
	}
	stmt := a.db.DB.WithContext(ctx).
		Table(adminTable).
		Where("id = ?", adminID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (a *admin) UpdateSignInAt(ctx context.Context, adminID string) error {
	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		admin, err := a.get(ctx, tx, adminID)
		if err != nil {
			return err
		}

		now := a.now()
		params := map[string]interface{}{
			"last_sign_in_at": now,
			"updated_at":      now,
		}
		if admin.FirstSignInAt.IsZero() {
			params["first_sign_in_at"] = now
		}
		err = tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", adminID).
			Updates(params).Error
		return err
	})
	return dbError(err)
}

func (a *admin) get(
	ctx context.Context,
	tx *gorm.DB,
	adminID string,
	fields ...string,
) (*entity.Admin, error) {
	var admin *entity.Admin

	stmt := a.db.Statement(ctx, tx, adminTable, fields...).
		Where("id = ?", adminID)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, err
	}
	if err := a.fill(ctx, tx, admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (a *admin) fill(ctx context.Context, tx *gorm.DB, admins ...*entity.Admin) error {
	var groups entity.AdminGroupUsers

	ids := entity.Admins(admins).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := a.db.Statement(ctx, tx, adminGroupUserTable).
		Where("admin_id IN (?)", ids).
		Where("expired_at IS NULL OR expired_at > ?", jst.Now())

	if err := stmt.Find(&groups).Error; err != nil {
		return err
	}

	entity.Admins(admins).Fill(groups.GroupByAdminID())
	return nil
}
