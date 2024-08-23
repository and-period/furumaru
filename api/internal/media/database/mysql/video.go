package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
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

func (v *video) List(ctx context.Context, params *database.ListVideosParams, fields ...string) (entity.Videos, error) {
	// TODO: 詳細の実装
	return entity.Videos{}, nil
}

func (v *video) Count(ctx context.Context, params *database.ListVideosParams) (int64, error) {
	// TODO: 詳細の実装
	return 0, nil
}

func (v *video) Get(ctx context.Context, videoID string, fields ...string) (*entity.Video, error) {
	// TODO: 詳細の実装
	return &entity.Video{}, nil
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
