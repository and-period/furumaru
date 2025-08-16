package entity

import "time"

// OrderMetadata - 注文付加情報
type OrderMetadata struct {
	OrderID        string    `gorm:"primaryKey;<-:create"` // 注文履歴ID
	PickupAt       time.Time `gorm:"default:null"`         // 受け取り日時
	PickupLocation string    `gorm:"default:null"`         // 受け取り場所
	CreatedAt      time.Time `gorm:"<-:create"`            // 作成日時
	UpdatedAt      time.Time `gorm:""`                     // 更新日時
}

type MultiOrderMetadata []*OrderMetadata

type NewOrderPickupMetadataParams struct {
	OrderID        string
	PickupAt       time.Time
	PickupLocation string
}

func NewOrderPickupMetadata(params *NewOrderPickupMetadataParams) *OrderMetadata {
	return &OrderMetadata{
		OrderID:        params.OrderID,
		PickupAt:       params.PickupAt,
		PickupLocation: params.PickupLocation,
	}
}

func (ms MultiOrderMetadata) MapByOrderID() map[string]*OrderMetadata {
	res := make(map[string]*OrderMetadata, len(ms))
	for _, m := range ms {
		res[m.OrderID] = m
	}
	return res
}
