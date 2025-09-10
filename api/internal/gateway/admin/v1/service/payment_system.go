package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// PaymentSystemStatus - 決済システム状態
type PaymentSystemStatus types.PaymentSystemStatus

type PaymentSystem struct {
	types.PaymentSystem
}

type PaymentSystems []*PaymentSystem

func NewPaymentSystemStatus(status entity.PaymentSystemStatus) PaymentSystemStatus {
	switch status {
	case entity.PaymentSystemStatusInUse:
		return PaymentSystemStatus(types.PaymentSystemStatusInUse)
	case entity.PaymentSystemStatusOutage:
		return PaymentSystemStatus(types.PaymentSystemStatusOutage)
	default:
		return PaymentSystemStatus(types.PaymentSystemStatusUnknown)
	}
}

func (s PaymentSystemStatus) StoreEntity() entity.PaymentSystemStatus {
	switch types.PaymentSystemStatus(s) {
	case types.PaymentSystemStatusInUse:
		return entity.PaymentSystemStatusInUse
	case types.PaymentSystemStatusOutage:
		return entity.PaymentSystemStatusOutage
	default:
		return entity.PaymentSystemStatusUnknown
	}
}

func (s PaymentSystemStatus) Response() types.PaymentSystemStatus {
	return types.PaymentSystemStatus(s)
}

func NewPaymentSystem(system *entity.PaymentSystem) *PaymentSystem {
	return &PaymentSystem{
		PaymentSystem: types.PaymentSystem{
			MethodType: NewPaymentMethodType(system.MethodType).Response(),
			Status:     NewPaymentSystemStatus(system.Status).Response(),
			CreatedAt:  jst.Unix(system.CreatedAt),
			UpdatedAt:  jst.Unix(system.UpdatedAt),
		},
	}
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
