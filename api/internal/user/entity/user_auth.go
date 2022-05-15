package entity

import "github.com/and-period/marche/api/pkg/cognito"

// UserAuth - 購入者認証情報
type UserAuth struct {
	UserID       string // ユーザーID
	AccessToken  string // アクセストークン
	RefreshToken string // 更新トークン
	ExpiresIn    int32  // 有効期限
}

func NewUserAuth(userID string, rs *cognito.AuthResult) *UserAuth {
	return &UserAuth{
		UserID:       userID,
		AccessToken:  rs.AccessToken,
		RefreshToken: rs.RefreshToken,
		ExpiresIn:    rs.ExpiresIn,
	}
}
