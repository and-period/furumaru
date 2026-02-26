package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

// PaymentSystemStatus - 決済システム状態
type PaymentSystemStatus types.PaymentSystemStatus

// PaymentProviderType - 決済プロバイダー種別
type PaymentProviderType types.PaymentProviderType

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

func NewPaymentProviderType(providerType entity.PaymentProviderType) PaymentProviderType {
	switch providerType {
	case entity.PaymentProviderTypeKomoju:
		return PaymentProviderType(types.PaymentProviderTypeKomoju)
	case entity.PaymentProviderTypeStripe:
		return PaymentProviderType(types.PaymentProviderTypeStripe)
	default:
		return PaymentProviderType(types.PaymentProviderTypeUnknown)
	}
}

func (t PaymentProviderType) StoreEntity() entity.PaymentProviderType {
	switch types.PaymentProviderType(t) {
	case types.PaymentProviderTypeKomoju:
		return entity.PaymentProviderTypeKomoju
	case types.PaymentProviderTypeStripe:
		return entity.PaymentProviderTypeStripe
	default:
		return entity.PaymentProviderTypeUnknown
	}
}

func (t PaymentProviderType) Response() types.PaymentProviderType {
	return types.PaymentProviderType(t)
}

func NewPaymentSystem(system *entity.PaymentSystem) *PaymentSystem {
	return &PaymentSystem{
		PaymentSystem: types.PaymentSystem{
			MethodType:   NewPaymentMethodType(system.MethodType).Response(),
			ProviderType: NewPaymentProviderType(system.ProviderType).Response(),
			Status:       NewPaymentSystemStatus(system.Status).Response(),
			CreatedAt:    jst.Unix(system.CreatedAt),
			UpdatedAt:    jst.Unix(system.UpdatedAt),
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
