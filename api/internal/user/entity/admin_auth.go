package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

// AdminAuth - 管理者認証情報
type AdminAuth struct {
	AdminID      string    // 管理者ID
	Type         AdminType // 権限
	GroupIDs     []string  // グループID一覧
	AccessToken  string    // アクセストークン
	RefreshToken string    // 更新トークン
	ExpiresIn    int32     // 有効期限
}

func NewAdminAuth(admin *Admin, rs *cognito.AuthResult) *AdminAuth {
	return &AdminAuth{
		AdminID:      admin.ID,
		Type:         admin.Type,
		GroupIDs:     admin.GroupIDs,
		AccessToken:  rs.AccessToken,
		RefreshToken: rs.RefreshToken,
		ExpiresIn:    rs.ExpiresIn,
	}
}

// AdminAuthEvent - 管理者OAuthイベント
type AdminAuthEvent struct {
	AdminID      string    `dynamodbav:"admin_id"`            // 管理者ID
	ProviderType string    `dynamodbav:"provider_type"`       // プロバイダー種別
	Nonce        string    `dynamodbav:"nonce"`               // セキュア文字列（リプレイアタック対策）
	ExpiredAt    time.Time `dynamodbav:"expired_at,unixtime"` // 有効期限
	CreatedAt    time.Time `dynamodbav:"created_at"`          // 登録日時
	UpdatedAt    time.Time `dynamodbav:"updated_at"`          // 更新日時
}

type AdminAuthEventParams struct {
	AdminID      string
	ProviderType string
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
		"admin_id":      e.AdminID,
		"provider_type": e.ProviderType,
	}
}
