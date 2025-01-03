package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const producerTable = "producers"

type producer struct {
	database.Producer
	db  *mysql.Client
	now func() time.Time
}

func NewProducer(db *mysql.Client, mysql database.Producer) database.Producer {
	return &producer{
		Producer: mysql,
		db:       db,
		now:      jst.Now,
	}
}

type listProducersParams database.ListProducersParams

func (p listProducersParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.Name != "" {
		stmt = stmt.Where("`username` LIKE ?", "%"+p.Name+"%").
			Or("`profile` LIKE ?", "%"+p.Name+"%")
	}
	stmt = stmt.Order("updated_at DESC")
	return stmt
}

func (p listProducersParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (p *producer) List(
	ctx context.Context, params *database.ListProducersParams, fields ...string,
) (entity.Producers, error) {
	var producers entity.Producers

	prm := listProducersParams(*params)

	stmt := p.db.Statement(ctx, p.db.DB, producerTable, fields...)
	stmt = prm.stmt(stmt)
	stmt = prm.pagination(stmt)

	if err := stmt.Find(&producers).Error; err != nil {
		return nil, dbError(err)
	}
	if err := p.fill(ctx, p.db.DB, producers...); err != nil {
		return nil, dbError(err)
	}
	return producers, nil
}

func (p *producer) Count(ctx context.Context, params *database.ListProducersParams) (int64, error) {
	prm := listProducersParams(*params)

	total, err := p.db.Count(ctx, p.db.DB, &entity.Producer{}, prm.stmt)
	return total, dbError(err)
}

func (p *producer) fill(ctx context.Context, tx *gorm.DB, producers ...*entity.Producer) error {
	var admins entity.Admins

	ids := entity.Producers(producers).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := p.db.Statement(ctx, tx, adminTable).Unscoped().Where("id IN (?)", ids)
	if err := stmt.Find(&admins).Error; err != nil {
		return err
	}
	// TODO: 管理者グループID一覧を取得する処理を追加
	if err := admins.Fill(nil); err != nil {
		return err
	}
	return entity.Producers(producers).Fill(admins.Map())
}
