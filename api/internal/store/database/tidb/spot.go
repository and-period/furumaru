package tidb

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const spotTable = "spots"

type spot struct {
	database.Spot
	db  *mysql.Client
	now func() time.Time
}

func NewSpot(db *mysql.Client, mysql database.Spot) database.Spot {
	return &spot{
		Spot: mysql,
		db:   db,
		now:  jst.Now,
	}
}

type listSpotsParams database.ListSpotsParams

func (p listSpotsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", p.Name)).
			Or("`description` LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
	}
	if p.ExcludeApproved {
		stmt = stmt.Where("approved = ?", false)
	}
	if p.ExcludeDisabled {
		stmt = stmt.Where("approved = ?", true)
	}
	return stmt
}

func (p listSpotsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (s *spot) List(ctx context.Context, params *database.ListSpotsParams, fields ...string) (entity.Spots, error) {
	var spots entity.Spots

	prm := listSpotsParams(*params)

	stmt := s.db.Statement(ctx, s.db.DB, spotTable, fields...)
	stmt = prm.stmt(stmt)
	stmt = prm.pagination(stmt)

	if err := stmt.Find(&spots).Error; err != nil {
		return nil, dbError(err)
	}
	if err := spots.Fill(); err != nil {
		return nil, dbError(err)
	}
	return spots, nil
}

func (s *spot) Count(ctx context.Context, params *database.ListSpotsParams) (int64, error) {
	prm := listSpotsParams(*params)

	total, err := s.db.Count(ctx, s.db.DB, &entity.Spot{}, prm.stmt)
	return total, dbError(err)
}
