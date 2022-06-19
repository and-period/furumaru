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

const coordinatorTable = "coordinators"

var coordinatorFields = []string{
	"id", "email", "phone_number",
	"lastname", "firstname", "lastname_kana", "firstname_kana",
	"company_name", "store_name", "thumbnail_url", "header_url",
	"twitter_account", "instagram_account", "facebook_account",
	"postal_code", "prefecture", "city", "address_line1", "address_line2",
	"created_at", "updated_at", "deleted_at",
}

type coordinator struct {
	db  *database.Client
	now func() time.Time
}

func NewCoordinator(db *database.Client) Coordinator {
	return &coordinator{
		db:  db,
		now: jst.Now,
	}
}

func (a *coordinator) List(
	ctx context.Context, params *ListCoordinatorsParams, fields ...string,
) (entity.Coordinators, error) {
	var coordinators entity.Coordinators
	if len(fields) == 0 {
		fields = coordinatorFields
	}

	stmt := a.db.DB.WithContext(ctx).Table(coordinatorTable).Select(fields)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&coordinators).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return coordinators, nil
}

func (a *coordinator) Get(
	ctx context.Context, coordinatorID string, fields ...string,
) (*entity.Coordinator, error) {
	var coordinator *entity.Coordinator
	if len(fields) == 0 {
		fields = coordinatorFields
	}

	stmt := a.db.DB.WithContext(ctx).
		Table(coordinatorTable).Select(fields).
		Where("id = ?", coordinatorID)

	if err := stmt.First(&coordinator).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return coordinator, nil
}

func (a *coordinator) Create(
	ctx context.Context, auth *entity.AdminAuth, coordinator *entity.Coordinator,
) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := a.now()
		auth.CreatedAt, auth.UpdatedAt = now, now
		coordinator.CreatedAt, coordinator.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(adminAuthTable).Create(&auth).Error
		if err != nil {
			return nil, err
		}
		err = tx.WithContext(ctx).Table(coordinatorTable).Create(&coordinator).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (a *coordinator) UpdateEmail(ctx context.Context, coordinatorID, email string) error {
	_, err := a.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var current *entity.Coordinator
		err := tx.WithContext(ctx).
			Table(coordinatorTable).Select("id").
			Where("id = ?", coordinatorID).
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
			Table(coordinatorTable).
			Where("id = ?", coordinatorID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}
