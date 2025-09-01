package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Cart struct {
	response.Cart
}

type Carts []*Cart

type CartItem struct {
	response.CartItem
}

type CartItems []*CartItem

func NewCart(basket *entity.CartBasket) *Cart {
	return &Cart{
		Cart: response.Cart{
			Number:        basket.BoxNumber,
			Type:          NewShippingType(basket.BoxType).Response(),
			Size:          NewShippingSize(basket.BoxSize).Response(),
			Rate:          basket.BoxRate,
			Items:         NewCartItems(basket.Items).Response(),
			CoordinatorID: basket.CoordinatorID,
		},
	}
}

func (c *Cart) Response() *response.Cart {
	return &c.Cart
}

func NewCarts(cart *entity.Cart) Carts {
	if cart == nil {
		return Carts{}
	}
	res := make(Carts, len(cart.Baskets))
	for i := range cart.Baskets {
		res[i] = NewCart(cart.Baskets[i])
	}
	return res
}

func (cs Carts) Response() []*response.Cart {
	res := make([]*response.Cart, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}

func NewCartItem(item *entity.CartItem) *CartItem {
	return &CartItem{
		CartItem: response.CartItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		},
	}
}

func (i *CartItem) Response() *response.CartItem {
	return &i.CartItem
}

func NewCartItems(items entity.CartItems) CartItems {
	res := make(CartItems, len(items))
	for i := range items {
		res[i] = NewCartItem(items[i])
	}
	return res
}

func (is CartItems) Response() []*response.CartItem {
	res := make([]*response.CartItem, len(is))
	for i := range is {
		res[i] = is[i].Response()
	}
	return res
}
