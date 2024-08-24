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

const experienceTypeTable = "experience_types"

type experienceType struct {
	db  *mysql.Client
	now func() time.Time
}

func newExperienceType(db *mysql.Client) database.ExperienceType {
	return &experienceType{
		db:  db,
		now: jst.Now,
	}
}

type listExperienceTypesParams database.ListExperienceTypesParams

func (p listExperienceTypesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("MATCH (`name`, `description`) AGAINST (? IN NATURAL LANGUAGE MODE)", p.Name)
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

func (t *experienceType) MultiGet(ctx context.Context, experienceIDs []string, fields ...string) (entity.ExperienceTypes, error) {
	var types entity.ExperienceTypes

	stmt := t.db.Statement(ctx, t.db.DB, experienceTypeTable, fields...).
		Where("id IN (?)", experienceIDs)

	err := stmt.Find(&types).Error
	return types, dbError(err)
}

func (t *experienceType) Get(ctx context.Context, experienceID string, fields ...string) (*entity.ExperienceType, error) {
	experienceType, err := t.get(ctx, t.db.DB, experienceID, fields...)
	return experienceType, dbError(err)
}

func (t *experienceType) Create(ctx context.Context, experience *entity.ExperienceType) error {
	now := t.now()
	experience.CreatedAt, experience.UpdatedAt = now, now

	err := t.db.DB.WithContext(ctx).Table(experienceTypeTable).Create(experience).Error
	return dbError(err)
}

func (t *experienceType) Update(ctx context.Context, experienceID string, params *database.UpdateExperienceTypeParams) error {
	updates := map[string]interface{}{
		"name":       params.Name,
		"updated_at": t.now(),
	}
	stmt := t.db.DB.WithContext(ctx).Table(experienceTypeTable).Where("id = ?", experienceID)
	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (t *experienceType) Delete(ctx context.Context, experienceID string) error {
	stmt := t.db.DB.WithContext(ctx).Table(experienceTypeTable).Where("id = ?", experienceID)
	err := stmt.Delete(&entity.ExperienceType{}).Error
	return dbError(err)
}

func (t *experienceType) get(ctx context.Context, db *gorm.DB, experienceID string, fields ...string) (*entity.ExperienceType, error) {
	var experience *entity.ExperienceType

	stmt := t.db.Statement(ctx, db, experienceTypeTable, fields...).
		Where("id = ?", experienceID)

	if err := stmt.First(&experience).Error; err != nil {
		return nil, err
	}
	return experience, nil
}
