package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
	set "github.com/and-period/furumaru/api/pkg/set/v2"
)

type OrderItem struct {
	response.OrderItem
	orderID string
}

type OrderItems []*OrderItem

func NewOrderItem(item *entity.OrderItem) *OrderItem {
	return &OrderItem{
		OrderItem: response.OrderItem{
			ProductID: item.ProductID,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Weight:    NewProductWeight(item.Weight, item.WeightUnit),
		},
		orderID: item.OrderID,
	}
}

func (i *OrderItem) Fill(product *Product) {
	if product != nil {
		i.Name = product.Name
		i.Media = product.Media
	}
}

func (i *OrderItem) Response() *response.OrderItem {
	return &i.OrderItem
}

func NewOrderItems(items entity.OrderItems) OrderItems {
	res := make(OrderItems, len(items))
	for i := range items {
		res[i] = NewOrderItem(items[i])
	}
	return res
}

func (is OrderItems) Fill(products map[string]*Product) {
	for i := range is {
		is[i].Fill(products[is[i].ProductID])
	}
}

func (is OrderItems) ProductIDs() []string {
	return set.UniqBy(is, func(i *OrderItem) string {
		return i.ProductID
	})
}

func (is OrderItems) Response() []*response.OrderItem {
	res := make([]*response.OrderItem, len(is))
	for i := range is {
		res[i] = is[i].Response()
	}
	return res
}
