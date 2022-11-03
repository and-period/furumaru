package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const producerTable = "producers"

var producerFields = []string{
	"admin_id", "coordinator_id", "phone_number", "store_name", "thumbnail_url", "header_url",
	"postal_code", "prefecture", "city", "address_line1", "address_line2",
	"created_at", "updated_at", "deleted_at",
}

type producer struct {
	db  *database.Client
	now func() time.Time
}

func NewProducer(db *database.Client) Producer {
	return &producer{
		db:  db,
		now: jst.Now,
	}
}

func (p *producer) List(
	ctx context.Context, params *ListProducersParams, fields ...string,
) (entity.Producers, error) {
	var producers entity.Producers
	if len(fields) == 0 {
		fields = producerFields
	}

	stmt := p.db.DB.WithContext(ctx).Table(producerTable).Select(fields)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&producers).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	if err := p.fill(ctx, p.db.DB, producers...); err != nil {
		return nil, exception.InternalError(err)
	}
	return producers, nil
}

func (p *producer) Count(ctx context.Context, params *ListProducersParams) (int64, error) {
	var total int64

	stmt := p.db.DB.WithContext(ctx).Table(producerTable).Select("COUNT(*)")
	stmt = params.stmt(stmt)

	err := stmt.Count(&total).Error
	return total, exception.InternalError(err)
}

func (p *producer) MultiGet(
	ctx context.Context, producerIDs []string, fields ...string,
) (entity.Producers, error) {
	producers, err := p.multiGet(ctx, p.db.DB, producerIDs, fields...)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	if err := p.fill(ctx, p.db.DB, producers...); err != nil {
		return nil, exception.InternalError(err)
	}
	return producers, nil
}

func (p *producer) Get(
	ctx context.Context, producerID string, fields ...string,
) (*entity.Producer, error) {
	producer, err := p.get(ctx, p.db.DB, producerID, fields...)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	if err := p.fill(ctx, p.db.DB, producer); err != nil {
		return nil, exception.InternalError(err)
	}
	return producer, nil
}

func (p *producer) Create(
	ctx context.Context, producer *entity.Producer, auth func(ctx context.Context) error,
) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := p.now()
		producer.Admin.CreatedAt, producer.Admin.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(adminTable).Create(&producer.Admin).Error; err != nil {
			return nil, err
		}
		producer.CreatedAt, producer.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(producerTable).Create(&producer).Error; err != nil {
			return nil, err
		}
		return nil, auth(ctx)
	})
	return exception.InternalError(err)
}

func (p *producer) Update(ctx context.Context, producerID string, params *UpdateProducerParams) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := p.get(ctx, tx, producerID); err != nil {
			return nil, err
		}

		now := p.now()
		adminParams := map[string]interface{}{
			"lastname":       params.Lastname,
			"firstname":      params.Firstname,
			"lastname_kana":  params.LastnameKana,
			"firstname_kana": params.FirstnameKana,
			"updated_at":     now,
		}
		producerParams := map[string]interface{}{
			"store_name":    params.StoreName,
			"thumbnail_url": params.ThumbnailURL,
			"header_url":    params.HeaderURL,
			"phone_number":  params.PhoneNumber,
			"postal_code":   params.PostalCode,
			"city":          params.City,
			"address_line1": params.AddressLine1,
			"address_line2": params.AddressLine2,
			"updated_at":    now,
		}

		err := tx.WithContext(ctx).
			Table(adminTable).
			Where("id = ?", producerID).
			Updates(adminParams).Error
		if err != nil {
			return nil, err
		}
		err = tx.WithContext(ctx).
			Table(producerTable).
			Where("admin_id = ?", producerID).
			Updates(producerParams).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (p *producer) UpdateRelationship(ctx context.Context, coordinatorID string, producerIDs ...string) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var id *string
		if coordinatorID != "" {
			id = &coordinatorID
		}

		params := map[string]interface{}{
			"coordinator_id": id,
			"updated_at":     p.now(),
		}
		err := tx.WithContext(ctx).
			Table(producerTable).
			Where("admin_id IN (?)", producerIDs).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (p *producer) Delete(ctx context.Context, producerID string, auth func(ctx context.Context) error) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := p.get(ctx, tx, producerID); err != nil {
			return nil, err
		}

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
			return nil, err
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
			return nil, err
		}
		return nil, auth(ctx)
	})
	return exception.InternalError(err)
}

func (p *producer) multiGet(
	ctx context.Context, tx *gorm.DB, producerIDs []string, fields ...string,
) (entity.Producers, error) {
	var producers entity.Producers
	if len(fields) == 0 {
		fields = producerFields
	}

	err := tx.WithContext(ctx).
		Table(producerTable).Select(fields).
		Where("admin_id IN (?)", producerIDs).
		Find(&producers).Error
	return producers, err
}

func (p *producer) get(
	ctx context.Context, tx *gorm.DB, producerID string, fields ...string,
) (*entity.Producer, error) {
	var producer *entity.Producer
	if len(fields) == 0 {
		fields = producerFields
	}

	err := tx.WithContext(ctx).
		Table(producerTable).Select(fields).
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

	stmt := tx.WithContext(ctx).
		Table(adminTable).Select(adminFields).
		Where("id IN (?)", ids)
	if err := stmt.Find(&admins).Error; err != nil {
		return err
	}

	adminMap := admins.Map()

	for i, p := range producers {
		admin, ok := adminMap[p.AdminID]
		if !ok {
			admin = &entity.Admin{}
		}

		producers[i].Fill(admin)
	}
	return nil
}
