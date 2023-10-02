package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const adddressTable = "addresses"

type address struct {
	db  *mysql.Client
	now func() time.Time
}

func newAddress(db *mysql.Client) database.Address {
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
	return addresses, dbError(err)
}
