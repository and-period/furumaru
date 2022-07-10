package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
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

func (a *adminAuth) GetByAdminID(
	ctx context.Context, adminID string, fields ...string,
) (*entity.AdminAuth, error) {
	var auth *entity.AdminAuth
	if len(fields) == 0 {
		fields = adminAuthFields
	}

	stmt := a.db.DB.WithContext(ctx).
		Table(adminAuthTable).Select(fields).
		Where("admin_id = ?", adminID)

	if err := stmt.First(&auth).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return auth, nil
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
