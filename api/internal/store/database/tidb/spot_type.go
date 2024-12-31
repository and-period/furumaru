package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const spotTypeTable = "spot_types"

type spotType struct {
	database.SpotType
	db  *mysql.Client
	now func() time.Time
}

func NewSpotType(db *mysql.Client, mysql database.SpotType) database.SpotType {
	return &spotType{
		SpotType: mysql,
		db:       db,
		now:      jst.Now,
	}
}

type listSpotTypesParams database.ListSpotTypesParams

func (p listSpotTypesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`name` LIKE ?", "%"+p.Name+"%")
	}
	return stmt
}

func (p listSpotTypesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (t *spotType) List(ctx context.Context, params *database.ListSpotTypesParams, fields ...string) (entity.SpotTypes, error) {
	var types entity.SpotTypes

	p := listSpotTypesParams(*params)

	stmt := t.db.Statement(ctx, t.db.DB, spotTypeTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	err := stmt.Find(&types).Error
	return types, dbError(err)
}

func (t *spotType) Count(ctx context.Context, params *database.ListSpotTypesParams) (int64, error) {
	p := listSpotTypesParams(*params)

	total, err := t.db.Count(ctx, t.db.DB, &entity.SpotType{}, p.stmt)
	return total, dbError(err)
}
