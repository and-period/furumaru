package tidb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/datatypes"
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
	internal, err := newInternalReceivedQueues(queues)
	if err != nil {
		return dbError(err)
	}
	err = q.db.DB.WithContext(ctx).Table(receivedQueueTable).Create(&internal).Error
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

	return internal.entity()
}

type internalReceivedQueue struct {
	entity.ReceivedQueue `gorm:"embedded"`
	UserIDsJSON          datatypes.JSON `gorm:"default:null;column:user_ids"` // 送信先ユーザーID一覧(JSON)
}

type internalReceivedQueues []*internalReceivedQueue

func newInternalReceivedQueue(queue *entity.ReceivedQueue) (*internalReceivedQueue, error) {
	userIDs, err := json.Marshal(queue.UserIDs)
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal received queue user ids: %w", err)
	}
	internal := &internalReceivedQueue{
		ReceivedQueue: *queue,
		UserIDsJSON:   userIDs,
	}
	return internal, nil
}

func (q *internalReceivedQueue) entity() (*entity.ReceivedQueue, error) {
	if err := q.unmarshalUserIDs(); err != nil {
		return nil, err
	}
	return &q.ReceivedQueue, nil
}

func (q *internalReceivedQueue) unmarshalUserIDs() error {
	if q.UserIDsJSON == nil {
		return nil
	}
	var userIDs []string
	if err := json.Unmarshal(q.UserIDsJSON, &userIDs); err != nil {
		return fmt.Errorf("tidb: failed to unmarshal received queue user ids: %w", err)
	}
	q.ReceivedQueue.UserIDs = userIDs
	return nil
}

func newInternalReceivedQueues(queues entity.ReceivedQueues) (internalReceivedQueues, error) {
	res := make(internalReceivedQueues, len(queues))
	for i := range queues {
		internal, err := newInternalReceivedQueue(queues[i])
		if err != nil {
			return nil, err
		}
		res[i] = internal
	}
	return res, nil
}

func (qs internalReceivedQueues) entities() (entity.ReceivedQueues, error) {
	res := make(entity.ReceivedQueues, len(qs))
	for i := range qs {
		q, err := qs[i].entity()
		if err != nil {
			return nil, err
		}
		res[i] = q
	}
	return res, nil
}
