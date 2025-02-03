package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/cognito"
)

var errInvalidAuthUsername = errors.New("entity: invalid username")

// AdminAuthProviderType - 管理者認証プロバイダ種別
type AdminAuthProviderType int32

const (
	AdminAuthProviderTypeUnknown AdminAuthProviderType = 0
	AdminAuthProviderTypeGoogle  AdminAuthProviderType = 1 // Google認証
)

func (t AdminAuthProviderType) ToCognito() cognito.ProviderType {
	switch t {
	case AdminAuthProviderTypeGoogle:
		return cognito.ProviderTypeGoogle
	default:
		return cognito.ProviderTypeUnknown
	}
}

// AdminAuthProvider - 管理者認証プロバイダ
type AdminAuthProvider struct {
	AdminID      string                `gorm:"primaryKey;<-:create"` // 管理者ID
	ProviderType AdminAuthProviderType `gorm:"primaryKey;<-:create"` // プロバイダ種別
	AccountID    string                `gorm:"default:null"`         // アカウントID
	Email        string                `gorm:"default:null"`         // メールアドレス
	CreatedAt    time.Time             `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time             `gorm:""`                     // 更新日時
}

type AdminAuthProviders []*AdminAuthProvider

type AdminAuthProviderParams struct {
	AdminID      string
	ProviderType AdminAuthProviderType
	Auth         *cognito.AuthUser
}

func NewAdminAuthProvider(params *AdminAuthProviderParams) (*AdminAuthProvider, error) {
	// Cognitoユーザー名の形式）#{Provier名}_${アカウントID}
	strs := strings.Split(params.Auth.Username, "_")
	if len(strs) != 2 {
		return nil, errInvalidAuthUsername
	}
	return &AdminAuthProvider{
		AdminID:      params.AdminID,
		ProviderType: params.ProviderType,
		AccountID:    strs[1],
		Email:        params.Auth.Email,
	}, nil
}
