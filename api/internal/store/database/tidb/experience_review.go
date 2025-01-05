package tidb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const experienceReviewTable = "experience_reviews"

type experienceReview struct {
	db  *mysql.Client
	now func() time.Time
}

func NewExperienceReview(db *mysql.Client) database.ExperienceReview {
	return &experienceReview{
		db:  db,
		now: time.Now,
	}
}

type listExperienceReviewsParams database.ListExperienceReviewsParams

func (p listExperienceReviewsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.ExperienceID != "" {
		stmt = stmt.Where("experience_id = ?", p.ExperienceID)
	}
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
	}
	if len(p.Rates) > 0 {
		stmt = stmt.Where("rate IN (?)", p.Rates)
	}
	stmt = stmt.Order("created_at DESC")
	return stmt
}

func (r *experienceReview) List(
	ctx context.Context, params *database.ListExperienceReviewsParams, fields ...string,
) (entity.ExperienceReviews, string, error) {
	var reviews entity.ExperienceReviews

	p := listExperienceReviewsParams(*params)

	stmt := r.db.Statement(ctx, r.db.DB, experienceReviewTable, fields...)
	stmt = p.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(int(p.Limit) + 1)
	}
	if params.NextToken != "" {
		nsec, err := strconv.ParseInt(params.NextToken, 10, 64)
		if err != nil {
			return nil, "", fmt.Errorf("database: failed to parse next token: %s: %w", err.Error(), database.ErrInvalidArgument)
		}
		stmt = stmt.Where("created_at <= ?", time.Unix(0, nsec))
	}

	if err := stmt.Find(&reviews).Error; err != nil {
		return nil, "", dbError(err)
	}
	if err := r.fill(ctx, r.db.DB, reviews...); err != nil {
		return nil, "", dbError(err)
	}

	var nextToken string
	if len(reviews) > int(params.Limit) {
		nextToken = strconv.FormatInt(reviews[params.Limit].CreatedAt.UnixNano(), 10)
		reviews = reviews[:len(reviews)-1]
	}
	return reviews, nextToken, nil
}

func (r *experienceReview) Get(ctx context.Context, reviewID string, fields ...string) (*entity.ExperienceReview, error) {
	var review *entity.ExperienceReview

	stmt := r.db.Statement(ctx, r.db.DB, experienceReviewTable, fields...).
		Where("id = ?", reviewID)

	if err := stmt.First(&review).Error; err != nil {
		return nil, dbError(err)
	}
	if err := r.fill(ctx, r.db.DB, review); err != nil {
		return nil, dbError(err)
	}
	return review, nil
}

func (r *experienceReview) Create(ctx context.Context, review *entity.ExperienceReview) error {
	now := r.now()
	review.CreatedAt, review.UpdatedAt = now, now

	err := r.db.DB.WithContext(ctx).Table(experienceReviewTable).Create(&review).Error
	return dbError(err)
}

func (r *experienceReview) Update(ctx context.Context, reviewID string, params *database.UpdateExperienceReviewParams) error {
	updates := map[string]interface{}{
		"rate":       params.Rate,
		"title":      params.Title,
		"comment":    params.Comment,
		"updated_at": r.now(),
	}
	stmt := r.db.DB.WithContext(ctx).Table(experienceReviewTable).Where("id = ?", reviewID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (r *experienceReview) Delete(ctx context.Context, experienceReviewID string) error {
	stmt := r.db.DB.WithContext(ctx).Table(experienceReviewTable).Where("id = ?", experienceReviewID)

	err := stmt.Delete(&entity.ExperienceReview{}).Error
	return dbError(err)
}

func (r *experienceReview) Aggregate(
	ctx context.Context, params *database.AggregateExperienceReviewsParams,
) (entity.AggregatedExperienceReviews, error) {
	var reviews entity.AggregatedExperienceReviews

	fields := []string{
		"experience_id",
		"COUNT(*) as count",
		"AVG(rate) as average",
		"SUM(CASE WHEN rate = 1 THEN 1 ELSE 0 END) AS rate1",
		"SUM(CASE WHEN rate = 2 THEN 1 ELSE 0 END) AS rate2",
		"SUM(CASE WHEN rate = 3 THEN 1 ELSE 0 END) AS rate3",
		"SUM(CASE WHEN rate = 4 THEN 1 ELSE 0 END) AS rate4",
		"SUM(CASE WHEN rate = 5 THEN 1 ELSE 0 END) AS rate5",
	}

	stmt := r.db.Statement(ctx, r.db.DB, experienceReviewTable, fields...).
		Where("experience_id IN (?)", params.ExperienceIDs).
		Group("experience_id")

	err := stmt.Scan(&reviews).Error
	return reviews, dbError(err)
}

func (r *experienceReview) fill(ctx context.Context, tx *gorm.DB, reviews ...*entity.ExperienceReview) error {
	var reactions entity.AggregatedExperienceReviewReactions

	ids := entity.ExperienceReviews(reviews).IDs()
	if len(ids) == 0 {
		return nil
	}

	fields := []string{
		"review_id",
		"reaction_type",
		"COUNT(*) as total",
	}

	stmt := r.db.Statement(ctx, tx, experienceReviewReactionTable, fields...).
		Where("review_id IN (?)", ids).
		Group("review_id, reaction_type")

	if err := stmt.Find(&reactions).Error; err != nil {
		return err
	}
	if len(reactions) == 0 {
		return nil
	}
	entity.ExperienceReviews(reviews).SetReactions(reactions.GroupByReviewID())

	return nil
}
