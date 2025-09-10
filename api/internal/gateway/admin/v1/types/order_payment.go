package types

// PaymentMethodType - 決済手段
type PaymentMethodType int32

const (
	PaymentMethodTypeUnknown      PaymentMethodType = 0
	PaymentMethodTypeCash         PaymentMethodType = 1  // 代引支払い
	PaymentMethodTypeCreditCard   PaymentMethodType = 2  // クレジットカード決済
	PaymentMethodTypeKonbini      PaymentMethodType = 3  // コンビニ決済
	PaymentMethodTypeBankTransfer PaymentMethodType = 4  // 銀行振込決済
	PaymentMethodTypePayPay       PaymentMethodType = 5  // QR決済（PayPay）
	PaymentMethodTypeLinePay      PaymentMethodType = 6  // QR決済（LINE Pay）
	PaymentMethodTypeMerpay       PaymentMethodType = 7  // QR決済（メルペイ）
	PaymentMethodTypeRakutenPay   PaymentMethodType = 8  // QR決済（楽天ペイ）
	PaymentMethodTypeAUPay        PaymentMethodType = 9  // QR決済（au PAY）
	PaymentMethodTypeNone         PaymentMethodType = 10 // 決済無し
	PaymentMethodTypePaidy        PaymentMethodType = 11 // ペイディ（Paidy）
	PaymentMethodTypePayEasy      PaymentMethodType = 12 // ペイジー（PayEasy）
)

// PaymentStatus - 支払い状況
type PaymentStatus int32

const (
	PaymentStatusUnknown    PaymentStatus = 0
	PaymentStatusUnpaid     PaymentStatus = 1 // 未支払い
	PaymentStatusAuthorized PaymentStatus = 2 // オーソリ済み
	PaymentStatusPaid       PaymentStatus = 3 // 支払い済み
	PaymentStatusCanceled   PaymentStatus = 4 // キャンセル済み
	PaymentStatusFailed     PaymentStatus = 5 // 失敗
)
