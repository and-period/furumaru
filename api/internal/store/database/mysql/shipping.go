package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const shippingTable = "shippings"

type shipping struct {
	db  *mysql.Client
	now func() time.Time
}

func newShipping(db *mysql.Client) database.Shipping {
	return &shipping{
		db:  db,
		now: jst.Now,
	}
}

type listShippingsParams database.ListShippingsParams

func (p listShippingsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("`%s` ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("`%s` DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

func (p listShippingsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (s *shipping) List(ctx context.Context, params *database.ListShippingsParams, fields ...string) (entity.Shippings, error) {
	var shippings entity.Shippings

	p := listShippingsParams(*params)

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&shippings).Error; err != nil {
		return nil, dbError(err)
	}
	if err := shippings.Fill(); err != nil {
		return nil, dbError(err)
	}
	return shippings, nil
}

func (s *shipping) Count(ctx context.Context, params *database.ListShippingsParams) (int64, error) {
	p := listShippingsParams(*params)

	total, err := s.db.Count(ctx, s.db.DB, &entity.Shipping{}, p.stmt)
	return total, dbError(err)
}

func (s *shipping) MultiGet(ctx context.Context, shippingIDs []string, fields ...string) (entity.Shippings, error) {
	var shippings entity.Shippings

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...).
		Where("id IN (?)", shippingIDs)
	if err := stmt.Find(&shippings).Error; err != nil {
		return nil, dbError(err)
	}
	if err := shippings.Fill(); err != nil {
		return nil, dbError(err)
	}
	return shippings, nil
}

func (s *shipping) Get(ctx context.Context, shoppingID string, fields ...string) (*entity.Shipping, error) {
	shopping, err := s.get(ctx, s.db.DB, shoppingID, fields...)
	return shopping, dbError(err)
}

func (s *shipping) Create(ctx context.Context, shipping *entity.Shipping) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		if shipping.IsDefault {
			updates := map[string]interface{}{
				"is_default": false,
				"updated_at": s.now(),
			}
			stmt := tx.WithContext(ctx).Table(shippingTable).Where("coordinator_id = ?", shipping.CoordinatorID)
			if err := stmt.Updates(updates).Error; err != nil {
				return err
			}
		}

		if err := shipping.FillJSON(); err != nil {
			return err
		}

		now := s.now()
		shipping.CreatedAt, shipping.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(shippingTable).Create(&shipping).Error
		return err
	})
	return dbError(err)
}

func (s *shipping) Update(ctx context.Context, shippingID string, params *database.UpdateShippingParams) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		current, err := s.get(ctx, tx, shippingID)
		if err != nil {
			return err
		}

		if params.IsDefault {
			updates := map[string]interface{}{
				"is_default": false,
				"updated_at": s.now(),
			}
			stmt := tx.WithContext(ctx).Table(shippingTable).Where("coordinator_id = ?", current.CoordinatorID)
			if err := stmt.Updates(updates).Error; err != nil {
				return err
			}
		}

		updates := map[string]interface{}{
			"name":                params.Name,
			"is_default":          params.IsDefault,
			"box60_refrigerated":  params.Box60Refrigerated,
			"box60_frozen":        params.Box60Frozen,
			"box80_refrigerated":  params.Box80Refrigerated,
			"box80_frozen":        params.Box80Frozen,
			"box100_refrigerated": params.Box100Refrigerated,
			"box100_frozen":       params.Box100Frozen,
			"has_free_shipping":   params.HasFreeShipping,
			"free_shipping_rates": params.FreeShippingRates,
			"updated_at":          s.now(),
		}
		if len(params.Box60Rates) > 0 {
			rates, err := params.Box60Rates.Marshal()
			if err != nil {
				return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
			}
			updates["box60_rates"] = rates
		}
		if len(params.Box80Rates) > 0 {
			rates, err := params.Box80Rates.Marshal()
			if err != nil {
				return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
			}
			updates["box80_rates"] = rates
		}
		if len(params.Box100Rates) > 0 {
			rates, err := params.Box100Rates.Marshal()
			if err != nil {
				return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
			}
			updates["box100_rates"] = rates
		}
		err = tx.WithContext(ctx).
			Table(shippingTable).
			Where("id = ?", shippingID).
			Updates(updates).Error
		return err
	})
	return dbError(err)
}

func (s *shipping) Delete(ctx context.Context, shippingID string) error {
	stmt := s.db.DB.WithContext(ctx).
		Table(shippingTable).
		Where("id = ?", shippingID)

	err := stmt.Delete(&entity.Shipping{}).Error
	return dbError(err)
}

func (s *shipping) get(
	ctx context.Context, tx *gorm.DB, shippingID string, fields ...string,
) (*entity.Shipping, error) {
	var shipping *entity.Shipping

	err := s.db.Statement(ctx, tx, shippingTable, fields...).
		Where("id = ?", shippingID).
		First(&shipping).Error
	if err != nil {
		return nil, err
	}
	if err := shipping.Fill(); err != nil {
		return nil, err
	}
	return shipping, nil
}
