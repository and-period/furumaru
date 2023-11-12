package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
)

// OrderItem - 注文商品情報
type OrderItem struct {
	FulfillmentID     string    `gorm:"primaryKey;<-:create"` // 注文配送ID
	ProductRevisionID int64     `gorm:"primaryKey;<-:create"` // 商品ID
	OrderID           string    `gorm:""`                     // 注文履歴ID
	Quantity          int64     `gorm:""`                     // 購入数量
	CreatedAt         time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time `gorm:""`                     // 更新日時
}

type OrderItems []*OrderItem

type NewOrderItemParams struct {
	OrderID       string
	FulfillmentID string
	Item          *CartItem
	Product       *Product
}

type NewOrderItemsParams struct {
	OrderID     string
	Fulfillment *OrderFulfillment
	Items       CartItems
	Products    map[string]*Product
}

func NewOrderItem(params *NewOrderItemParams) *OrderItem {
	return &OrderItem{
		FulfillmentID:     params.FulfillmentID,
		ProductRevisionID: params.Product.ProductRevision.ID,
		OrderID:           params.OrderID,
		Quantity:          params.Item.Quantity,
	}
}

func NewOrderItems(params *NewOrderItemsParams) (OrderItems, error) {
	res := make(OrderItems, len(params.Items))
	for i, item := range params.Items {
		product, ok := params.Products[item.ProductID]
		if !ok {
			return nil, errNotFoundProduct
		}
		iparams := &NewOrderItemParams{
			OrderID:       params.OrderID,
			FulfillmentID: params.Fulfillment.ID,
			Item:          item,
			Product:       product,
		}
		res[i] = NewOrderItem(iparams)
	}
	return res, nil
}

func (is OrderItems) ProductRevisionIDs() []int64 {
	return set.UniqBy(is, func(i *OrderItem) int64 {
		return i.ProductRevisionID
	})
}

func (is OrderItems) GroupByFulfillmentID() map[string]OrderItems {
	res := make(map[string]OrderItems, len(is))
	for _, i := range is {
		if _, ok := res[i.FulfillmentID]; !ok {
			res[i.FulfillmentID] = make(OrderItems, 0, len(is))
		}
		res[i.FulfillmentID] = append(res[i.FulfillmentID], i)
	}
	return res
}

func (is OrderItems) GroupByOrderID() map[string]OrderItems {
	res := make(map[string]OrderItems, len(is))
	for _, i := range is {
		if _, ok := res[i.OrderID]; !ok {
			res[i.OrderID] = make(OrderItems, 0, len(is))
		}
		res[i.OrderID] = append(res[i.OrderID], i)
	}
	return res
}
