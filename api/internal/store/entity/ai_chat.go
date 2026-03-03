package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// AiChatSession - AIチャットセッション
type AiChatSession struct {
	ID        string    `gorm:"primaryKey;<-:create"` // セッションID
	AdminID   string    `gorm:""`                     // 管理者ID
	ProductID string    `gorm:""`                     // 関連商品ID
	Title     string    `gorm:""`                     // セッションタイトル
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type AiChatSessions []*AiChatSession

type NewAiChatSessionParams struct {
	AdminID   string
	ProductID string
	Title     string
}

func NewAiChatSession(params *NewAiChatSessionParams) *AiChatSession {
	return &AiChatSession{
		ID:        uuid.Base58Encode(uuid.New()),
		AdminID:   params.AdminID,
		ProductID: params.ProductID,
		Title:     params.Title,
	}
}

// AiChatMessage - AIチャットメッセージ
type AiChatMessage struct {
	ID        string    `gorm:"primaryKey;<-:create"` // メッセージID
	SessionID string    `gorm:""`                     // セッションID
	Role      string    `gorm:""`                     // ロール (user, assistant)
	Content   string    `gorm:""`                     // メッセージ内容 (JSON)
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
}

type AiChatMessages []*AiChatMessage

type NewAiChatMessageParams struct {
	SessionID string
	Role      string
	Content   string
}

func NewAiChatMessage(params *NewAiChatMessageParams) *AiChatMessage {
	return &AiChatMessage{
		ID:        uuid.Base58Encode(uuid.New()),
		SessionID: params.SessionID,
		Role:      params.Role,
		Content:   params.Content,
	}
}
