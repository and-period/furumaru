package mysql

import (
	"context"
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
	productTable         = "products"
	productRevisionTable = "product_revisions"
)

type product struct {
	db  *mysql.Client
	now func() time.Time
}

func newProduct(db *mysql.Client) database.Product {
	return &product{
		db:  db,
		now: jst.Now,
	}
}

type listProductsParams database.ListProductsParams

func (p listProductsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("name LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
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
	if p.OnlyPublished {
		stmt = stmt.Where("public = ?", true)
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
	var products entity.Products

	prm := listProductsParams(*params)

	stmt := p.db.Statement(ctx, p.db.DB, productTable, fields...)
	stmt = prm.stmt(stmt)
	stmt = prm.pagination(stmt)

	if err := stmt.Find(&products).Error; err != nil {
		return nil, dbError(err)
	}
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
		if err := product.FillJSON(); err != nil {
			return err
		}

		now := p.now()

		product.CreatedAt, product.UpdatedAt = now, now
		product.ProductRevision.CreatedAt, product.ProductRevision.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(productTable).Create(&product).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Table(productRevisionTable).Create(&product.ProductRevision).Error
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
		tagIDs, err := entity.ProductMarshalTagIDs(params.TagIDs)
		if err != nil {
			return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
		}
		points, err := entity.ProductMarshalRecommendedPoints(params.RecommendedPoints)
		if err != nil {
			return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
		}

		updates := map[string]interface{}{
			"product_type_id":     params.TypeID,
			"product_tag_ids":     tagIDs,
			"name":                params.Name,
			"description":         params.Description,
			"recommended_points":  points,
			"public":              params.Public,
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
			media, err := params.Media.Marshal()
			if err != nil {
				return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
			}
			updates["media"] = media
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

func (p *product) UpdateMedia(
	ctx context.Context, productID string, set func(media entity.MultiProductMedia) bool,
) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		product, err := p.get(ctx, tx, productID, "id", "media")
		if err != nil {
			return err
		}
		if exists := set(product.Media); !exists {
			return fmt.Errorf("database: media is non-existent: %w", database.ErrFailedPrecondition)
		}

		buf, err := product.Media.Marshal()
		if err != nil {
			return err
		}
		params := map[string]interface{}{
			"media":      datatypes.JSON(buf),
			"updated_at": p.now(),
		}

		err = tx.WithContext(ctx).
			Table(productTable).
			Where("id = ?", productID).
			Updates(params).Error
		return err
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
	var products entity.Products

	stmt := p.db.Statement(ctx, tx, productTable, fields...).
		Where("id IN (?)", productIDs)

	if err := stmt.Find(&products).Error; err != nil {
		return nil, err
	}
	if err := p.fill(ctx, tx, products...); err != nil {
		return nil, err
	}
	return products, nil
}

func (p *product) get(ctx context.Context, tx *gorm.DB, productID string, fields ...string) (*entity.Product, error) {
	var product *entity.Product

	stmt := p.db.Statement(ctx, tx, productTable, fields...).
		Where("id = ?", productID)

	if err := stmt.First(&product).Error; err != nil {
		return nil, err
	}
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
		return nil
	}
	if len(revisions) == 0 {
		return nil
	}
	return entity.Products(products).Fill(revisions.MapByProductID(), p.now())
}
