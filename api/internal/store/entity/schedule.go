package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// Schedule - 開催スケジュール
type Schedule struct {
	ID            string         `gorm:"primaryKey;<-:create"` // テンプレートID
	CoordinatorID string         `gorm:""`                     // 仲介者ID
	ShippingID    string         `gorm:""`                     // 配送設定ID
	Title         string         `gorm:""`                     // タイトル
	Description   string         `gorm:""`                     // 説明
	ThumbnailURL  string         `gorm:""`                     // サムネイルURL
	StartAt       time.Time      `gorm:""`                     // 開催開始日時
	EndAt         time.Time      `gorm:""`                     // 開催終了日時
	Canceled      bool           `gorm:""`                     // 開催中止フラグ
	CreatedAt     time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time      `gorm:""`                     // 更新日時
	DeletedAt     gorm.DeletedAt `gorm:"default:null"`
}

type Schedules []*Schedule

type NewScheduleParams struct {
	CoordinatorID string
	ShippingID    string
	Title         string
	Description   string
	ThumbnailURL  string
	StartAt       time.Time
	EndAt         time.Time
}

func NewSchedule(params *NewScheduleParams) *Schedule {
	return &Schedule{
		ID:            uuid.Base58Encode(uuid.New()),
		CoordinatorID: params.CoordinatorID,
		ShippingID:    params.ShippingID,
		Title:         params.Title,
		Description:   params.Description,
		ThumbnailURL:  params.ThumbnailURL,
		StartAt:       params.StartAt,
		EndAt:         params.EndAt,
	}
}
