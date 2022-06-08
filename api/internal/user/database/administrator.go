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
	"id", "email", "phone_number",
	"lastname", "firstname", "lastname_kana", "firstname_kana",
	"created_at", "updated_at", "deleted_at",
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
	return administrators, nil
}

func (a *administrator) Get(
	ctx context.Context, administratorID string, fields ...string,
) (*entity.Administrator, error) {
	var administrator *entity.Administrator
	if len(fields) == 0 {
		fields = administratorFields
	}

	stmt := a.db.DB.WithContext(ctx).
		Table(administratorTable).Select(fields).
		Where("id = ?", administratorID)

	if err := stmt.First(&administrator).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return administrator, nil
}

func (a *administrator) Create(
	ctx context.Context, auth *entity.AdminAuth, administrator *entity.Administrator,
) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := a.now()
		auth.CreatedAt, auth.UpdatedAt = now, now
		administrator.CreatedAt, administrator.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(adminAuthTable).Create(&auth).Error
		if err != nil {
			return nil, err
		}
		err = tx.WithContext(ctx).Table(administratorTable).Create(&administrator).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (a *administrator) UpdateEmail(ctx context.Context, administratorID, email string) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var current *entity.Administrator
		err := tx.WithContext(ctx).
			Table(administratorTable).Select("id").
			Where("id = ?", administratorID).
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
			Table(administratorTable).
			Where("id = ?", administratorID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}
