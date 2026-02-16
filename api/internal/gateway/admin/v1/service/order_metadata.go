package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type OrderMetadata struct {
	types.OrderMetadata
	orderID string
}

func NewOrderMetadata(metadata *entity.OrderMetadata) *OrderMetadata {
	return &OrderMetadata{
		OrderMetadata: types.OrderMetadata{
			OrderRequest:   metadata.OrderRequest,
			PickupAt:       jst.Unix(metadata.PickupAt),
			PickupLocation: metadata.PickupLocation,
		},
		orderID: metadata.OrderID,
	}
}

func (m *OrderMetadata) Response() *types.OrderMetadata {
	return &m.OrderMetadata
}
