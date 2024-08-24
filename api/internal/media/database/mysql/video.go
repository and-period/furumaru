package mysql

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
	db  *mysql.Client
	now func() time.Time
}

func newVideo(db *mysql.Client) database.Video {
	return &video{
		db:  db,
		now: jst.Now,
	}
}

type listVideosParams database.ListVideosParams

func (p listVideosParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("MATCH (`title`, `description`) AGAINST (? IN NATURAL LANGUAGE MODE)", p.Name)
	}
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
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
	stmt = p.stmt(stmt)
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

	total, err := v.db.Count(ctx, v.db.DB, &entity.Video{}, p.stmt)
	return total, dbError(err)
}

func (v *video) Get(ctx context.Context, videoID string, fields ...string) (*entity.Video, error) {
	video, err := v.get(ctx, v.db.DB, videoID, fields...)
	return video, dbError(err)
}

func (v *video) Create(ctx context.Context, video *entity.Video) error {
	// TODO: 詳細の実装
	return nil
}

func (v *video) Update(ctx context.Context, videoID string, params *database.UpdateVideoParams) error {
	// TODO: 詳細の実装
	return nil
}

func (v *video) Delete(ctx context.Context, videoID string) error {
	// TODO: 詳細の実装
	return nil
}

func (v *video) get(ctx context.Context, tx *gorm.DB, videoID string, fields ...string) (*entity.Video, error) {
	var video *entity.Video

	stmt := v.db.Statement(ctx, tx, videoTable, fields...).Where("id = ?", videoID)

	if err := stmt.First(&video).Error; err != nil {
		return nil, err
	}
	if err := v.fill(ctx, tx, video); err != nil {
		return nil, err
	}
	return video, nil
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
