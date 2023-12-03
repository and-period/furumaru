package entity

import (
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// Order - 注文履歴情報
type Order struct {
	OrderPayment      `gorm:"-"`
	OrderFulfillments `gorm:"-"`
	OrderItems        `gorm:"-"`
	ID                string         `gorm:"primaryKey;<-:create"` // 注文履歴ID
	UserID            string         `gorm:""`                     // ユーザーID
	CoordinatorID     string         `gorm:""`                     // 注文受付担当者ID
	PromotionID       string         `gorm:"default:null"`         // プロモーションID
	ShippingMessage   string         `gorm:"default:null"`         // 発送時のメッセージ
	CreatedAt         time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time      `gorm:""`                     // 更新日時
	CompletedAt       time.Time      `gorm:"default:null"`         // 対応完了日時
	DeletedAt         gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Orders []*Order

// AggregatedOrder - 注文履歴集計情報
type AggregatedOrder struct {
	UserID     string // ユーザーID
	OrderCount int64  // 注文合計回数
	Subtotal   int64  // 購入合計金額
	Discount   int64  // 割引合計金額
}

type AggregatedOrders []*AggregatedOrder

type NewOrderParams struct {
	CoordinatorID     string
	Customer          *entity.User
	BillingAddress    *entity.Address
	ShippingAddress   *entity.Address
	Shipping          *Shipping
	Baskets           CartBaskets
	Products          Products
	PaymentMethodType PaymentMethodType
	Promotion         *Promotion
}

func NewOrder(params *NewOrderParams) (*Order, error) {
	var promotionID string
	if params.Promotion != nil {
		promotionID = params.Promotion.ID
	}
	orderID := uuid.Base58Encode(uuid.New())
	pparams := &NewOrderPaymentParams{
		OrderID:    orderID,
		Address:    params.BillingAddress,
		MethodType: params.PaymentMethodType,
		Baskets:    params.Baskets,
		Products:   params.Products,
		Shipping:   params.Shipping,
		Promotion:  params.Promotion,
	}
	payment, err := NewOrderPayment(pparams)
	if err != nil {
		return nil, err
	}
	fparams := &NewOrderFulfillmentsParams{
		OrderID:  orderID,
		Address:  params.ShippingAddress,
		Baskets:  params.Baskets,
		Products: params.Products.Map(),
	}
	fulfillments, items, err := NewOrderFulfillments(fparams)
	if err != nil {
		return nil, err
	}
	return &Order{
		OrderPayment:      *payment,
		OrderFulfillments: fulfillments,
		OrderItems:        items,
		ID:                orderID,
		UserID:            params.Customer.ID,
		CoordinatorID:     params.CoordinatorID,
		PromotionID:       promotionID,
	}, nil
}

func (o *Order) Fill(payment *OrderPayment, fulfillments OrderFulfillments, items OrderItems) {
	o.OrderPayment = *payment
	o.OrderFulfillments = fulfillments
	o.OrderItems = items
}

func (os Orders) IDs() []string {
	res := make([]string, len(os))
	for i := range os {
		res[i] = os[i].ID
	}
	return res
}

func (os Orders) UserIDs() []string {
	return set.UniqBy(os, func(o *Order) string {
		return o.UserID
	})
}

func (os Orders) CoordinatorIDs() []string {
	return set.UniqBy(os, func(o *Order) string {
		return o.CoordinatorID
	})
}

func (os Orders) PromotionIDs() []string {
	res := set.NewEmpty[string](len(os))
	for i := range os {
		if os[i].PromotionID == "" {
			continue
		}
		res.Add(os[i].PromotionID)
	}
	return res.Slice()
}

func (os Orders) AddressRevisionIDs() []int64 {
	res := set.NewEmpty[int64](len(os) * 2) // payment + fulfillment
	for i := range os {
		res.Add(os[i].OrderPayment.AddressRevisionID)
		res.Add(os[i].OrderFulfillments.AddressRevisionIDs()...)
	}
	return res.Slice()
}

func (os Orders) ProductRevisionIDs() []int64 {
	res := set.NewEmpty[int64](len(os))
	for i := range os {
		res.Add(os[i].ProductRevisionIDs()...)
	}
	return res.Slice()
}

func (os Orders) Fill(payments map[string]*OrderPayment, fulfillments map[string]OrderFulfillments, items map[string]OrderItems) {
	for _, o := range os {
		payment, ok := payments[o.ID]
		if !ok {
			payment = &OrderPayment{}
		}
		o.Fill(payment, fulfillments[o.ID], items[o.ID])
	}
}

func (os AggregatedOrders) Map() map[string]*AggregatedOrder {
	res := make(map[string]*AggregatedOrder, len(os))
	for _, o := range os {
		res[o.UserID] = o
	}
	return res
}
