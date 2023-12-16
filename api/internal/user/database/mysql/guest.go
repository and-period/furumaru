package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const guestTable = "guests"

type guest struct {
	db  *mysql.Client
	now func() time.Time
}

func newGuest(db *mysql.Client) database.Guest {
	return &guest{
		db:  db,
		now: jst.Now,
	}
}

func (g *guest) Delete(ctx context.Context, userID string) error {
	err := g.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := g.now()
		guestParams := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
		}
		err := tx.WithContext(ctx).
			Table(guestTable).
			Where("user_id = ?", userID).
			Updates(guestParams).Error
		if err != nil {
			return err
		}
		userParams := map[string]interface{}{
			"updated_at": now,
			"deleted_at": now,
		}
		err = tx.WithContext(ctx).
			Table(userTable).
			Where("id = ?", userID).
			Updates(userParams).Error
		return err
	})
	return dbError(err)
}
