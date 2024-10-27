package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const (
	videoTable           = "videos"
	videoProductTable    = "video_products"
	videoExperienceTable = "video_experiences"
)

type video struct {
	database.Video
	db  *mysql.Client
	now func() time.Time
}

func newVideo(db *mysql.Client, mysql database.Video) database.Video {
	return &video{
		Video: mysql,
		db:    db,
		now:   jst.Now,
	}
}

type listVideosParams database.ListVideosParams

func (p listVideosParams) stmt(stmt *gorm.DB, now time.Time) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`title` LIKE ?", "%"+p.Name+"%").
			Or("`description` LIKE ?", "%"+p.Name+"%")
	}
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.OnlyPublished {
		stmt = stmt.Where("public = ? AND published_at <= ?", true, now)
	}
	if p.OnlyDisplayProduct {
		stmt = stmt.Where("display_product = ?", true)
	}
	if p.OnlyDisplayExperience {
		stmt = stmt.Where("display_experience = ?", true)
	}
	if p.ExcludeLimited {
		stmt = stmt.Where("limited = ?", false)
	}
	return stmt.Order("published_at DESC")
}

func (p listVideosParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (v *video) List(ctx context.Context, params *database.ListVideosParams, fields ...string) (entity.Videos, error) {
	var videos entity.Videos

	p := listVideosParams(*params)

	stmt := v.db.Statement(ctx, v.db.DB, videoTable, fields...)
	stmt = p.stmt(stmt, v.now())
	stmt = p.pagination(stmt)

	if err := stmt.Find(&videos).Error; err != nil {
		return nil, dbError(err)
	}
	if err := v.fill(ctx, v.db.DB, videos...); err != nil {
		return nil, dbError(err)
	}
	return videos, nil
}

func (v *video) Count(ctx context.Context, params *database.ListVideosParams) (int64, error) {
	p := listVideosParams(*params)

	fn := func(stmt *gorm.DB) *gorm.DB { return p.stmt(stmt, v.now()) }
	total, err := v.db.Count(ctx, v.db.DB, &entity.Video{}, fn)
	return total, dbError(err)
}

func (v *video) fill(ctx context.Context, tx *gorm.DB, videos ...*entity.Video) error {
	var (
		products    entity.VideoProducts
		experiences entity.VideoExperiences
	)

	ids := entity.Videos(videos).IDs()
	if len(ids) == 0 {
		return nil
	}

	stmt := v.db.Statement(ctx, tx, videoProductTable).Where("video_id IN (?)", ids)
	if err := stmt.Find(&products).Error; err != nil {
		return err
	}
	stmt = v.db.Statement(ctx, tx, videoExperienceTable).Where("video_id IN (?)", ids)
	if err := stmt.Find(&experiences).Error; err != nil {
		return err
	}

	entity.Videos(videos).Fill(products.GroupByVideoID(), experiences.GroupByVideoID(), v.now())
	return nil
}
