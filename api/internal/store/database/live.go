package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const liveTable = "lives"

type live struct {
	db  *database.Client
	now func() time.Time
}

func NewLive(db *database.Client) Live {
	return &live{
		db:  db,
		now: jst.Now,
	}
}

func (l *live) ListByScheduleID(ctx context.Context, scheduleID string, fields ...string) (entity.Lives, error) {
	var lives entity.Lives

	stmt := l.db.Statement(ctx, l.db.DB, liveTable, fields...).Where("schedule_id = ?", scheduleID)

	if err := stmt.Find(&lives).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := l.fill(ctx, l.db.DB, lives...); err != nil {
		return nil, exception.InternalError(err)
	}
	return lives, nil
}

func (l *live) Get(ctx context.Context, liveID string, fields ...string) (*entity.Live, error) {
	live, err := l.get(ctx, l.db.DB, liveID, fields...)
	return live, exception.InternalError(err)
}

func (l *live) Create(ctx context.Context, live *entity.Live) error {
	err := l.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := l.now()
		live.CreatedAt, live.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(liveTable).Create(&live).Error; err != nil {
			return err
		}

		products := entity.NewLiveProducts(live.ID, live.ProductIDs)
		return l.replaceProducts(ctx, tx, live.ID, products)
	})
	return exception.InternalError(err)
}

func (l *live) Update(ctx context.Context, liveID string, params *UpdateLiveParams) error {
	err := l.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.get(ctx, tx, liveID); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"comment":    params.Comment,
			"start_at":   params.StartAt,
			"end_at":     params.EndAt,
			"updated_at": l.now(),
		}

		err := tx.WithContext(ctx).
			Table(liveTable).
			Where("id = ?", liveID).
			Updates(updates).Error
		if err != nil {
			return err
		}

		products := entity.NewLiveProducts(liveID, params.ProductIDs)
		return l.replaceProducts(ctx, tx, liveID, products)
	})
	return exception.InternalError(err)
}

func (l *live) Delete(ctx context.Context, liveID string) error {
	err := l.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := l.get(ctx, tx, liveID); err != nil {
			return err
		}

		err := tx.WithContext(ctx).
			Table(liveTable).
			Where("id = ?", liveID).
			Delete(&entity.Live{}).Error
		return err
	})
	return exception.InternalError(err)
}

func (l *live) get(ctx context.Context, tx *gorm.DB, liveID string, fields ...string) (*entity.Live, error) {
	var live *entity.Live

	stmt := l.db.Statement(ctx, tx, liveTable, fields...).Where("id = ?", liveID)

	if err := stmt.First(&live).Error; err != nil {
		return nil, err
	}
	if err := l.fill(ctx, tx, live); err != nil {
		return nil, err
	}
	return live, nil
}

func (l *live) fill(ctx context.Context, tx *gorm.DB, lives ...*entity.Live) error {
	var products entity.LiveProducts

	ids := entity.Lives(lives).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := l.db.Statement(ctx, tx, liveProductTable).Where("live_id IN (?)", ids)
	if err := stmt.Find(&products).Error; err != nil {
		return err
	}
	entity.Lives(lives).Fill(products.GroupByLiveID())
	return nil
}

func (l *live) replaceProducts(ctx context.Context, tx *gorm.DB, liveID string, products entity.LiveProducts) error {
	// 不要なレコードを削除
	stmt := tx.WithContext(ctx).
		Where("live_id = ?", liveID).
		Where("product_id NOT IN (?)", products.ProductIDs())
	if err := stmt.Delete(&entity.LiveProduct{}).Error; err != nil {
		return err
	}

	// レコードの登録/更新
	if len(products) == 0 {
		return nil
	}
	for _, product := range products {
		params := map[string]interface{}{
			"updated_at": l.now(),
		}
		conds := clause.OnConflict{
			Columns:   []clause.Column{{Name: "live_id"}, {Name: "product_id"}},
			DoUpdates: clause.Assignments(params),
		}
		if err := tx.WithContext(ctx).Omit(clause.Associations).Clauses(conds).Create(&product).Error; err != nil {
			return err
		}
	}
	return nil
}
