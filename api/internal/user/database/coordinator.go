package database

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const coordinatorTable = "coordinators"

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

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&coordinators).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := c.fill(ctx, c.db.DB, coordinators...); err != nil {
		return nil, exception.InternalError(err)
	}
	return coordinators, nil
}

func (c *coordinator) Count(ctx context.Context, _ *ListCoordinatorsParams) (int64, error) {
	total, err := c.db.Count(ctx, c.db.DB, &entity.Coordinator{}, nil)
	return total, exception.InternalError(err)
}

func (c *coordinator) MultiGet(
	ctx context.Context, coordinatorIDs []string, fields ...string,
) (entity.Coordinators, error) {
	var coordinators entity.Coordinators

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...).
		Where("admin_id IN (?)", coordinatorIDs)

	if err := stmt.Find(&coordinators).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := c.fill(ctx, c.db.DB, coordinators...); err != nil {
		return nil, exception.InternalError(err)
	}
	return coordinators, nil
}

func (c *coordinator) Get(
	ctx context.Context, coordinatorID string, fields ...string,
) (*entity.Coordinator, error) {
	coordinator, err := c.get(ctx, c.db.DB, coordinatorID, fields...)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	if err := c.fill(ctx, c.db.DB, coordinator); err != nil {
		return nil, exception.InternalError(err)
	}
	return coordinator, nil
}

func (c *coordinator) Create(
	ctx context.Context, coordinator *entity.Coordinator, auth func(ctx context.Context) error,
) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		if err := coordinator.FillJSON(); err != nil {
			return err
		}
		now := c.now()
		coordinator.Admin.CreatedAt, coordinator.Admin.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(adminTable).Create(&coordinator.Admin).Error; err != nil {
			return err
		}
		coordinator.CreatedAt, coordinator.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(coordinatorTable).Create(&coordinator).Error; err != nil {
			return err
		}
		return auth(ctx)
	})
	return exception.InternalError(err)
}

func (c *coordinator) Update(ctx context.Context, coordinatorID string, params *UpdateCoordinatorParams) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := c.get(ctx, tx, coordinatorID); err != nil {
			return err
		}

		now := c.now()
		adminParams := map[string]interface{}{
			"lastname":       params.Lastname,
			"firstname":      params.Firstname,
			"lastname_kana":  params.LastnameKana,
			"firstname_kana": params.FirstnameKana,
			"updated_at":     now,
		}
		coordinatorParams := map[string]interface{}{
			"marche_name":         params.MarcheName,
			"username":            params.Username,
			"profile":             params.Profile,
			"thumbnail_url":       params.ThumbnailURL,
			"header_url":          params.HeaderURL,
			"promotion_video_url": params.PromotionVideoURL,
			"bonus_video_url":     params.BonusVideoURL,
			"instagram_id":        params.InstagramID,
			"facebook_id":         params.FacebookID,
			"phone_number":        params.PhoneNumber,
			"postal_code":         params.PostalCode,
			"city":                params.City,
			"address_line1":       params.AddressLine1,
			"address_line2":       params.AddressLine2,
			"updated_at":          now,
		}
		if len(params.ProductTypeIDs) > 0 {
			productTypeIDs, err := entity.CoordinatorMarshalProductTypeIDs(params.ProductTypeIDs)
			if err != nil {
				return fmt.Errorf("database: %w: %s", exception.ErrInvalidArgument, err.Error())
			}
			coordinatorParams["product_type_ids"] = datatypes.JSON(productTypeIDs)
		}

		err := tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", coordinatorID).
			Updates(adminParams).Error
		if err != nil {
			return err
		}
		err = tx.WithContext(ctx).
			Table(coordinatorTable).
			Where("admin_id = ?", coordinatorID).
			Updates(coordinatorParams).Error
		return err
	})
	return exception.InternalError(err)
}

func (c *coordinator) UpdateThumbnails(ctx context.Context, coordinatorID string, thumbnails common.Images) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		coordinator, err := c.get(ctx, tx, coordinatorID, "thumbnail_url")
		if err != nil {
			return err
		}
		if coordinator.ThumbnailURL == "" {
			return fmt.Errorf("database: thumbnail url is empty: %w", exception.ErrFailedPrecondition)
		}

		buf, err := thumbnails.Marshal()
		if err != nil {
			return err
		}
		params := map[string]interface{}{
			"thumbnails": datatypes.JSON(buf),
			"updated_at": c.now(),
		}

		err = tx.WithContext(ctx).
			Table(coordinatorTable).
			Where("admin_id = ?", coordinatorID).
			Updates(params).Error
		return err
	})
	return exception.InternalError(err)
}

func (c *coordinator) UpdateHeaders(ctx context.Context, coordinatorID string, headers common.Images) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		coordinator, err := c.get(ctx, tx, coordinatorID, "header_url")
		if err != nil {
			return err
		}
		if coordinator.HeaderURL == "" {
			return fmt.Errorf("database: header url is empty: %w", exception.ErrFailedPrecondition)
		}

		buf, err := headers.Marshal()
		if err != nil {
			return err
		}
		params := map[string]interface{}{
			"headers":    datatypes.JSON(buf),
			"updated_at": c.now(),
		}

		err = tx.WithContext(ctx).
			Table(coordinatorTable).
			Where("admin_id = ?", coordinatorID).
			Updates(params).Error
		return err
	})
	return exception.InternalError(err)
}

func (c *coordinator) Delete(ctx context.Context, coordinatorID string, auth func(ctx context.Context) error) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := c.get(ctx, tx, coordinatorID); err != nil {
			return err
		}

		now := c.now()
		coordinatorParams := map[string]interface{}{
			"updated_at": now,
			"deleted_at": now,
		}
		err := tx.WithContext(ctx).
			Table(coordinatorTable).
			Where("admin_id = ?", coordinatorID).
			Updates(coordinatorParams).Error
		if err != nil {
			return err
		}
		adminParams := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		err = tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", coordinatorID).
			Updates(adminParams).Error
		if err != nil {
			return err
		}
		return auth(ctx)
	})
	return exception.InternalError(err)
}

func (c *coordinator) get(
	ctx context.Context, tx *gorm.DB, coordinatorID string, fields ...string,
) (*entity.Coordinator, error) {
	var coordinator *entity.Coordinator

	err := c.db.Statement(ctx, tx, coordinatorTable, fields...).
		Where("admin_id = ?", coordinatorID).
		First(&coordinator).Error
	return coordinator, err
}

func (c *coordinator) fill(ctx context.Context, tx *gorm.DB, coordinators ...*entity.Coordinator) error {
	var admins entity.Admins

	ids := entity.Coordinators(coordinators).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := c.db.Statement(ctx, tx, adminTable).
		Where("id IN (?)", ids)
	if err := stmt.Find(&admins).Error; err != nil {
		return err
	}

	adminMap := admins.Map()

	for i, c := range coordinators {
		admin, ok := adminMap[c.AdminID]
		if !ok {
			admin = &entity.Admin{}
		}

		if err := coordinators[i].Fill(admin); err != nil {
			return err
		}
	}
	return nil
}
