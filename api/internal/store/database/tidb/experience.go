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

const (
	experienceTable         = "experiences"
	experienceRevisionTable = "experience_revisions"
)

type experience struct {
	database.Experience
	db  *mysql.Client
	now func() time.Time
}

func newExperience(db *mysql.Client, mysql database.Experience) database.Experience {
	return &experience{
		Experience: mysql,
		db:         db,
		now:        jst.Now,
	}
}

type listExperiencesParams database.ListExperiencesParams

func (p listExperiencesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`title` LIKE ?", "%"+p.Name+"%").
			Or("`description` LIKE ?", "%"+p.Name+"%")
	}
	if p.HostPrefecture > 0 {
		stmt = stmt.Where("host_prefecture = ?", p.HostPrefecture)
	}
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.ProducerID != "" {
		stmt = stmt.Where("producer_id = ?", p.ProducerID)
	}
	if p.OnlyPublished {
		stmt = stmt.Where("public = ?", true).Where("deleted_at IS NULL")
	}
	if !p.EndAtGte.IsZero() {
		stmt = stmt.Where("end_at >= ?", p.EndAtGte)
	}
	if !p.ExcludeDeleted {
		stmt = stmt.Unscoped()
	}
	return stmt.Order("start_at DESC")
}

func (p listExperiencesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (e *experience) List(ctx context.Context, params *database.ListExperiencesParams, fields ...string) (entity.Experiences, error) {
	var experiences entity.Experiences

	p := listExperiencesParams(*params)

	stmt := e.db.Statement(ctx, e.db.DB, experienceTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&experiences).Error; err != nil {
		return nil, dbError(err)
	}
	if err := e.fill(ctx, e.db.DB, experiences...); err != nil {
		return nil, dbError(err)
	}
	return experiences, nil
}

func (e *experience) Count(ctx context.Context, params *database.ListExperiencesParams) (int64, error) {
	p := listExperiencesParams(*params)

	total, err := e.db.Count(ctx, e.db.DB, &entity.Experience{}, p.stmt)
	return total, dbError(err)
}

func (e *experience) fill(ctx context.Context, tx *gorm.DB, experiences ...*entity.Experience) error {
	var revisions entity.ExperienceRevisions

	ids := entity.Experiences(experiences).IDs()
	if len(ids) == 0 {
		return nil
	}

	sub := tx.Table(experienceRevisionTable).
		Select("MAX(id)").
		Where("experience_id IN (?)", ids).
		Group("experience_id")
	stmt := e.db.Statement(ctx, tx, experienceRevisionTable).
		Where("id IN (?)", sub)

	if err := stmt.Find(&revisions).Error; err != nil {
		return err
	}
	if len(revisions) == 0 {
		return nil
	}
	return entity.Experiences(experiences).Fill(revisions.MapByExperienceID(), e.now())
}
