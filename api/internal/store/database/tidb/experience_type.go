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

const experienceTypeTable = "experience_types"

type experienceType struct {
	database.ExperienceType
	db  *mysql.Client
	now func() time.Time
}

func NewExperienceType(db *mysql.Client, mysql database.ExperienceType) database.ExperienceType {
	return &experienceType{
		ExperienceType: mysql,
		db:             db,
		now:            jst.Now,
	}
}

type listExperienceTypesParams database.ListExperienceTypesParams

func (p listExperienceTypesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`name` LIKE ?", "%"+p.Name+"%").
			Or("`description` LIKE ?", "%"+p.Name+"%")
	}
	return stmt
}

func (p listExperienceTypesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (t *experienceType) List(
	ctx context.Context, params *database.ListExperienceTypesParams, fields ...string,
) (entity.ExperienceTypes, error) {
	var types entity.ExperienceTypes

	p := listExperienceTypesParams(*params)

	stmt := t.db.Statement(ctx, t.db.DB, experienceTypeTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	err := stmt.Find(&types).Error
	return types, dbError(err)
}

func (t *experienceType) Count(ctx context.Context, params *database.ListExperienceTypesParams) (int64, error) {
	p := listExperienceTypesParams(*params)

	total, err := t.db.Count(ctx, t.db.DB, &entity.ExperienceType{}, p.stmt)
	return total, dbError(err)
}
