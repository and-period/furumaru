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

const promotionTable = "promotions"

type promotion struct {
	db  *mysql.Client
	now func() time.Time
}

func newPromotion(db *mysql.Client) database.Promotion {
	return &promotion{
		db:  db,
		now: jst.Now,
	}
}

type listPromotionsParams database.ListPromotionsParams

func (p listPromotionsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Title != "" {
		stmt = stmt.Where("MATCH (`title`) AGAINST (? IN NATURAL LANGUAGE MODE)", p.Title)
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

func (p *promotion) MultiGet(ctx context.Context, promotionIDs []string, fields ...string) (entity.Promotions, error) {
	var promotions entity.Promotions

	stmt := p.db.Statement(ctx, p.db.DB, promotionTable, fields...).
		Where("id IN (?)", promotionIDs)

	if err := stmt.Find(&promotions).Error; err != nil {
		return nil, dbError(err)
	}
	promotions.Fill(p.now())
	return promotions, nil
}

func (p *promotion) Get(ctx context.Context, promotionID string, fields ...string) (*entity.Promotion, error) {
	promotion, err := p.get(ctx, p.db.DB, promotionID, fields...)
	return promotion, dbError(err)
}

func (p *promotion) GetByCode(ctx context.Context, code string, fields ...string) (*entity.Promotion, error) {
	var promotion *entity.Promotion

	stmt := p.db.Statement(ctx, p.db.DB, promotionTable, fields...).
		Where("code = ?", code)

	if err := stmt.First(&promotion).Error; err != nil {
		return nil, dbError(err)
	}
	promotion.Fill(p.now())
	return promotion, nil
}

func (p *promotion) Create(ctx context.Context, promotion *entity.Promotion) error {
	now := p.now()
	promotion.CreatedAt, promotion.UpdatedAt = now, now

	err := p.db.DB.WithContext(ctx).Table(promotionTable).Create(&promotion).Error
	return dbError(err)
}

func (p *promotion) Update(ctx context.Context, promotionID string, params *database.UpdatePromotionParams) error {
	updates := map[string]interface{}{
		"title":         params.Title,
		"description":   params.Description,
		"public":        params.Public,
		"discount_type": params.DiscountType,
		"discount_rate": params.DiscountRate,
		"code":          params.Code,
		"code_type":     params.CodeType,
		"start_at":      params.StartAt,
		"end_at":        params.EndAt,
		"updated_at":    p.now(),
	}
	stmt := p.db.DB.WithContext(ctx).
		Table(promotionTable).
		Where("id = ?", promotionID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (p *promotion) Delete(ctx context.Context, promotionID string) error {
	stmt := p.db.DB.WithContext(ctx).
		Table(promotionTable).
		Where("id = ?", promotionID)

	err := stmt.Delete(&entity.Promotion{}).Error
	return dbError(err)
}

func (p *promotion) get(
	ctx context.Context, tx *gorm.DB, promotionID string, fields ...string,
) (*entity.Promotion, error) {
	var promotion *entity.Promotion

	stmt := p.db.Statement(ctx, tx, promotionTable, fields...).
		Where("id = ?", promotionID)

	if err := stmt.First(&promotion).Error; err != nil {
		return nil, err
	}
	promotion.Fill(p.now())
	return promotion, nil
}
