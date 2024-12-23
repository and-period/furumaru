package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm/clause"
)

const experienceReviewReactionTable = "experience_review_reactions"

type experienceReviewReaction struct {
	db  *mysql.Client
	now func() time.Time
}

func newExperienceReviewReaction(db *mysql.Client) database.ExperienceReviewReaction {
	return &experienceReviewReaction{
		db:  db,
		now: time.Now,
	}
}

func (r *experienceReviewReaction) Upsert(ctx context.Context, reaction *entity.ExperienceReviewReaction) error {
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

func (r *experienceReviewReaction) Delete(ctx context.Context, experienceReviewID, userID string) error {
	stmt := r.db.DB.WithContext(ctx).Where("review_id = ? AND user_id = ?", experienceReviewID, userID)

	err := stmt.Delete(&entity.ExperienceReviewReaction{}).Error
	return dbError(err)
}

func (r *experienceReviewReaction) GetUserReactions(ctx context.Context, experienceID, userID string) (entity.ExperienceReviewReactions, error) {
	var reactions entity.ExperienceReviewReactions

	stmt := r.db.Statement(ctx, r.db.DB, experienceReviewReactionTable).
		Joins("JOIN experience_reviews ON experience_review_reactions.review_id = experience_reviews.id").
		Where("experience_reviews.experience_id = ?", experienceID).
		Where("experience_review_reactions.user_id = ?", userID)

	err := stmt.Find(&reactions).Error
	return reactions, dbError(err)
}
