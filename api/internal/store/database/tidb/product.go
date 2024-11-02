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
	database.Product
	db  *mysql.Client
	now func() time.Time
}

func newProduct(db *mysql.Client, mysql database.Product) database.Product {
	return &product{
		Product: mysql,
		db:      db,
		now:     jst.Now,
	}
}

type listProductsParams database.ListProductsParams

func (p listProductsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", p.Name)).
			Or("`description` LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
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
		stmt = stmt.Where("public = ?", true).Where("deleted_at IS NULL")
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
	return entity.Products(products).Fill(revisions.MapByProductID(), p.now())
}
