package mysql

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const broadcastCommentTable = "broadcast_comments"

type broadcastComment struct {
	db  *mysql.Client
	now func() time.Time
}

func newBroadcastComment(db *mysql.Client) database.BroadcastComment {
	return &broadcastComment{
		db:  db,
		now: jst.Now,
	}
}

func (c *broadcastComment) List(
	ctx context.Context,
	params *database.ListBroadcastCommentsParams,
	fields ...string,
) (entity.BroadcastComments, string, error) {
	var comments entity.BroadcastComments

	stmt := c.db.Statement(ctx, c.db.DB, broadcastCommentTable, fields...).
		Where("broadcast_id = ?", params.BroadcastID).
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
	for i := range params.Orders {
		var value string
		if params.Orders[i].OrderByASC {
			value = fmt.Sprintf("`%s` ASC", params.Orders[i].Key)
		} else {
			value = fmt.Sprintf("`%s` DESC", params.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	if len(params.Orders) == 0 {
		stmt = stmt.Order("created_at DESC")
	}

	if err := stmt.Find(&comments).Error; err != nil {
		return nil, "", dbError(err)
	}
	var nextToken string
	if len(comments) > int(params.Limit) {
		nextToken = strconv.FormatInt(comments[params.Limit].CreatedAt.UnixNano(), 10)
		comments = comments[:params.Limit]
	}
	return comments, nextToken, nil
}

func (c *broadcastComment) Create(ctx context.Context, comment *entity.BroadcastComment) error {
	now := c.now()
	comment.CreatedAt, comment.UpdatedAt = now, now

	err := c.db.DB.WithContext(ctx).Table(broadcastCommentTable).Create(&comment).Error
	return dbError(err)
}
