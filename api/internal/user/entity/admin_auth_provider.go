package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/cognito"
)

var (
	ErrInvalidAdminAuthUsername     = errors.New("entity: invalid admin auth username")
	ErrInvalidAdminAuthProviderType = errors.New("entity: invalid admin auth provider type")
)

// AdminAuthProviderType - 管理者認証プロバイダ種別
type AdminAuthProviderType int32

const (
	AdminAuthProviderTypeUnknown AdminAuthProviderType = 0
	AdminAuthProviderTypeGoogle  AdminAuthProviderType = 1 // Google認証
	AdminAuthProviderTypeLINE    AdminAuthProviderType = 2 // LINE認証
)

func (t AdminAuthProviderType) ToCognito() cognito.ProviderType {
	switch t {
	case AdminAuthProviderTypeGoogle:
		return cognito.ProviderTypeGoogle
	case AdminAuthProviderTypeLINE:
		return cognito.ProviderTypeLINE
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
	if strs := strings.Split(params.Auth.Username, "_"); len(strs) != 2 {
		return nil, ErrInvalidAdminAuthUsername // CognitoユーザーはAuthProvider作成不要
	}
	identity := findAdminAuthIdentity(params.Auth, params.ProviderType)
	if identity == nil {
		return nil, ErrInvalidAdminAuthProviderType // 対象のプロバイダについての連携情報が存在しない
	}
	return &AdminAuthProvider{
		AdminID:      params.AdminID,
		ProviderType: params.ProviderType,
		AccountID:    identity.UserID,
		Email:        params.Auth.Email,
	}, nil
}

func findAdminAuthIdentity(user *cognito.AuthUser, providerType AdminAuthProviderType) *cognito.AuthUserIdentity {
	target := providerType.ToCognito()
	for _, identity := range user.Identities {
		if identity.ProviderType == target {
			return identity
		}
	}
	return nil
}
