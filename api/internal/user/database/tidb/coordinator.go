package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	apmysql "github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const coordinatorTable = "coordinators"

type coordinator struct {
	database.Coordinator
	db  *apmysql.Client
	now func() time.Time
}

func NewCoordinator(db *apmysql.Client, mysql database.Coordinator) database.Coordinator {
	return &coordinator{
		Coordinator: mysql,
		db:          db,
		now:         jst.Now,
	}
}

type listCoordinatorsParams database.ListCoordinatorsParams

func (p listCoordinatorsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`username` LIKE ?", "%"+p.Name+"%").
			Or("`marche_name` LIKE ?", "%"+p.Name+"%").
			Or("`profile` LIKE ?", "%"+p.Name+"%")
	}
	stmt = stmt.Order("updated_at DESC")
	return stmt
}

func (p listCoordinatorsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (c *coordinator) List(
	ctx context.Context, params *database.ListCoordinatorsParams, fields ...string,
) (entity.Coordinators, error) {
	var coordinators entity.Coordinators

	p := listCoordinatorsParams(*params)

	stmt := c.db.Statement(ctx, c.db.DB, coordinatorTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&coordinators).Error; err != nil {
		return nil, dbError(err)
	}
	if err := c.fill(ctx, c.db.DB, coordinators...); err != nil {
		return nil, dbError(err)
	}
	return coordinators, nil
}

func (c *coordinator) Count(ctx context.Context, params *database.ListCoordinatorsParams) (int64, error) {
	p := listCoordinatorsParams(*params)

	total, err := c.db.Count(ctx, c.db.DB, &entity.Coordinator{}, p.stmt)
	return total, dbError(err)
}

func (c *coordinator) fill(ctx context.Context, tx *gorm.DB, coordinators ...*entity.Coordinator) error {
	var admins entity.Admins

	ids := entity.Coordinators(coordinators).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := c.db.Statement(ctx, tx, adminTable).Unscoped().Where("id IN (?)", ids)
	if err := stmt.Find(&admins).Error; err != nil {
		return err
	}
	// TODO: 管理者グループID一覧を取得する処理を追加
	if err := admins.Fill(nil); err != nil {
		return err
	}
	return entity.Coordinators(coordinators).Fill(admins.Map())
}
