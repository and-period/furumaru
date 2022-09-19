package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// Schedule - 開催スケジュール
type Schedule struct {
	ID          string         `gorm:"primaryKey;<-:create"` // テンプレートID
	Title       string         `gorm:""`                     // タイトル
	Description string         `gorm:""`                     // 説明
	ThumnailURL string         `gorm:""`                     // サムネイルURL
	StartAt     time.Time      `gorm:""`                     // 開催開始日時
	EndAt       time.Time      `gorm:""`                     // 開催終了日時
	Canceled    bool           `gorm:""`                     // 開催中止フラグ
	CreatedAt   time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time      `gorm:""`                     // 更新日時
	DeletedAt   gorm.DeletedAt `gorm:"default:null"`
}

type NewScheduleParams struct {
	Title       string
	Description string
	ThumnailURL string
	StartAt     time.Time
	EndAt       time.Time
}

func NewSchedule(params *NewScheduleParams) *Schedule {
	return &Schedule{
		ID:          uuid.Base58Encode(uuid.New()),
		Title:       params.Title,
		Description: params.Description,
		ThumnailURL: params.ThumnailURL,
		StartAt:     params.StartAt,
		EndAt:       params.EndAt,
	}
}
