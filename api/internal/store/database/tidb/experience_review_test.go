package tidb

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExperienceReview(t *testing.T) {
	assert.NotNil(t, NewExperienceReview(nil))
}

func TestExperienceReview_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	experienceType := testExperienceType("type-id", "野菜", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	p := testExperience("experience-id", "type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ExperienceRevision).Error
	require.NoError(t, err)

	reviews := make(entity.ExperienceReviews, 3)
	reviews[0] = testExperienceReview("review-id01", "experience-id", "user-id01", now().Add(-time.Hour))
	reviews[1] = testExperienceReview("review-id02", "experience-id", "user-id02", now())
	reviews[2] = testExperienceReview("review-id03", "experience-id", "user-id03", now().Add(time.Hour))
	err = db.DB.Create(&reviews).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListExperienceReviewsParams
	}
	type want struct {
		reviews entity.ExperienceReviews
		token   string
		err     error
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
				params: &database.ListExperienceReviewsParams{
					ExperienceID: "experience-id",
					Limit:        1,
					NextToken:    strconv.FormatInt(now().UnixNano(), 10),
				},
			},
			want: want{
				reviews: reviews[1:2],
				token:   strconv.FormatInt(now().Add(-time.Hour).UnixNano(), 10),
				err:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &experienceReview{db: db, now: now}
			actual, token, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.token, token)
			assert.Equal(t, tt.want.reviews, actual)
		})
	}
}

func TestExperienceReview_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
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
	}
	type want struct {
		review *entity.ExperienceReview
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

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &experienceReview{db: db, now: now}
			actual, err := db.Get(ctx, tt.args.reviewID)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.review, actual)
		})
	}
}

func TestExperienceReview_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	experienceType := testExperienceType("type-id", "野菜", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	p := testExperience("experience-id", "type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ExperienceRevision).Error
	require.NoError(t, err)

	type args struct {
		review *entity.ExperienceReview
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
				review: testExperienceReview("review-id", "experience-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				review := testExperienceReview("review-id", "experience-id", "user-id", now())
				err := db.DB.Create(&review).Error
				require.NoError(t, err)
			},
			args: args{
				review: testExperienceReview("review-id", "experience-id", "user-id", now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			err := delete(ctx, experienceReviewTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experienceReview{db: db, now: now}
			err = db.Create(ctx, tt.args.review)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperienceReview_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	experienceType := testExperienceType("type-id", "野菜", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	p := testExperience("experience-id", "type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ExperienceRevision).Error
	require.NoError(t, err)

	type args struct {
		reviewID string
		params   *database.UpdateExperienceReviewParams
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
				review := testExperienceReview("review-id", "experience-id", "user-id", now())
				err := db.DB.Create(&review).Error
				require.NoError(t, err)
			},
			args: args{
				reviewID: "review-id",
				params: &database.UpdateExperienceReviewParams{
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
			ctx := t.Context()
			err := delete(ctx, experienceReviewTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experienceReview{db: db, now: now}
			err = db.Update(ctx, tt.args.reviewID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperienceReview_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	experienceType := testExperienceType("type-id", "野菜", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	p := testExperience("experience-id", "type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ExperienceRevision).Error
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
				review := testExperienceReview("review-id", "experience-id", "user-id", now())
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
			ctx := t.Context()
			err := delete(ctx, experienceReviewTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &experienceReview{db: db, now: now}
			err = db.Delete(ctx, tt.args.reviewID)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestExperienceReview_Aggregate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := dbClient
	now := func() time.Time {
		return current
	}
	err := deleteAll(t.Context())
	require.NoError(t, err)

	experienceType := testExperienceType("type-id", "野菜", now())
	err = db.DB.Create(&experienceType).Error
	require.NoError(t, err)
	p := testExperience("experience-id", "type-id", "shop-id", "coordinator-id", "producer-id", 1, now())
	err = db.DB.Table(experienceTable).Create(&p).Error
	require.NoError(t, err)
	err = db.DB.Create(&p.ExperienceRevision).Error
	require.NoError(t, err)

	reviews := make(entity.ExperienceReviews, 4)
	reviews[0] = testExperienceReview("review-id01", "experience-id", "user-id01", now())
	reviews[0].Rate = 1
	reviews[1] = testExperienceReview("review-id02", "experience-id", "user-id02", now())
	reviews[1].Rate = 5
	reviews[2] = testExperienceReview("review-id03", "experience-id", "user-id03", now())
	reviews[2].Rate = 3
	reviews[3] = testExperienceReview("review-id04", "experience-id", "user-id04", now())
	reviews[3].Rate = 1
	err = db.DB.Create(&reviews).Error
	require.NoError(t, err)

	type args struct {
		params *database.AggregateExperienceReviewsParams
	}
	type want struct {
		reviews entity.AggregatedExperienceReviews
		err     error
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
				params: &database.AggregateExperienceReviewsParams{
					ExperienceIDs: []string{"experience-id"},
				},
			},
			want: want{
				reviews: entity.AggregatedExperienceReviews{
					{
						ExperienceID: "experience-id",
						Count:        4,
						Average:      2.5,
						Rate1:        2,
						Rate2:        0,
						Rate3:        1,
						Rate4:        0,
						Rate5:        1,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			tt.setup(ctx, t, db)

			db := &experienceReview{db: db, now: now}
			actual, err := db.Aggregate(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.reviews, actual)
		})
	}
}

func testExperienceReview(id, experienceID, userID string, now time.Time) *entity.ExperienceReview {
	return &entity.ExperienceReview{
		ID:           id,
		ExperienceID: experienceID,
		UserID:       userID,
		Rate:         5,
		Title:        "おすすめの商品です",
		Comment:      "とても良い商品でした",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
