package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const producerTable = "producers"

type producer struct {
	db  *mysql.Client
	now func() time.Time
}

func newProducer(db *mysql.Client) database.Producer {
	return &producer{
		db:  db,
		now: jst.Now,
	}
}

type listProducersParams database.ListProducersParams

func (p listProducersParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.OnlyUnrelated {
		stmt = stmt.Where("coordinator_id IS NULL")
	}
	if p.Username != "" {
		stmt = stmt.Where("username LIKE ?", fmt.Sprintf("%%%s%%", p.Username))
	}
	stmt = stmt.Order("updated_at DESC")
	return stmt
}

func (p listProducersParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (p *producer) List(
	ctx context.Context, params *database.ListProducersParams, fields ...string,
) (entity.Producers, error) {
	var producers entity.Producers

	prm := listProducersParams(*params)

	stmt := p.db.Statement(ctx, p.db.DB, producerTable, fields...)
	stmt = prm.stmt(stmt)
	stmt = prm.pagination(stmt)

	if err := stmt.Find(&producers).Error; err != nil {
		return nil, dbError(err)
	}
	if err := p.fill(ctx, p.db.DB, producers...); err != nil {
		return nil, dbError(err)
	}
	return producers, nil
}

func (p *producer) Count(ctx context.Context, params *database.ListProducersParams) (int64, error) {
	prm := listProducersParams(*params)

	total, err := p.db.Count(ctx, p.db.DB, &entity.Producer{}, prm.stmt)
	return total, dbError(err)
}

func (p *producer) MultiGet(
	ctx context.Context, producerIDs []string, fields ...string,
) (entity.Producers, error) {
	producers, err := p.multiGet(ctx, p.db.DB, producerIDs, fields...)
	if err != nil {
		return nil, dbError(err)
	}
	if err := p.fill(ctx, p.db.DB, producers...); err != nil {
		return nil, dbError(err)
	}
	return producers, nil
}

func (p *producer) Get(
	ctx context.Context, producerID string, fields ...string,
) (*entity.Producer, error) {
	producer, err := p.get(ctx, p.db.DB, producerID, fields...)
	if err != nil {
		return nil, dbError(err)
	}
	if err := p.fill(ctx, p.db.DB, producer); err != nil {
		return nil, dbError(err)
	}
	return producer, nil
}

func (p *producer) Create(
	ctx context.Context, producer *entity.Producer, auth func(ctx context.Context) error,
) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := p.now()
		producer.Admin.CreatedAt, producer.Admin.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(adminTable).Create(&producer.Admin).Error; err != nil {
			return err
		}
		producer.CreatedAt, producer.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(producerTable).Create(&producer).Error; err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (p *producer) Update(ctx context.Context, producerID string, params *database.UpdateProducerParams) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := p.now()
		adminParams := map[string]interface{}{
			"lastname":       params.Lastname,
			"firstname":      params.Firstname,
			"lastname_kana":  params.LastnameKana,
			"firstname_kana": params.FirstnameKana,
			"updated_at":     now,
		}
		producerParams := map[string]interface{}{
			"username":            params.Username,
			"profile":             params.Profile,
			"thumbnail_url":       params.ThumbnailURL,
			"header_url":          params.HeaderURL,
			"promotion_video_url": params.PromotionVideoURL,
			"bonus_video_url":     params.BonusVideoURL,
			"instagram_id":        params.InstagramID,
			"facebook_id":         params.FacebookID,
			"phone_number":        params.PhoneNumber,
			"postal_code":         params.PostalCode,
			"prefecture":          params.PrefectureCode,
			"city":                params.City,
			"address_line1":       params.AddressLine1,
			"address_line2":       params.AddressLine2,
			"updated_at":          now,
		}

		err := tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", producerID).
			Updates(adminParams).Error
		if err != nil {
			return err
		}
		err = tx.WithContext(ctx).
			Table(producerTable).
			Where("admin_id = ?", producerID).
			Updates(producerParams).Error
		return err
	})
	return dbError(err)
}

func (p *producer) UpdateThumbnails(ctx context.Context, producerID string, thumbnails common.Images) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		producer, err := p.get(ctx, tx, producerID, "thumbnail_url")
		if err != nil {
			return err
		}
		if producer.ThumbnailURL == "" {
			return fmt.Errorf("database: thumbnail url is empty: %w", database.ErrFailedPrecondition)
		}

		buf, err := thumbnails.Marshal()
		if err != nil {
			return err
		}
		params := map[string]interface{}{
			"thumbnails": datatypes.JSON(buf),
			"updated_at": p.now(),
		}

		err = tx.WithContext(ctx).
			Table(producerTable).
			Where("admin_id = ?", producerID).
			Updates(params).Error
		return err
	})
	return dbError(err)
}

func (p *producer) UpdateHeaders(ctx context.Context, producerID string, headers common.Images) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		producer, err := p.get(ctx, tx, producerID, "header_url")
		if err != nil {
			return err
		}
		if producer.HeaderURL == "" {
			return fmt.Errorf("database: header url is empty: %w", database.ErrFailedPrecondition)
		}

		buf, err := headers.Marshal()
		if err != nil {
			return err
		}
		params := map[string]interface{}{
			"headers":    datatypes.JSON(buf),
			"updated_at": p.now(),
		}

		err = tx.WithContext(ctx).
			Table(producerTable).
			Where("admin_id = ?", producerID).
			Updates(params).Error
		return err
	})
	return dbError(err)
}

func (p *producer) UpdateRelationship(ctx context.Context, coordinatorID string, producerIDs ...string) error {
	var id *string
	if coordinatorID != "" {
		id = &coordinatorID
	}

	params := map[string]interface{}{
		"coordinator_id": id,
		"updated_at":     p.now(),
	}
	stmt := p.db.DB.WithContext(ctx).
		Table(producerTable).
		Where("admin_id IN (?)", producerIDs)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (p *producer) Delete(ctx context.Context, producerID string, auth func(ctx context.Context) error) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := p.now()
		producerParams := map[string]interface{}{
			"updated_at": now,
			"deleted_at": now,
		}
		err := tx.WithContext(ctx).
			Table(producerTable).
			Where("admin_id = ?", producerID).
			Updates(producerParams).Error
		if err != nil {
			return err
		}
		adminParams := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		err = tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", producerID).
			Updates(adminParams).Error
		if err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (p *producer) AggregateByCoordinatorID(
	ctx context.Context, coordinatorIDs []string,
) (map[string]int64, error) {
	fields := []string{
		"coordinator_id",
		"COUNT(*) AS total",
	}

	stmt := p.db.Statement(ctx, p.db.DB, producerTable, fields...).
		Where("coordinator_id IN (?)", coordinatorIDs).
		Where("deleted_at IS NULL").
		Group("coordinator_id")

	rows, err := stmt.Rows()
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()

	res := make(map[string]int64, len(coordinatorIDs))
	for rows.Next() {
		var (
			coordinatorID string
			total         int64
		)
		if err := rows.Scan(&coordinatorID, &total); err != nil {
			return nil, dbError(err)
		}
		res[coordinatorID] = total
	}
	return res, nil
}

func (p *producer) multiGet(
	ctx context.Context, tx *gorm.DB, producerIDs []string, fields ...string,
) (entity.Producers, error) {
	var producers entity.Producers

	err := p.db.Statement(ctx, tx, producerTable, fields...).
		Where("admin_id IN (?)", producerIDs).
		Find(&producers).Error
	return producers, err
}

func (p *producer) get(
	ctx context.Context, tx *gorm.DB, producerID string, fields ...string,
) (*entity.Producer, error) {
	var producer *entity.Producer

	err := p.db.Statement(ctx, tx, producerTable, fields...).
		Where("admin_id = ?", producerID).
		First(&producer).Error
	return producer, err
}

func (p *producer) fill(ctx context.Context, tx *gorm.DB, producers ...*entity.Producer) error {
	var admins entity.Admins

	ids := entity.Producers(producers).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := p.db.Statement(ctx, tx, adminTable).Where("id IN (?)", ids)
	if err := stmt.Find(&admins).Error; err != nil {
		return err
	}
	admins.Fill()
	return entity.Producers(producers).Fill(admins.Map())
}
