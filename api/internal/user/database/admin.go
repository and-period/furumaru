package database

import (
	"context"
	"time"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/pkg/jst"
	"gorm.io/gorm"
)

const adminTable = "admins"

var adminFields = []string{
	"id", "cognito_id", "email", "role", "thumbnail_url",
	"lastname", "firstname", "lastname_kana", "firstname_kana",
	"created_at", "updated_at", "deleted_at",
}

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

func (a *admin) List(ctx context.Context, params *ListAdminsParams, fields ...string) (entity.Admins, error) {
	var admins entity.Admins
	if len(fields) == 0 {
		fields = adminFields
	}

	stmt := a.db.DB.WithContext(ctx).Table(adminTable).Select(fields)
	if len(params.Roles) > 0 {
		stmt = stmt.Where("role IN (?)", params.Roles)
	}
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&admins).Error
	return admins, dbError(err)
}

func (a *admin) MultiGet(ctx context.Context, adminIDs []string, fields ...string) (entity.Admins, error) {
	var admins entity.Admins
	if len(fields) == 0 {
		fields = adminFields
	}

	stmt := a.db.DB.WithContext(ctx).
		Table(adminTable).Select(fields).
		Where("id IN (?)", adminIDs)

	err := stmt.Find(&admins).Error
	return admins, dbError(err)
}

func (a *admin) Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error) {
	var admin *entity.Admin
	if len(fields) == 0 {
		fields = adminFields
	}

	stmt := a.db.DB.WithContext(ctx).
		Table(adminTable).Select(fields).
		Where("id = ?", adminID)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, dbError(err)
	}
	return admin, nil
}

func (a *admin) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error) {
	var admin *entity.Admin
	if len(fields) == 0 {
		fields = adminFields
	}

	stmt := a.db.DB.WithContext(ctx).
		Table(adminTable).Select(fields).
		Where("cognito_id = ?", cognitoID)

	if err := stmt.First(&admin).Error; err != nil {
		return nil, dbError(err)
	}
	return admin, nil
}

func (a *admin) Create(ctx context.Context, admin *entity.Admin) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := a.now()
		admin.CreatedAt, admin.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(adminTable).Create(&admin).Error
		return nil, err
	})
	return dbError(err)
}

func (a *admin) UpdateEmail(ctx context.Context, adminID, email string) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var current *entity.Admin
		err := tx.WithContext(ctx).
			Table(adminTable).Select("id").
			Where("id = ?", adminID).
			First(&current).Error
		if err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"id":         current.ID,
			"email":      email,
			"updated_at": a.now(),
		}
		err = tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", current.ID).
			Updates(params).Error
		return nil, err
	})
	return dbError(err)
}
