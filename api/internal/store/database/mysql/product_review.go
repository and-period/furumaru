package mysql

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

const productReviewTable = "product_reviews"

type productReview struct {
	db  *mysql.Client
	now func() time.Time
}

func newProductReview(db *mysql.Client) database.ProductReview {
	return &productReview{
		db:  db,
		now: time.Now,
	}
}

type listProductReviewsParams database.ListProductReviewsParams

func (p listProductReviewsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.ProductID != "" {
		stmt = stmt.Where("product_id = ?", p.ProductID)
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

func (r *productReview) List(
	ctx context.Context, params *database.ListProductReviewsParams, fields ...string,
) (entity.ProductReviews, string, error) {
	var reviews entity.ProductReviews

	p := listProductReviewsParams(*params)

	stmt := r.db.Statement(ctx, r.db.DB, productReviewTable, fields...)
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

func (r *productReview) Get(ctx context.Context, reviewID string, fields ...string) (*entity.ProductReview, error) {
	var review *entity.ProductReview

	stmt := r.db.Statement(ctx, r.db.DB, productReviewTable, fields...).
		Where("id = ?", reviewID)

	if err := stmt.First(&review).Error; err != nil {
		return nil, dbError(err)
	}
	if err := r.fill(ctx, r.db.DB, review); err != nil {
		return nil, dbError(err)
	}
	return review, nil
}

func (r *productReview) Create(ctx context.Context, review *entity.ProductReview) error {
	now := r.now()
	review.CreatedAt, review.UpdatedAt = now, now

	err := r.db.DB.WithContext(ctx).Table(productReviewTable).Create(&review).Error
	return dbError(err)
}

func (r *productReview) Update(ctx context.Context, reviewID string, params *database.UpdateProductReviewParams) error {
	updates := map[string]interface{}{
		"rate":       params.Rate,
		"title":      params.Title,
		"comment":    params.Comment,
		"updated_at": r.now(),
	}
	stmt := r.db.DB.WithContext(ctx).Table(productReviewTable).Where("id = ?", reviewID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (r *productReview) Delete(ctx context.Context, productReviewID string) error {
	stmt := r.db.DB.WithContext(ctx).Table(productReviewTable).Where("id = ?", productReviewID)

	err := stmt.Delete(&entity.ProductReview{}).Error
	return dbError(err)
}

func (r *productReview) Aggregate(
	ctx context.Context, params *database.AggregateProductReviewsParams,
) (entity.AggregatedProductReviews, error) {
	var reviews entity.AggregatedProductReviews

	fields := []string{
		"product_id",
		"COUNT(*) as count",
		"AVG(rate) as average",
		"SUM(CASE WHEN rate = 1 THEN 1 ELSE 0 END) AS rate1",
		"SUM(CASE WHEN rate = 2 THEN 1 ELSE 0 END) AS rate2",
		"SUM(CASE WHEN rate = 3 THEN 1 ELSE 0 END) AS rate3",
		"SUM(CASE WHEN rate = 4 THEN 1 ELSE 0 END) AS rate4",
		"SUM(CASE WHEN rate = 5 THEN 1 ELSE 0 END) AS rate5",
	}

	stmt := r.db.Statement(ctx, r.db.DB, productReviewTable, fields...).
		Where("product_id IN (?)", params.ProductIDs).
		Group("product_id")

	err := stmt.Scan(&reviews).Error
	return reviews, dbError(err)
}

func (r *productReview) fill(ctx context.Context, tx *gorm.DB, reviews ...*entity.ProductReview) error {
	var reactions entity.AggregatedProductReviewReactions

	ids := entity.ProductReviews(reviews).IDs()
	if len(ids) == 0 {
		return nil
	}

	fields := []string{
		"review_id",
		"reaction_type",
		"COUNT(*) as total",
	}

	stmt := r.db.Statement(ctx, tx, productReviewReactionTable, fields...).
		Where("preview_id IN (?)", ids).
		Group("review_id, reaction_type")

	if err := stmt.Find(&reactions).Error; err != nil {
		return err
	}
	if len(reactions) == 0 {
		return nil
	}
	entity.ProductReviews(reviews).SetReactions(reactions.GroupByReviewID())

	return nil
}
