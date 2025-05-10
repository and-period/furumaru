package tidb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/datatypes"
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

type listShippingsParams database.ListShippingsParams

func (p listShippingsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.ShopID != "" {
		stmt = stmt.Where("shop_id = ?", p.ShopID)
	}
	if len(p.CoordinatorIDs) > 0 {
		stmt = stmt.Where("coordinator_id IN (?)", p.CoordinatorIDs)
	}
	if p.OnlyInUse {
		stmt = stmt.Where("in_use = ?", true)
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

	prm := listShippingsParams(*params)

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...)
	stmt = prm.stmt(stmt)
	stmt = prm.pagination(stmt)

	if err := stmt.Find(&shippings).Error; err != nil {
		return nil, dbError(err)
	}
	if err := s.fill(ctx, s.db.DB, shippings...); err != nil {
		return nil, dbError(err)
	}
	return shippings, nil
}

func (s *shipping) ListByCoordinatorIDs(ctx context.Context, coordinatorIDs []string, fields ...string) (entity.Shippings, error) {
	var shippings entity.Shippings

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...).
		Where("coordinator_id IN (?)", coordinatorIDs).
		Where("in_use = ?", true)

	if err := stmt.Find(&shippings).Error; err != nil {
		return nil, dbError(err)
	}
	if err := s.fill(ctx, s.db.DB, shippings...); err != nil {
		return nil, dbError(err)
	}
	return shippings, nil
}

func (s *shipping) Count(ctx context.Context, params *database.ListShippingsParams) (int64, error) {
	prm := listShippingsParams(*params)

	total, err := s.db.Count(ctx, s.db.DB, &entity.Shipping{}, prm.stmt)
	return total, dbError(err)
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
	var internal internalShippingRevisions

	stmt := s.db.Statement(ctx, s.db.DB, shippingRevisionTable).
		Where("id IN (?)", revisionIDs)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	if len(internal) == 0 {
		return entity.Shippings{}, nil
	}
	revisions, err := internal.entities()
	if err != nil {
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

func (s *shipping) Get(ctx context.Context, shippingID string, fields ...string) (*entity.Shipping, error) {
	var shipping *entity.Shipping

	stmt := s.db.Statement(ctx, s.db.DB, shippingTable, fields...).
		Where("id = ?", shippingID)

	if err := stmt.First(&shipping).Error; err != nil {
		return nil, dbError(err)
	}
	if err := s.fill(ctx, s.db.DB, shipping); err != nil {
		return nil, dbError(err)
	}
	return shipping, nil
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
		Where("coordinator_id = ?", coordinatorID).
		Where("in_use = ?", true)

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
		now := s.now()

		// 既存で使用中の配送設定がある場合、それを無効化する
		if shipping.InUse {
			params := map[string]interface{}{
				"in_use":     false,
				"updated_at": now,
			}
			stmt := tx.WithContext(ctx).
				Table(shippingTable).
				Where("shop_id = ?", shipping.ShopID).
				Where("in_use = ?", true)
			if err := stmt.Updates(params).Error; err != nil {
				return err
			}
		}

		// 配送設定を登録
		shipping.CreatedAt, shipping.UpdatedAt = now, now
		shipping.ShippingRevision.CreatedAt, shipping.ShippingRevision.UpdatedAt = now, now

		internal, err := newInternalShippingRevision(&shipping.ShippingRevision)
		if err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Table(shippingTable).Create(&shipping).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Table(shippingRevisionTable).Create(&internal).Error
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

	now := s.now()
	revision.CreatedAt, revision.UpdatedAt = now, now

	internal, err := newInternalShippingRevision(revision)
	if err != nil {
		return fmt.Errorf("tidb: %w: %s", database.ErrInvalidArgument, err.Error())
	}

	err = s.db.DB.WithContext(ctx).Table(shippingRevisionTable).Create(&internal).Error
	return dbError(err)
}

func (s *shipping) UpdateInUse(ctx context.Context, shopID, shippingID string) error {
	return s.db.Transaction(ctx, func(tx *gorm.DB) error {
		current, err := s.get(ctx, tx, shippingID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if current != nil {
			// 別店舗の情報は更新できない
			if current.ShopID != shopID {
				return fmt.Errorf("tidb: shop id mismatch: %w", database.ErrFailedPrecondition)
			}
			// すでに使用中の場合は何もしない
			if current.InUse {
				return nil
			}
		}

		oldParams := map[string]interface{}{
			"in_use":     false,
			"updated_at": s.now(),
		}
		stmt := tx.WithContext(ctx).
			Table(shippingTable).
			Where("shop_id = ?", shopID).
			Where("in_use = ?", true)
		if err := stmt.Updates(oldParams).Error; err != nil {
			return err
		}

		newParams := map[string]interface{}{
			"in_use":     true,
			"updated_at": s.now(),
		}
		stmt = tx.WithContext(ctx).
			Table(shippingTable).
			Where("id = ?", shippingID)
		return stmt.Updates(newParams).Error
	})
}

func (s *shipping) Delete(ctx context.Context, shippingID string) error {
	err := s.db.Transaction(ctx, func(tx *gorm.DB) error {
		if shippingID == entity.DefaultShippingID {
			return fmt.Errorf("tidb: default shipping cannot be deleted: %w", database.ErrPermissionDenied)
		}

		current, err := s.get(ctx, tx, shippingID)
		if err != nil {
			return err
		}
		if current.InUse {
			return fmt.Errorf("tidb: shipping is in use: %w", database.ErrFailedPrecondition)
		}

		params := map[string]interface{}{
			"deleted_at": s.now(),
		}
		stmt := tx.WithContext(ctx).
			Table(shippingTable).
			Where("id = ?", shippingID)

		return stmt.Updates(params).Error
	})
	return dbError(err)
}

func (s *shipping) get(ctx context.Context, tx *gorm.DB, shippingID string, fields ...string) (*entity.Shipping, error) {
	var shipping *entity.Shipping

	stmt := s.db.Statement(ctx, tx, shippingTable, fields...).
		Where("id = ?", shippingID)

	if err := stmt.First(&shipping).Error; err != nil {
		return nil, err
	}
	if err := s.fill(ctx, tx, shipping); err != nil {
		return nil, err
	}
	return shipping, nil
}

func (s *shipping) fill(ctx context.Context, tx *gorm.DB, shippings ...*entity.Shipping) error {
	var internal internalShippingRevisions

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

	if err := stmt.Find(&internal).Error; err != nil {
		return err
	}
	if len(internal) == 0 {
		return nil
	}
	revisions, err := internal.entities()
	if err != nil {
		return err
	}

	entity.Shippings(shippings).Fill(revisions.MapByShippingID())
	return nil
}

type internalShippingRevision struct {
	entity.ShippingRevision `gorm:"embedded"`
	Box60RatesJSON          datatypes.JSON `gorm:"default:null;column:box60_rates"`  // 箱サイズ60の通常便配送料一覧(JSON)
	Box80RatesJSON          datatypes.JSON `gorm:"default:null;column:box80_rates"`  // 箱サイズ80の通常便配送料一覧(JSON)
	Box100RatesJSON         datatypes.JSON `gorm:"default:null;column:box100_rates"` // 箱サイズ100の通常便配送料一覧(JSON)
}

type internalShippingRevisions []*internalShippingRevision

func newInternalShippingRevision(revision *entity.ShippingRevision) (*internalShippingRevision, error) {
	box60Rates, err := revision.Box60Rates.Marshal()
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal box60 rates: %w", err)
	}
	box80Rates, err := revision.Box80Rates.Marshal()
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal box80 rates: %w", err)
	}
	box100Rates, err := revision.Box100Rates.Marshal()
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal box100 rates: %w", err)
	}
	internal := &internalShippingRevision{
		ShippingRevision: *revision,
		Box60RatesJSON:   box60Rates,
		Box80RatesJSON:   box80Rates,
		Box100RatesJSON:  box100Rates,
	}
	return internal, nil
}

func (r *internalShippingRevision) entity() (*entity.ShippingRevision, error) {
	if err := r.unmarshalBox60Rates(); err != nil {
		return nil, err
	}
	if err := r.unmarshalBox80Rates(); err != nil {
		return nil, err
	}
	if err := r.unmarshalBox100Rates(); err != nil {
		return nil, err
	}
	return &r.ShippingRevision, nil
}

func (r *internalShippingRevision) unmarshalBox60Rates() error {
	if r == nil || r.Box60RatesJSON == nil {
		return nil
	}
	var rates entity.ShippingRates
	if err := json.Unmarshal(r.Box60RatesJSON, &rates); err != nil {
		return fmt.Errorf("tidb: failed to unmarshal box60 rates: %w", err)
	}
	r.ShippingRevision.Box60Rates = rates
	return nil
}

func (r *internalShippingRevision) unmarshalBox80Rates() error {
	if r == nil || r.Box80RatesJSON == nil {
		return nil
	}
	var rates entity.ShippingRates
	if err := json.Unmarshal(r.Box80RatesJSON, &rates); err != nil {
		return fmt.Errorf("tidb: failed to unmarshal box80 rates: %w", err)
	}
	r.ShippingRevision.Box80Rates = rates
	return nil
}

func (r *internalShippingRevision) unmarshalBox100Rates() error {
	if r == nil || r.Box100RatesJSON == nil {
		return nil
	}
	var rates entity.ShippingRates
	if err := json.Unmarshal(r.Box100RatesJSON, &rates); err != nil {
		return fmt.Errorf("tidb: failed to unmarshal box100 rates: %w", err)
	}
	r.ShippingRevision.Box100Rates = rates
	return nil
}

func (rs internalShippingRevisions) entities() (entity.ShippingRevisions, error) {
	res := make(entity.ShippingRevisions, len(rs))
	for i := range rs {
		r, err := rs[i].entity()
		if err != nil {
			return nil, err
		}
		res[i] = r
	}
	return res, nil
}
