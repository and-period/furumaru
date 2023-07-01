package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const promotionTable = "promotions"

type promotion struct {
	db  *database.Client
	now func() time.Time
}

func NewPromotion(db *database.Client) Promotion {
	return &promotion{
		db:  db,
		now: jst.Now,
	}
}

func (p *promotion) List(ctx context.Context, params *ListPromotionsParams, fields ...string) (entity.Promotions, error) {
	var promotions entity.Promotions

	stmt := p.db.Statement(ctx, p.db.DB, promotionTable, fields...)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&promotions).Error
	return promotions, exception.InternalError(err)
}

func (p *promotion) Count(ctx context.Context, params *ListPromotionsParams) (int64, error) {
	total, err := p.db.Count(ctx, p.db.DB, &entity.Promotion{}, params.stmt)
	return total, exception.InternalError(err)
}

func (p *promotion) MultiGet(ctx context.Context, promotionIDs []string, fields ...string) (entity.Promotions, error) {
	var promotions entity.Promotions

	stmt := p.db.Statement(ctx, p.db.DB, promotionTable, fields...).
		Where("id IN (?)", promotionIDs)

	err := stmt.Find(&promotions).Error
	return promotions, exception.InternalError(err)
}

func (p *promotion) Get(ctx context.Context, promotionID string, fields ...string) (*entity.Promotion, error) {
	promotion, err := p.get(ctx, p.db.DB, promotionID, fields...)
	return promotion, exception.InternalError(err)
}

func (p *promotion) Create(ctx context.Context, promotion *entity.Promotion) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := p.now()
		promotion.CreatedAt, promotion.UpdatedAt = now, now

		return tx.WithContext(ctx).Table(promotionTable).Create(&promotion).Error
	})
	return exception.InternalError(err)
}

func (p *promotion) Update(ctx context.Context, promotionID string, params *UpdatePromotionParams) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := p.get(ctx, tx, promotionID); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"title":         params.Title,
			"description":   params.Description,
			"public":        params.Public,
			"published_at":  params.PublishedAt,
			"discount_type": params.DiscountType,
			"discount_rate": params.DiscountRate,
			"code":          params.Code,
			"code_type":     params.CodeType,
			"start_at":      params.StartAt,
			"end_at":        params.EndAt,
			"updated_at":    p.now(),
		}
		err := tx.WithContext(ctx).
			Table(promotionTable).
			Where("id = ?", promotionID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
}

func (p *promotion) Delete(ctx context.Context, promotionID string) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := p.get(ctx, tx, promotionID); err != nil {
			return err
		}

		err := tx.WithContext(ctx).
			Table(promotionTable).
			Where("id = ?", promotionID).
			Delete(&entity.Promotion{}).Error
		return err
	})
	return exception.InternalError(err)
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
	return promotion, nil
}
