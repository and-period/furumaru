package entity

import (
	"time"

	"gorm.io/datatypes"
)

// VideoStatus - オンデマンド配信状況
type VideoStatus int32

const (
	VideoStatusUnknown   VideoStatus = 0
	VideoStatusPrivate   VideoStatus = 1 // 非公開
	VideoStatusWaiting   VideoStatus = 2 // 公開前
	VideoStatusLimited   VideoStatus = 3 // 限定公開
	VideoStatusPublished VideoStatus = 4 // 公開済み
)

// Video - オンデマンド配信情報
type Video struct {
	ID                string         `gorm:"primaryKey;<-:create"`             // オンデマンド動画ID
	CoordinatorID     string         `gorm:""`                                 // コーディネータID
	CategoryIDs       []string       `gorm:"-"`                                // カテゴリID一覧
	CategoryIDsJSON   datatypes.JSON `gorm:"default:null;column:category_ids"` // カテゴリID一覧(JSON)
	ProductIDs        []string       `gorm:"-"`                                // 商品ID一覧
	ProductIDsJSON    datatypes.JSON `gorm:"default:null;column:product_ids"`  // 商品ID一覧(JSON)
	Title             string         `gorm:""`                                 // タイトル
	Description       string         `gorm:""`                                 // 説明
	Status            VideoStatus    `gorm:"-"`                                // 配信状況
	ThumbnailURL      string         `gorm:""`                                 // サムネイルURL
	OriginalVideoURL  string         `gorm:""`                                 // オリジナル動画URL
	ProcessedVideoURL string         `gorm:""`                                 // 加工済み動画URL
	Public            bool           `gorm:""`                                 // 公開設定
	Limited           bool           `gorm:""`                                 // 限定公開設定
	PublishedAt       time.Time      `gorm:""`                                 // 公開日時
	CreatedAt         time.Time      `gorm:"<-:create"`                        // 作成日時
	UpdatedAt         time.Time      `gorm:""`                                 // 更新日時
}
