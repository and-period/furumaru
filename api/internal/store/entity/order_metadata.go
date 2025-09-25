package entity

import (
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
)

// OrderMetadata - 注文付加情報
type OrderMetadata struct {
	OrderID         string    `gorm:"primaryKey;<-:create"` // 注文履歴ID
	OrderRequest    string    `gorm:"default:null"`         // 注文時リクエスト
	PickupAt        time.Time `gorm:"default:null"`         // 受け取り日時
	PickupLocation  string    `gorm:"default:null"`         // 受け取り場所
	ShippingMessage string    `gorm:"default:null"`         // 発送時メッセージ
	CreatedAt       time.Time `gorm:"<-:create"`            // 作成日時
	UpdatedAt       time.Time `gorm:""`                     // 更新日時
}

type MultiOrderMetadata []*OrderMetadata

type NewOrderMetadataParams struct {
	OrderID         string
	Pickup          bool
	ShippingAddress *entity.Address
	ShippingMessage string
	PickupAt        time.Time
	PickupLocation  string
	OrderRequest    string
}

func NewOrderMetadata(params *NewOrderMetadataParams) *OrderMetadata {
	res := &OrderMetadata{
		OrderID:      params.OrderID,
		OrderRequest: params.OrderRequest,
	}
	switch {
	case params.Pickup:
		res.PickupAt = params.PickupAt
		res.PickupLocation = params.PickupLocation
	case params.ShippingAddress != nil:
		res.ShippingMessage = params.ShippingMessage
	}
	return res
}

func (ms MultiOrderMetadata) MapByOrderID() map[string]*OrderMetadata {
	res := make(map[string]*OrderMetadata, len(ms))
	for _, m := range ms {
		res[m.OrderID] = m
	}
	return res
}
