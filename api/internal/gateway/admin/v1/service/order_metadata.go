package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
)

type OrderMetadata struct {
	response.OrderMetadata
	orderID string
}

func NewOrderMetadata(metadata *entity.OrderMetadata) *OrderMetadata {
	return &OrderMetadata{
		OrderMetadata: response.OrderMetadata{
			PickupAt:       jst.Unix(metadata.PickupAt),
			PickupLocation: metadata.PickupLocation,
		},
		orderID: metadata.OrderID,
	}
}

func (m *OrderMetadata) Response() *response.OrderMetadata {
	return &m.OrderMetadata
}
