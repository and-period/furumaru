package entity

import "time"

// オンデマンド配信関連商品情報
type VideoProduct struct {
	VideoID   string    `gorm:"primaryKey;<-:create"` // オンデマンド動画ID
	ProductID string    `gorm:"primaryKey;<-:create"` // 商品ID
	Priority  int64     `gorm:"default:0"`            // 表示優先度
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type VideoProducts []*VideoProduct
