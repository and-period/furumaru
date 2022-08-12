package entity

import (
	"errors"
	"time"

	"github.com/and-period/furumaru/api/pkg/cognito"
)

var errInvalidAdminRole = errors.New("entity: invalid admin role")

// AdminRole - 管理者権限
type AdminRole int32

const (
	AdminRoleUnknown       AdminRole = 0
	AdminRoleAdministrator AdminRole = 1 // 管理者
	AdminRoleCoordinator   AdminRole = 2 // 仲介者
	AdminRoleProducer      AdminRole = 3 // 生産者
)

// AdminAuth - 管理者認証情報
type AdminAuth struct {
	AdminID      string    `gorm:"primaryKey;<-:create"` // 管理者ID
	CognitoID    string    `gorm:"<-:create"`            // 管理者ID (Cognito用)
	Role         AdminRole `gorm:"<-:create"`            // 権限
	Device       string    `gorm:""`                     // デバイストークン(Push通知用)
	AccessToken  string    `gorm:"-"`                    // アクセストークン
	RefreshToken string    `gorm:"-"`                    // 更新トークン
	ExpiresIn    int32     `gorm:"-"`                    // 有効期限
	CreatedAt    time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time `gorm:""`                     // 更新日時
}

type AdminAuths []*AdminAuth

func NewAdminRole(role int32) (AdminRole, error) {
	res := AdminRole(role)
	if err := res.Validate(); err != nil {
		return AdminRoleUnknown, err
	}
	return res, nil
}

func (r AdminRole) Validate() error {
	switch r {
	case AdminRoleAdministrator, AdminRoleCoordinator, AdminRoleProducer:
		return nil
	default:
		return errInvalidAdminRole
	}
}

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
