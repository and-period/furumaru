package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

// PaymentSystemStatus - 決済システム状態
type PaymentSystemStatus int32

const (
	PaymentSystemStatusUnknown PaymentSystemStatus = 0
	PaymentSystemStatusInUse   PaymentSystemStatus = 1 // 利用可能
	PaymentSystemStatusOutage  PaymentSystemStatus = 2 // 停止中
)

type PaymentSystem struct {
	types.PaymentSystem
}

type PaymentSystems []*PaymentSystem

func NewPaymentSystemStatus(status entity.PaymentSystemStatus) PaymentSystemStatus {
	switch status {
	case entity.PaymentSystemStatusInUse:
		return PaymentSystemStatusInUse
	case entity.PaymentSystemStatusOutage:
		return PaymentSystemStatusOutage
	default:
		return PaymentSystemStatusUnknown
	}
}

func (s PaymentSystemStatus) StoreEntity() entity.PaymentSystemStatus {
	switch s {
	case PaymentSystemStatusInUse:
		return entity.PaymentSystemStatusInUse
	case PaymentSystemStatusOutage:
		return entity.PaymentSystemStatusOutage
	default:
		return entity.PaymentSystemStatusUnknown
	}
}

func (s PaymentSystemStatus) Response() types.PaymentStatus {
	return types.PaymentStatus(s)
}

func NewPaymentSystem(system *entity.PaymentSystem) *PaymentSystem {
	return &PaymentSystem{
		PaymentSystem: types.PaymentSystem{
			MethodType: NewPaymentMethodType(system.MethodType).Response(),
			Status:     NewPaymentSystemStatus(system.Status).Response(),
		},
	}
}

func (s *PaymentSystem) InService() bool {
	return PaymentSystemStatus(s.Status) == PaymentSystemStatusInUse
}

func (s *PaymentSystem) Response() *types.PaymentSystem {
	return &s.PaymentSystem
}

func NewPaymentSystems(systems entity.PaymentSystems) PaymentSystems {
	res := make(PaymentSystems, len(systems))
	for i := range systems {
		res[i] = NewPaymentSystem(systems[i])
	}
	return res
}

func (ss PaymentSystems) Response() []*types.PaymentSystem {
	res := make([]*types.PaymentSystem, len(ss))
	for i := range ss {
		res[i] = ss[i].Response()
	}
	return res
}
