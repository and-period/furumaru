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
	ID              string         `gorm:"primaryKey;<-:create"` // テンプレートID
	ShopID          string         `gorm:"default:null"`         // 店舗ID
	CoordinatorID   string         `gorm:""`                     // コーディネータID
	Status          ScheduleStatus `gorm:"-"`                    // 開催状況
	Title           string         `gorm:""`                     // タイトル
	Description     string         `gorm:""`                     // 説明
	ThumbnailURL    string         `gorm:""`                     // サムネイルURL
	ImageURL        string         `gorm:""`                     // ふた絵URL
	OpeningVideoURL string         `gorm:""`                     // オープニング動画URL
	Public          bool           `gorm:""`                     // 公開フラグ
	Approved        bool           `gorm:""`                     // 承認フラグ
	ApprovedAdminID string         `gorm:""`                     // 承認した管理者ID
	StartAt         time.Time      `gorm:""`                     // 開催開始日時
	EndAt           time.Time      `gorm:""`                     // 開催終了日時
	CreatedAt       time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt       time.Time      `gorm:""`                     // 更新日時
	DeletedAt       gorm.DeletedAt `gorm:"default:null"`
}

type Schedules []*Schedule

type NewScheduleParams struct {
	CoordinatorID   string
	Title           string
	Description     string
	ThumbnailURL    string
	ImageURL        string
	OpeningVideoURL string
	Public          bool
	StartAt         time.Time
	EndAt           time.Time
}

func NewSchedule(params *NewScheduleParams) *Schedule {
	return &Schedule{
		ID:              uuid.Base58Encode(uuid.New()),
		CoordinatorID:   params.CoordinatorID,
		Title:           params.Title,
		Description:     params.Description,
		ThumbnailURL:    params.ThumbnailURL,
		ImageURL:        params.ImageURL,
		OpeningVideoURL: params.OpeningVideoURL,
		Public:          params.Public,
		Approved:        true, // デフォルトは承認済みにしておく
		ApprovedAdminID: "",
		StartAt:         params.StartAt,
		EndAt:           params.EndAt,
	}
}

func (s *Schedule) Fill(now time.Time) error {
	s.SetStatus(now)
	return nil
}

func (s *Schedule) SetStatus(now time.Time) {
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

func (s *Schedule) Published() bool {
	if s == nil {
		return false
	}
	return s.Public && s.Approved
}

func (ss Schedules) Fill(now time.Time) error {
	for i := range ss {
		if err := ss[i].Fill(now); err != nil {
			return err
		}
	}
	return nil
}

func (ss Schedules) IDs() []string {
	res := make([]string, len(ss))
	for i := range ss {
		res[i] = ss[i].ID
	}
	return res
}

func (ss Schedules) CoordinatorIDs() []string {
	return set.UniqBy(ss, func(s *Schedule) string {
		return s.CoordinatorID
	})
}
