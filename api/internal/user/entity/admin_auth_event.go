package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// AdminAuthEvent - 管理者OAuthイベント
type AdminAuthEvent struct {
	AdminID      string                `dynamodbav:"admin_id"`            // 管理者ID
	ProviderType AdminAuthProviderType `dynamodbav:"provider_type"`       // プロバイダ種別
	Nonce        string                `dynamodbav:"nonce"`               // セキュア文字列（リプレイアタック対策）
	ExpiredAt    time.Time             `dynamodbav:"expired_at,unixtime"` // 有効期限
	CreatedAt    time.Time             `dynamodbav:"created_at"`          // 登録日時
	UpdatedAt    time.Time             `dynamodbav:"updated_at"`          // 更新日時
}

type AdminAuthEventParams struct {
	AdminID      string
	ProviderType AdminAuthProviderType
	Now          time.Time
	TTL          time.Duration
}

func NewAdminAuthEvent(params *AdminAuthEventParams) *AdminAuthEvent {
	return &AdminAuthEvent{
		AdminID:      params.AdminID,
		ProviderType: params.ProviderType,
		Nonce:        uuid.Base58Encode(uuid.New()),
		ExpiredAt:    params.Now.Add(params.TTL),
		CreatedAt:    params.Now,
		UpdatedAt:    params.Now,
	}
}

func (e *AdminAuthEvent) TableName() string {
	return "admin-auth-events"
}

func (e *AdminAuthEvent) PrimaryKey() map[string]interface{} {
	return map[string]interface{}{
		"admin_id": e.AdminID,
	}
}
