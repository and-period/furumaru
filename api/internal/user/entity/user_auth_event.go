package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// UserAuthEvent - 購入者OAuthイベント
type UserAuthEvent struct {
	SessionID    string               `dynamodbav:"session_id"`          // セッションID
	ProviderType UserAuthProviderType `dynamodbav:"provider_type"`       // プロバイダ種別
	Nonce        string               `dynamodbav:"nonce"`               // セキュア文字列（リプレイアタック対策）
	ExpiredAt    time.Time            `dynamodbav:"expired_at,unixtime"` // 有効期限
	CreatedAt    time.Time            `dynamodbav:"created_at"`          // 登録日時
	UpdatedAt    time.Time            `dynamodbav:"updated_at"`          // 更新日時
}

type UserAuthEventParams struct {
	SessionID    string
	ProviderType UserAuthProviderType
	Now          time.Time
	TTL          time.Duration
}

func NewUserAuthEvent(params *UserAuthEventParams) *UserAuthEvent {
	return &UserAuthEvent{
		SessionID:    params.SessionID,
		ProviderType: params.ProviderType,
		Nonce:        uuid.Base58Encode(uuid.New()),
		ExpiredAt:    params.Now.Add(params.TTL),
		CreatedAt:    params.Now,
		UpdatedAt:    params.Now,
	}
}

func (e *UserAuthEvent) TableName() string {
	return "user-auth-events"
}

func (e *UserAuthEvent) PrimaryKey() map[string]interface{} {
	return map[string]interface{}{
		"session_id": e.SessionID,
	}
}
