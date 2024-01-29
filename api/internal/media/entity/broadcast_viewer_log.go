package entity

import "time"

// BroadcastViewerLog - ライブ視聴履歴情報
type BroadcastViewerLog struct {
	BroadcastID string    `gorm:"primaryKey;<-:create"` // ライブ配信ID
	SessionID   string    `gorm:"primaryKey;<-:create"` // セッションID
	CreatedAt   time.Time `gorm:"primaryKey;<-:create"` // 登録日時
	UserID      string    `gorm:"default:null"`         // ユーザーID
	UserAgent   string    `gorm:""`                     // ユーザーエージェント
	ClientIP    string    `gorm:""`                     // 接続元IPアドレス
	UpdatedAt   time.Time `gorm:""`                     // 更新日時
}

type BroadcastViewerLogParams struct {
	BroadcastID string
	SessionID   string
	UserID      string
	UserAgent   string
	ClientIP    string
}

func NewBroadcastViewerLog(params *BroadcastViewerLogParams) *BroadcastViewerLog {
	return &BroadcastViewerLog{
		BroadcastID: params.BroadcastID,
		SessionID:   params.SessionID,
		UserID:      params.UserID,
		UserAgent:   params.UserAgent,
		ClientIP:    params.ClientIP,
	}
}
