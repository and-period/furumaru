package entity

import (
	"strings"
	"time"
)

// Admin - 管理者共通情報
type Admin struct {
	ID            string    // 管理者ID
	Role          AdminRole // 管理者権限
	Lastname      string    // 姓
	Firstname     string    // 名
	LastnameKana  string    // 姓(かな)
	FirstnameKana string    // 名(かな)
	Email         string    // メールアドレス
	CreatedAt     time.Time // 登録日時
	UpdatedAt     time.Time // 更新日時
}

type Admins []*Admin

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
