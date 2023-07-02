package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// ScheduleStatus - 開催状況
type ScheduleStatus int32

const (
	ScheduleStatusUnknown    ScheduleStatus = 0
	ScheduleStatusPrivate    ScheduleStatus = 1 // 非公開
	ScheduleStatusInProgress ScheduleStatus = 2 // 管理者承認前
	ScheduleStatusWaiting    ScheduleStatus = 3 // 開催前
	ScheduleStatusLive       ScheduleStatus = 4 // 開催中
	ScheduleStatusClosed     ScheduleStatus = 5 // 開催終了
)

// Schedule - 開催スケジュール
type Schedule struct {
	ID                   string         `gorm:"primaryKey;<-:create"` // テンプレートID
	CoordinatorID        string         `gorm:""`                     // コーディネータID
	ShippingID           string         `gorm:""`                     // 配送設定ID
	Status               ScheduleStatus `gorm:"-"`                    // 開催状況
	Title                string         `gorm:""`                     // タイトル
	Description          string         `gorm:""`                     // 説明
	ThumbnailURL         string         `gorm:""`                     // サムネイルURL
	OpeningVideoURL      string         `gorm:""`                     // オープニング動画URL
	IntermissionVideoURL string         `gorm:""`                     // 幕間動画URL
	Public               bool           `gorm:""`                     // 公開フラグ
	Approved             bool           `gorm:""`                     // 承認フラグ
	ApprovedAdminID      string         `gorm:""`                     // 承認した管理者ID
	StartAt              time.Time      `gorm:""`                     // 開催開始日時
	EndAt                time.Time      `gorm:""`                     // 開催終了日時
	CreatedAt            time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt            time.Time      `gorm:""`                     // 更新日時
	DeletedAt            gorm.DeletedAt `gorm:"default:null"`
}

type Schedules []*Schedule

type NewScheduleParams struct {
	CoordinatorID        string
	ShippingID           string
	Title                string
	Description          string
	ThumbnailURL         string
	OpeningVideoURL      string
	IntermissionVideoURL string
	StartAt              time.Time
	EndAt                time.Time
}

func NewSchedule(params *NewScheduleParams) *Schedule {
	return &Schedule{
		ID:                   uuid.Base58Encode(uuid.New()),
		CoordinatorID:        params.CoordinatorID,
		ShippingID:           params.ShippingID,
		Title:                params.Title,
		Description:          params.Description,
		ThumbnailURL:         params.ThumbnailURL,
		OpeningVideoURL:      params.OpeningVideoURL,
		IntermissionVideoURL: params.IntermissionVideoURL,
		Approved:             false,
		ApprovedAdminID:      "",
		StartAt:              params.StartAt,
		EndAt:                params.EndAt,
	}
}

func (s *Schedule) Fill(now time.Time) {
	switch {
	case !s.Approved:
		s.Status = ScheduleStatusInProgress
	case !s.Public:
		s.Status = ScheduleStatusPrivate
	case now.Before(s.StartAt):
		s.Status = ScheduleStatusWaiting
	case now.Before(s.EndAt):
		s.Status = ScheduleStatusLive
	default:
		s.Status = ScheduleStatusClosed
	}
}

func (ss Schedules) Fill(now time.Time) {
	for i := range ss {
		ss[i].Fill(now)
	}
}

func (ss Schedules) CoordinatorIDs() []string {
	return set.UniqBy(ss, func(s *Schedule) string {
		return s.CoordinatorID
	})
}

func (ss Schedules) ShippingIDs() []string {
	return set.UniqBy(ss, func(s *Schedule) string {
		return s.ShippingID
	})
}
