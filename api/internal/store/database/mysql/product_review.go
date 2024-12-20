package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
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

func (r *productReview) Get(ctx context.Context, reviewID string, fields ...string) (*entity.ProductReview, error) {
	var review *entity.ProductReview

	stmt := r.db.Statement(ctx, r.db.DB, productReviewTable, fields...).
		Where("id = ?", reviewID)

	if err := stmt.First(&review).Error; err != nil {
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
