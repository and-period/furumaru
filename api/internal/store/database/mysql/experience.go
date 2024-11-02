package mysql

import (
	"context"
	"errors"
	"fmt"
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
	db  *mysql.Client
	now func() time.Time
}

func newExperience(db *mysql.Client) database.Experience {
	return &experience{
		db:  db,
		now: jst.Now,
	}
}

type listExperiencesParams database.ListExperiencesParams

func (p listExperiencesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("MATCH (`title`, `description`) AGAINST (? IN NATURAL LANGUAGE MODE)", p.Name)
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
	var internal internalExperiences

	p := listExperiencesParams(*params)

	stmt := e.db.Statement(ctx, e.db.DB, experienceTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	experiences := internal.entities()

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

func (e *experience) MultiGet(ctx context.Context, experienceIDs []string, fields ...string) (entity.Experiences, error) {
	experiences, err := e.multiGet(ctx, e.db.DB, experienceIDs, fields...)
	return experiences, dbError(err)
}

func (e *experience) MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Experiences, error) {
	var revisions entity.ExperienceRevisions

	stmt := e.db.Statement(ctx, e.db.DB, experienceRevisionTable).
		Where("id IN (?)", revisionIDs)

	if err := stmt.Find(&revisions).Error; err != nil {
		return nil, dbError(err)
	}
	if len(revisions) == 0 {
		return entity.Experiences{}, nil
	}

	experiences, err := e.multiGet(ctx, e.db.DB, revisions.ExperienceIDs(), fields...)
	if err != nil {
		return nil, err
	}
	if len(experiences) == 0 {
		return entity.Experiences{}, nil
	}

	res, err := revisions.Merge(experiences.Map())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *experience) Get(ctx context.Context, experienceID string, fields ...string) (*entity.Experience, error) {
	experience, err := e.get(ctx, e.db.DB, experienceID, fields...)
	return experience, dbError(err)
}

func (e *experience) Create(ctx context.Context, experience *entity.Experience) error {
	err := e.db.Transaction(ctx, func(tx *gorm.DB) error {
		if err := experience.FillJSON(); err != nil {
			return err
		}

		now := e.now()

		experience.CreatedAt, experience.UpdatedAt = now, now
		experience.ExperienceRevision.CreatedAt, experience.ExperienceRevision.UpdatedAt = now, now

		internal := newInternalExperience(experience)

		if err := tx.WithContext(ctx).Table(experienceTable).Create(&internal).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Table(experienceRevisionTable).Create(&internal.ExperienceRevision).Error
	})
	return dbError(err)
}

func (e *experience) Update(ctx context.Context, experienceID string, params *database.UpdateExperienceParams) error {
	now := e.now()
	rparams := &entity.NewExperienceRevisionParams{
		ExperienceID:          experienceID,
		PriceAdult:            params.PriceAdult,
		PriceJuniorHighSchool: params.PriceJuniorHighSchool,
		PriceElementarySchool: params.PriceElementarySchool,
		PricePreschool:        params.PricePreschool,
		PriceSenior:           params.PriceSenior,
	}
	revision := entity.NewExperienceRevision(rparams)

	err := e.db.Transaction(ctx, func(tx *gorm.DB) error {
		media, err := params.Media.Marshal()
		if err != nil {
			return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
		}
		points, err := entity.ExperienceMarshalRecommendedPoints(params.RecommendedPoints)
		if err != nil {
			return fmt.Errorf("database: %w: %s", database.ErrInvalidArgument, err.Error())
		}
		openTime, err := jst.ParseFromHHMM(params.BusinessOpenTime)
		if err != nil {
			return fmt.Errorf("entity: invalid business open time: %w", err)
		}
		closeTime, err := jst.ParseFromHHMM(params.BusinessCloseTime)
		if err != nil {
			return fmt.Errorf("entity: invalid business close time: %w", err)
		}
		if !openTime.Before(closeTime) {
			return errors.New("entity: invalid business time")
		}

		updates := map[string]interface{}{
			"experience_type_id":  params.TypeID,
			"title":               params.Title,
			"description":         params.Description,
			"public":              params.Public,
			"sold_out":            params.SoldOut,
			"media":               nil,
			"recommended_points":  points,
			"promotion_video_url": params.PromotionVideoURL,
			"duration":            params.Duration,
			"direction":           params.Direction,
			"business_open_time":  params.BusinessOpenTime,
			"business_close_time": params.BusinessCloseTime,
			"host_postal_code":    params.HostPostalCode,
			"host_prefecture":     params.HostPrefectureCode,
			"host_city":           params.HostCity,
			"host_address_line1":  params.HostAddressLine1,
			"host_address_line2":  params.HostAddressLine2,
			"host_geolocation":    newExperienceHostGeolocation(params.HostLongitude, params.HostLatitude),
			"start_at":            params.StartAt,
			"end_at":              params.EndAt,
			"updated_at":          now,
		}
		if len(media) > 0 {
			updates["media"] = media
		}

		stmt := tx.WithContext(ctx).Table(experienceTable).Where("id = ?", experienceID)
		if err := stmt.Updates(updates).Error; err != nil {
			return err
		}

		revision.CreatedAt, revision.UpdatedAt = now, now
		return tx.WithContext(ctx).Table(experienceRevisionTable).Create(&revision).Error
	})
	return dbError(err)
}

func (e *experience) Delete(ctx context.Context, experienceID string) error {
	params := map[string]interface{}{
		"deleted_at": e.now(),
	}
	stmt := e.db.DB.WithContext(ctx).Table(experienceTable).Where("id = ?", experienceID)
	err := stmt.Updates(params).Error
	return dbError(err)
}

func (e *experience) multiGet(ctx context.Context, tx *gorm.DB, experienceIDs []string, fields ...string) (entity.Experiences, error) {
	var internal internalExperiences

	stmt := e.db.Statement(ctx, tx, experienceTable, fields...).
		Where("id IN (?)", experienceIDs)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, err
	}
	experiences := internal.entities()

	if err := e.fill(ctx, tx, experiences...); err != nil {
		return nil, err
	}
	return experiences, nil
}

func (e *experience) get(ctx context.Context, tx *gorm.DB, experienceID string, fields ...string) (*entity.Experience, error) {
	var internal *internalExperience

	stmt := e.db.Statement(ctx, tx, experienceTable, fields...).
		Where("id = ?", experienceID)

	if err := stmt.First(&internal).Error; err != nil {
		return nil, err
	}
	experience := internal.entity()

	if err := e.fill(ctx, tx, experience); err != nil {
		return nil, err
	}
	return experience, nil
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

type internalExperience struct {
	entity.Experience `gorm:"embedded"`
	HostGeolocation   mysql.Geometry `gorm:""`  // 開催場所(座標情報)
	HostLongitude     float64        `gorm:"-"` // 開催場所(座標情報:経度)
	HostLatitude      float64        `gorm:"-"` // 開催場所(座標情報:緯度)
}

type internalExperiences []*internalExperience

func newInternalExperience(experience *entity.Experience) *internalExperience {
	return &internalExperience{
		Experience:      *experience,
		HostGeolocation: newExperienceHostGeolocation(experience.HostLongitude, experience.HostLatitude),
		HostLongitude:   experience.HostLongitude,
		HostLatitude:    experience.HostLatitude,
	}
}

func newExperienceHostGeolocation(longitude, latitude float64) mysql.Geometry {
	return mysql.Geometry{
		X: longitude,
		Y: latitude,
	}
}

func (e *internalExperience) entity() *entity.Experience {
	if e == nil {
		return nil
	}
	e.Experience.HostLongitude = e.HostGeolocation.X
	e.Experience.HostLatitude = e.HostGeolocation.Y
	return &e.Experience
}

func (es internalExperiences) entities() entity.Experiences {
	res := make(entity.Experiences, len(es))
	for i := range es {
		res[i] = es[i].entity()
	}
	return res
}
