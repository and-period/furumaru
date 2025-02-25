package entity

import (
	"time"
)

// ScheduleType - 通知スケジュール種別
type ScheduleType int32

const (
	ScheduleTypeUnknown                 ScheduleType = 0
	ScheduleTypeNotification            ScheduleType = 1 // お知らせ通知
	ScheduleTypeStartLive               ScheduleType = 2 // ライブ配信開始通知
	ScheduleTypeReviewProductRequest    ScheduleType = 3 // 商品レビュー依頼通知
	ScheduleTypeReviewExperienceRequest ScheduleType = 4 // 体験レビュー依頼通知
)

var ScheduleTypes = []ScheduleType{
	ScheduleTypeNotification,
	ScheduleTypeStartLive,
	ScheduleTypeReviewProductRequest,
	ScheduleTypeReviewExperienceRequest,
}

// ScheduleStatus - 通知スケジュール実行状態
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
	Deadline    time.Time      `gorm:"default:null"`         // 送信締め切り日時
	CreatedAt   time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time      `gorm:""`                     // 更新日時
}

type Schedules []*Schedule

type NewScheduleParams struct {
	MessageType ScheduleType
	MessageID   string
	SentAt      time.Time
	Deadline    time.Time
}

func NewSchedule(params *NewScheduleParams) *Schedule {
	return &Schedule{
		MessageType: params.MessageType,
		MessageID:   params.MessageID,
		Status:      ScheduleStatusWaiting,
		Count:       0,
		SentAt:      params.SentAt,
		Deadline:    params.Deadline,
	}
}

func (s *Schedule) Executable(now time.Time) bool {
	if !s.Deadline.IsZero() && s.Deadline.Before(now) {
		// 通知予定時間を過ぎてしまっている
		return false
	}
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
	if !s.Deadline.IsZero() && s.Deadline.Before(now) {
		// 通知予定時間を過ぎてしまっている
		return true
	}
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
