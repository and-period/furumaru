package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const cartActionLogTable = "cart_action_logs"

type cartActionLog struct {
	db  *mysql.Client
	now func() time.Time
}

func newCartActionLog(db *mysql.Client) database.CartActionLog {
	return &cartActionLog{
		db:  db,
		now: jst.Now,
	}
}

func (c *cartActionLog) Create(ctx context.Context, log *entity.CartActionLog) error {
	now := c.now()
	log.CreatedAt, log.UpdatedAt = now, now

	err := c.db.DB.WithContext(ctx).Table(cartActionLogTable).Create(&log).Error
	return dbError(err)
}
