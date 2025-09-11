package types

// FulfillmentStatus - 配送状況
type FulfillmentStatus int32

const (
	FulfillmentStatusUnknown     FulfillmentStatus = 0
	FulfillmentStatusUnfulfilled FulfillmentStatus = 1 // 未発送
	FulfillmentStatusFulfilled   FulfillmentStatus = 2 // 発送済み
)

// ShippingCarrier - 配送会社
type ShippingCarrier int32

const (
	ShippingCarrierUnknown ShippingCarrier = 0
	ShippingCarrierYamato  ShippingCarrier = 1 // ヤマト運輸
	ShippingCarrierSagawa  ShippingCarrier = 2 // 佐川急便
)

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

// OrderFulfillment - 配送情報
type OrderFulfillment struct {
	FulfillmentID   string            `json:"fulfillmentId"`   // 配送情報ID
	TrackingNumber  string            `json:"trackingNumber"`  // 伝票番号
	Status          FulfillmentStatus `json:"status"`          // 配送状況
	ShippingCarrier ShippingCarrier   `json:"shippingCarrier"` // 配送会社
	ShippingType    ShippingType      `json:"shippingType"`    // 配送方法
	BoxNumber       int64             `json:"boxNumber"`       // 箱の通番
	BoxSize         ShippingSize      `json:"boxSize"`         // 箱の大きさ
	BoxRate         int64             `json:"boxRate"`         // 箱の占有率
	ShippedAt       int64             `json:"shippedAt"`       // 配送日時
	*Address                          // 配送先情報
}
