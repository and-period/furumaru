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

const storeTable = "stores"

var storeFields = []string{
	"id", "name", "thumbnail_url",
	"created_at", "updated_at", "deleted_at",
}

type store struct {
	db  *database.Client
	now func() time.Time
}

func NewStore(db *database.Client) Store {
	return &store{
		db:  db,
		now: jst.Now,
	}
}

func (s *store) List(ctx context.Context, params *ListStoresParams, fields ...string) (entity.Stores, error) {
	var stores entity.Stores
	if len(fields) == 0 {
		fields = storeFields
	}

	stmt := s.db.DB.WithContext(ctx).Table(storeTable).Select(fields)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&stores).Error
	return stores, exception.InternalError(err)
}

func (s *store) Get(ctx context.Context, storeID int64, fields ...string) (*entity.Store, error) {
	var store *entity.Store
	if len(fields) == 0 {
		fields = storeFields
	}

	stmt := s.db.DB.WithContext(ctx).
		Table(storeTable).Select(fields).
		Where("id = ?", storeID)

	if err := stmt.First(&store).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return store, nil
}

func (s *store) Create(ctx context.Context, store *entity.Store) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := s.now()
		store.CreatedAt, store.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(storeTable).Create(&store).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (s *store) Update(ctx context.Context, storeID int64, name, thumbnailURL string) error {
	_, err := s.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var current *entity.Store
		err := tx.WithContext(ctx).
			Table(storeTable).Select("id").
			Where("id = ?", storeID).
			First(&current).Error
		if err != nil {
			return nil, err
		}

		params := map[string]interface{}{
			"id":            current.ID,
			"name":          name,
			"thumbnail_url": thumbnailURL,
			"updated_at":    s.now(),
		}
		err = tx.WithContext(ctx).
			Table(storeTable).
			Where("id = ?", current.ID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}
