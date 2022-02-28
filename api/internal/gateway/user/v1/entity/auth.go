package entity

import (
	"github.com/and-period/marche/api/internal/gateway/entity"
	"github.com/and-period/marche/api/internal/gateway/util"
)

// Auth - 認証情報
type Auth struct {
	UserID       string `json:"userId"`       // ユーザーID
	AccessToken  string `json:"accessToken"`  // アクセストークン
	RefreshToken string `json:"refreshToken"` // 更新トークン
	ExpiresIn    int32  `json:"expiresIn"`    // 有効期限
	TokenType    string `json:"tokenType"`    // トークン種別
}

func NewAuth(auth *entity.UserAuth) *Auth {
	return &Auth{
		UserID:       auth.UserId,
		AccessToken:  auth.AccessToken,
		RefreshToken: auth.RefreshToken,
		ExpiresIn:    auth.ExpiresIn,
		TokenType:    util.AuthTokenType,
	}
}
