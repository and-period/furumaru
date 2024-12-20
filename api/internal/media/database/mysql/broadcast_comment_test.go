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

func TestBroadcastComment(t *testing.T) {
	assert.NotNil(t, newBroadcastComment(nil))
}

func TestBroadcastComment_List(t *testing.T) {
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

	broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
	err = db.DB.Create(&broadcast).Error
	require.NoError(t, err)

	comments := make(entity.BroadcastComments, 2)
	comments[0] = testBroadcastComment("comment-id01", "broadcast-id", "user-id", now().Add(-time.Minute))
	comments[1] = testBroadcastComment("comment-id02", "broadcast-id", "user-id", now())
	err = db.DB.Create(&comments).Error
	require.NoError(t, err)

	type args struct {
		params *database.ListBroadcastCommentsParams
	}
	type want struct {
		comments entity.BroadcastComments
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
				params: &database.ListBroadcastCommentsParams{
					BroadcastID:  "broadcast-id",
					WithDisabled: false,
					CreatedAtGte: now().Add(-time.Hour),
					CreatedAtLt:  now().Add(time.Hour),
					Limit:        1,
					NextToken:    "",
				},
			},
			want: want{
				comments: comments[:1],
				token:    strconv.FormatInt(comments[1].CreatedAt.UnixNano(), 10),
				err:      nil,
			},
		},
		{
			name:  "success with next token",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {},
			args: args{
				params: &database.ListBroadcastCommentsParams{
					BroadcastID:  "broadcast-id",
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tt.setup(ctx, t, db)

			db := &broadcastComment{db: db, now: now}
			actual, token, err := db.List(ctx, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.comments, actual)
			assert.Equal(t, tt.want.token, token)
		})
	}
}

func TestBroadcastComment_Create(t *testing.T) {
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

	broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
	err = db.DB.Create(&broadcast).Error
	require.NoError(t, err)

	type args struct {
		comment *entity.BroadcastComment
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
				comment: testBroadcastComment("comment-id", "broadcast-id", "user-id", now()),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "already exists",
			setup: func(ctx context.Context, t *testing.T, db *mysql.Client) {
				comment := testBroadcastComment("comment-id", "broadcast-id", "user-id", now())
				err := db.DB.Create(&comment).Error
				require.NoError(t, err)
			},
			args: args{
				comment: testBroadcastComment("comment-id", "broadcast-id", "user-id", now()),
			},
			want: want{
				err: database.ErrAlreadyExists,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, broadcastCommentTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &broadcastComment{db: db, now: now}
			err = db.Create(ctx, tt.args.comment)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func TestBroadcastComment_Update(t *testing.T) {
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

	broadcast := testBroadcast("broadcast-id", "schedule-id", "coordinator-id", now())
	err = db.DB.Create(&broadcast).Error
	require.NoError(t, err)

	type args struct {
		commentID string
		params    *database.UpdateBroadcastCommentParams
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
				comment := testBroadcastComment("comment-id", "broadcast-id", "user-id", now())
				err = db.DB.Create(&comment).Error
				require.NoError(t, err)
			},
			args: args{
				commentID: "comment-id",
				params: &database.UpdateBroadcastCommentParams{
					Disabled: true,
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := delete(ctx, broadcastCommentTable)
			require.NoError(t, err)

			tt.setup(ctx, t, db)

			db := &broadcastComment{db: db, now: now}
			err = db.Update(ctx, tt.args.commentID, tt.args.params)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func testBroadcastComment(commentID, broadcastID, userID string, now time.Time) *entity.BroadcastComment {
	return &entity.BroadcastComment{
		ID:          commentID,
		BroadcastID: broadcastID,
		UserID:      userID,
		Content:     "こんにちは",
		Disabled:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
