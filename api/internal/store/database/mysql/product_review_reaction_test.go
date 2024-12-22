package mysql

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

func TestProductReviewReaction(t *testing.T) {
	assert.NotNil(t, newProductReviewReaction(nil))
}

func TestProductReviewReaction_Upsert(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)
	p := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
	err = db.DB.Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ProductRevision).Error
	require.NoError(t, err)

	review := testProductReview("review-id", "product-id", "user-id", now())
	err = db.DB.Create(&review).Error
	require.NoError(t, err)

	type args struct {
		reaction *entity.ProductReviewReaction
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
				reaction: testProductReviewReaction("review-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "success update",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				reaction := testProductReviewReaction("review-id", "user-id", now())
				err := db.DB.Create(&reaction).Error
				require.NoError(t, err)
			},
			args: args{
				reaction: testProductReviewReaction("review-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, productReviewReactionTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productReviewReaction{db: db, now: now}
			err = db.Upsert(ctx, tt.args.reaction)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestProductReviewReaction_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)
	p := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
	err = db.DB.Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ProductRevision).Error
	require.NoError(t, err)

	review := testProductReview("review-id", "product-id", "user-id", now())
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
				reaction := testProductReviewReaction("review-id", "user-id", now())
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, productReviewReactionTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productReviewReaction{db: db, now: now}
			err = db.Delete(ctx, tt.args.reviewID, tt.args.userID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestProductReviewReaction_GetUserReactions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}

	err := deleteAll(ctx)
	require.NoError(t, err)

	category := testCategory("category-id", "野菜", now())
	err = db.DB.Create(&category).Error
	require.NoError(t, err)
	productType := testProductType("type-id", "category-id", "野菜", now())
	err = db.DB.Create(&productType).Error
	require.NoError(t, err)
	productTag := testProductTag("tag-id", "贈答品", now())
	err = db.DB.Create(&productTag).Error
	require.NoError(t, err)
	p := testProduct("product-id", "type-id", "category-id", "coordinator-id", "producer-id", []string{"tag-id"}, 1, now())
	err = db.DB.Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ProductRevision).Error
	require.NoError(t, err)

	review := testProductReview("review-id", "product-id", "user-id", now())
	err = db.DB.Create(&review).Error
	require.NoError(t, err)

	reaction := testProductReviewReaction("review-id", "user-id", now())
	err = db.DB.Create(&reaction).Error
	require.NoError(t, err)

	type args struct {
		productID string
		userID    string
	}
	type want struct {
		reactions entity.ProductReviewReactions
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
				productID: "product-id",
				userID:    "user-id",
			},
			want: want{
				reactions: entity.ProductReviewReactions{reaction},
				err:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &productReviewReaction{db: db, now: now}
			actual, err := db.GetUserReactions(ctx, tt.args.productID, tt.args.userID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.reactions, actual)
		})
	}
}

func testProductReviewReaction(reviewID, userID string, now time.Time) *entity.ProductReviewReaction {
	return &entity.ProductReviewReaction{
		ReviewID:     reviewID,
		UserID:       userID,
		ReactionType: entity.ProductReviewReactionTypeLike,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
