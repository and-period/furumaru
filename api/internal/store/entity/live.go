package entity

import (
	"time"

	set "github.com/and-period/furumaru/api/pkg/set/v2"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

type LiveStatus int32

const (
	LiveStatusUnknown  LiveStatus = 0
	LiveStatusWaiting  LiveStatus = 1 // 開始前
	LiveStatusOpened   LiveStatus = 2 // 配信中
	LiveStatusClosed   LiveStatus = 3 // 配信終了
	LiveStatusCanceled LiveStatus = 4 // 配信中止
)

// ライブ配信情報
type Live struct {
	LiveProducts `gorm:"-"`
	ID           string         `gorm:"primaryKey;<-:create"` // ライブ配信ID
	ScheduleID   string         `gorm:""`                     // 開催スケジュールID
	ProducerID   string         `gorm:""`                     // 生産者ID
	Title        string         `gorm:""`                     // タイトル
	Description  string         `gorm:""`                     // 説明
	Status       LiveStatus     `gorm:"-"`                    // 配信ステータス
	Published    bool           `gorm:""`                     // 配信公開フラグ
	Canceled     bool           `gorm:""`                     // 配信中止フラグ
	StartAt      time.Time      `gorm:""`                     // 配信開始日時
	EndAt        time.Time      `gorm:""`                     // 配信終了日時
	ChannelArn   string         `gorm:""`                     // チャンネルArn
	CreatedAt    time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time      `gorm:""`                     // 更新日時
	DeletedAt    gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Lives []*Live

type NewLiveParams struct {
	ScheduleID  string
	ProducerID  string
	Title       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
}

func NewLive(params *NewLiveParams) *Live {
	return &Live{
		ID:          uuid.Base58Encode(uuid.New()),
		ScheduleID:  params.ScheduleID,
		ProducerID:  params.ProducerID,
		Title:       params.Title,
		Description: params.Description,
		StartAt:     params.StartAt,
		EndAt:       params.EndAt,
	}
}

func (l *Live) Fill(products LiveProducts, now time.Time) {
	l.LiveProducts = products
	l.Status = l.status(now)
}

func (l *Live) status(now time.Time) LiveStatus {
	if l.Canceled {
		return LiveStatusCanceled
	}
	if !l.Published {
		return LiveStatusWaiting
	}
	switch {
	case now.Before(l.StartAt):
		return LiveStatusWaiting
	case now.Before(l.EndAt):
		return LiveStatusOpened
	default:
		return LiveStatusClosed
	}
}

func (ls Lives) ProducerIDs() []string {
	return set.UniqBy(ls, func(l *Live) string {
		return l.ProducerID
	})
}

func (ls Lives) Fill(products map[string]LiveProducts, now time.Time) {
	for i := range ls {
		ls[i].Fill(products[ls[i].ID], now)
	}
}
