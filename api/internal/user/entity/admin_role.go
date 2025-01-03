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

// RelatedAdminGroup - 管理者グループと管理者の紐付け情報
type RelatedAdminGroup struct {
	AdminID   string    `gorm:"primaryKey;<-:create"` // 管理者ID
	GroupID   string    `gorm:"primaryKey;<-:create"` // 管理者グループID
	ExpiredAt time.Time `gorm:"default:null"`         // 有効期限
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type RelatedAdminGroups []RelatedAdminGroup

func (gs RelatedAdminGroups) GroupIDs() []string {
	res := make([]string, len(gs))
	for i := range gs {
		res[i] = gs[i].GroupID
	}
	return res
}

func (gs RelatedAdminGroups) GroupByAdminID() map[string]RelatedAdminGroups {
	res := make(map[string]RelatedAdminGroups, len(gs))
	for _, g := range gs {
		if _, ok := res[g.AdminID]; !ok {
			res[g.AdminID] = make(RelatedAdminGroups, 0)
		}
		res[g.AdminID] = append(res[g.AdminID], g)
	}
	return res
}

// AdminGroup - 管理者グループ情報
type AdminGroup struct {
	ID             string         `gorm:"primaryKey;<-:create"` // 管理者グループID
	Type           AdminType      `gorm:"<-:create"`            // 管理者グループ種別
	Name           string         `gorm:""`                     // 管理者グループ名
	Description    string         `gorm:""`                     // 説明
	CreatedAdminID string         `gorm:""`                     // 登録者ID
	UpdatedAdminID string         `gorm:""`                     // 更新者ID
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

// AdminGroupRole - 管理者グループと管理者権限の紐付け情報
type AdminGroupRole struct {
	AdminGroupID   string         `gorm:"primaryKey;<-:create"` // 管理者グループID
	RoleID         string         `gorm:"primaryKey;<-:create"` // 管理者権限ID
	CreatedAdminID string         `gorm:""`                     // 登録者ID
	UpdatedAdminID string         `gorm:""`                     // 更新者ID
	ExpiredAt      time.Time      `gorm:"default:null"`         // 有効期限
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type AdminGroupRoles []AdminGroupRole

func (rs AdminGroupRoles) Write(w *csv.Writer) error {
	for _, r := range rs {
		record := []string{
			"g",            // レコードタイプ
			r.AdminGroupID, // 管理者グループID
			r.RoleID,       // 管理者権限ID
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// AdminRole - 管理者権限情報
type AdminRole struct {
	ID             string         `gorm:"primaryKey;<-:create"` // 管理者権限ID
	Name           string         `gorm:""`                     // 管理者権限名
	Note           string         `gorm:""`                     // 備考
	CreatedAdminID string         `gorm:""`                     // 登録者ID
	UpdatedAdminID string         `gorm:""`                     // 更新者ID
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type AdminRoles []AdminRole

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

type AdminPolicies []AdminPolicy

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
	RoleID         string    `gorm:"primaryKey;<-:create"` // 管理者権限ID
	PolicyID       string    `gorm:"primaryKey;<-:create"` // 管理者ポリシーID
	CreatedAdminID string    `gorm:""`                     // 登録者ID
	UpdatedAdminID string    `gorm:""`                     // 更新者ID
	CreatedAt      time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time `gorm:""`                     // 更新日時
}

type AdminRolePolicies []AdminRolePolicy

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
