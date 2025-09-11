package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

type Cart struct {
	types.Cart
}

type Carts []*Cart

type CartItem struct {
	types.CartItem
}

type CartItems []*CartItem

func NewCart(basket *entity.CartBasket) *Cart {
	return &Cart{
		Cart: types.Cart{
			Number:        basket.BoxNumber,
			Type:          NewShippingType(basket.BoxType).Response(),
			Size:          NewShippingSize(basket.BoxSize).Response(),
			Rate:          basket.BoxRate,
			Items:         NewCartItems(basket.Items).Response(),
			CoordinatorID: basket.CoordinatorID,
		},
	}
}

func (c *Cart) Response() *types.Cart {
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

func (cs Carts) Response() []*types.Cart {
	res := make([]*types.Cart, len(cs))
	for i := range cs {
		res[i] = cs[i].Response()
	}
	return res
}

func NewCartItem(item *entity.CartItem) *CartItem {
	return &CartItem{
		CartItem: types.CartItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		},
	}
}

func (i *CartItem) Response() *types.CartItem {
	return &i.CartItem
}

func NewCartItems(items entity.CartItems) CartItems {
	res := make(CartItems, len(items))
	for i := range items {
		res[i] = NewCartItem(items[i])
	}
	return res
}

func (is CartItems) Response() []*types.CartItem {
	res := make([]*types.CartItem, len(is))
	for i := range is {
		res[i] = is[i].Response()
	}
	return res
}
