package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const threadTable = "threads"

type thread struct {
	db  *mysql.Client
	now func() time.Time
}

func NewThread(db *mysql.Client) database.Thread {
	return &thread{
		db:  db,
		now: jst.Now,
	}
}

type listThreadsParams database.ListThreadsParams

func (p listThreadsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.ContactID != "" {
		stmt = stmt.Where("contact_id = ?", p.ContactID)
	}
	return stmt
}

func (p listThreadsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (t *thread) List(ctx context.Context, params *database.ListThreadsParams, fields ...string) (entity.Threads, error) {
	var threads entity.Threads

	p := listThreadsParams(*params)

	stmt := t.db.Statement(ctx, t.db.DB, threadTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&threads).Error; err != nil {
		return nil, dbError(err)
	}
	return threads, nil
}

func (t *thread) Count(ctx context.Context, params *database.ListThreadsParams) (int64, error) {
	p := listThreadsParams(*params)

	total, err := t.db.Count(ctx, t.db.DB, &entity.Thread{}, p.stmt)
	return total, dbError(err)
}

func (t *thread) Get(ctx context.Context, threadID string, fields ...string) (*entity.Thread, error) {
	thread, err := t.get(ctx, t.db.DB, threadID, fields...)
	return thread, dbError(err)
}

func (t *thread) Create(ctx context.Context, thread *entity.Thread) error {
	now := t.now()
	thread.CreatedAt, thread.UpdatedAt = now, now

	err := t.db.DB.WithContext(ctx).Table(threadTable).Create(&thread).Error
	return dbError(err)
}

func (t *thread) Update(ctx context.Context, threadID string, params *database.UpdateThreadParams) error {
	updates := map[string]interface{}{
		"content":    params.Content,
		"user_id":    params.UserID,
		"user_type":  params.UserType,
		"updated_at": t.now(),
	}
	stmt := t.db.DB.WithContext(ctx).
		Table(threadTable).
		Where("id = ?", threadID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (t *thread) Delete(ctx context.Context, threadID string) error {
	params := map[string]interface{}{
		"deleted_at": t.now(),
	}
	stmt := t.db.DB.WithContext(ctx).
		Table(threadTable).
		Where("id = ?", threadID)

	err := stmt.Updates(params).Error
	return dbError(err)
}

func (t *thread) get(
	ctx context.Context, tx *gorm.DB, threadID string, fields ...string,
) (*entity.Thread, error) {
	var thread *entity.Thread

	stmt := t.db.Statement(ctx, tx, threadTable, fields...).
		Where("id = ?", threadID)

	if err := stmt.First(&thread).Error; err != nil {
		return nil, err
	}
	return thread, nil
}
