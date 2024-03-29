package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type OrderItem struct {
	response.OrderItem
	orderID string
}

type OrderItems []*OrderItem

func NewOrderItem(item *entity.OrderItem, product *Product) *OrderItem {
	var (
		productID string
		price     int64
	)
	if product != nil {
		productID = product.ID
		price = product.Price
	}
	return &OrderItem{
		OrderItem: response.OrderItem{
			FulfillmentID: item.FulfillmentID,
			ProductID:     productID,
			Price:         price,
			Quantity:      item.Quantity,
		},
		orderID: item.OrderID,
	}
}

func (i *OrderItem) Response() *response.OrderItem {
	return &i.OrderItem
}

func NewOrderItems(items entity.OrderItems, products map[int64]*Product) OrderItems {
	res := make(OrderItems, len(items))
	for i, v := range items {
		res[i] = NewOrderItem(v, products[v.ProductRevisionID])
	}
	return res
}

func (is OrderItems) Response() []*response.OrderItem {
	res := make([]*response.OrderItem, len(is))
	for i := range is {
		res[i] = is[i].Response()
	}
	return res
}
