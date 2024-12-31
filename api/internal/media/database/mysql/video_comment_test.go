package mysql

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestVideoComment(t *testing.T) {
	assert.NotNil(t, NewVideoComment(nil))
}

func TestVideoComment_List(t *testing.T) {
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

	vide := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	err = db.DB.Create(&vide).Error
	require.NoError(t, err)

	comments := make(entity.VideoComments, 2)
	comments[0] = testVideoComment("comment-id01", "video-id", "user-id", now().Add(-time.Minute))
	comments[1] = testVideoComment("comment-id02", "video-id", "user-id", now())
	err = db.DB.Create(&comments).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListVideoCommentsParams
	}
	type want struct {
		comments entity.VideoComments
		token    string
		err      error
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
				params: &database.ListVideoCommentsParams{
					VideoID:      "video-id",
					WithDisabled: false,
					CreatedAtGte: now().Add(-time.Hour),
					CreatedAtLt:  now().Add(time.Hour),
					Limit:        1,
					NextToken:    "",
				},
			},
			want: want{
				comments: comments[1:],
				token:    strconv.FormatInt(comments[0].CreatedAt.UnixNano(), 10),
				err:      nil,
			},
		},
		{
			name:  "success with next token",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListVideoCommentsParams{
					VideoID:      "video-id",
					WithDisabled: false,
					CreatedAtGte: now().Add(-time.Hour),
					CreatedAtLt:  now().Add(time.Hour),
					Limit:        2,
					NextToken:    strconv.FormatInt(comments[1].CreatedAt.UnixNano(), 10),
				},
			},
			want: want{
				comments: comments[1:],
				token:    "",
				err:      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &videoComment{db: db, now: now}
			comments, token, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.comments, comments)
			assert.Equal(t, tt.want.token, token)
		})
	}
}

func TestVideoComment_Create(t *testing.T) {
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

	vide := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	err = db.DB.Create(&vide).Error
	require.NoError(t, err)

	type args struct {
		comment *entity.VideoComment
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
				comment: testVideoComment("comment-id", "video-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "duplicate",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				comment := testVideoComment("comment-id", "video-id", "user-id", now())
				err := db.DB.Create(&comment).Error
				require.NoError(t, err)
			},
			args: args{
				comment: testVideoComment("comment-id", "video-id", "user-id", now()),
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

			err := delete(ctx, videoCommentTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &videoComment{db: db, now: now}
			err = db.Create(ctx, tt.args.comment)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestVideoComment_Update(t *testing.T) {
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

	vide := testVideo("video-id", "coordinator-id", []string{"product-id"}, []string{"experience-id"}, now())
	err = db.DB.Create(&vide).Error
	require.NoError(t, err)

	type args struct {
		commentID string
		params    *database.UpdateVideoCommentParams
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
				comment := testVideoComment("comment-id", "video-id", "user-id", now())
				err := db.DB.Create(&comment).Error
				require.NoError(t, err)
			},
			args: args{
				commentID: "comment-id",
				params: &database.UpdateVideoCommentParams{
					Disabled: false,
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

			err := delete(ctx, videoCommentTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &videoComment{db: db, now: now}
			err = db.Update(ctx, tt.args.commentID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testVideoComment(commentID, videoID, userID string, now time.Time) *entity.VideoComment {
	return &entity.VideoComment{
		ID:        commentID,
		VideoID:   videoID,
		UserID:    userID,
		Content:   "とても面白いですね",
		Disabled:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
