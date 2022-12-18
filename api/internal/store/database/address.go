package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
)

const adddressTable = "addresses"

type address struct {
	db  *database.Client
	now func() time.Time
}

func NewAddress(db *database.Client) Address {
	return &address{
		db:  db,
		now: jst.Now,
	}
}

func (a *address) MultiGet(ctx context.Context, addressIDs []string, fields ...string) (entity.Addresses, error) {
	var addresses entity.Addresses

	err := a.db.Statement(ctx, a.db.DB, adddressTable, fields...).
		Where("id IN (?)", addressIDs).
		Find(&addresses).Error
	return addresses, exception.InternalError(err)
}
