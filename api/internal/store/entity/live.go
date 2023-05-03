package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
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
	LiveProducts   `gorm:"-"`
	ID             string         `gorm:"primaryKey;<-:create"` // ライブ配信ID
	ScheduleID     string         `gorm:""`                     // 開催スケジュールID
	ProducerID     string         `gorm:""`                     // 生産者ID
	Title          string         `gorm:""`                     // タイトル
	Description    string         `gorm:""`                     // 説明
	Status         LiveStatus     `gorm:""`                     // 配信ステータス
	StartAt        time.Time      `gorm:""`                     // 配信開始日時
	EndAt          time.Time      `gorm:""`                     // 配信終了日時
	ChannelArn     string         `gorm:"default:null"`         // チャンネルArn
	StreamKeyArn   string         `gorm:"default:null"`         // ストリームキーArn
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
	ChannelName    string         `gorm:"-"`                    // チャンネル名
	IngestEndpoint string         `gorm:"-"`                    // 配信エンドポイント
	StreamKey      string         `gorm:"-"`                    // ストリームキー
	StreamID       string         `gorm:"-"`                    // ストリームID
	PlaybackURL    string         `gorm:"-"`                    // 再生用URL
	ViewerCount    int64          `gorm:"-"`                    // 視聴者数
}

type Lives []*Live

type NewLiveParams struct {
	ScheduleID  string
	ProducerID  string
	Title       string
	Description string
	Status      LiveStatus
	StartAt     time.Time
	EndAt       time.Time
}

type FillLiveIvsParams struct {
	ChannelName    string
	IngestEndpoint string
	StreamKey      string
	PlaybackURL    string
	StreamID       string
	ViewerCount    int64
}

func NewLive(params *NewLiveParams) *Live {
	return &Live{
		ID:          uuid.Base58Encode(uuid.New()),
		ScheduleID:  params.ScheduleID,
		ProducerID:  params.ProducerID,
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		StartAt:     params.StartAt,
		EndAt:       params.EndAt,
	}
}

func (l *Live) Fill(products LiveProducts) {
	l.LiveProducts = products
}

func (l *Live) FillIVS(params FillLiveIvsParams) {
	l.ChannelName = params.ChannelName
	l.IngestEndpoint = params.IngestEndpoint
	l.StreamKey = params.StreamKey
	l.PlaybackURL = params.PlaybackURL
	l.StreamID = params.StreamID
	l.ViewerCount = params.ViewerCount
}

func (ls Lives) IDs() []string {
	return set.UniqBy(ls, func(l *Live) string {
		return l.ID
	})
}

func (ls Lives) ProducerIDs() []string {
	return set.UniqBy(ls, func(l *Live) string {
		return l.ProducerID
	})
}

func (ls Lives) Fill(products map[string]LiveProducts) {
	for i := range ls {
		ls[i].Fill(products[ls[i].ID])
	}
}
