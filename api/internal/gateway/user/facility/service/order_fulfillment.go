package service

import "github.com/and-period/furumaru/api/internal/store/entity"

// ShippingSize - 配送時の箱の大きさ
type ShippingSize int32

const (
	ShippingSizeUnknown ShippingSize = 0
	ShippingSize60      ShippingSize = 1 // 箱のサイズ:60
	ShippingSize80      ShippingSize = 2 // 箱のサイズ:80
	ShippingSize100     ShippingSize = 3 // 箱のサイズ:100
)

// ShippingType - 配送方法
type ShippingType int32

const (
	ShippingTypeUnknown ShippingType = 0
	ShippingTypeNormal  ShippingType = 1 // 常温・冷蔵便
	ShippingTypeFrozen  ShippingType = 2 // 冷凍便
	ShippingTypePickup  ShippingType = 3 // 店舗受取
)

func NewShippingSize(size entity.ShippingSize) ShippingSize {
	switch size {
	case entity.ShippingSize60:
		return ShippingSize60
	case entity.ShippingSize80:
		return ShippingSize80
	case entity.ShippingSize100:
		return ShippingSize100
	default:
		return ShippingSizeUnknown
	}
}

func (s ShippingSize) Response() int32 {
	return int32(s)
}

func NewShippingType(typ entity.ShippingType) ShippingType {
	switch typ {
	case entity.ShippingTypeNormal:
		return ShippingTypeNormal
	case entity.ShippingTypeFrozen:
		return ShippingTypeFrozen
	case entity.ShippingTypePickup:
		return ShippingTypePickup
	default:
		return ShippingTypeUnknown
	}
}

func (t ShippingType) Response() int32 {
	return int32(t)
}
