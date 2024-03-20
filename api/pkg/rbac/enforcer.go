//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package rbac

import (
	"fmt"

	casbin "github.com/casbin/casbin/v2"
)

type Enforcer interface {
	Enforce(rvals ...interface{}) (bool, error)                       // 引数で渡された情報と一致するかの検証
	GetRolesForUser(group string, domain ...string) ([]string, error) // 引数で渡されたグループ名に対する検証ルールを返す
}

// NewEnforcer - 認可情報の検証用クライアントの生成
func NewEnforcer(modelPath, policyPath string) (Enforcer, error) {
	if modelPath == "" || policyPath == "" {
		return &casbin.Enforcer{}, nil
	}
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, fmt.Errorf("casbin: failed to new enforcer: %w", err)
	}
	if err := enforcer.InitWithFile(modelPath, policyPath); err != nil {
		return nil, fmt.Errorf("casbin: failed to init enforcer: %w", err)
	}
	enforcer.EnableLog(false)
	return enforcer, nil
}
