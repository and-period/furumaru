package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
)

// 配送ステータス
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

// OrderFulfillment - 注文配送情報
type OrderFulfillment struct {
	ID                string            `gorm:"primaryKey;<-:create"` // 注文配送ID
	OrderID           string            `gorm:""`                     // 注文履歴ID
	AddressRevisionID int64             `gorm:""`                     // 配送先情報ID
	Status            FulfillmentStatus `gorm:""`                     // 配送ステータス
	TrackingNumber    string            `gorm:"default:null"`         // 配送伝票番号
	ShippingCarrier   ShippingCarrier   `gorm:""`                     // 配送会社
	ShippingMethod    DeliveryType      `gorm:""`                     // 配送方法
	BoxNumber         int64             `gorm:""`                     // 箱の通番
	BoxSize           ShippingSize      `gorm:""`                     // 箱の大きさ
	ShippedAt         time.Time         `gorm:"default:null"`         // 配送日時
	CreatedAt         time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time         `gorm:""`                     // 更新日時
}

type OrderFulfillments []*OrderFulfillment

func (fs OrderFulfillments) AddressRevisionIDs() []int64 {
	return set.UniqBy(fs, func(f *OrderFulfillment) int64 {
		return f.AddressRevisionID
	})
}

func (fs OrderFulfillments) GroupByOrderID() map[string]OrderFulfillments {
	res := make(map[string]OrderFulfillments, len(fs))
	for _, f := range fs {
		if _, ok := res[f.OrderID]; !ok {
			res[f.OrderID] = make(OrderFulfillments, 0, len(fs))
		}
		res[f.OrderID] = append(res[f.OrderID], f)
	}
	return res
}
