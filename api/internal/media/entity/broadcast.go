package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

type BroadcastOrderBy string

const (
	BroadcastOrderByUpdatedAt BroadcastOrderBy = "updated_at"
)

// BroadcastType - ライブ配信種別
type BroadcastType int32

const (
	BroadcastTypeUnknown   BroadcastType = 0
	BroadcastTypeNormal    BroadcastType = 1 // 通常配信
	BroadcastTypeRehearsal BroadcastType = 2 // リハーサル配信
)

// BroadcastStatus - ライブ配信状況
type BroadcastStatus int32

const (
	BroadcastStatusUnknown  BroadcastStatus = 0
	BroadcastStatusDisabled BroadcastStatus = 1 // リソース未作成
	BroadcastStatusWaiting  BroadcastStatus = 2 // リソース作成/削除中
	BroadcastStatusIdle     BroadcastStatus = 3 // 停止中
	BroadcastStatusActive   BroadcastStatus = 4 // 配信中
)

// Broadcast - ライブ配信情報
type Broadcast struct {
	ID                        string          `gorm:"primaryKey;<-:create"` // ライブ配信ID
	ScheduleID                string          `gorm:"default:null"`         // 開催スケジュールID
	CoordinatorID             string          `gorm:""`                     // コーディネータID
	Type                      BroadcastType   `gorm:""`                     // ライブ配信種別
	Status                    BroadcastStatus `gorm:""`                     // ライブ配信状況
	InputURL                  string          `gorm:""`                     // ライブ配信URL(入力)
	OutputURL                 string          `gorm:""`                     // ライブ配信URL(出力)
	ArchiveURL                string          `gorm:""`                     // アーカイブ配信URL
	ArchiveFixed              bool            `gorm:""`                     // アーカイブ映像を編集したか
	CloudFrontDistributionArn string          `gorm:"default:null"`         // CloudFrontディストリビューションARN
	MediaLiveChannelArn       string          `gorm:"default:null"`         // MediaLiveチャンネルARN
	MediaLiveChannelID        string          `gorm:"default:null"`         // MediaLiveチャンネルID
	MediaLiveRTMPInputArn     string          `gorm:"default:null"`         // MediaLiveインプットARN(RTMP)
	MediaLiveRTMPInputName    string          `gorm:"default:null"`         // MediaLiveインプット名(RTMP)
	MediaLiveMP4InputArn      string          `gorm:"default:null"`         // MediaLiveインプットARN(MP4)
	MediaLiveMP4InputName     string          `gorm:"default:null"`         // MediaLiveインプット名(MP4)
	MediaStoreContainerArn    string          `gorm:"default:null"`         // MediaStoreコンテナARN
	YoutubeStreamURL          string          `gorm:"default:null"`         // YouTube配信URL
	YoutubeStreamKey          string          `gorm:"default:null"`         // YouTubeストリームキー
	CreatedAt                 time.Time       `gorm:"<-:create"`            // 登録日時
	UpdatedAt                 time.Time       `gorm:""`                     // 更新日時
}

type Broadcasts []*Broadcast

type NewBroadcastParams struct {
	ScheduleID    string
	CoordinatorID string
}

func NewBroadcast(params *NewBroadcastParams) *Broadcast {
	return &Broadcast{
		ID:            uuid.Base58Encode(uuid.New()),
		ScheduleID:    params.ScheduleID,
		CoordinatorID: params.CoordinatorID,
		Type:          BroadcastTypeNormal,
		Status:        BroadcastStatusDisabled,
	}
}

func (bs Broadcasts) ScheduleIDs() []string {
	return set.UniqBy(bs, func(b *Broadcast) string {
		return b.ScheduleID
	})
}
