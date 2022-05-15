package entity

import "github.com/and-period/marche/api/pkg/cognito"

// AdminAuth - 管理者認証情報
type AdminAuth struct {
	AdminID      string    // 管理者ID
	Role         AdminRole // 権限
	AccessToken  string    // アクセストークン
	RefreshToken string    // 更新トークン
	ExpiresIn    int32     // 有効期限
}

func NewAdminAuth(adminID string, role AdminRole, rs *cognito.AuthResult) *AdminAuth {
	return &AdminAuth{
		AdminID:      adminID,
		Role:         role,
		AccessToken:  rs.AccessToken,
		RefreshToken: rs.RefreshToken,
		ExpiresIn:    rs.ExpiresIn,
	}
}
