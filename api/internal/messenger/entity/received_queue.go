package entity

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// ReceivedQueue - 通知キュー管理
type ReceivedQueue struct {
	ID          string         `gorm:"primaryKey;<-:create"`         // 通知キューID
	EventType   EventType      `gorm:""`                             // 通知種別
	UserType    UserType       `gorm:""`                             // 送信先ユーザー種別
	UserIDs     []string       `gorm:"-"`                            // 送信先ユーザーID一覧
	UserIDsJSON datatypes.JSON `gorm:"default:null;column:user_ids"` // 送信先ユーザーID一覧(JSON)
	Done        bool           `gorm:""`                             // 完了フラグ
	CreatedAt   time.Time      `gorm:"<-:create"`                    // 登録日時
	UpdatedAt   time.Time      `gorm:""`                             // 更新日時
}

func NewReceivedQueue(payload *WorkerPayload) *ReceivedQueue {
	return &ReceivedQueue{
		ID:        payload.QueueID,
		EventType: payload.EventType,
		UserType:  payload.UserType,
		UserIDs:   payload.UserIDs,
		Done:      false,
	}
}

func (q *ReceivedQueue) Fill() error {
	var userIDs []string
	if err := json.Unmarshal(q.UserIDsJSON, &userIDs); err != nil {
		return err
	}
	q.UserIDs = userIDs
	return nil
}

func (q *ReceivedQueue) FillJSON() error {
	v, err := json.Marshal(q.UserIDs)
	if err != nil {
		return err
	}
	q.UserIDsJSON = datatypes.JSON(v)
	return nil
}
