package entity

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/pkg/set"
)

var bascketWeightLimits = map[ShippingSize]int64{
	ShippingSize60:  2e3,  //  2kg =  2,000g
	ShippingSize80:  5e3,  //  5kg =  5,000g
	ShippingSize100: 10e3, // 10kg = 10,000g
}

// Cart - カート情報
type Cart struct {
	SessionID string      `dynamodbav:"session_id"`          // セッションID
	Baskets   CartBaskets `dynamodbav:"baskets"`             // 買い物かご一覧
	ExpiredAt time.Time   `dynamodbav:"expired_at,unixtime"` // 有効期限
	CreatedAt time.Time   `dynamodbav:"created_at"`          // 登録日時
	UpdatedAt time.Time   `dynamodbav:"updated_at"`          // 更新日時
}

// CartBasket - 買い物かご情報
type CartBasket struct {
	BoxNumber     int64        `dynamodbav:"box_number"` // 箱の通番
	BoxType       DeliveryType `dynamodbav:"box_type"`   // 箱の種別
	BoxSize       ShippingSize `dynamodbav:"box_size"`   // 箱のサイズ
	Items         CartItems    `dynamodbav:"items"`      // 商品一覧
	coordinatorID string       `dynamodbav:"-"`          // コーディネータID
}

type CartBaskets []*CartBasket

// CartItem - 買い物かご
type CartItem struct {
	ProductID string `dynamodbav:"product_id"` // 商品ID
	Quantity  int64  `dynamodbav:"quantity"`   // 数量
}

type CartItems []*CartItem

// cartItem - 買い物かごの分割グループ
type cartGroup struct {
	key           string
	coordinatorID string
	boxType       DeliveryType
	products      Products
}

type cartGroups []*cartGroup

// Refresh - カート内の整理
func (c *Cart) Refresh(products Products) error {
	baskets, err := refreshCart(c.Baskets, products.Map())
	if err != nil {
		return err
	}
	c.Baskets = baskets
	return nil
}

// MergeByProductID - 商品IDを基に、買い物かご内の商品を統合
func (bs CartBaskets) MergeByProductID() CartItems {
	items := make(map[string]int64, len(bs))
	for i := range bs {
		for _, item := range bs[i].Items {
			items[item.ProductID] += item.Quantity
		}
	}
	res := make(CartItems, 0, len(items))
	for productID, quantity := range items {
		item := &CartItem{
			ProductID: productID,
			Quantity:  quantity,
		}
		res = append(res, item)
	}
	return res
}

// AdjustItems - 商品在庫数に合わせて、買い物かご内の商品数量を調整
func (bs CartBaskets) AdjustItems(products map[string]*Product) CartItems {
	items := bs.MergeByProductID()
	res := make(CartItems, 0, len(items))
	for _, item := range items {
		product, ok := products[item.ProductID]
		if !ok {
			continue
		}
		if item.Quantity > product.Inventory {
			item.Quantity = product.Inventory
		}
		res = append(res, item)
	}
	return res
}

func (bs CartBaskets) ProductIDs() []string {
	set := set.NewEmpty[string](len(bs))
	for i := range bs {
		set.Add(bs[i].Items.ProductIDs()...)
	}
	return set.Slice()
}

func NewCartItem(productID string, quantity int64) *CartItem {
	return &CartItem{
		ProductID: productID,
		Quantity:  quantity,
	}
}

func NewCartItems(items map[string]int64) CartItems {
	res := make(CartItems, 0, len(items))
	for productID, quantity := range items {
		res = append(res, NewCartItem(productID, quantity))
	}
	return res
}

func NewCartItemsWithProducts(products Products) CartItems {
	items := make(map[string]int64, len(products))
	for _, p := range products {
		items[p.ID]++
	}
	return NewCartItems(items)
}

func (is CartItems) ProductIDs() []string {
	return set.UniqBy(is, func(i *CartItem) string {
		return i.ProductID
	})
}

func (is CartItems) divide(products map[string]*Product) (cartGroups, error) {
	citems := is.groupByCartBasketKey(products)
	res := make(cartGroups, 0, len(citems))
	for key, items := range citems {
		coordinatorID, boxType, err := parseCartBasketKey(key)
		if err != nil {
			return nil, err
		}
		group := &cartGroup{
			key:           key,
			coordinatorID: coordinatorID,
			boxType:       boxType,
			products:      items,
		}
		res = append(res, group)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].key <= res[j].key
	})
	return res, nil
}

func (is CartItems) groupByCartBasketKey(products map[string]*Product) map[string]Products {
	res := make(map[string]Products, len(products))
	for _, item := range is {
		product, ok := products[item.ProductID]
		if !ok {
			continue
		}
		key := generateCartBasketKey(product.CoordinatorID, product.DeliveryType)
		if _, ok := res[key]; !ok {
			res[key] = make(Products, 0, item.Quantity)
		}
		for i := 0; i < int(item.Quantity); i++ {
			res[key] = append(res[key], product)
		}
	}
	return res
}

func generateCartBasketKey(coordinatorID string, typ DeliveryType) string {
	return fmt.Sprintf("%s:%d", coordinatorID, typ)
}

func parseCartBasketKey(key string) (string, DeliveryType, error) {
	strs := strings.Split(key, ":")
	if len(strs) != 2 {
		return "", DeliveryTypeUnknown, errors.New("invalid cart basket key format")
	}
	typ, err := strconv.ParseInt(strs[1], 10, 64)
	if err != nil {
		return "", DeliveryTypeUnknown, err
	}
	return strs[0], DeliveryType(typ), nil
}

func refreshCart(baskets CartBaskets, products map[string]*Product) (CartBaskets, error) {
	// 商品在庫数に合わせて買い物かご内の商品数量を調整
	items := baskets.AdjustItems(products)

	// 整理前の商品一覧の作成（同じ商品は、数量分レコード作成）
	groups, err := items.divide(products)
	if err != nil {
		return nil, err
	}

	boxNumber := int64(1)
	res := make(CartBaskets, 0, len(baskets))

	// 買い物かごの整理(再作成)
	for _, group := range groups {
		stacks := group.products

		// 箱の占有率が大きいもの順に並び替え
		sort.Slice(stacks, func(i, j int) bool {
			return stacks[i].Box100Rate >= stacks[j].Box100Rate
		})

		for {
			// 商品がまだあるかの検証
			if len(stacks) == 0 {
				break
			}

			// 箱に詰める商品と詰めない商品を分割
			picked, rest := pickCartBasket(stacks)

			// 箱の作成
			basket := &CartBasket{
				BoxNumber:     boxNumber,
				BoxType:       group.boxType,
				BoxSize:       calcShippingSize(picked),
				Items:         NewCartItemsWithProducts(picked),
				coordinatorID: group.coordinatorID,
			}
			res = append(res, basket)

			boxNumber++
			stacks = rest
		}
	}

	return res, nil
}

// pickCartBasket - 箱に詰める商品と詰めない商品を分割
func pickCartBasket(products Products) (Products, Products) {
	freeRate := int64(100)
	freeWeight := bascketWeightLimits[ShippingSize100]

	picked := make(Products, 0, len(products)) // 箱に入ったもの
	rest := make(Products, 0, len(products))   // 箱に入らなかったもの
	for i, product := range products {
		// すでにいっぱいの場合、次の箱に詰めるようにする
		if freeRate == 0 || freeWeight == 0 {
			rest = append(rest, products[i:]...)
			break
		}
		weight := product.WeightGram()
		// 箱に入るかの検証
		if product.Box100Rate > freeRate || weight > freeWeight {
			rest = append(rest, product)
			continue
		}
		picked = append(picked, product)
		freeRate -= product.Box100Rate
		freeWeight -= weight
	}

	return picked, rest
}

// calcShippingSize - 買い物かごの占有率を見て、かごの大きさを決定
func calcShippingSize(products Products) ShippingSize {
	weight := products.WeightGram()
	if products.Box60Rate() <= 100 && weight <= bascketWeightLimits[ShippingSize60] {
		return ShippingSize60
	}
	if products.Box80Rate() <= 100 && weight <= bascketWeightLimits[ShippingSize80] {
		return ShippingSize80
	}
	return ShippingSize100
}
