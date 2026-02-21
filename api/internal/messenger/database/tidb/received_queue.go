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

const receivedQueueTable = "received_queues"

type receivedQueue struct {
	db  *mysql.Client
	now func() time.Time
}

func NewReceivedQueue(db *mysql.Client) database.ReceivedQueue {
	return &receivedQueue{
		db:  db,
		now: jst.Now,
	}
}

func (q *receivedQueue) Get(
	ctx context.Context, queueID string, typ entity.NotifyType, fields ...string,
) (*entity.ReceivedQueue, error) {
	queue, err := q.get(ctx, q.db.DB, queueID, typ, fields...)
	return queue, dbError(err)
}

func (q *receivedQueue) MultiCreate(ctx context.Context, queues ...*entity.ReceivedQueue) error {
	now := q.now()
	for _, queue := range queues {
		queue.CreatedAt, queue.UpdatedAt = now, now
	}
	internal := newInternalReceivedQueues(queues)
	err := q.db.DB.WithContext(ctx).Table(receivedQueueTable).Create(&internal).Error
	return dbError(err)
}

func (q *receivedQueue) UpdateDone(ctx context.Context, queueID string, typ entity.NotifyType, done bool) error {
	updates := map[string]interface{}{
		"done":       done,
		"updated_at": q.now(),
	}
	stmt := q.db.DB.WithContext(ctx).
		Table(receivedQueueTable).
		Where("id = ?", queueID).
		Where("notify_type = ?", typ)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (q *receivedQueue) get(
	ctx context.Context, tx *gorm.DB, queueID string, typ entity.NotifyType, fields ...string,
) (*entity.ReceivedQueue, error) {
	var internal *internalReceivedQueue

	err := q.db.Statement(ctx, tx, receivedQueueTable, fields...).
		Where("id = ?", queueID).
		Where("notify_type = ?", typ).
		First(&internal).Error
	if err != nil {
		return nil, err
	}

	return internal.entity(), nil
}

type internalReceivedQueue struct {
	entity.ReceivedQueue `gorm:"embedded"`
	UserIDsJSON          mysql.JSONColumn[[]string] `gorm:"default:null;column:user_ids"` // 送信先ユーザーID一覧(JSON)
}

type internalReceivedQueues []*internalReceivedQueue

func newInternalReceivedQueue(queue *entity.ReceivedQueue) *internalReceivedQueue {
	return &internalReceivedQueue{
		ReceivedQueue: *queue,
		UserIDsJSON:   mysql.NewJSONColumn(queue.UserIDs),
	}
}

func (q *internalReceivedQueue) entity() *entity.ReceivedQueue {
	q.ReceivedQueue.UserIDs = q.UserIDsJSON.Val
	return &q.ReceivedQueue
}

func newInternalReceivedQueues(queues entity.ReceivedQueues) internalReceivedQueues {
	res := make(internalReceivedQueues, len(queues))
	for i := range queues {
		res[i] = newInternalReceivedQueue(queues[i])
	}
	return res
}
