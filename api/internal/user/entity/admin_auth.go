package entity

import "github.com/and-period/marche/api/pkg/cognito"

// AdminAuth - 管理者認証情報
type AdminAuth struct {
	AdminID      string
	Role         AdminRole
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
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
