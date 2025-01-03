package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const producerTable = "producers"

type producer struct {
	db  *mysql.Client
	now func() time.Time
}

func NewProducer(db *mysql.Client) database.Producer {
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
	if p.Name != "" {
		stmt = stmt.Where("MATCH (`username`, `profile`) AGAINST (? IN NATURAL LANGUAGE MODE)", p.Name)
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

func (p *producer) MultiGetWithDeleted(
	ctx context.Context, producerIDs []string, fields ...string,
) (entity.Producers, error) {
	var producers entity.Producers

	stmt := p.db.Statement(ctx, p.db.DB, producerTable, fields...).
		Where("admin_id IN (?)", producerIDs).
		Unscoped()

	if err := stmt.Find(&producers).Error; err != nil {
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

func (p *producer) GetWithDeleted(
	ctx context.Context, producerID string, fields ...string,
) (*entity.Producer, error) {
	var producer *entity.Producer

	stmt := p.db.Statement(ctx, p.db.DB, producerTable, fields...).
		Where("admin_id = ?", producerID).
		Unscoped()

	if err := stmt.First(&producer).Error; err != nil {
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
			"lastname":       mysql.NullString(params.Lastname),
			"firstname":      mysql.NullString(params.Firstname),
			"lastname_kana":  mysql.NullString(params.LastnameKana),
			"firstname_kana": mysql.NullString(params.FirstnameKana),
			"email":          mysql.NullString(params.Email),
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
			"phone_number":        mysql.NullString(params.PhoneNumber),
			"postal_code":         mysql.NullString(params.PostalCode),
			"prefecture":          mysql.NullInt(params.PrefectureCode),
			"city":                mysql.NullString(params.City),
			"address_line1":       mysql.NullString(params.AddressLine1),
			"address_line2":       mysql.NullString(params.AddressLine2),
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

func (p *producer) Delete(ctx context.Context, producerID string, auth func(ctx context.Context) error) error {
	err := p.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := p.now()
		updates := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		stmt := tx.WithContext(ctx).Table(adminTable).Where("id = ?", producerID)
		if err := stmt.Updates(updates).Error; err != nil {
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
		"producers.coordinator_id AS coordinator_id",
		"COUNT(producers.admin_id) AS total",
	}

	stmt := p.db.Statement(ctx, p.db.DB, producerTable, fields...).
		Joins("INNER JOIN admins ON admins.id = producers.admin_id").
		Where("producers.coordinator_id IN (?)", coordinatorIDs).
		Where("admins.deleted_at IS NULL").
		Group("producers.coordinator_id")

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

	stmt := p.db.Statement(ctx, tx, adminTable).Unscoped().Where("id IN (?)", ids)
	if err := stmt.Find(&admins).Error; err != nil {
		return err
	}
	// TODO: 管理者グループID一覧を取得する処理を追加
	if err := admins.Fill(nil); err != nil {
		return err
	}
	return entity.Producers(producers).Fill(admins.Map())
}
