package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
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
	ScheduleID                string          `gorm:""`                     // 開催スケジュールID
	Status                    BroadcastStatus `gorm:""`                     // ライブ配信状況
	InputURL                  string          `gorm:""`                     // ライブ配信URL(入力)
	OutputURL                 string          `gorm:""`                     // ライブ配信URL(出力)
	ArchiveURL                string          `gorm:""`                     // アーカイブ配信URL
	CloudFrontDistributionArn string          `gorm:"default:null"`         // CloudFrontディストリビューションARN
	MediaLiveChannelArn       string          `gorm:"default:null"`         // MediaLiveチャンネルARN
	MediaLiveChannelID        string          `gorm:"default:null"`         // MediaLiveチャンネルID
	MediaLiveRTMPInputArn     string          `gorm:"default:null"`         // MediaLiveインプットARN(RTMP)
	MediaLiveRTMPInputName    string          `gorm:"default:null"`         // MediaLiveインプット名(RTMP)
	MediaLiveMP4InputArn      string          `gorm:"default:null"`         // MediaLiveインプットARN(MP4)
	MediaLiveMP4InputName     string          `gorm:"default:null"`         // MediaLiveインプット名(MP4)
	MediaStoreContainerArn    string          `gorm:"default:null"`         // MediaStoreコンテナARN
	CreatedAt                 time.Time       `gorm:"<-:create"`            // 登録日時
	UpdatedAt                 time.Time       `gorm:""`                     // 更新日時
}

type Broadcasts []*Broadcast

type NewBroadcastParams struct {
	ScheduleID string
}

func NewBroadcast(params *NewBroadcastParams) *Broadcast {
	return &Broadcast{
		ID:         uuid.Base58Encode(uuid.New()),
		ScheduleID: params.ScheduleID,
		Status:     BroadcastStatusDisabled,
	}
}
