package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const liveTable = "lives"

type live struct {
	db  *mysql.Client
	now func() time.Time
}

func newLive(db *mysql.Client) database.Live {
	return &live{
		db:  db,
		now: jst.Now,
	}
}

type listLivesParams database.ListLivesParams

func (p listLivesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if len(p.ScheduleIDs) > 0 {
		stmt = stmt.Where("schedule_id IN (?)", p.ScheduleIDs)
	}
	if p.ProducerID != "" {
		stmt = stmt.Where("producer_id = ?", p.ProducerID)
	}
	stmt = stmt.Order("start_at ASC, end_at ASC")
	return stmt
}

func (p listLivesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (l *live) List(ctx context.Context, params *database.ListLivesParams, fields ...string) (entity.Lives, error) {
	var lives entity.Lives

	p := listLivesParams(*params)

	stmt := l.db.Statement(ctx, l.db.DB, liveTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&lives).Error; err != nil {
		return nil, dbError(err)
	}
	if err := l.fill(ctx, l.db.DB, lives...); err != nil {
		return nil, dbError(err)
	}
	return lives, nil
}

func (l *live) Count(ctx context.Context, params *database.ListLivesParams) (int64, error) {
	p := listLivesParams(*params)

	total, err := l.db.Count(ctx, l.db.DB, &entity.Live{}, p.stmt)
	return total, dbError(err)
}

func (l *live) Get(ctx context.Context, liveID string, fields ...string) (*entity.Live, error) {
	live, err := l.get(ctx, l.db.DB, liveID, fields...)
	return live, dbError(err)
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
	return dbError(err)
}

func (l *live) Update(ctx context.Context, liveID string, params *database.UpdateLiveParams) error {
	err := l.db.Transaction(ctx, func(tx *gorm.DB) error {
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
	return dbError(err)
}

func (l *live) Delete(ctx context.Context, liveID string) error {
	stmt := l.db.DB.WithContext(ctx).
		Table(liveTable).
		Where("id = ?", liveID)

	err := stmt.Delete(&entity.Live{}).Error
	return dbError(err)
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
