package tidb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	apmysql "github.com/and-period/furumaru/api/pkg/mysql"
	"golang.org/x/sync/errgroup"
	"gorm.io/datatypes"
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
			Or("`marche_name` LIKE ?", "%"+p.Name+"%").
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
	var internal internalCoordinators

	p := listCoordinatorsParams(*params)

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	coordinators, err := internal.entities()
	if err != nil {
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
	var internal internalCoordinators

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...).
		Where("admin_id IN (?)", coordinatorIDs)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	coordinators, err := internal.entities()
	if err != nil {
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
	var internal internalCoordinators

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...).
		Where("admin_id IN (?)", coordinatorIDs).
		Unscoped()

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	coordinators, err := internal.entities()
	if err != nil {
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
	var internal *internalCoordinator

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...).
		Where("admin_id = ?", coordinatorID).
		Unscoped()

	if err := stmt.First(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	coordinator, err := internal.entity()
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
		internal, err := newInternalCoordinator(coordinator)
		if err != nil {
			return err
		}
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
		if err := tx.WithContext(ctx).Table(coordinatorTable).Create(&internal).Error; err != nil {
			return err
		}
		if len(groups) > 0 {
			if err := tx.WithContext(ctx).Table(adminGroupUserTable).Create(&groups).Error; err != nil {
				return err
			}
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (c *coordinator) Update(ctx context.Context, coordinatorID string, params *database.UpdateCoordinatorParams) error {
	err := c.db.Transaction(ctx, func(tx *gorm.DB) error {
		productTypeIDs, err := json.Marshal(params.ProductTypeIDs)
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
	var internal *internalCoordinator

	stmt := c.db.Statement(ctx, tx, coordinatorTable, fields...).
		Where("admin_id = ?", coordinatorID)

	if err := stmt.First(&internal).Error; err != nil {
		return nil, err
	}
	coordinator, err := internal.entity()
	if err != nil {
		return nil, err
	}

	return coordinator, err
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

type internalCoordinator struct {
	entity.Coordinator `gorm:"embedded"`
	ProductTypeIDsJSON datatypes.JSON `gorm:"default:null;column:product_type_ids"` // 取り扱い品目ID一覧(JSON)
}

type internalCoordinators []*internalCoordinator

func newInternalCoordinator(coordinator *entity.Coordinator) (*internalCoordinator, error) {
	typeIDs, err := json.Marshal(coordinator.ProductTypeIDs)
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal product type IDs: %w", err)
	}
	internal := &internalCoordinator{
		Coordinator:        *coordinator,
		ProductTypeIDsJSON: typeIDs,
	}
	return internal, nil
}

func (c *internalCoordinator) entity() (*entity.Coordinator, error) {
	if err := c.unmarshalProductTypeIDs(); err != nil {
		return nil, err
	}
	return &c.Coordinator, nil
}

func (c *internalCoordinator) unmarshalProductTypeIDs() error {
	if c == nil || c.ProductTypeIDsJSON == nil {
		return nil
	}
	var ids []string
	if err := json.Unmarshal(c.ProductTypeIDsJSON, &ids); err != nil {
		return fmt.Errorf("tidb: failed to unmarshal product type IDs: %w", err)
	}
	c.ProductTypeIDs = ids
	return nil
}

func (cs internalCoordinators) entities() (entity.Coordinators, error) {
	res := make(entity.Coordinators, len(cs))
	for i := range cs {
		c, err := cs[i].entity()
		if err != nil {
			return nil, err
		}
		res[i] = c
	}
	return res, nil
}
