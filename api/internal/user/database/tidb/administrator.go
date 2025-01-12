package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

const administratorTable = "administrators"

type administrator struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdministrator(db *mysql.Client) database.Administrator {
	return &administrator{
		db:  db,
		now: jst.Now,
	}
}

func (a *administrator) List(
	ctx context.Context, params *database.ListAdministratorsParams, fields ...string,
) (entity.Administrators, error) {
	var administrators entity.Administrators

	stmt := a.db.Statement(ctx, a.db.DB, administratorTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&administrators).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, administrators...); err != nil {
		return nil, dbError(err)
	}
	return administrators, nil
}

func (a *administrator) Count(ctx context.Context, _ *database.ListAdministratorsParams) (int64, error) {
	total, err := a.db.Count(ctx, a.db.DB, &entity.Administrator{}, nil)
	return total, dbError(err)
}

func (a *administrator) MultiGet(
	ctx context.Context, administratorIDs []string, fields ...string,
) (entity.Administrators, error) {
	var administrators entity.Administrators

	stmt := a.db.Statement(ctx, a.db.DB, administratorTable, fields...).
		Where("admin_id IN (?)", administratorIDs)

	if err := stmt.Find(&administrators).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, administrators...); err != nil {
		return nil, dbError(err)
	}
	return administrators, nil
}

func (a *administrator) Get(
	ctx context.Context, administratorID string, fields ...string,
) (*entity.Administrator, error) {
	administrator, err := a.get(ctx, a.db.DB, administratorID, fields...)
	if err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, administrator); err != nil {
		return nil, dbError(err)
	}
	return administrator, nil
}

func (a *administrator) Create(
	ctx context.Context, administrator *entity.Administrator, auth func(ctx context.Context) error,
) error {
	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		params := &entity.NewAdminGroupUsersParams{
			AdminID:  administrator.AdminID,
			GroupIDs: administrator.GroupIDs,
		}
		groups := entity.NewAdminGroupUsers(params)

		now := a.now()
		administrator.CreatedAt, administrator.UpdatedAt = now, now
		administrator.Admin.CreatedAt, administrator.Admin.UpdatedAt = now, now
		for _, group := range groups {
			group.CreatedAt, group.UpdatedAt = now, now
		}

		if err := tx.WithContext(ctx).Table(adminTable).Create(&administrator.Admin).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Table(administratorTable).Create(&administrator).Error; err != nil {
			return err
		}
		if len(groups) > 0 {
			if err := tx.WithContext(ctx).Table(adminGroupUserTable).Create(&groups).Error; err != nil {
				return err
			}
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (a *administrator) Update(ctx context.Context, administratorID string, params *database.UpdateAdministratorParams) error {
	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := a.now()
		adminParams := map[string]interface{}{
			"lastname":       params.Lastname,
			"firstname":      params.Firstname,
			"lastname_kana":  params.LastnameKana,
			"firstname_kana": params.FirstnameKana,
			"updated_at":     now,
		}
		administratorParams := map[string]interface{}{
			"phone_number": params.PhoneNumber,
			"updated_at":   now,
		}

		err := tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", administratorID).
			Updates(adminParams).Error
		if err != nil {
			return err
		}
		err = tx.WithContext(ctx).
			Table(administratorTable).
			Where("admin_id = ?", administratorID).
			Updates(administratorParams).Error
		return err
	})
	return dbError(err)
}

func (a *administrator) Delete(
	ctx context.Context, administratorID string, auth func(ctx context.Context) error,
) error {
	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := a.now()
		updates := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		stmt := tx.WithContext(ctx).Table(adminTable).Where("id = ?", administratorID)
		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (a *administrator) get(
	ctx context.Context, tx *gorm.DB, administratorID string, fields ...string,
) (*entity.Administrator, error) {
	var administrator *entity.Administrator

	err := a.db.Statement(ctx, tx, administratorTable, fields...).
		Where("admin_id = ?", administratorID).
		First(&administrator).Error
	return administrator, err
}

func (a *administrator) fill(ctx context.Context, tx *gorm.DB, administrators ...*entity.Administrator) error {
	var (
		admins entity.Admins
		groups entity.AdminGroupUsers
	)

	ids := entity.Administrators(administrators).IDs()
	if len(ids) == 0 {
		return nil
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		stmt := a.db.Statement(ectx, tx, adminTable).
			Unscoped().
			Where("id IN (?)", ids)
		return stmt.Find(&admins).Error
	})
	eg.Go(func() error {
		stmt := a.db.Statement(ectx, tx, adminGroupUserTable).
			Where("admin_id IN (?)", ids).
			Where("expired_at IS NULL OR expired_at > ?", jst.Now())
		return stmt.Find(&groups).Error
	})
	if err := eg.Wait(); err != nil {
		return err
	}

	entity.Administrators(administrators).Fill(admins.Map(), groups.GroupByAdminID())
	return nil
}
