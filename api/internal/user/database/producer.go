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
	"id", "email", "phone_number",
	"lastname", "firstname", "lastname_kana", "firstname_kana",
	"store_name", "thumbnail_url", "header_url",
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

	err := stmt.Find(&producers).Error
	return producers, exception.InternalError(err)
}

func (p *producer) Count(ctx context.Context, params *ListProducersParams) (int64, error) {
	var total int64

	stmt := p.db.DB.WithContext(ctx).Table(producerTable).Select("COUNT(*)")

	err := stmt.Count(&total).Error
	return total, exception.InternalError(err)
}

func (p *producer) MultiGet(
	ctx context.Context, producerIDs []string, fields ...string,
) (entity.Producers, error) {
	var producers entity.Producers
	if len(fields) == 0 {
		fields = producerFields
	}

	err := p.db.DB.WithContext(ctx).
		Table(producerTable).Select(fields).
		Where("id IN (?)", producerIDs).
		Find(&producers).Error
	return producers, exception.InternalError(err)
}

func (p *producer) Get(
	ctx context.Context, producerID string, fields ...string,
) (*entity.Producer, error) {
	producer, err := p.get(ctx, p.db.DB, producerID, fields...)
	return producer, exception.InternalError(err)
}

func (p *producer) Create(
	ctx context.Context, auth *entity.AdminAuth, producer *entity.Producer,
) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := p.now()
		auth.CreatedAt, auth.UpdatedAt = now, now
		producer.CreatedAt, producer.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(adminAuthTable).Create(&auth).Error
		if err != nil {
			return nil, err
		}
		err = tx.WithContext(ctx).Table(producerTable).Create(&producer).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (p *producer) Update(ctx context.Context, producerID string, params *UpdateProducerParams) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := p.get(ctx, tx, producerID); err != nil {
			return nil, err
		}

		updates := map[string]interface{}{
			"lastname":       params.Lastname,
			"firstname":      params.Firstname,
			"lastname_kana":  params.LastnameKana,
			"firstname_kana": params.FirstnameKana,
			"store_name":     params.StoreName,
			"thumbnail_url":  params.ThumbnailURL,
			"header_url":     params.HeaderURL,
			"phone_number":   params.PhoneNumber,
			"postal_code":    params.PostalCode,
			"city":           params.City,
			"address_line1":  params.AddressLine1,
			"address_line2":  params.AddressLine2,
			"updated_at":     p.now(),
		}
		err := tx.WithContext(ctx).
			Table(producerTable).
			Where("id = ?", producerID).
			Updates(updates).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (p *producer) UpdateEmail(ctx context.Context, producerID, email string) error {
	_, err := p.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := p.get(ctx, tx, producerID); err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"email":      email,
			"updated_at": p.now(),
		}
		err := tx.WithContext(ctx).
			Table(producerTable).
			Where("id = ?", producerID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
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
		Where("id = ?", producerID).
		First(&producer).Error
	return producer, err
}
