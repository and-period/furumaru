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

const receivedQueueTable = "received_queues"

type receivedQueue struct {
	db  *database.Client
	now func() time.Time
}

func NewReceivedQueue(db *database.Client) ReceivedQueue {
	return &receivedQueue{
		db:  db,
		now: jst.Now,
	}
}

func (q *receivedQueue) Get(ctx context.Context, queueID string, fields ...string) (*entity.ReceivedQueue, error) {
	queue, err := q.get(ctx, q.db.DB, queueID, fields...)
	return queue, exception.InternalError(err)
}

func (q *receivedQueue) Create(ctx context.Context, queue *entity.ReceivedQueue) error {
	err := q.db.Transaction(ctx, func(tx *gorm.DB) error {
		if err := queue.FillJSON(); err != nil {
			return err
		}

		now := q.now()
		queue.CreatedAt, queue.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(receivedQueueTable).Create(&queue).Error
		return err
	})
	return exception.InternalError(err)
}

func (q *receivedQueue) UpdateDone(ctx context.Context, queueID string, done bool) error {
	err := q.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := q.get(ctx, tx, queueID); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"done":       done,
			"updated_at": q.now(),
		}
		err := tx.WithContext(ctx).
			Table(receivedQueueTable).
			Where("id = ?", queueID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
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
