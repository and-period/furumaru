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

func (t *thread) Get(ctx context.Context, threadID string, fields ...string) (*entity.Thread, error) {
	thread, err := t.get(ctx, t.db.DB, threadID, fields...)
	return thread, exception.InternalError(err)
}

func (t *thread) listByContactID(
	ctx context.Context, tx *gorm.DB, contactID string, fields ...string,
) (entity.Threads, error) {
	var threads entity.Threads

	err := t.db.Statement(ctx, tx, threadTable, fields...).
		Where("contact_id = ?", contactID).
		Find(&threads).Error
	if err != nil {
		return nil, err
	}
	return threads, nil
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
