package database

import (
	"context"
	"time"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/pkg/jst"
	"gorm.io/gorm"
)

const shopTable = "shops"

var shopFields = []string{
	"id", "cognito_id", "email",
	"lastname", "firstname", "lastname_kana", "firstname_kana",
	"created_at", "updated_at", "deleted_at",
}

type shop struct {
	db  *database.Client
	now func() time.Time
}

func NewShop(db *database.Client) Shop {
	return &shop{
		db:  db,
		now: jst.Now,
	}
}

func (s *shop) List(ctx context.Context, params *ListShopsParams, fields ...string) (entity.Shops, error) {
	var shops entity.Shops
	if len(fields) == 0 {
		fields = shopFields
	}

	stmt := s.db.DB.WithContext(ctx).Table(shopTable).Select(fields)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&shops).Error
	return shops, dbError(err)
}

func (s *shop) MultiGet(ctx context.Context, shopIDs []string, fields ...string) (entity.Shops, error) {
	var shops entity.Shops
	if len(fields) == 0 {
		fields = shopFields
	}

	stmt := s.db.DB.WithContext(ctx).
		Table(shopTable).Select(fields).
		Where("id IN (?)", shopIDs)

	err := stmt.Find(&shops).Error
	return shops, dbError(err)
}

func (s *shop) Get(ctx context.Context, shopID string, fields ...string) (*entity.Shop, error) {
	var shop *entity.Shop
	if len(fields) == 0 {
		fields = shopFields
	}

	stmt := s.db.DB.WithContext(ctx).
		Table(shopTable).Select(fields).
		Where("id = ?", shopID)

	if err := stmt.First(&shop).Error; err != nil {
		return nil, dbError(err)
	}
	return shop, nil
}

func (s *shop) Create(ctx context.Context, shop *entity.Shop) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := s.now()
		shop.CreatedAt, shop.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(shopTable).Create(&shop).Error
		return nil, err
	})
	return dbError(err)
}
