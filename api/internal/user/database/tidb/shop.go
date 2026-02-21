package tidb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
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

type listShopsParams database.ListShopsParams

func (p listShopsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if len(p.CoordinatorIDs) > 0 {
		stmt = stmt.Where("coordinator_id IN (?)", p.CoordinatorIDs)
	}
	if len(p.ProducerIDs) > 0 {
		stmt = stmt.Where("id IN (SELECT DISTINCT(shop_id) FROM shop_producers WHERE producer_id IN (?))", p.ProducerIDs)
	}
	return stmt
}

func (p listShopsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (s *shop) List(ctx context.Context, params *database.ListShopsParams, fields ...string) (entity.Shops, error) {
	var internal internalShops

	p := listShopsParams(*params)

	stmt := s.db.Statement(ctx, s.db.DB, shopTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	shops := internal.entities()

	if err := s.fill(ctx, s.db.DB, shops...); err != nil {
		return nil, dbError(err)
	}
	return shops, nil
}

func (s *shop) Count(ctx context.Context, params *database.ListShopsParams) (int64, error) {
	p := listShopsParams(*params)

	total, err := s.db.Count(ctx, s.db.DB, &entity.Shop{}, p.stmt)
	return total, dbError(err)
}

func (s *shop) MultiGet(ctx context.Context, shopIDs []string, fields ...string) (entity.Shops, error) {
	var internal internalShops

	stmt := s.db.Statement(ctx, s.db.DB, shopTable, fields...).Where("id IN (?)", shopIDs)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	shops := internal.entities()

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
	shop := internal.entity()

	if err := s.fill(ctx, s.db.DB, shop); err != nil {
		return nil, dbError(err)
	}
	return shop, nil
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

func (s *shop) RemoveProductType(ctx context.Context, productTypeID string) error {
	sub := gorm.Expr("JSON_REMOVE(product_type_ids, JSON_UNQUOTE(JSON_SEARCH(product_type_ids, 'one', ?)))", productTypeID)

	stmt := s.db.DB.WithContext(ctx).
		Table(shopTable).
		Where("JSON_SEARCH(product_type_ids, 'one', ?) IS NOT NULL", productTypeID)

	err := stmt.Update("product_type_ids", sub).Error
	return dbError(err)
}

func (s *shop) ListProducers(ctx context.Context, params *database.ListShopProducersParams) ([]string, error) {
	var producerIDs []string

	fields := []string{
		"DISTINCT(producer_id) AS producer_id",
	}
	stmt := s.db.Statement(ctx, s.db.DB, shopProducerTable, fields...).
		Where("shop_id = ?", params.ShopID).
		Order("producer_id")

	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Scan(&producerIDs).Error
	return producerIDs, dbError(err)
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
	shop := internal.entity()

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
	ProductTypeIDsJSON mysql.JSONColumn[[]string]       `gorm:"default:null;column:product_type_ids"` // 取り扱い商品種別ID一覧
	BusinessDaysJSON   mysql.JSONColumn[[]time.Weekday] `gorm:"default:null;column:business_days"`    // 営業曜日(発送可能日)一覧(JSON)
}

type internalShops []*internalShop

func newInternalShop(shop *entity.Shop) *internalShop {
	return &internalShop{
		Shop: entity.Shop{
			ID:            shop.ID,
			CoordinatorID: shop.CoordinatorID,
			Name:          shop.Name,
			Activated:     shop.Activated,
			CreatedAt:     shop.CreatedAt,
			UpdatedAt:     shop.UpdatedAt,
			DeletedAt:     shop.DeletedAt,
		},
		ProductTypeIDsJSON: mysql.NewJSONColumn(shop.ProductTypeIDs),
		BusinessDaysJSON:   mysql.NewJSONColumn(shop.BusinessDays),
	}
}

func (s *internalShop) entity() *entity.Shop {
	s.Shop.ProductTypeIDs = s.ProductTypeIDsJSON.Val
	s.Shop.BusinessDays = s.BusinessDaysJSON.Val
	return &s.Shop
}

func (ss internalShops) entities() entity.Shops {
	res := make(entity.Shops, len(ss))
	for i := range ss {
		res[i] = ss[i].entity()
	}
	return res
}
