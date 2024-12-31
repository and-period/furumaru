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

const productTypeTable = "product_types"

type productType struct {
	database.ProductType
	db  *mysql.Client
	now func() time.Time
}

func NewProductType(db *mysql.Client, mysql database.ProductType) database.ProductType {
	return &productType{
		ProductType: mysql,
		db:          db,
		now:         jst.Now,
	}
}

type listProductTypesParams database.ListProductTypesParams

func (p listProductTypesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.CategoryID != "" {
		stmt = stmt.Where("category_id = ?", p.CategoryID)
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
		stmt = stmt.Order("category_id ASC, name ASC")
	}
	return stmt
}

func (p listProductTypesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (t *productType) List(
	ctx context.Context, params *database.ListProductTypesParams, fields ...string,
) (entity.ProductTypes, error) {
	var productTypes entity.ProductTypes

	p := listProductTypesParams(*params)

	stmt := t.db.Statement(ctx, t.db.DB, productTypeTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&productTypes).Error; err != nil {
		return nil, dbError(err)
	}
	return productTypes, nil
}

func (t *productType) Count(ctx context.Context, params *database.ListProductTypesParams) (int64, error) {
	p := listProductTypesParams(*params)

	total, err := t.db.Count(ctx, t.db.DB, &entity.ProductType{}, p.stmt)
	return total, dbError(err)
}
