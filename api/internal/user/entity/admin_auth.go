package entity

import (
	"github.com/and-period/furumaru/api/pkg/cognito"
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
