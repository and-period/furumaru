package entity

import (
	"time"
)

type ScheduleType int32

const (
	ScheduleTypeUnknown      ScheduleType = 0
	ScheduleTypeNotification ScheduleType = 1 // お知らせ通知
)

var ScheduleTypes = []ScheduleType{
	ScheduleTypeNotification,
}

type ScheduleStatus int32

const (
	ScheduleStatusWaiting    ScheduleStatus = 0 // 実行前
	ScheduleStatusProcessing ScheduleStatus = 1 // 実行中
	ScheduleStatusDone       ScheduleStatus = 2 // 完了
	ScheduleStatusCanceled   ScheduleStatus = 3 // 中止
)

// Schedule - 通知スケジュール管理
type Schedule struct {
	MessageType ScheduleType   `gorm:"primaryKey;<-:create"` // 通知種別
	MessageID   string         `gorm:"primaryKey;<-:create"` // 通知ID
	Status      ScheduleStatus `gorm:""`                     // 実行ステータス
	Count       int64          `gorm:""`                     // 実行回数
	SentAt      time.Time      `gorm:""`                     // 送信日時
	CreatedAt   time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time      `gorm:""`                     // 更新日時
}

type Schedules []*Schedule

func NewSchedule(messageType ScheduleType, messageID string, sentAt time.Time) *Schedule {
	return &Schedule{
		MessageType: messageType,
		MessageID:   messageID,
		Status:      ScheduleStatusWaiting,
		Count:       0,
		SentAt:      sentAt,
	}
}

func (s *Schedule) Executable(now time.Time) bool {
	switch s.Status {
	case ScheduleStatusWaiting:
		if now.Before(s.SentAt) {
			// 通知予定時刻より前の場合、まだ実行しない
			return false
		}
		return true
	case ScheduleStatusDone, ScheduleStatusCanceled:
		return false
	case ScheduleStatusProcessing:
		if now.Before(s.UpdatedAt.Add(10 * time.Minute)) {
			// 前回実行開始から10分経っていない場合は最実行しない
			return false
		}
		return s.Count < 2 // 1度のみ再実行させる
	default:
		return false
	}
}

func (s *Schedule) ShouldCancel(now time.Time) bool {
	switch s.Status {
	case ScheduleStatusWaiting, ScheduleStatusDone, ScheduleStatusCanceled:
		return false
	case ScheduleStatusProcessing:
		if now.Before(s.UpdatedAt.Add(10 * time.Minute)) {
			// 前回実行開始から10分経っていない場合は処理を続ける
			return false
		}
		return s.Count > 1 // 再実行後も終わっていない場合は処理を中止する
	default:
		return false
	}
}

func (ss Schedules) Map() map[string]*Schedule {
	res := make(map[string]*Schedule, len(ss))
	for _, s := range ss {
		res[s.MessageID] = s
	}
	return res
}
