package entity

import (
	"errors"
	"strings"
	"time"
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

// Admin - 管理者共通情報
type Admin struct {
	ID            string    `gorm:"primaryKey;<-:create"` // 管理者ID
	CognitoID     string    `gorm:"<-:create"`            // 管理者ID (Cognito用)
	Role          AdminRole `gorm:"<-:create"`            // 管理者権限
	Lastname      string    `gorm:""`                     // 姓
	Firstname     string    `gorm:""`                     // 名
	LastnameKana  string    `gorm:""`                     // 姓(かな)
	FirstnameKana string    `gorm:""`                     // 名(かな)
	Email         string    `gorm:""`                     // メールアドレス
	Device        string    `gorm:""`                     // デバイストークン(Push通知用)
	CreatedAt     time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time `gorm:""`                     // 更新日時
}

type Admins []*Admin

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

func NewAdminFromAdministrator(administrator *Administrator) *Admin {
	return &Admin{
		ID:            administrator.ID,
		Role:          AdminRoleAdministrator,
		Lastname:      administrator.Lastname,
		Firstname:     administrator.Firstname,
		LastnameKana:  administrator.LastnameKana,
		FirstnameKana: administrator.FirstnameKana,
		Email:         administrator.Email,
		CreatedAt:     administrator.CreatedAt,
		UpdatedAt:     administrator.UpdatedAt,
	}
}

func NewAdminFromCoordinator(coordinator *Coordinator) *Admin {
	return &Admin{
		ID:            coordinator.ID,
		Role:          AdminRoleCoordinator,
		Lastname:      coordinator.Lastname,
		Firstname:     coordinator.Firstname,
		LastnameKana:  coordinator.LastnameKana,
		FirstnameKana: coordinator.FirstnameKana,
		Email:         coordinator.Email,
		CreatedAt:     coordinator.CreatedAt,
		UpdatedAt:     coordinator.UpdatedAt,
	}
}

func NewAdminFromProducer(producer *Producer) *Admin {
	return &Admin{
		ID:            producer.ID,
		Role:          AdminRoleProducer,
		Lastname:      producer.Lastname,
		Firstname:     producer.Firstname,
		LastnameKana:  producer.LastnameKana,
		FirstnameKana: producer.FirstnameKana,
		Email:         producer.Email,
		CreatedAt:     producer.CreatedAt,
		UpdatedAt:     producer.UpdatedAt,
	}
}

func (a *Admin) Name() string {
	return strings.TrimSpace(strings.Join([]string{a.Lastname, a.Firstname}, " "))
}

func NewAdminsFromAdministrators(administrators Administrators) Admins {
	res := make(Admins, len(administrators))
	for i := range administrators {
		res[i] = NewAdminFromAdministrator(administrators[i])
	}
	return res
}

func NewAdminsFromCoordinators(coordinators Coordinators) Admins {
	res := make(Admins, len(coordinators))
	for i := range coordinators {
		res[i] = NewAdminFromCoordinator(coordinators[i])
	}
	return res
}

func NewAdminsFromProducers(producers Producers) Admins {
	res := make(Admins, len(producers))
	for i := range producers {
		res[i] = NewAdminFromProducer(producers[i])
	}
	return res
}
