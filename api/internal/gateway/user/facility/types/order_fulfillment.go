package types

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
