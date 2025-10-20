package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	apmysql "github.com/and-period/furumaru/api/pkg/mysql"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

const coordinatorTable = "coordinators"

type coordinator struct {
	db  *apmysql.Client
	now func() time.Time
}

func NewCoordinator(db *apmysql.Client) database.Coordinator {
	return &coordinator{
		db:  db,
		now: jst.Now,
	}
}

type listCoordinatorsParams database.ListCoordinatorsParams

func (p listCoordinatorsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`username` LIKE ?", "%"+p.Name+"%").
			Or("`profile` LIKE ?", "%"+p.Name+"%")
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

func (c *coordinator) MultiGetWithDeleted(
	ctx context.Context, coordinatorIDs []string, fields ...string,
) (entity.Coordinators, error) {
	var coordinators entity.Coordinators

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...).
		Where("admin_id IN (?)", coordinatorIDs).
		Unscoped()

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

func (c *coordinator) GetWithDeleted(
	ctx context.Context, coordinatorID string, fields ...string,
) (*entity.Coordinator, error) {
	var coordinator *entity.Coordinator

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...).
		Where("admin_id = ?", coordinatorID).
		Unscoped()

	if err := stmt.First(&coordinator).Error; err != nil {
		return nil, dbError(err)
	}
	if err := c.fill(ctx, c.db.DB, coordinator); err != nil {
		return nil, dbError(err)
	}
	return coordinator, nil
}

func (c *coordinator) Create(
	ctx context.Context, coordinator *entity.Coordinator, shop *entity.Shop, auth func(ctx context.Context) error,
) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		params := &entity.NewAdminGroupUsersParams{
			AdminID:  coordinator.AdminID,
			GroupIDs: coordinator.GroupIDs,
		}
		groups := entity.NewAdminGroupUsers(params)

		now := c.now()
		coordinator.CreatedAt, coordinator.UpdatedAt = now, now
		coordinator.Admin.CreatedAt, coordinator.Admin.UpdatedAt = now, now
		for _, group := range groups {
			group.CreatedAt, group.UpdatedAt = now, now
		}

		if err := tx.WithContext(ctx).Table(adminTable).Create(&coordinator.Admin).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Table(coordinatorTable).Create(&coordinator).Error; err != nil {
			return err
		}

		if len(groups) > 0 {
			if err := tx.WithContext(ctx).Table(adminGroupUserTable).Create(&groups).Error; err != nil {
				return err
			}
		}

		if err := tx.WithContext(ctx).Table(shopTable).Create(&shop).Error; err != nil {
			return err
		}

		return auth(ctx)
	})
	return dbError(err)
}

func (c *coordinator) Update(ctx context.Context, coordinatorID string, params *database.UpdateCoordinatorParams) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := c.now()
		adminParams := map[string]interface{}{
			"lastname":       params.Lastname,
			"firstname":      params.Firstname,
			"lastname_kana":  params.LastnameKana,
			"firstname_kana": params.FirstnameKana,
			"updated_at":     now,
		}
		coordinatorParams := map[string]interface{}{
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
			"updated_at":          now,
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
	return dbError(err)
}

func (c *coordinator) Delete(ctx context.Context, coordinatorID string, auth func(ctx context.Context) error) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := c.now()
		aupdates := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		stmt := tx.WithContext(ctx).Table(adminTable).Where("id = ?", coordinatorID)
		if err := stmt.Updates(aupdates).Error; err != nil {
			return err
		}
		cupdate := map[string]interface{}{
			"updated_at": now,
			"deleted_at": now,
		}
		stmt = tx.WithContext(ctx).Table(coordinatorTable).Where("admin_id = ?", coordinatorID)
		if err := stmt.Updates(cupdate).Error; err != nil {
			return err
		}
		supdate := map[string]interface{}{
			"activated":  false,
			"updated_at": now,
			"deleted_at": now,
		}
		stmt = tx.WithContext(ctx).Table(shopTable).Where("coordinator_id = ?", coordinatorID)
		if err := stmt.Updates(supdate).Error; err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (c *coordinator) get(
	ctx context.Context, tx *gorm.DB, coordinatorID string, fields ...string,
) (*entity.Coordinator, error) {
	var coordinator *entity.Coordinator

	stmt := c.db.Statement(ctx, tx, coordinatorTable, fields...).
		Where("admin_id = ?", coordinatorID)

	if err := stmt.First(&coordinator).Error; err != nil {
		return nil, err
	}
	return coordinator, nil
}

func (c *coordinator) fill(ctx context.Context, tx *gorm.DB, coordinators ...*entity.Coordinator) error {
	var (
		admins entity.Admins
		groups entity.AdminGroupUsers
	)

	ids := entity.Coordinators(coordinators).IDs()
	if len(ids) == 0 {
		return nil
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		stmt := c.db.Statement(ectx, tx, adminTable).
			Unscoped().
			Where("id IN (?)", ids)
		return stmt.Find(&admins).Error
	})
	eg.Go(func() error {
		stmt := c.db.Statement(ectx, tx, adminGroupUserTable).
			Where("admin_id IN (?)", ids).
			Where("expired_at IS NULL OR expired_at > ?", jst.Now())
		return stmt.Find(&groups).Error
	})
	if err := eg.Wait(); err != nil {
		return err
	}

	entity.Coordinators(coordinators).Fill(admins.Map(), groups.GroupByAdminID())
	return nil
}
