package tidb

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExperienceReviewReaction(t *testing.T) {
	assert.NotNil(t, NewExperienceReviewReaction(nil))
}

func TestExperienceReviewReaction_Upsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	experienceType := testExperienceType("type-id", "野菜", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	p := testExperience("experience-id", "type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ExperienceRevision).Error
	require.NoError(t, err)

	review := testExperienceReview("review-id", "experience-id", "user-id", now())
	err = db.DB.Create(&review).Error
	require.NoError(t, err)

	type args struct {
		reaction *entity.ExperienceReviewReaction
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success create",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				reaction: testExperienceReviewReaction("review-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success update",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				reaction := testExperienceReviewReaction("review-id", "user-id", now())
				err := db.DB.Create(&reaction).Error
				require.NoError(t, err)
			},
			args: args{
				reaction: testExperienceReviewReaction("review-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, experienceReviewReactionTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experienceReviewReaction{db: db, now: now}
			err = db.Upsert(ctx, tt.args.reaction)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperienceReviewReaction_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	experienceType := testExperienceType("type-id", "野菜", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	p := testExperience("experience-id", "type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ExperienceRevision).Error
	require.NoError(t, err)

	review := testExperienceReview("review-id", "experience-id", "user-id", now())
	err = db.DB.Create(&review).Error
	require.NoError(t, err)

	type args struct {
		reviewID string
		userID   string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name: "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				reaction := testExperienceReviewReaction("review-id", "user-id", now())
				err := db.DB.Create(&reaction).Error
				require.NoError(t, err)
			},
			args: args{
				reviewID: "review-id",
				userID:   "user-id",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, experienceReviewReactionTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experienceReviewReaction{db: db, now: now}
			err = db.Delete(ctx, tt.args.reviewID, tt.args.userID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperienceReviewReaction_GetUserReactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	shop := testShop("shop-id", "coordinator-id", []string{}, []string{}, now())
	err = db.DB.Table(shopTable).Create(&shop).Error
	require.NoError(t, err)

	experienceType := testExperienceType("type-id", "野菜", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	p := testExperience("experience-id", "type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ExperienceRevision).Error
	require.NoError(t, err)

	review := testExperienceReview("review-id", "experience-id", "user-id", now())
	err = db.DB.Create(&review).Error
	require.NoError(t, err)

	reaction := testExperienceReviewReaction("review-id", "user-id", now())
	err = db.DB.Create(&reaction).Error
	require.NoError(t, err)

	type args struct {
		experienceID string
		userID       string
	}
	type want struct {
		reactions entity.ExperienceReviewReactions
		err       error
	}
	tests := []struct {
		name  string
		setup func(ctx context.Context, t *testing.T, db *mysql.Client)
		args  args
		want  want
	}{
		{
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				experienceID: "experience-id",
				userID:       "user-id",
			},
			want: want{
				reactions: entity.ExperienceReviewReactions{reaction},
				err:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &experienceReviewReaction{db: db, now: now}
			actual, err := db.GetUserReactions(ctx, tt.args.experienceID, tt.args.userID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.reactions, actual)
		})
	}
}

func testExperienceReviewReaction(reviewID, userID string, now time.Time) *entity.ExperienceReviewReaction {
	return &entity.ExperienceReviewReaction{
		ReviewID:     reviewID,
		UserID:       userID,
		ReactionType: entity.ExperienceReviewReactionTypeLike,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
