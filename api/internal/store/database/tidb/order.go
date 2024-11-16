package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const orderTable = "orders"

type order struct {
	database.Order
	db  *mysql.Client
	now func() time.Time
}

func newOrder(db *mysql.Client, mysql database.Order) database.Order {
	return &order{
		Order: mysql,
		db:    db,
		now:   jst.Now,
	}
}

type listOrdersParams database.ListOrdersParams

func (p listOrdersParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
	}
	if len(p.Types) > 0 {
		stmt = stmt.Where("type IN (?)", p.Types)
	}
	if len(p.Statuses) > 0 {
		stmt = stmt.Where("status IN (?)", p.Statuses)
	}
	return stmt
}

func (p listOrdersParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (o *order) ListUserIDs(ctx context.Context, params *database.ListOrdersParams) ([]string, int64, error) {
	var userIDs []string
	var total int64

	p := listOrdersParams(*params)

	stmt := o.db.Statement(ctx, o.db.DB, orderTable, "DISTINCT(user_id)")
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)
	if err := stmt.Find(&userIDs).Error; err != nil {
		return nil, 0, dbError(err)
	}

	stmt = o.db.Statement(ctx, o.db.DB, orderTable, "COUNT(DISTINCT(user_id))")
	stmt = p.stmt(stmt)
	if err := stmt.Count(&total).Error; err != nil {
		return nil, 0, dbError(err)
	}

	return userIDs, total, nil
}
