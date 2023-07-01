package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const threadTable = "threads"

type thread struct {
	db  *database.Client
	now func() time.Time
}

func NewThread(db *database.Client) Thread {
	return &thread{
		db:  db,
		now: jst.Now,
	}
}

func (t *thread) ListByContactID(
	ctx context.Context, params *ListThreadsByContactIDParams, fields ...string,
) (entity.Threads, error) {
	var threads entity.Threads

	stmt := t.db.Statement(ctx, t.db.DB, threadTable, fields...)
	stmt = params.stmt(stmt)
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}
	if err := stmt.Find(&threads).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return threads, nil
}

func (t *thread) Count(ctx context.Context, params *ListThreadsByContactIDParams) (int64, error) {
	total, err := t.db.Count(ctx, t.db.DB, &entity.Thread{}, params.stmt)
	return total, exception.InternalError(err)
}

func (t *thread) Create(ctx context.Context, thread *entity.Thread) error {
	err := t.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := t.now()
		thread.CreatedAt, thread.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(threadTable).Create(&thread).Error
		return err
	})
	return exception.InternalError(err)
}

func (t *thread) Get(ctx context.Context, threadID string, fields ...string) (*entity.Thread, error) {
	thread, err := t.get(ctx, t.db.DB, threadID, fields...)
	return thread, exception.InternalError(err)
}

func (t *thread) Update(ctx context.Context, threadID string, params *UpdateThreadParams) error {
	err := t.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := t.get(ctx, tx, threadID); err != nil {
			return err
		}
		updates := map[string]interface{}{
			"content":    params.Content,
			"user_id":    params.UserID,
			"user_type":  params.UserType,
			"updated_at": t.now(),
		}
		err := tx.WithContext(ctx).
			Table(threadTable).
			Where("id = ?", threadID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
}

func (t *thread) Delete(ctx context.Context, threadID string) error {
	err := t.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := t.get(ctx, tx, threadID); err != nil {
			return err
		}

		params := map[string]interface{}{
			"deleted_at": t.now(),
		}
		err := tx.WithContext(ctx).
			Table(threadTable).
			Where("id = ?", threadID).
			Updates(params).Error
		return err
	})
	return exception.InternalError(err)
}

func (t *thread) get(
	ctx context.Context, tx *gorm.DB, threadID string, fields ...string,
) (*entity.Thread, error) {
	var thread *entity.Thread

	err := t.db.Statement(ctx, tx, threadTable, fields...).
		Where("id = ?", threadID).
		First(&thread).Error
	return thread, err
}
