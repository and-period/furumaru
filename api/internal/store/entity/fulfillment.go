package entity

import (
	"time"

	"gorm.io/gorm"
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

// Fulfillment - 注文配送情報
type Fulfillment struct {
	OrderID         string          `gorm:"primaryKey;<-:create"` // 注文履歴ID
	AddressID       string          `gorm:""`                     // 配送先情報ID
	TrackingNumber  string          `gorm:"default:null"`         // 配送伝票番号
	ShippingCarrier ShippingCarrier `gorm:""`                     // 配送会社
	ShippingMethod  DeliveryType    `gorm:""`                     // 配送方法
	BoxSize         ShippingSize    `gorm:""`                     // 箱の大きさ
	CreatedAt       time.Time       `gorm:"<-:create"`            // 登録日時
	UpdatedAt       time.Time       `gorm:""`                     // 更新日時
	DeletedAt       gorm.DeletedAt  `gorm:"default:null"`         // 削除日時
}

type Fulfillments []*Fulfillment

func (fs Fulfillments) MapByOrderID() map[string]*Fulfillment {
	res := make(map[string]*Fulfillment, len(fs))
	for _, f := range fs {
		res[f.OrderID] = f
	}
	return res
}
