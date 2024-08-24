package entity

import (
	"time"
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
	VideoProducts    `gorm:"-"`
	VideoExperiences `gorm:"-"`
	ID               string      `gorm:"primaryKey;<-:create"` // オンデマンド動画ID
	CoordinatorID    string      `gorm:""`                     // コーディネータID
	ProductIDs       []string    `gorm:"-"`                    // 商品ID一覧
	ExperienceIDs    []string    `gorm:"-"`                    // 体験ID一覧
	Title            string      `gorm:""`                     // タイトル
	Description      string      `gorm:""`                     // 説明
	Status           VideoStatus `gorm:"-"`                    // 配信状況
	ThumbnailURL     string      `gorm:""`                     // サムネイルURL
	VideoURL         string      `gorm:""`                     // 動画URL
	Public           bool        `gorm:""`                     // 公開設定
	Limited          bool        `gorm:""`                     // 限定公開設定
	PublishedAt      time.Time   `gorm:""`                     // 公開日時
	CreatedAt        time.Time   `gorm:"<-:create"`            // 作成日時
	UpdatedAt        time.Time   `gorm:""`                     // 更新日時
}

type Videos []*Video
