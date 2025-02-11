package tidb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	shopTable         = "shops"
	shopProducerTable = "shop_producers"
)

type shop struct {
	db  *mysql.Client
	now func() time.Time
}

func NewShop(db *mysql.Client) database.Shop {
	return &shop{
		db:  db,
		now: time.Now,
	}
}

func (s *shop) ListByProducerID(ctx context.Context, producerID string, fields ...string) (entity.Shops, error) {
	var internal internalShops

	stmt := s.db.Statement(ctx, s.db.DB, shopTable, fields...).
		Joins("JOIN shop_producers ON shop_producers.shop_id = shops.id").
		Where("shop_producers.producer_id = ?", producerID).
		Where("shops.deleted_at IS NULL")

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	shops, err := internal.entities()
	if err != nil {
		return nil, dbError(err)
	}

	if err := s.fill(ctx, s.db.DB, shops...); err != nil {
		return nil, dbError(err)
	}
	return shops, nil
}

func (s *shop) Get(ctx context.Context, shopID string, fields ...string) (*entity.Shop, error) {
	shop, err := s.get(ctx, s.db.DB, shopID, fields...)
	return shop, dbError(err)
}

func (s *shop) GetByCoordinatorID(ctx context.Context, coordinatorID string, fields ...string) (*entity.Shop, error) {
	var internal *internalShop

	stmt := s.db.Statement(ctx, s.db.DB, shopTable, fields...).
		Where("coordinator_id = ?", coordinatorID)

	if err := stmt.First(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	shop, err := internal.entity()
	if err != nil {
		return nil, dbError(err)
	}

	if err := s.fill(ctx, s.db.DB, shop); err != nil {
		return nil, dbError(err)
	}
	return shop, nil
}

func (s *shop) Create(ctx context.Context, shop *entity.Shop) error {
	now := s.now()
	shop.CreatedAt, shop.UpdatedAt = now, now

	internal, err := newInternalShop(shop)
	if err != nil {
		return err
	}

	err = s.db.DB.WithContext(ctx).Table(shopTable).Create(&internal).Error
	return dbError(err)
}

func (s *shop) Update(ctx context.Context, shopID string, params *database.UpdateShopParams) error {
	productTypeIDs, err := json.Marshal(params.ProductTypeIDs)
	if err != nil {
		return fmt.Errorf("tidb: failed to marshal product type ids: %w", err)
	}
	businessDays, err := json.Marshal(params.BusinessDays)
	if err != nil {
		return fmt.Errorf("tidb: failed to marshal business days: %w", err)
	}

	updates := map[string]interface{}{
		"name":             params.Name,
		"product_type_ids": productTypeIDs,
		"business_days":    businessDays,
		"updated_at":       s.now(),
	}
	stmt := s.db.DB.WithContext(ctx).Table(shopTable).Where("id = ?", shopID)

	err = stmt.Updates(updates).Error
	return dbError(err)
}

func (s *shop) Delete(ctx context.Context, shopID string) error {
	updates := map[string]interface{}{
		"activated":  false,
		"updated_at": s.now(),
		"deleted_at": s.now(),
	}
	stmt := s.db.DB.WithContext(ctx).Table(shopTable).Where("id = ?", shopID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (s *shop) RemoveProductType(ctx context.Context, productTypeID string) error {
	sub := gorm.Expr("JSON_REMOVE(product_type_ids, JSON_UNQUOTE(JSON_SEARCH(product_type_ids, 'one', ?)))", productTypeID)

	stmt := s.db.DB.WithContext(ctx).
		Table(shopTable).
		Where("JSON_SEARCH(product_type_ids, 'one', ?) IS NOT NULL", productTypeID)

	err := stmt.Update("product_type_ids", sub).Error
	return dbError(err)
}

func (s *shop) RelateProducer(ctx context.Context, shopID, producerID string) error {
	now := s.now()
	producer := &entity.ShopProducer{
		ShopID:     shopID,
		ProducerID: producerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	updates := map[string]interface{}{
		"updated_at": s.now(),
	}
	stmt := s.db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "shop_id"}, {Name: "producer_id"}},
		DoUpdates: clause.Assignments(updates),
	})

	err := stmt.Create(&producer).Error
	return dbError(err)
}

func (s *shop) UnrelateProducer(ctx context.Context, shopID, producerID string) error {
	stmt := s.db.DB.WithContext(ctx).Where("shop_id = ? AND producer_id = ?", shopID, producerID)

	err := stmt.Delete(&entity.ShopProducer{}).Error
	return dbError(err)
}

func (s *shop) get(ctx context.Context, tx *gorm.DB, shopID string, fields ...string) (*entity.Shop, error) {
	var internal *internalShop

	stmt := s.db.Statement(ctx, tx, shopTable, fields...).Where("id = ?", shopID)

	if err := stmt.First(&internal).Error; err != nil {
		return nil, err
	}
	shop, err := internal.entity()
	if err != nil {
		return nil, err
	}

	if err := s.fill(ctx, tx, shop); err != nil {
		return nil, err
	}
	return shop, nil
}

func (s *shop) fill(ctx context.Context, tx *gorm.DB, shops ...*entity.Shop) error {
	var producers entity.ShopProducers

	ids := entity.Shops(shops).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := s.db.Statement(ctx, tx, shopProducerTable).Where("shop_id IN ?", ids)

	if err := stmt.Find(&producers).Error; err != nil {
		return err
	}
	if len(producers) == 0 {
		return nil
	}
	entity.Shops(shops).Fill(producers.GroupByShopID())
	return nil
}

type internalShop struct {
	entity.Shop        `gorm:"embedded"`
	ProductTypeIDsJSON datatypes.JSON `gorm:"default:null;column:product_type_ids"` // 取り扱い商品種別ID一覧
	BusinessDaysJSON   datatypes.JSON `gorm:"default:null;column:business_days"`    // 営業曜日(発送可能日)一覧(JSON)
}

type internalShops []*internalShop

func newInternalShop(shop *entity.Shop) (*internalShop, error) {
	productTypeIDs, err := json.Marshal(shop.ProductTypeIDs)
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal product type ids: %w", err)
	}
	businessDays, err := json.Marshal(shop.BusinessDays)
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal business days: %w", err)
	}
	internal := &internalShop{
		Shop: entity.Shop{
			ID:            shop.ID,
			CoordinatorID: shop.CoordinatorID,
			Name:          shop.Name,
			Activated:     shop.Activated,
			CreatedAt:     shop.CreatedAt,
			UpdatedAt:     shop.UpdatedAt,
			DeletedAt:     shop.DeletedAt,
		},
		ProductTypeIDsJSON: productTypeIDs,
		BusinessDaysJSON:   businessDays,
	}
	return internal, nil
}

func (s *internalShop) entity() (*entity.Shop, error) {
	if err := s.unmarshalProductTypeIDs(); err != nil {
		return nil, err
	}
	if err := s.unmarshalBusinessDays(); err != nil {
		return nil, err
	}
	return &s.Shop, nil
}

func (s *internalShop) unmarshalProductTypeIDs() error {
	if s == nil || s.ProductTypeIDsJSON == nil {
		return nil
	}
	var productTypeIDs []string
	if err := json.Unmarshal(s.ProductTypeIDsJSON, &productTypeIDs); err != nil {
		return fmt.Errorf("tidb: failed to unmarshal product type ids: %w", err)
	}
	s.ProductTypeIDs = productTypeIDs
	return nil
}

func (s *internalShop) unmarshalBusinessDays() error {
	if s == nil || s.BusinessDaysJSON == nil {
		return nil
	}
	var businessDays []time.Weekday
	if err := json.Unmarshal(s.BusinessDaysJSON, &businessDays); err != nil {
		return fmt.Errorf("tidb: failed to unmarshal business days: %w", err)
	}
	s.BusinessDays = businessDays
	return nil
}

func (ss internalShops) entities() (entity.Shops, error) {
	res := make(entity.Shops, len(ss))
	for i := range ss {
		s, err := ss[i].entity()
		if err != nil {
			return nil, err
		}
		res[i] = s
	}
	return res, nil
}
