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

// PaymentMethodType - 決済手段
type PaymentMethodType int32

const (
	PaymentMethodTypeUnknown     PaymentMethodType = 0
	PaymentMethodTypeCash        PaymentMethodType = 1 // 代引支払い
	PaymentMethodTypeCreditCard  PaymentMethodType = 2 // クレジットカード決済
	PaymentMethodTypeKonbini     PaymentMethodType = 3 // コンビニ決済
	PaymentMethodTypeBankTranser PaymentMethodType = 4 // 銀行振込決済
	PaymentMethodTypePayPay      PaymentMethodType = 5 // QR決済（PayPay）
	PaymentMethodTypeLinePay     PaymentMethodType = 6 // QR決済（LINE Pay）
	PaymentMethodTypeMerpay      PaymentMethodType = 7 // QR決済（メルペイ）
	PaymentMethodTypeRakutenPay  PaymentMethodType = 8 // QR決済（楽天ペイ）
	PaymentMethodTypeAUPay       PaymentMethodType = 9 // QR決済（au PAY）
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
