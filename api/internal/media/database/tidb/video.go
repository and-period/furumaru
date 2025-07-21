package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func NewVideo(db *mysql.Client) database.Video {
	return &video{
		db:  db,
		now: jst.Now,
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

func (v *video) List(
	ctx context.Context,
	params *database.ListVideosParams,
	fields ...string,
) (entity.Videos, error) {
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

func (v *video) ListByProductID(
	ctx context.Context,
	productID string,
	fields ...string,
) (entity.Videos, error) {
	var videos entity.Videos

	sub := v.db.DB.
		Table(videoProductTable).
		Select("video_id").
		Where("product_id = ?", productID)
	stmt := v.db.DB.WithContext(ctx).
		Table(videoTable).
		Where("id IN (?)", sub)

	if err := stmt.Find(&videos).Error; err != nil {
		return nil, dbError(err)
	}
	if err := v.fill(ctx, v.db.DB, videos...); err != nil {
		return nil, dbError(err)
	}
	return videos, nil
}

func (v *video) ListByExperienceID(
	ctx context.Context,
	experienceID string,
	fields ...string,
) (entity.Videos, error) {
	var videos entity.Videos

	sub := v.db.DB.
		Table(videoExperienceTable).
		Select("video_id").
		Where("experience_id = ?", experienceID)
	stmt := v.db.DB.WithContext(ctx).
		Table(videoTable).
		Where("id IN (?)", sub)

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

func (v *video) Get(ctx context.Context, videoID string, fields ...string) (*entity.Video, error) {
	video, err := v.get(ctx, v.db.DB, videoID, fields...)
	return video, dbError(err)
}

func (v *video) Create(ctx context.Context, video *entity.Video) error {
	err := v.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := v.now()
		video.CreatedAt, video.UpdatedAt = now, now
		video.VideoProducts = entity.NewVideoProducts(video.ID, video.ProductIDs)
		video.VideoExperiences = entity.NewVideoExperiences(video.ID, video.ExperienceIDs)

		if err := tx.WithContext(ctx).Table(videoTable).Create(&video).Error; err != nil {
			return err
		}
		if err := v.replaceProducts(ctx, tx, video.ID, video.VideoProducts); err != nil {
			return err
		}
		if err := v.replaceExperiences(ctx, tx, video.ID, video.VideoExperiences); err != nil {
			return err
		}
		return nil
	})
	return dbError(err)
}

func (v *video) Update(
	ctx context.Context,
	videoID string,
	params *database.UpdateVideoParams,
) error {
	err := v.db.Transaction(ctx, func(tx *gorm.DB) error {
		products := entity.NewVideoProducts(videoID, params.ProductIDs)
		experiences := entity.NewVideoExperiences(videoID, params.ExperienceIDs)

		updates := map[string]interface{}{
			"title":              params.Title,
			"description":        params.Description,
			"thumbnail_url":      params.ThumbnailURL,
			"video_url":          params.VideoURL,
			"public":             params.Public,
			"limited":            params.Limited,
			"display_product":    params.DisplayProduct,
			"display_experience": params.DisplayExperience,
			"published_at":       params.PublishedAt,
			"updated_at":         v.now(),
		}
		stmt := tx.WithContext(ctx).Table(videoTable).Where("id = ?", videoID)
		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}

		if err := v.replaceProducts(ctx, tx, videoID, products); err != nil {
			return err
		}
		if err := v.replaceExperiences(ctx, tx, videoID, experiences); err != nil {
			return err
		}
		return nil
	})
	return dbError(err)
}

func (v *video) Delete(ctx context.Context, videoID string) error {
	stmt := v.db.DB.WithContext(ctx).Table(videoTable).Where("id = ?", videoID)
	err := stmt.Delete(&entity.Video{}).Error
	return dbError(err)
}

func (v *video) get(
	ctx context.Context,
	tx *gorm.DB,
	videoID string,
	fields ...string,
) (*entity.Video, error) {
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

func (v *video) replaceProducts(
	ctx context.Context,
	tx *gorm.DB,
	videoID string,
	products entity.VideoProducts,
) error {
	// 不要なレコードを削除
	stmt := tx.WithContext(ctx).Where("video_id = ?", videoID)
	if len(products.ProductIDs()) > 0 {
		stmt = stmt.Where("product_id NOT IN (?)", products.ProductIDs())
	}
	if err := stmt.Delete(&entity.VideoProduct{}).Error; err != nil {
		return err
	}

	// レコードの登録/更新
	if len(products) == 0 {
		return nil
	}
	for _, product := range products {
		params := map[string]interface{}{
			"video_id":   product.VideoID,
			"product_id": product.ProductID,
			"priority":   product.Priority,
			"updated_at": v.now(),
		}
		conds := clause.OnConflict{
			Columns:   []clause.Column{{Name: "video_id"}, {Name: "product_id"}},
			DoUpdates: clause.Assignments(params),
		}
		if err := tx.WithContext(ctx).Omit(clause.Associations).Clauses(conds).Create(&product).Error; err != nil {
			return err
		}
	}
	return nil
}

func (v *video) replaceExperiences(
	ctx context.Context,
	tx *gorm.DB,
	videoID string,
	experiences entity.VideoExperiences,
) error {
	// 不要なレコードを削除
	stmt := tx.WithContext(ctx).Where("video_id = ?", videoID)
	if len(experiences.ExperienceIDs()) > 0 {
		stmt = stmt.Where("experience_id NOT IN (?)", experiences.ExperienceIDs())
	}
	if err := stmt.Delete(&entity.VideoExperience{}).Error; err != nil {
		return err
	}

	// レコードの登録/更新
	if len(experiences) == 0 {
		return nil
	}
	for _, experience := range experiences {
		params := map[string]interface{}{
			"video_id":      experience.VideoID,
			"experience_id": experience.ExperienceID,
			"priority":      experience.Priority,
			"updated_at":    v.now(),
		}
		conds := clause.OnConflict{
			Columns:   []clause.Column{{Name: "video_id"}, {Name: "experience_id"}},
			DoUpdates: clause.Assignments(params),
		}
		if err := tx.WithContext(ctx).Omit(clause.Associations).Clauses(conds).Create(&experience).Error; err != nil {
			return err
		}
	}
	return nil
}
