package entity

import "time"

// NotifyType - 通知種別
type NotifyType int32

const (
	NotifyTypeUnknown NotifyType = 0
	NotifyTypeEmail   NotifyType = 1 // メール通知
	NotifyTypeMessage NotifyType = 2 // メッセージ通知
	NotifyTypePush    NotifyType = 3 // プッシュ通知
	NotifyTypeReport  NotifyType = 4 // システムレポート
)

// ReceivedQueue - 通知キュー管理
type ReceivedQueue struct {
	ID         string     `gorm:"primaryKey;<-:create"` // 通知キューID
	NotifyType NotifyType `gorm:"primaryKey;<-:create"` // 通知種別
	EventType  EventType  `gorm:""`                     // 実行種別
	UserType   UserType   `gorm:""`                     // 送信先ユーザー種別
	UserIDs    []string   `gorm:"-"`                    // 送信先ユーザーID一覧
	Done       bool       `gorm:""`                     // 完了フラグ
	CreatedAt  time.Time  `gorm:"<-:create"`            // 登録日時
	UpdatedAt  time.Time  `gorm:""`                     // 更新日時
}

type ReceivedQueues []*ReceivedQueue

func NewReceivedQueue(payload *WorkerPayload, notifyType NotifyType) *ReceivedQueue {
	return &ReceivedQueue{
		ID:         payload.QueueID,
		NotifyType: notifyType,
		EventType:  payload.EventType,
		UserType:   payload.UserType,
		UserIDs:    payload.UserIDs,
		Done:       false,
	}
}

func NewReceivedQueues(payload *WorkerPayload) ReceivedQueues {
	const max = 4
	res := make(ReceivedQueues, 0, max)
	if payload.Email != nil {
		res = append(res, NewReceivedQueue(payload, NotifyTypeEmail))
	}
	if payload.Message != nil {
		res = append(res, NewReceivedQueue(payload, NotifyTypeMessage))
	}
	if payload.Push != nil {
		res = append(res, NewReceivedQueue(payload, NotifyTypePush))
	}
	if payload.Report != nil {
		res = append(res, NewReceivedQueue(payload, NotifyTypeReport))
	}
	return res
}
