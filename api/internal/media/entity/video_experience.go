package entity

import "time"

// オンデマンド配信関連体験情報
type VideoExperience struct {
	VideoID      string    `gorm:"primaryKey;<-:create"` // オンデマンド動画ID
	ExperienceID string    `gorm:"primaryKey;<-:create"` // 体験ID
	Priority     int64     `gorm:"default:0"`            // 表示優先度
	CreatedAt    time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time `gorm:""`                     // 更新日時
}

type VideoExperiences []*VideoExperience
