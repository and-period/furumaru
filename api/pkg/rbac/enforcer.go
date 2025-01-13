//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package rbac

import (
	"errors"
	"fmt"
	"os"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
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

func NewEnforcerFromString(modelText, policyText string) (Enforcer, error) {
	if modelText == "" || policyText == "" {
		return nil, errors.New("casbin: invalid argument")
	}
	model, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, fmt.Errorf("casbin: failed to new model: %w", err)
	}
	file, err := os.CreateTemp("", "casbin-policy-*.csv")
	if err != nil {
		return nil, fmt.Errorf("casbin: failed to create temp file: %w", err)
	}
	defer os.Remove(file.Name())
	if _, err := file.WriteString(policyText); err != nil {
		return nil, fmt.Errorf("casbin: failed to write policy: %w", err)
	}
	enforcer, err := casbin.NewEnforcer(model, fileadapter.NewAdapter(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("casbin: failed to new enforcer: %w", err)
	}
	enforcer.EnableLog(false)
	return enforcer, nil
}
