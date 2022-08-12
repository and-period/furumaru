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

const adminAuthTable = "admin_auths"

var adminAuthFields = []string{
	"admin_id", "cognito_id", "role", "created_at", "updated_at",
}

type adminAuth struct {
	db  *database.Client
	now func() time.Time
}

func NewAdminAuth(db *database.Client) AdminAuth {
	return &adminAuth{
		db:  db,
		now: jst.Now,
	}
}

func (a *adminAuth) MultiGet(
	ctx context.Context, adminIDs []string, fields ...string,
) (entity.AdminAuths, error) {
	var auth entity.AdminAuths
	if len(fields) == 0 {
		fields = adminAuthFields
	}

	err := a.db.DB.WithContext(ctx).
		Table(adminAuthTable).Select(fields).
		Where("admin_id IN (?)", adminIDs).
		Find(&auth).Error
	return auth, exception.InternalError(err)
}

func (a *adminAuth) GetByAdminID(
	ctx context.Context, adminID string, fields ...string,
) (*entity.AdminAuth, error) {
	auth, err := a.getByAdminID(ctx, a.db.DB, adminID, fields...)
	return auth, exception.InternalError(err)
}

func (a *adminAuth) GetByCognitoID(
	ctx context.Context, cognitoID string, fields ...string,
) (*entity.AdminAuth, error) {
	var auth *entity.AdminAuth
	if len(fields) == 0 {
		fields = adminAuthFields
	}

	stmt := a.db.DB.WithContext(ctx).
		Table(adminAuthTable).Select(fields).
		Where("cognito_id = ?", cognitoID)

	if err := stmt.First(&auth).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return auth, nil
}

func (a *adminAuth) UpdateDevice(ctx context.Context, adminID, device string) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := a.getByAdminID(ctx, tx, adminID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"device":     device,
			"updated_at": a.now(),
		}
		err := tx.WithContext(ctx).
			Table(adminAuthTable).
			Where("admin_id = ?", adminID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (a *adminAuth) getByAdminID(
	ctx context.Context, tx *gorm.DB, adminID string, fields ...string,
) (*entity.AdminAuth, error) {
	var auth *entity.AdminAuth
	if len(fields) == 0 {
		fields = adminAuthFields
	}

	err := tx.WithContext(ctx).
		Table(adminAuthTable).Select(fields).
		Where("admin_id = ?", adminID).
		First(&auth).Error
	return auth, err
}
