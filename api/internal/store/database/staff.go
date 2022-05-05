package database

import (
	"context"
	"time"

	"github.com/and-period/marche/api/internal/store/entity"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/pkg/jst"
)

const staffTable = "staffs"

var staffFields = []string{
	"store_id", "user_id", "role",
	"created_at", "updated_at",
}

type staff struct {
	db  *database.Client
	now func() time.Time
}

func NewStaff(db *database.Client) Staff {
	return &staff{
		db:  db,
		now: jst.Now,
	}
}

func (s *staff) ListByStoreID(ctx context.Context, storeID int64, fields ...string) (entity.Staffs, error) {
	var staffs entity.Staffs
	if len(fields) == 0 {
		fields = staffFields
	}

	stmt := s.db.DB.WithContext(ctx).
		Table(staffTable).Select(fields).
		Where("store_id = ?", storeID)

	err := stmt.Find(&staffs).Error
	return staffs, dbError(err)
}
