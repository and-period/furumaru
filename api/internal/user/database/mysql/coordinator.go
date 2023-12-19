package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const coordinatorTable = "coordinators"

type coordinator struct {
	db  *mysql.Client
	now func() time.Time
}

func newCoordinator(db *mysql.Client) database.Coordinator {
	return &coordinator{
		db:  db,
		now: jst.Now,
	}
}

type listCoordinatorsParams database.ListCoordinatorsParams

func (p listCoordinatorsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Username != "" {
		stmt = stmt.Where("username LIKE ?", fmt.Sprintf("%%%s%%", p.Username))
	}
	stmt = stmt.Order("updated_at DESC")
	return stmt
}

func (p listCoordinatorsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (c *coordinator) List(
	ctx context.Context, params *database.ListCoordinatorsParams, fields ...string,
) (entity.Coordinators, error) {
	var coordinators entity.Coordinators

	p := listCoordinatorsParams(*params)

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&coordinators).Error; err != nil {
		return nil, dbError(err)
	}
	if err := c.fill(ctx, c.db.DB, coordinators...); err != nil {
		return nil, dbError(err)
	}
	return coordinators, nil
}

func (c *coordinator) Count(ctx context.Context, params *database.ListCoordinatorsParams) (int64, error) {
	p := listCoordinatorsParams(*params)

	total, err := c.db.Count(ctx, c.db.DB, &entity.Coordinator{}, p.stmt)
	return total, dbError(err)
}

func (c *coordinator) MultiGet(
	ctx context.Context, coordinatorIDs []string, fields ...string,
) (entity.Coordinators, error) {
	var coordinators entity.Coordinators

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...).
		Where("admin_id IN (?)", coordinatorIDs)

	if err := stmt.Find(&coordinators).Error; err != nil {
		return nil, dbError(err)
	}
	if err := c.fill(ctx, c.db.DB, coordinators...); err != nil {
		return nil, dbError(err)
	}
	return coordinators, nil
}

func (c *coordinator) Get(
	ctx context.Context, coordinatorID string, fields ...string,
) (*entity.Coordinator, error) {
	coordinator, err := c.get(ctx, c.db.DB, coordinatorID, fields...)
	if err != nil {
		return nil, dbError(err)
	}
	if err := c.fill(ctx, c.db.DB, coordinator); err != nil {
		return nil, dbError(err)
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
	return dbError(err)
}

func (c *coordinator) Update(ctx context.Context, coordinatorID string, params *database.UpdateCoordinatorParams) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		productTypeIDs, err := entity.CoordinatorMarshalProductTypeIDs(params.ProductTypeIDs)
		if err != nil {
			return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
		}
		businessDays, err := entity.CoordinatorMarshalBusinessDays(params.BusinessDays)
		if err != nil {
			return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
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
			"product_type_ids":    productTypeIDs,
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
			"prefecture":          params.PrefectureCode,
			"city":                params.City,
			"address_line1":       params.AddressLine1,
			"address_line2":       params.AddressLine2,
			"business_days":       businessDays,
			"updated_at":          now,
		}

		err = tx.WithContext(ctx).
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
	return dbError(err)
}

func (c *coordinator) UpdateThumbnails(ctx context.Context, coordinatorID string, thumbnails common.Images) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		coordinator, err := c.get(ctx, tx, coordinatorID, "thumbnail_url")
		if err != nil {
			return err
		}
		if coordinator.ThumbnailURL == "" {
			return fmt.Errorf("database: thumbnail url is empty: %w", database.ErrFailedPrecondition)
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
	return dbError(err)
}

func (c *coordinator) UpdateHeaders(ctx context.Context, coordinatorID string, headers common.Images) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		coordinator, err := c.get(ctx, tx, coordinatorID, "header_url")
		if err != nil {
			return err
		}
		if coordinator.HeaderURL == "" {
			return fmt.Errorf("database: header url is empty: %w", database.ErrFailedPrecondition)
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
	return dbError(err)
}

func (c *coordinator) Delete(ctx context.Context, coordinatorID string, auth func(ctx context.Context) error) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := c.now()
		updates := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		stmt := tx.WithContext(ctx).Table(adminTable).Where("id = ?", coordinatorID)
		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (c *coordinator) RemoveProductTypeID(ctx context.Context, productTypeID string) error {
	sub := gorm.Expr("JSON_REMOVE(product_type_ids, JSON_UNQUOTE(JSON_SEARCH(product_type_ids, 'one', ?)))", productTypeID)

	err := c.db.DB.WithContext(ctx).
		Table(coordinatorTable).
		Where("JSON_SEARCH(product_type_ids, 'one', ?) IS NOT NULL", productTypeID).
		Update("product_type_ids", sub).Error
	return dbError(err)
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

	stmt := c.db.Statement(ctx, tx, adminTable).Unscoped().Where("id IN (?)", ids)
	if err := stmt.Find(&admins).Error; err != nil {
		return err
	}
	admins.Fill()
	return entity.Coordinators(coordinators).Fill(admins.Map())
}
