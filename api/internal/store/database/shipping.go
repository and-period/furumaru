package database

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const shippingTable = "shippings"

var shippingFields = []string{
	"id", "name",
	"box60_rates", "box60_refrigerated", "box60_frozen",
	"box80_rates", "box80_refrigerated", "box80_frozen",
	"box100_rates", "box100_refrigerated", "box100_frozen",
	"has_free_shipping", "free_shipping_rates", "created_at", "updated_at",
}

type shipping struct {
	db  *database.Client
	now func() time.Time
}

func NewShipping(db *database.Client) Shipping {
	return &shipping{
		db:  db,
		now: jst.Now,
	}
}

func (s *shipping) List(ctx context.Context, params *ListShippingsParams, fields ...string) (entity.Shippings, error) {
	var shippings entity.Shippings
	if len(fields) == 0 {
		fields = shippingFields
	}

	stmt := s.db.DB.WithContext(ctx).Table(shippingTable).Select(fields)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&shippings).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := shippings.Fill(); err != nil {
		return nil, exception.InternalError(err)
	}
	return shippings, nil
}

func (s *shipping) Count(ctx context.Context, params *ListShippingsParams) (int64, error) {
	var total int64

	stmt := s.db.DB.WithContext(ctx).Table(shippingTable).Select("COUNT(*)")

	err := stmt.Find(&total).Error
	return total, exception.InternalError(err)
}

func (s *shipping) MultiGet(ctx context.Context, shippingIDs []string, fields ...string) (entity.Shippings, error) {
	var shippings entity.Shippings
	if len(fields) == 0 {
		fields = shippingFields
	}

	err := s.db.DB.WithContext(ctx).
		Table(shippingTable).Select(fields).
		Where("id IN (?)", shippingIDs).
		Find(&shippings).Error
	if err := shippings.Fill(); err != nil {
		return nil, err
	}
	return shippings, exception.InternalError(err)
}

func (s *shipping) Get(ctx context.Context, shoppingID string, fields ...string) (*entity.Shipping, error) {
	shopping, err := s.get(ctx, s.db.DB, shoppingID, fields...)
	return shopping, exception.InternalError(err)
}

func (s *shipping) Create(ctx context.Context, shipping *entity.Shipping) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if err := shipping.FillJSON(); err != nil {
			return nil, err
		}

		now := s.now()
		shipping.CreatedAt, shipping.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(shippingTable).Create(&shipping).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (s *shipping) Update(ctx context.Context, shippingID string, params *UpdateShippingParams) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := s.get(ctx, tx, shippingID); err != nil {
			return nil, err
		}

		updates := map[string]interface{}{
			"name":                params.Name,
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
				return nil, fmt.Errorf("database: %w: %s", exception.ErrInvalidArgument, err.Error())
			}
			updates["box60_rates"] = rates
		}
		if len(params.Box80Rates) > 0 {
			rates, err := params.Box80Rates.Marshal()
			if err != nil {
				return nil, fmt.Errorf("database: %w: %s", exception.ErrInvalidArgument, err.Error())
			}
			updates["box80_rates"] = rates
		}
		if len(params.Box100Rates) > 0 {
			rates, err := params.Box100Rates.Marshal()
			if err != nil {
				return nil, fmt.Errorf("database: %w: %s", exception.ErrInvalidArgument, err.Error())
			}
			updates["box100_rates"] = rates
		}
		err := tx.WithContext(ctx).
			Table(shippingTable).
			Where("id = ?", shippingID).
			Updates(updates).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (s *shipping) Delete(ctx context.Context, shippingID string) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := s.get(ctx, tx, shippingID); err != nil {
			return nil, err
		}

		err := tx.WithContext(ctx).
			Table(shippingTable).
			Where("id = ?", shippingID).
			Delete(&entity.Shipping{}).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (s *shipping) get(
	ctx context.Context, tx *gorm.DB, shippingID string, fields ...string,
) (*entity.Shipping, error) {
	var shipping *entity.Shipping
	if len(fields) == 0 {
		fields = shippingFields
	}

	err := tx.WithContext(ctx).
		Table(shippingTable).Select(fields).
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
