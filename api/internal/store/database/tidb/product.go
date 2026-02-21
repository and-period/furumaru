package tidb

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
	productTable         = "products"
	productRevisionTable = "product_revisions"
)

type product struct {
	db  *mysql.Client
	now func() time.Time
}

func NewProduct(db *mysql.Client) database.Product {
	return &product{
		db:  db,
		now: jst.Now,
	}
}

type listProductsParams database.ListProductsParams

func (p listProductsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", p.Name)).
			Or("`description` LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.ShopID != "" {
		stmt = stmt.Where("shop_id = ?", p.ShopID)
	}
	if p.ProducerID != "" {
		stmt = stmt.Where("producer_id = ?", p.ProducerID)
	}
	if len(p.ProducerIDs) > 0 {
		stmt = stmt.Where("producer_id IN (?)", p.ProducerIDs)
	}
	if len(p.ProductTypeIDs) > 0 {
		stmt = stmt.Where("product_type_id IN (?)", p.ProductTypeIDs)
	}
	if p.ProductTagID != "" {
		stmt = stmt.Where("JSON_SEARCH(product_tag_ids, 'all', ?) IS NOT NULL", p.ProductTagID)
	}
	if len(p.Scopes) > 0 {
		stmt = stmt.Where("scope IN (?)", p.Scopes)
	}
	if !p.EndAtGte.IsZero() {
		stmt = stmt.Where("end_at >= ?", p.EndAtGte)
	}
	if !p.ExcludeDeleted {
		stmt = stmt.Unscoped()
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	if len(p.Orders) == 0 {
		stmt = stmt.Order("start_at DESC")
	}
	return stmt
}

func (p listProductsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (p *product) List(ctx context.Context, params *database.ListProductsParams, fields ...string) (entity.Products, error) {
	var internal internalProducts

	prm := listProductsParams(*params)

	stmt := p.db.Statement(ctx, p.db.DB, productTable, fields...)
	stmt = prm.stmt(stmt)
	stmt = prm.pagination(stmt)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	products := internal.entities()

	if err := p.fill(ctx, p.db.DB, products...); err != nil {
		return nil, dbError(err)
	}
	return products, nil
}

func (p *product) Count(ctx context.Context, params *database.ListProductsParams) (int64, error) {
	prm := listProductsParams(*params)

	total, err := p.db.Count(ctx, p.db.DB, &entity.Product{}, prm.stmt)
	return total, dbError(err)
}

func (p *product) MultiGet(ctx context.Context, productIDs []string, fields ...string) (entity.Products, error) {
	products, err := p.multiGet(ctx, p.db.DB, productIDs, fields...)
	return products, dbError(err)
}

func (p *product) MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Products, error) {
	var revisions entity.ProductRevisions

	stmt := p.db.Statement(ctx, p.db.DB, productRevisionTable).
		Where("id IN (?)", revisionIDs)

	if err := stmt.Find(&revisions).Error; err != nil {
		return nil, dbError(err)
	}
	if len(revisions) == 0 {
		return entity.Products{}, nil
	}

	products, err := p.multiGet(ctx, p.db.DB, revisions.ProductIDs(), fields...)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return entity.Products{}, nil
	}

	res, err := revisions.Merge(products.Map())
	if err != nil {
		return nil, dbError(err)
	}
	return res, nil
}

func (p *product) Get(ctx context.Context, productID string, fields ...string) (*entity.Product, error) {
	product, err := p.get(ctx, p.db.DB, productID, fields...)
	return product, dbError(err)
}

func (p *product) Create(ctx context.Context, product *entity.Product) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := p.now()

		product.CreatedAt, product.UpdatedAt = now, now
		product.ProductRevision.CreatedAt, product.ProductRevision.UpdatedAt = now, now

		internal := newInternalProduct(product)

		if err := tx.WithContext(ctx).Table(productTable).Create(&internal).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Table(productRevisionTable).Create(&internal.ProductRevision).Error
	})
	return dbError(err)
}

func (p *product) Update(ctx context.Context, productID string, params *database.UpdateProductParams) error {
	now := p.now()
	rparams := &entity.NewProductRevisionParams{
		ProductID: productID,
		Price:     params.Price,
		Cost:      params.Cost,
	}
	revision := entity.NewProductRevision(rparams)

	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"product_type_id":     params.TypeID,
			"product_tag_ids":     mysql.NewJSONColumn(params.TagIDs),
			"name":                params.Name,
			"description":         params.Description,
			"media":               nil,
			"recommended_points":  mysql.NewJSONColumn(params.RecommendedPoints),
			"scope":               params.Scope,
			"inventory":           params.Inventory,
			"weight":              params.Weight,
			"weight_unit":         params.WeightUnit,
			"item":                params.Item,
			"item_unit":           params.ItemUnit,
			"item_description":    params.ItemDescription,
			"expiration_date":     params.ExpirationDate,
			"storage_method_type": params.StorageMethodType,
			"delivery_type":       params.DeliveryType,
			"box60_rate":          params.Box60Rate,
			"box80_rate":          params.Box80Rate,
			"box100_rate":         params.Box100Rate,
			"origin_prefecture":   params.OriginPrefectureCode,
			"origin_city":         params.OriginCity,
			"start_at":            params.StartAt,
			"end_at":              params.EndAt,
			"updated_at":          p.now(),
		}
		if len(params.Media) > 0 {
			updates["media"] = mysql.NewJSONColumn(params.Media)
		}

		stmt := tx.WithContext(ctx).Table(productTable).Where("id = ?", productID)
		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}

		revision.CreatedAt, revision.UpdatedAt = now, now
		return tx.WithContext(ctx).Table(productRevisionTable).Create(&revision).Error
	})
	return dbError(err)
}

func (p *product) DecreaseInventory(ctx context.Context, revisionID, quantity int64) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		var product *entity.Product

		stmt := p.db.Statement(ctx, tx, productTable, "products.*").
			Joins("INNER JOIN product_revisions ON products.id = product_revisions.product_id").
			Where("product_revisions.id = ?", revisionID)
		if err := stmt.First(&product).Error; err != nil {
			return err
		}
		if product.Inventory == 0 {
			return nil
		}

		inventory := product.Inventory - quantity
		if inventory < 0 {
			inventory = 0
		}
		params := map[string]interface{}{
			"inventory":  inventory,
			"updated_at": p.now(),
		}

		err := tx.WithContext(ctx).
			Table(productTable).
			Where("id = ?", product.ID).
			Updates(params).Error
		return err
	})
	return dbError(err)
}

func (p *product) Delete(ctx context.Context, productID string) error {
	params := map[string]interface{}{
		"deleted_at": p.now(),
	}
	stmt := p.db.DB.WithContext(ctx).
		Table(productTable).
		Where("id = ?", productID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (p *product) multiGet(ctx context.Context, tx *gorm.DB, productIDs []string, fields ...string) (entity.Products, error) {
	var internal internalProducts

	stmt := p.db.Statement(ctx, tx, productTable, fields...).Unscoped().Where("id IN (?)", productIDs)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, err
	}
	products := internal.entities()

	if err := p.fill(ctx, tx, products...); err != nil {
		return nil, err
	}
	return products, nil
}

func (p *product) get(ctx context.Context, tx *gorm.DB, productID string, fields ...string) (*entity.Product, error) {
	var internal *internalProduct

	stmt := p.db.Statement(ctx, tx, productTable, fields...).Unscoped().Where("id = ?", productID)

	if err := stmt.First(&internal).Error; err != nil {
		return nil, err
	}
	product := internal.entity()

	if err := p.fill(ctx, tx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (p *product) fill(ctx context.Context, tx *gorm.DB, products ...*entity.Product) error {
	var revisions entity.ProductRevisions

	ids := entity.Products(products).IDs()
	if len(ids) == 0 {
		return nil
	}

	sub := tx.Table(productRevisionTable).
		Select("MAX(id)").
		Where("product_id IN (?)", ids).
		Group("product_id")
	stmt := p.db.Statement(ctx, tx, productRevisionTable).
		Where("id IN (?)", sub)

	if err := stmt.Find(&revisions).Error; err != nil {
		return err
	}
	if len(revisions) == 0 {
		return nil
	}
	entity.Products(products).Fill(revisions.MapByProductID(), p.now())
	return nil
}

type internalProduct struct {
	entity.Product        `gorm:"embedded"`
	TagIDsJSON            mysql.JSONColumn[[]string]                  `gorm:"default:null;column:product_tag_ids"`    // 商品タグID一覧(JSON)
	MediaJSON             mysql.JSONColumn[entity.MultiProductMedia]  `gorm:"default:null;column:media"`              // メディア一覧(JSON)
	RecommendedPointsJSON mysql.JSONColumn[[]string]                  `gorm:"default:null;column:recommended_points"` // おすすめポイント一覧(JSON)
}

type internalProducts []*internalProduct

func newInternalProduct(product *entity.Product) *internalProduct {
	return &internalProduct{
		Product:               *product,
		TagIDsJSON:            mysql.NewJSONColumn(product.TagIDs),
		MediaJSON:             mysql.NewJSONColumn(product.Media),
		RecommendedPointsJSON: mysql.NewJSONColumn(product.RecommendedPoints),
	}
}

func (p *internalProduct) entity() *entity.Product {
	e := p.Product
	e.TagIDs = p.TagIDsJSON.Val
	e.Media = p.MediaJSON.Val
	e.RecommendedPoints = p.RecommendedPointsJSON.Val
	return &e
}



func (ps internalProducts) entities() entity.Products {
	res := make(entity.Products, len(ps))
	for i := range ps {
		res[i] = ps[i].entity()
	}
	return res
}
