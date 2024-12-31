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

const (
	shippingTable         = "shippings"
	shippingRevisionTable = "shipping_revisions"
)

type shipping struct {
	db  *mysql.Client
	now func() time.Time
}

func NewShipping(db *mysql.Client) database.Shipping {
	return &shipping{
		db:  db,
		now: jst.Now,
	}
}

func (s *shipping) ListByCoordinatorIDs(ctx context.Context, coordinatorIDs []string, fields ...string) (entity.Shippings, error) {
	var shippings entity.Shippings

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...).
		Where("coordinator_id IN (?)", coordinatorIDs)

	if err := stmt.Find(&shippings).Error; err != nil {
		return nil, dbError(err)
	}
	if err := s.fill(ctx, s.db.DB, shippings...); err != nil {
		return nil, dbError(err)
	}
	return shippings, nil
}

func (s *shipping) MultiGet(ctx context.Context, shippingIDs []string, fields ...string) (entity.Shippings, error) {
	var shippings entity.Shippings

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...).
		Where("id IN (?)", shippingIDs)

	if err := stmt.Find(&shippings).Error; err != nil {
		return nil, dbError(err)
	}
	if err := s.fill(ctx, s.db.DB, shippings...); err != nil {
		return nil, dbError(err)
	}
	return shippings, nil
}

func (s *shipping) MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Shippings, error) {
	var revisions entity.ShippingRevisions

	stmt := s.db.Statement(ctx, s.db.DB, shippingRevisionTable).
		Where("id IN (?)", revisionIDs)

	if err := stmt.Find(&revisions).Error; err != nil {
		return nil, dbError(err)
	}
	if len(revisions) == 0 {
		return entity.Shippings{}, nil
	}
	if err := revisions.Fill(); err != nil {
		return nil, dbError(err)
	}

	shippings, err := s.MultiGet(ctx, revisions.ShippingIDs(), fields...)
	if err != nil {
		return nil, err
	}
	if len(shippings) == 0 {
		return entity.Shippings{}, nil
	}

	res, err := revisions.Merge(shippings.Map())
	if err != nil {
		return nil, dbError(err)
	}
	return res, nil
}

func (s *shipping) GetDefault(ctx context.Context, fields ...string) (*entity.Shipping, error) {
	var shipping *entity.Shipping

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...).
		Where("id = ?", entity.DefaultShippingID)

	if err := stmt.First(&shipping).Error; err != nil {
		return nil, dbError(err)
	}
	if err := s.fill(ctx, s.db.DB, shipping); err != nil {
		return nil, dbError(err)
	}
	return shipping, nil
}

func (s *shipping) GetByCoordinatorID(ctx context.Context, coordinatorID string, fields ...string) (*entity.Shipping, error) {
	var shipping *entity.Shipping

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...).
		Where("coordinator_id = ?", coordinatorID)

	if err := stmt.First(&shipping).Error; err != nil {
		return nil, dbError(err)
	}
	if err := s.fill(ctx, s.db.DB, shipping); err != nil {
		return nil, dbError(err)
	}
	return shipping, nil
}

func (s *shipping) Create(ctx context.Context, shipping *entity.Shipping) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		if err := shipping.ShippingRevision.FillJSON(); err != nil {
			return err
		}

		now := s.now()

		shipping.CreatedAt, shipping.UpdatedAt = now, now
		shipping.ShippingRevision.CreatedAt, shipping.ShippingRevision.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(shippingTable).Create(&shipping).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Table(shippingRevisionTable).Create(&shipping.ShippingRevision).Error
	})
	return dbError(err)
}

func (s *shipping) Update(ctx context.Context, shippingID string, params *database.UpdateShippingParams) error {
	rparams := &entity.NewShippingRevisionParams{
		ShippingID:        shippingID,
		Box60Rates:        params.Box60Rates,
		Box60Frozen:       params.Box60Frozen,
		Box80Rates:        params.Box80Rates,
		Box80Frozen:       params.Box80Frozen,
		Box100Rates:       params.Box100Rates,
		Box100Frozen:      params.Box100Frozen,
		HasFreeShipping:   params.HasFreeShipping,
		FreeShippingRates: params.FreeShippingRates,
	}
	revision := entity.NewShippingRevision(rparams)
	if err := revision.FillJSON(); err != nil {
		return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
	}

	now := s.now()
	revision.CreatedAt, revision.UpdatedAt = now, now
	err := s.db.DB.WithContext(ctx).Table(shippingRevisionTable).Create(&revision).Error
	return dbError(err)
}

func (s *shipping) fill(ctx context.Context, tx *gorm.DB, shippings ...*entity.Shipping) error {
	var revisions entity.ShippingRevisions

	ids := entity.Shippings(shippings).IDs()
	if len(ids) == 0 {
		return nil
	}

	sub := tx.Table(shippingRevisionTable).
		Select("MAX(id)").
		Where("shipping_id IN (?)", ids).
		Group("shipping_id")
	stmt := s.db.Statement(ctx, tx, shippingRevisionTable).
		Where("id IN (?)", sub)

	if err := stmt.Find(&revisions).Error; err != nil {
		return err
	}
	if len(revisions) == 0 {
		return nil
	}
	if err := revisions.Fill(); err != nil {
		return err
	}

	entity.Shippings(shippings).Fill(revisions.MapByShippingID())
	return nil
}
