package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type ShippingSize types.ShippingSize

type ShippingType types.ShippingType

func NewShippingSize(size entity.ShippingSize) ShippingSize {
	switch size {
	case entity.ShippingSize60:
		return ShippingSize(types.ShippingSize60)
	case entity.ShippingSize80:
		return ShippingSize(types.ShippingSize80)
	case entity.ShippingSize100:
		return ShippingSize(types.ShippingSize100)
	default:
		return ShippingSize(types.ShippingSizeUnknown)
	}
}

func (s ShippingSize) Response() types.ShippingSize {
	return types.ShippingSize(s)
}

func NewShippingType(typ entity.ShippingType) ShippingType {
	switch typ {
	case entity.ShippingTypeNormal:
		return ShippingType(types.ShippingTypeNormal)
	case entity.ShippingTypeFrozen:
		return ShippingType(types.ShippingTypeFrozen)
	case entity.ShippingTypePickup:
		return ShippingType(types.ShippingTypePickup)
	default:
		return ShippingType(types.ShippingTypeUnknown)
	}
}

func (t ShippingType) Response() types.ShippingType {
	return types.ShippingType(t)
}
