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

const promotionTable = "promotions"

type promotion struct {
	database.Promotion
	db  *mysql.Client
	now func() time.Time
}

func NewPromotion(db *mysql.Client, mysql database.Promotion) database.Promotion {
	return &promotion{
		Promotion: mysql,
		db:        db,
		now:       jst.Now,
	}
}

type listPromotionsParams database.ListPromotionsParams

func (p listPromotionsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Title != "" {
		stmt = stmt.Where("`title` LIKE ?", fmt.Sprintf("%%%s%%", p.Title))
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

func (p listPromotionsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (p *promotion) List(ctx context.Context, params *database.ListPromotionsParams, fields ...string) (entity.Promotions, error) {
	var promotions entity.Promotions

	prm := listPromotionsParams(*params)

	stmt := p.db.Statement(ctx, p.db.DB, promotionTable, fields...)
	stmt = prm.stmt(stmt)
	stmt = prm.pagination(stmt)

	if err := stmt.Find(&promotions).Error; err != nil {
		return nil, dbError(err)
	}
	promotions.Fill(p.now())
	return promotions, nil
}

func (p *promotion) Count(ctx context.Context, params *database.ListPromotionsParams) (int64, error) {
	prm := listPromotionsParams(*params)

	total, err := p.db.Count(ctx, p.db.DB, &entity.Promotion{}, prm.stmt)
	return total, dbError(err)
}
