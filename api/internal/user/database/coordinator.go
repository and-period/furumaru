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

func (c *coordinator) List(
	ctx context.Context, params *ListCoordinatorsParams, fields ...string,
) (entity.Coordinators, error) {
	var coordinators entity.Coordinators
	if len(fields) == 0 {
		fields = coordinatorFields
	}

	stmt := c.db.DB.WithContext(ctx).Table(coordinatorTable).Select(fields)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&coordinators).Error
	return coordinators, exception.InternalError(err)
}

func (c *coordinator) Count(ctx context.Context, params *ListCoordinatorsParams) (int64, error) {
	var total int64

	stmt := c.db.DB.WithContext(ctx).Table(coordinatorTable).Select("COUNT(*)")

	err := stmt.Count(&total).Error
	return total, exception.InternalError(err)
}

func (c *coordinator) MultiGet(
	ctx context.Context, coordinatorIDs []string, fields ...string,
) (entity.Coordinators, error) {
	var coordinators entity.Coordinators
	if len(fields) == 0 {
		fields = coordinatorFields
	}

	err := c.db.DB.WithContext(ctx).
		Table(coordinatorTable).Select(fields).
		Where("id IN (?)", coordinatorIDs).
		Find(&coordinators).Error
	return coordinators, exception.InternalError(err)
}

func (c *coordinator) Get(
	ctx context.Context, coordinatorID string, fields ...string,
) (*entity.Coordinator, error) {
	coordinator, err := c.get(ctx, c.db.DB, coordinatorID, fields...)
	return coordinator, exception.InternalError(err)
}

func (c *coordinator) Create(
	ctx context.Context, auth *entity.AdminAuth, coordinator *entity.Coordinator,
) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := c.now()
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

func (c *coordinator) Update(ctx context.Context, coordinatorID string, params *UpdateCoordinatorParams) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := c.get(ctx, tx, coordinatorID); err != nil {
			return nil, err
		}

		updates := map[string]interface{}{
			"lastname":          params.Lastname,
			"firstname":         params.Firstname,
			"lastname_kana":     params.LastnameKana,
			"firstname_kana":    params.FirstnameKana,
			"company_name":      params.CompanyName,
			"store_name":        params.StoreName,
			"thumbnail_url":     params.ThumbnailURL,
			"header_url":        params.HeaderURL,
			"twitter_account":   params.TwitterAccount,
			"instagram_account": params.InstagramAccount,
			"facebook_account":  params.FacebookAccount,
			"phone_number":      params.PhoneNumber,
			"postal_code":       params.PostalCode,
			"city":              params.City,
			"address_line1":     params.AddressLine1,
			"address_line2":     params.AddressLine2,
			"updated_at":        c.now(),
		}
		err := tx.WithContext(ctx).
			Table(coordinatorTable).
			Where("id = ?", coordinatorID).
			Updates(updates).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (c *coordinator) UpdateEmail(ctx context.Context, coordinatorID, email string) error {
	_, err := c.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := c.get(ctx, tx, coordinatorID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"email":      email,
			"updated_at": c.now(),
		}
		err := tx.WithContext(ctx).
			Table(coordinatorTable).
			Where("id = ?", coordinatorID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (c *coordinator) get(
	ctx context.Context, tx *gorm.DB, coordinatorID string, fields ...string,
) (*entity.Coordinator, error) {
	var coordinator *entity.Coordinator
	if len(fields) == 0 {
		fields = coordinatorFields
	}

	err := tx.WithContext(ctx).
		Table(coordinatorTable).Select(fields).
		Where("id = ?", coordinatorID).
		First(&coordinator).Error
	return coordinator, err
}
