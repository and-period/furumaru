package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const receivedQueueTable = "received_queues"

type receivedQueue struct {
	db  *mysql.Client
	now func() time.Time
}

func newReceivedQueue(db *mysql.Client) database.ReceivedQueue {
	return &receivedQueue{
		db:  db,
		now: jst.Now,
	}
}

func (q *receivedQueue) Get(ctx context.Context, queueID string, fields ...string) (*entity.ReceivedQueue, error) {
	queue, err := q.get(ctx, q.db.DB, queueID, fields...)
	return queue, dbError(err)
}

func (q *receivedQueue) Create(ctx context.Context, queue *entity.ReceivedQueue) error {
	if err := queue.FillJSON(); err != nil {
		return err
	}

	now := q.now()
	queue.CreatedAt, queue.UpdatedAt = now, now

	err := q.db.DB.WithContext(ctx).Table(receivedQueueTable).Create(&queue).Error
	return dbError(err)
}

func (q *receivedQueue) UpdateDone(ctx context.Context, queueID string, done bool) error {
	updates := map[string]interface{}{
		"done":       done,
		"updated_at": q.now(),
	}
	stmt := q.db.DB.WithContext(ctx).
		Table(receivedQueueTable).
		Where("id = ?", queueID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (q *receivedQueue) get(
	ctx context.Context, tx *gorm.DB, queueID string, fields ...string,
) (*entity.ReceivedQueue, error) {
	var queue *entity.ReceivedQueue

	err := q.db.Statement(ctx, tx, receivedQueueTable, fields...).
		Where("id = ?", queueID).
		First(&queue).Error
	if err != nil {
		return nil, err
	}
	if err := queue.Fill(); err != nil {
		return nil, err
	}
	return queue, nil
}
