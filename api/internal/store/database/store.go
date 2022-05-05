package database

import (
	"context"
	"time"

	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/pkg/jst"
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
	return stores, dbError(err)
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
		return nil, dbError(err)
	}
	return store, nil
}
