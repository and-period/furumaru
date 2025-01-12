package entity

import (
	"encoding/csv"
	"sort"
	"time"

	"gorm.io/gorm"
)

// AdminType - 管理者種別
type AdminType int32

const (
	AdminTypeUnknown       AdminType = 0
	AdminTypeAdministrator AdminType = 1 // 管理者
	AdminTypeCoordinator   AdminType = 2 // コーディネータ
	AdminTypeProducer      AdminType = 3 // 生産者
)

// AdminPolicyAction - 管理者ポリシー操作
type AdminPolicyAction string

const (
	AdminPolicyActionAllow AdminPolicyAction = "allow"
	AdminPolicyActionDeny  AdminPolicyAction = "deny"
)

// AdminGroupUser - 管理者グループと管理者の紐付け情報
type AdminGroupUser struct {
	GroupID        string    `gorm:"primaryKey;<-:create"` // 管理者グループID
	AdminID        string    `gorm:"primaryKey;<-:create"` // 管理者ID
	CreatedAdminID string    `gorm:"default:null"`         // 登録者ID
	UpdatedAdminID string    `gorm:"default:null"`         // 更新者ID
	ExpiredAt      time.Time `gorm:"default:null"`         // 有効期限
	CreatedAt      time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time `gorm:""`                     // 更新日時
}

type AdminGroupUsers []*AdminGroupUser

func (us AdminGroupUsers) GroupIDs() []string {
	res := make([]string, len(us))
	for i := range us {
		res[i] = us[i].GroupID
	}
	return res
}

func (us AdminGroupUsers) GroupByAdminID() map[string]AdminGroupUsers {
	res := make(map[string]AdminGroupUsers, len(us))
	for _, u := range us {
		if _, ok := res[u.AdminID]; !ok {
			res[u.AdminID] = make(AdminGroupUsers, 0)
		}
		res[u.AdminID] = append(res[u.AdminID], u)
	}
	return res
}

// AdminGroup - 管理者グループ情報
type AdminGroup struct {
	ID             string         `gorm:"primaryKey;<-:create"` // 管理者グループID
	Type           AdminType      `gorm:"<-:create"`            // 管理者グループ種別
	Name           string         `gorm:""`                     // 管理者グループ名
	Description    string         `gorm:""`                     // 説明
	CreatedAdminID string         `gorm:"default:null"`         // 登録者ID
	UpdatedAdminID string         `gorm:"default:null"`         // 更新者ID
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type AdminGroups []*AdminGroup

// AdminGroupRole - 管理者グループと管理者権限の紐付け情報
type AdminGroupRole struct {
	GroupID        string         `gorm:"primaryKey;<-:create"` // 管理者グループID
	RoleID         string         `gorm:"primaryKey;<-:create"` // 管理者権限ID
	CreatedAdminID string         `gorm:"default:null"`         // 登録者ID
	UpdatedAdminID string         `gorm:"default:null"`         // 更新者ID
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type AdminGroupRoles []*AdminGroupRole

func (rs AdminGroupRoles) Write(w *csv.Writer) error {
	for _, r := range rs {
		record := []string{
			"g",       // レコードタイプ
			r.GroupID, // 管理者グループID
			r.RoleID,  // 管理者権限ID
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// AdminRole - 管理者権限情報
type AdminRole struct {
	ID          string    `gorm:"primaryKey;<-:create"` // 管理者権限ID
	Name        string    `gorm:""`                     // 管理者権限名
	Description string    `gorm:""`                     // 備考
	CreatedAt   time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time `gorm:""`                     // 更新日時
}

type AdminRoles []*AdminRole

// AdminPolicy - 管理者ポリシー情報
type AdminPolicy struct {
	ID          string            `gorm:"primaryKey;<-:create"` // 管理者ポリシーID
	Name        string            `gorm:""`                     // 管理者ポリシー名
	Description string            `gorm:""`                     // 説明
	Priority    int64             `gorm:""`                     // 優先度
	Path        string            `gorm:""`                     // マッチパターン - Path
	Method      string            `gorm:""`                     // マッチパターン - Method
	Action      AdminPolicyAction `gorm:""`                     // マッチパターン - Action
	CreatedAt   time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt   time.Time         `gorm:""`                     // 更新日時
}

type AdminPolicies []*AdminPolicy

func (ps AdminPolicies) SortByPriority() AdminPolicies {
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].Priority <= ps[j].Priority
	})
	return ps
}

func (ps AdminPolicies) Write(w *csv.Writer) error {
	for _, p := range ps.SortByPriority() {
		record := []string{
			"p",              // レコードタイプ
			p.ID,             // 管理者ポリシーID
			p.Path,           // マッチパターン - Path
			p.Method,         // マッチパターン - Method
			string(p.Action), // マッチパターン - Action
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// AdminRolePolicy - 管理者ロールと管理者ポリシーの紐付け情報
type AdminRolePolicy struct {
	RoleID    string    `gorm:"primaryKey;<-:create"` // 管理者権限ID
	PolicyID  string    `gorm:"primaryKey;<-:create"` // 管理者ポリシーID
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type AdminRolePolicies []*AdminRolePolicy

func (ps AdminRolePolicies) Write(w *csv.Writer) error {
	for _, p := range ps {
		record := []string{
			"g",        // レコードタイプ
			p.RoleID,   // 管理者権限ID
			p.PolicyID, // 管理者ポリシーID
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return nil
}
