package mysql

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
	db  *mysql.Client
	now func() time.Time
}

func NewSpotType(db *mysql.Client) database.SpotType {
	return &spotType{
		db:  db,
		now: jst.Now,
	}
}

type listSpotTypesParams database.ListSpotTypesParams

func (p listSpotTypesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("MATCH (`name`) AGAINST (? IN NATURAL LANGUAGE MODE)", p.Name)
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

func (t *spotType) MultiGet(ctx context.Context, spotTypeIDs []string, fields ...string) (entity.SpotTypes, error) {
	var types entity.SpotTypes

	stmt := t.db.Statement(ctx, t.db.DB, spotTypeTable, fields...)
	stmt = stmt.Where("id IN (?)", spotTypeIDs)

	err := stmt.Find(&types).Error
	return types, dbError(err)
}

func (t *spotType) Get(ctx context.Context, spotTypeID string, fields ...string) (*entity.SpotType, error) {
	var spotType *entity.SpotType

	stmt := t.db.Statement(ctx, t.db.DB, spotTypeTable, fields...).Where("id = ?", spotTypeID)

	if err := stmt.First(&spotType).Error; err != nil {
		return nil, dbError(err)
	}
	return spotType, nil
}

func (t *spotType) Create(ctx context.Context, spotType *entity.SpotType) error {
	now := t.now()
	spotType.CreatedAt, spotType.UpdatedAt = now, now

	err := t.db.DB.WithContext(ctx).Table(spotTypeTable).Create(spotType).Error
	return dbError(err)
}

func (t *spotType) Update(ctx context.Context, spotTypeID string, params *database.UpdateSpotTypeParams) error {
	updates := map[string]interface{}{
		"name":       params.Name,
		"updated_at": t.now(),
	}
	stmt := t.db.DB.WithContext(ctx).Table(spotTypeTable).Where("id = ?", spotTypeID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (t *spotType) Delete(ctx context.Context, spotTypeID string) error {
	stmt := t.db.DB.WithContext(ctx).Table(spotTypeTable).Where("id = ?", spotTypeID)

	err := stmt.Delete(&entity.SpotType{}).Error
	return dbError(err)
}
