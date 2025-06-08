package tidb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const videoCommentTable = "video_comments"

type videoComment struct {
	db  *mysql.Client
	now func() time.Time
}

func NewVideoComment(db *mysql.Client) database.VideoComment {
	return &videoComment{
		db:  db,
		now: time.Now,
	}
}

func (c *videoComment) List(
	ctx context.Context,
	params *database.ListVideoCommentsParams,
	fields ...string,
) (entity.VideoComments, string, error) {
	var comments entity.VideoComments

	stmt := c.db.Statement(ctx, c.db.DB, videoCommentTable, fields...).
		Where("video_id = ?", params.VideoID).
		Limit(int(params.Limit) + 1)

	if !params.WithDisabled {
		stmt = stmt.Where("disabled = ?", false)
	}
	if !params.CreatedAtGte.IsZero() {
		stmt = stmt.Where("created_at >= ?", params.CreatedAtGte)
	}
	if !params.CreatedAtLt.IsZero() {
		stmt = stmt.Where("created_at < ?", params.CreatedAtLt)
	}
	if params.NextToken != "" {
		nsec, err := strconv.ParseInt(params.NextToken, 10, 64)
		if err != nil {
			return nil, "", fmt.Errorf("database: failed to parse next token: %s: %w", err.Error(), database.ErrInvalidArgument)
		}
		stmt = stmt.Where("created_at >= ?", time.Unix(0, nsec))
	}
	stmt = stmt.Order("created_at ASC")

	if err := stmt.Find(&comments).Error; err != nil {
		return nil, "", dbError(err)
	}
	var nextToken string
	if len(comments) > int(params.Limit) {
		nextToken = strconv.FormatInt(comments[params.Limit].CreatedAt.UnixNano(), 10)
		comments = comments[:len(comments)-1]
	}
	return comments, nextToken, nil
}

func (c *videoComment) Create(ctx context.Context, comment *entity.VideoComment) error {
	now := c.now()
	comment.CreatedAt, comment.UpdatedAt = now, now

	err := c.db.DB.WithContext(ctx).Table(videoCommentTable).Create(comment).Error
	return dbError(err)
}

func (c *videoComment) Update(ctx context.Context, commentID string, params *database.UpdateVideoCommentParams) error {
	updates := map[string]interface{}{
		"disabled":   params.Disabled,
		"updated_at": c.now(),
	}
	stmt := c.db.DB.WithContext(ctx).Table(videoCommentTable).Where("id = ?", commentID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}
