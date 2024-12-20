package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductReview(t *testing.T) {
	assert.NotNil(t, newProductReview(nil))
}

func TestProductReview_Get(t *testing.T) {
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
	}
	type want struct {
		review *entity.ProductReview
		err    error
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
				reviewID: "review-id",
			},
			want: want{
				review: review,
				err:    nil,
			},
		},
		{
			name:  "not found",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				reviewID: "",
			},
			want: want{
				review: nil,
				err:    database.ErrNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &productReview{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.reviewID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.review, actual)
		})
	}
}

func TestProductReview_Create(t *testing.T) {
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

	type args struct {
		review *entity.ProductReview
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
			name:  "success",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				review: testProductReview("review-id", "product-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				review := testProductReview("review-id", "product-id", "user-id", now())
				err := db.DB.Create(&review).Error
				require.NoError(t, err)
			},
			args: args{
				review: testProductReview("review-id", "product-id", "user-id", now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, productReviewTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productReview{db: db, now: now}
			err = db.Create(ctx, tt.args.review)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestProductReview_Update(t *testing.T) {
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

	type args struct {
		reviewID string
		params   *database.UpdateProductReviewParams
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
				review := testProductReview("review-id", "product-id", "user-id", now())
				err := db.DB.Create(&review).Error
				require.NoError(t, err)
			},
			args: args{
				reviewID: "review-id",
				params: &database.UpdateProductReviewParams{
					Rate:    4,
					Title:   "おすすめの商品です",
					Comment: "とても良い商品でした",
				},
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

			err := delete(ctx, productReviewTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productReview{db: db, now: now}
			err = db.Update(ctx, tt.args.reviewID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestProductReview_Delete(t *testing.T) {
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

	type args struct {
		reviewID string
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
				review := testProductReview("review-id", "product-id", "user-id", now())
				err := db.DB.Create(&review).Error
				require.NoError(t, err)
			},
			args: args{
				reviewID: "review-id",
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

			err := delete(ctx, productReviewTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &productReview{db: db, now: now}
			err = db.Delete(ctx, tt.args.reviewID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testProductReview(id, productID, userID string, now time.Time) *entity.ProductReview {
	return &entity.ProductReview{
		ID:        id,
		ProductID: productID,
		UserID:    userID,
		Rate:      5,
		Title:     "おすすめの商品です",
		Comment:   "とても良い商品でした",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
