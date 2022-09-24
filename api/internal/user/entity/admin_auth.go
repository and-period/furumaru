package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/cognito"
	set "github.com/and-period/furumaru/api/pkg/set/v2"
)

// AdminAuth - 管理者認証情報
type AdminAuth struct {
	AdminID      string    `gorm:"primaryKey;<-:create"` // 管理者ID
	CognitoID    string    `gorm:"<-:create"`            // Deprecated: 管理者ID (Cognito用)
	Role         AdminRole `gorm:"<-:create"`            // 権限
	Device       string    `gorm:""`                     // Deprecated: デバイストークン(Push通知用)
	AccessToken  string    `gorm:"-"`                    // アクセストークン
	RefreshToken string    `gorm:"-"`                    // 更新トークン
	ExpiresIn    int32     `gorm:"-"`                    // 有効期限
	CreatedAt    time.Time `gorm:"<-:create"`            // Deprecated: 登録日時
	UpdatedAt    time.Time `gorm:""`                     // Deprecated: 更新日時
}

type AdminAuths []*AdminAuth

func NewAdminAuth(adminID, cognitoID string, role AdminRole) *AdminAuth {
	return &AdminAuth{
		AdminID:   adminID,
		CognitoID: cognitoID,
		Role:      role,
	}
}

func (a *AdminAuth) Fill(rs *cognito.AuthResult) {
	a.AccessToken = rs.AccessToken
	a.RefreshToken = rs.RefreshToken
	a.ExpiresIn = rs.ExpiresIn
}

func (as AdminAuths) GroupByRole() map[AdminRole]AdminAuths {
	const maxRoles = 4
	res := make(map[AdminRole]AdminAuths, maxRoles)
	for _, a := range as {
		if _, ok := res[a.Role]; !ok {
			res[a.Role] = make(AdminAuths, 0, len(as))
		}
		res[a.Role] = append(res[a.Role], a)
	}
	return res
}

func (as AdminAuths) AdminIDs() []string {
	res := make([]string, len(as))
	for i := range as {
		res[i] = as[i].AdminID
	}
	return res
}

func (as AdminAuths) Devices() []string {
	set := set.New[string](len(as))
	for i := range as {
		if as[i].Device == "" {
			continue
		}
		set.Add(as[i].Device)
	}
	return set.Slice()
}
