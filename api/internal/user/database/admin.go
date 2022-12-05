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

const adminTable = "admins"

type admin struct {
	db  *database.Client
	now func() time.Time
}

func NewAdmin(db *database.Client) Admin {
	return &admin{
		db:  db,
		now: jst.Now,
	}
}

func (a *admin) MultiGet(
	ctx context.Context, adminIDs []string, fields ...string,
) (entity.Admins, error) {
	var admins entity.Admins

	err := a.db.Statement(ctx, a.db.DB, adminTable, fields...).
		Where("id IN (?)", adminIDs).
		Find(&admins).Error
	return admins, exception.InternalError(err)
}

func (a *admin) Get(
	ctx context.Context, adminID string, fields ...string,
) (*entity.Admin, error) {
	admin, err := a.get(ctx, a.db.DB, adminID, fields...)
	return admin, exception.InternalError(err)
}

func (a *admin) GetByCognitoID(
	ctx context.Context, cognitoID string, fields ...string,
) (*entity.Admin, error) {
	var admin *entity.Admin

	stmt := a.db.Statement(ctx, a.db.DB, adminTable, fields...).
		Where("cognito_id = ?", cognitoID)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return admin, nil
}

func (a *admin) UpdateEmail(ctx context.Context, adminID, email string) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := a.get(ctx, tx, adminID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"email":      email,
			"updated_at": a.now(),
		}
		err := tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", adminID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (a *admin) UpdateDevice(ctx context.Context, adminID, device string) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := a.get(ctx, tx, adminID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"device":     device,
			"updated_at": a.now(),
		}
		err := tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", adminID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (a *admin) get(ctx context.Context, tx *gorm.DB, adminID string, fields ...string) (*entity.Admin, error) {
	var admin *entity.Admin

	err := a.db.Statement(ctx, tx, adminTable, fields...).
		Where("id = ?", adminID).
		First(&admin).Error
	return admin, err
}
