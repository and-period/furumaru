package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm/clause"
)

const productReviewReactionTable = "product_review_reactions"

type productReviewReaction struct {
	db  *mysql.Client
	now func() time.Time
}

func NewProductReviewReaction(db *mysql.Client) database.ProductReviewReaction {
	return &productReviewReaction{
		db:  db,
		now: time.Now,
	}
}

func (r *productReviewReaction) Upsert(ctx context.Context, reaction *entity.ProductReviewReaction) error {
	now := r.now()
	reaction.CreatedAt, reaction.UpdatedAt = now, now

	updates := map[string]interface{}{
		"reaction_type": reaction.ReactionType,
		"updated_at":    reaction.UpdatedAt,
	}
	stmt := r.db.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "review_id"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(updates),
	})

	err := stmt.Create(&reaction).Error
	return dbError(err)
}

func (r *productReviewReaction) Delete(ctx context.Context, productReviewID, userID string) error {
	stmt := r.db.DB.WithContext(ctx).Where("review_id = ? AND user_id = ?", productReviewID, userID)

	err := stmt.Delete(&entity.ProductReviewReaction{}).Error
	return dbError(err)
}

func (r *productReviewReaction) GetUserReactions(ctx context.Context, productID, userID string) (entity.ProductReviewReactions, error) {
	var reactions entity.ProductReviewReactions

	stmt := r.db.Statement(ctx, r.db.DB, productReviewReactionTable).
		Joins("JOIN product_reviews ON product_review_reactions.review_id = product_reviews.id").
		Where("product_reviews.product_id = ?", productID).
		Where("product_review_reactions.user_id = ?", userID)

	err := stmt.Find(&reactions).Error
	return reactions, dbError(err)
}
