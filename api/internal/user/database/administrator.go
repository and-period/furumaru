package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const administratorTable = "administrators"

var administratorFields = []string{
	"admin_id", "phone_number", "created_at", "updated_at", "deleted_at",
}

type administrator struct {
	db  *database.Client
	now func() time.Time
}

func NewAdministrator(db *database.Client) Administrator {
	return &administrator{
		db:  db,
		now: jst.Now,
	}
}

func (a *administrator) List(
	ctx context.Context, params *ListAdministratorsParams, fields ...string,
) (entity.Administrators, error) {
	var administrators entity.Administrators
	if len(fields) == 0 {
		fields = administratorFields
	}

	stmt := a.db.DB.WithContext(ctx).Table(administratorTable).Select(fields)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&administrators).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := a.fill(ctx, a.db.DB, administrators...); err != nil {
		return nil, exception.InternalError(err)
	}
	return administrators, nil
}

func (a *administrator) Count(ctx context.Context, params *ListAdministratorsParams) (int64, error) {
	var total int64

	stmt := a.db.DB.WithContext(ctx).Table(administratorTable).Select("COUNT(*)")

	err := stmt.Count(&total).Error
	return total, exception.InternalError(err)
}

func (a *administrator) MultiGet(
	ctx context.Context, administratorIDs []string, fields ...string,
) (entity.Administrators, error) {
	var administrators entity.Administrators
	if len(fields) == 0 {
		fields = administratorFields
	}

	stmt := a.db.DB.WithContext(ctx).
		Table(administratorTable).Select(fields).
		Where("admin_id IN (?)", administratorIDs)

	if err := stmt.Find(&administrators).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := a.fill(ctx, a.db.DB, administrators...); err != nil {
		return nil, exception.InternalError(err)
	}
	return administrators, nil
}

func (a *administrator) Get(
	ctx context.Context, administratorID string, fields ...string,
) (*entity.Administrator, error) {
	administrator, err := a.get(ctx, a.db.DB, administratorID, fields...)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	if err := a.fill(ctx, a.db.DB, administrator); err != nil {
		return nil, exception.InternalError(err)
	}
	return administrator, nil
}

func (a *administrator) Create(
	ctx context.Context, administrator *entity.Administrator, auth func(ctx context.Context) error,
) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := a.now()
		administrator.Admin.CreatedAt, administrator.Admin.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(adminTable).Create(&administrator.Admin).Error; err != nil {
			return nil, err
		}
		administrator.CreatedAt, administrator.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(administratorTable).Create(&administrator).Error; err != nil {
			return nil, err
		}
		return nil, auth(ctx)
	})
	return exception.InternalError(err)
}

func (a *administrator) Update(ctx context.Context, administratorID string, params *UpdateAdministratorParams) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := a.get(ctx, tx, administratorID); err != nil {
			return nil, err
		}

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
			return nil, err
		}
		err = tx.WithContext(ctx).
			Table(administratorTable).
			Where("admin_id = ?", administratorID).
			Updates(administratorParams).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (a *administrator) Delete(
	ctx context.Context, administratorID string, auth func(ctx context.Context) error,
) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := a.get(ctx, tx, administratorID); err != nil {
			return nil, err
		}

		now := a.now()
		administratorParams := map[string]interface{}{
			"updated_at": now,
			"deleted_at": now,
		}
		err := tx.WithContext(ctx).
			Table(administratorTable).
			Where("admin_id = ?", administratorID).
			Updates(administratorParams).Error
		if err != nil {
			return nil, err
		}
		adminParams := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		err = tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", administratorID).
			Updates(adminParams).Error
		if err != nil {
			return nil, err
		}
		return nil, auth(ctx)
	})
	return exception.InternalError(err)
}

func (a *administrator) get(
	ctx context.Context, tx *gorm.DB, administratorID string, fields ...string,
) (*entity.Administrator, error) {
	var administrator *entity.Administrator
	if len(fields) == 0 {
		fields = administratorFields
	}

	err := tx.WithContext(ctx).
		Table(administratorTable).Select(fields).
		Where("admin_id = ?", administratorID).
		First(&administrator).Error
	return administrator, err
}

func (a *administrator) fill(ctx context.Context, tx *gorm.DB, administrators ...*entity.Administrator) error {
	var admins entity.Admins

	ids := entity.Administrators(administrators).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := tx.WithContext(ctx).
		Table(adminTable).Select(adminFields).
		Where("id IN (?)", ids)
	if err := stmt.Find(&admins).Error; err != nil {
		return err
	}

	adminMap := admins.Map()

	for i, a := range administrators {
		admin, ok := adminMap[a.AdminID]
		if !ok {
			admin = &entity.Admin{}
		}

		administrators[i].Fill(admin)
	}
	return nil
}
