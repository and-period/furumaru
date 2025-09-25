package entity

import (
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/set"
	"gorm.io/gorm"
)

// OrderType - 注文種別
type OrderType int32

const (
	OrderTypeUnknown    OrderType = 0
	OrderTypeProduct    OrderType = 1 // 商品
	OrderTypeExperience OrderType = 2 // 体験
)

// OrderStatus - 注文ステータス
type OrderStatus int32

const (
	OrderStatusUnknown   OrderStatus = 0
	OrderStatusUnpaid    OrderStatus = 1 // 支払い待ち
	OrderStatusWaiting   OrderStatus = 2 // 受注待ち
	OrderStatusPreparing OrderStatus = 3 // 発送準備中
	OrderStatusShipped   OrderStatus = 4 // 発送完了
	OrderStatusCompleted OrderStatus = 5 // 完了
	OrderStatusCanceled  OrderStatus = 6 // キャンセル
	OrderStatusRefunded  OrderStatus = 7 // 返金
	OrderStatusFailed    OrderStatus = 8 // 失敗
)

// OrderShippingType - 発送方法
type OrderShippingType int32

const (
	OrderShippingTypeUnknown  OrderShippingType = 0
	OrderShippingTypeNone     OrderShippingType = 1 // 発送なし
	OrderShippingTypeStandard OrderShippingType = 2 // 通常配送
	OrderShippingTypePickup   OrderShippingType = 3 // 店舗受取
)

// AggregateOrderPeriodType - 注文集計期間種別
type AggregateOrderPeriodType string

const (
	AggregateOrderPeriodTypeDay   AggregateOrderPeriodType = "day"   // 日
	AggregateOrderPeriodTypeWeek  AggregateOrderPeriodType = "week"  // 週
	AggregateOrderPeriodTypeMonth AggregateOrderPeriodType = "month" // 月
)

// Order - 注文履歴情報
type Order struct {
	OrderPayment      `gorm:"-"`
	OrderFulfillments `gorm:"-"`
	OrderItems        `gorm:"-"`
	OrderExperience   `gorm:"-"`
	OrderMetadata     `gorm:"-"`
	ID                string         `gorm:"primaryKey;<-:create"` // 注文履歴ID
	UserID            string         `gorm:""`                     // ユーザーID
	SessionID         string         `gorm:""`                     // 注文時セッションID
	ShopID            string         `gorm:"default:null"`         // 店舗ID
	CoordinatorID     string         `gorm:""`                     // 注文受付担当者ID
	PromotionID       string         `gorm:"default:null"`         // プロモーションID
	ManagementID      int64          `gorm:""`                     // 管理番号
	Type              OrderType      `gorm:""`                     // 注文種別
	Status            OrderStatus    `gorm:""`                     // 注文ステータス
	CreatedAt         time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time      `gorm:""`                     // 更新日時
	CompletedAt       time.Time      `gorm:"default:null"`         // 対応完了日時
	DeletedAt         gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

type Orders []*Order

type NewProductOrderParams struct {
	OrderID           string
	SessionID         string
	ShopID            string
	CoordinatorID     string
	Customer          *entity.User
	BillingAddress    *entity.Address
	ShippingAddress   *entity.Address
	Shipping          *Shipping
	Baskets           CartBaskets
	Products          Products
	PaymentMethodType PaymentMethodType
	Promotion         *Promotion
	Pickup            bool
	PickupAt          time.Time
	PickupLocation    string
	OrderRequest      string
}

type NewExperienceOrderParams struct {
	OrderID               string
	SessionID             string
	ShopID                string
	CoordinatorID         string
	Customer              *entity.User
	BillingAddress        *entity.Address
	Experience            *Experience
	PaymentMethodType     PaymentMethodType
	Promotion             *Promotion
	AdultCount            int64
	JuniorHighSchoolCount int64
	ElementarySchoolCount int64
	PreschoolCount        int64
	SeniorCount           int64
	Transportation        string
	RequetsedDate         string
	RequetsedTime         string
	OrderRequest          string
}

func NewProductOrder(params *NewProductOrderParams) (*Order, error) {
	var promotionID string
	if params.Promotion != nil {
		promotionID = params.Promotion.ID
	}
	pparams := &NewProductOrderPaymentParams{
		OrderID:    params.OrderID,
		Pickup:     params.Pickup,
		Address:    params.BillingAddress,
		MethodType: params.PaymentMethodType,
		Baskets:    params.Baskets,
		Products:   params.Products,
		Shipping:   params.Shipping,
		Promotion:  params.Promotion,
	}
	payment, err := NewProductOrderPayment(pparams)
	if err != nil {
		return nil, err
	}
	fparams := &NewOrderFulfillmentsParams{
		OrderID:  params.OrderID,
		Pickup:   params.Pickup,
		Address:  params.ShippingAddress,
		Baskets:  params.Baskets,
		Products: params.Products.Map(),
	}
	fulfillments, items, err := NewOrderFulfillments(fparams)
	if err != nil {
		return nil, err
	}
	mparams := &NewOrderMetadataParams{
		OrderID:         params.OrderID,
		Pickup:          params.Pickup,
		ShippingAddress: params.ShippingAddress,
		ShippingMessage: "ご注文ありがとうございます！商品到着まで今しばらくお待ち下さい。",
		PickupAt:        params.PickupAt,
		PickupLocation:  params.PickupLocation,
		OrderRequest:    params.OrderRequest,
	}
	metadata := NewOrderMetadata(mparams)
	return &Order{
		OrderPayment:      *payment,
		OrderFulfillments: fulfillments,
		OrderItems:        items,
		OrderMetadata:     *metadata,
		ID:                params.OrderID,
		SessionID:         params.SessionID,
		UserID:            params.Customer.ID,
		ShopID:            params.ShopID,
		CoordinatorID:     params.CoordinatorID,
		PromotionID:       promotionID,
		Type:              OrderTypeProduct,
		Status:            OrderStatusUnpaid, // 初期ステータスは「支払い待ち」で登録
	}, nil
}

func NewExperienceOrder(params *NewExperienceOrderParams) (*Order, error) {
	var promotionID string
	if params.Promotion != nil {
		promotionID = params.Promotion.ID
	}
	pparams := &NewExperienceOrderPaymentParams{
		OrderID:               params.OrderID,
		MethodType:            params.PaymentMethodType,
		Address:               params.BillingAddress,
		Experience:            params.Experience,
		Promotion:             params.Promotion,
		AdultCount:            params.AdultCount,
		JuniorHighSchoolCount: params.JuniorHighSchoolCount,
		ElementarySchoolCount: params.ElementarySchoolCount,
		PreschoolCount:        params.PreschoolCount,
		SeniorCount:           params.SeniorCount,
	}
	payment, err := NewExperienceOrderPayment(pparams)
	if err != nil {
		return nil, err
	}
	eparams := &NewOrderExperienceParams{
		OrderID:               params.OrderID,
		Experience:            params.Experience,
		AdultCount:            params.AdultCount,
		JuniorHighSchoolCount: params.JuniorHighSchoolCount,
		ElementarySchoolCount: params.ElementarySchoolCount,
		PreschoolCount:        params.PreschoolCount,
		SeniorCount:           params.SeniorCount,
		Transportation:        params.Transportation,
		RequestedDate:         params.RequetsedDate,
		RequestedTime:         params.RequetsedTime,
	}
	experience, err := NewOrderExperience(eparams)
	if err != nil {
		return nil, err
	}
	mparams := &NewOrderMetadataParams{
		OrderID: params.OrderID,
	}
	metadata := NewOrderMetadata(mparams)
	return &Order{
		OrderPayment:    *payment,
		OrderExperience: *experience,
		OrderMetadata:   *metadata,
		ID:              params.OrderID,
		SessionID:       params.SessionID,
		UserID:          params.Customer.ID,
		ShopID:          params.ShopID,
		CoordinatorID:   params.CoordinatorID,
		PromotionID:     promotionID,
		Type:            OrderTypeExperience,
		Status:          OrderStatusUnpaid, // 初期ステータスは「支払い待ち」で登録
	}, nil
}

func (o *Order) Fill(
	payment *OrderPayment,
	fulfillments OrderFulfillments,
	items OrderItems,
	experience *OrderExperience,
	metadata *OrderMetadata,
) {
	if payment != nil {
		o.OrderPayment = *payment
	}
	if experience != nil {
		o.OrderExperience = *experience
	}
	if metadata != nil {
		o.OrderMetadata = *metadata
	}
	o.OrderFulfillments = fulfillments
	o.OrderItems = items
}

func (o *Order) SetPaymentStatus(status PaymentStatus) {
	switch status {
	case PaymentStatusPending:
		o.Status = OrderStatusUnpaid
	case PaymentStatusAuthorized:
		o.Status = OrderStatusWaiting
	case PaymentStatusCaptured:
		o.Status = OrderStatusPreparing
	case PaymentStatusCanceled:
		o.Status = OrderStatusCanceled
	case PaymentStatusRefunded:
		o.Status = OrderStatusRefunded
	case PaymentStatusFailed:
		o.Status = OrderStatusFailed
	case PaymentStatusExpired:
		o.Status = OrderStatusFailed
	default:
		o.Status = OrderStatusUnknown
	}
}

func (o *Order) SetFulfillmentStatus(fulfillmentID string, status FulfillmentStatus) {
	for _, f := range o.OrderFulfillments {
		if f.ID == fulfillmentID {
			f.Status = status
			break
		}
	}
	if o.Fulfilled() {
		o.Status = OrderStatusShipped
	} else {
		o.Status = OrderStatusPreparing
	}
}

func (o *Order) SetTransaction(transactionID string, now time.Time) {
	if o.Total > 0 {
		o.SetTransactionID(transactionID, now)
		return
	}
	// 金額が0円の場合は支払い処理が不要なため、トランザクションIDを詰めると同時に支払い完了とする
	o.Status = OrderStatusPreparing
	o.SetTransactionID(o.ID, now) // 支払い状態を取得できるよう、なにかしらの値を詰める
	o.MethodType = PaymentMethodTypeNone
	o.OrderPayment.Status = PaymentStatusCaptured
	o.PaidAt, o.CapturedAt = now, now
}

func (o *Order) Completed() bool {
	if o == nil {
		return false
	}
	switch o.Status {
	case OrderStatusCompleted, OrderStatusCanceled, OrderStatusRefunded, OrderStatusFailed:
		return true
	default:
		return false
	}
}

func (o *Order) Preservable() bool {
	if o == nil {
		return false
	}
	return o.Capturable() || o.Completable()
}

func (o *Order) Capturable() bool {
	if o == nil {
		return false
	}
	return o.Status == OrderStatusWaiting
}

func (o *Order) Completable() bool {
	if o == nil {
		return false
	}
	return o.Status == OrderStatusPreparing || o.Status == OrderStatusShipped
}

func (o *Order) Cancelable() bool {
	if o == nil {
		return false
	}
	return o.OrderPayment.Status == PaymentStatusPending || o.OrderPayment.Status == PaymentStatusAuthorized
}

func (o *Order) Refundable() bool {
	if o == nil {
		return false
	}
	return o.OrderPayment.Status == PaymentStatusCaptured
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
		res.Add(os[i].AddressRevisionID)
		res.Add(os[i].AddressRevisionIDs()...)
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

func (os Orders) ExperienceRevisionIDs() []int64 {
	res := set.NewEmpty[int64](len(os))
	for i := range os {
		if os[i].ExperienceRevisionID == 0 {
			continue
		}
		res.Add(os[i].ExperienceRevisionID)
	}
	return res.Slice()
}

func (os Orders) Fill(
	payments map[string]*OrderPayment,
	fulfillments map[string]OrderFulfillments,
	items map[string]OrderItems,
	experiences map[string]*OrderExperience,
	metadata map[string]*OrderMetadata,
) {
	for _, o := range os {
		payment, ok := payments[o.ID]
		if !ok {
			payment = &OrderPayment{}
		}
		experience, ok := experiences[o.ID]
		if !ok {
			experience = &OrderExperience{}
		}
		meta, ok := metadata[o.ID]
		if !ok {
			meta = &OrderMetadata{}
		}
		o.Fill(payment, fulfillments[o.ID], items[o.ID], experience, meta)
	}
}

// AggregatedOrder - 注文履歴集計情報
type AggregatedOrder struct {
	OrderCount    int64 // 注文合計回数
	UserCount     int64 // 注文ユーザー数
	SalesTotal    int64 // 購入合計金額
	DiscountTotal int64 // 割引合計金額
}

// AggregatedUserOrder - 注文履歴集計情報
type AggregatedUserOrder struct {
	UserID     string // ユーザーID
	OrderCount int64  // 注文合計回数
	Subtotal   int64  // 購入合計金額
	Discount   int64  // 割引合計金額
	Total      int64  // 支払合計金額
}

type AggregatedUserOrders []*AggregatedUserOrder

func (os AggregatedUserOrders) Map() map[string]*AggregatedUserOrder {
	res := make(map[string]*AggregatedUserOrder, len(os))
	for _, o := range os {
		res[o.UserID] = o
	}
	return res
}

// AggregatedOrderPayment - 支払い情報別集計情報
type AggregatedOrderPayment struct {
	PaymentMethodType PaymentMethodType // 支払い種別
	OrderCount        int64             // 注文合計回数
	UserCount         int64             // 注文ユーザー数
	SalesTotal        int64             // 購入合計金額
}

type AggregatedOrderPayments []*AggregatedOrderPayment

func (ps AggregatedOrderPayments) Map() map[PaymentMethodType]*AggregatedOrderPayment {
	res := make(map[PaymentMethodType]*AggregatedOrderPayment, len(ps))
	for _, p := range ps {
		res[p.PaymentMethodType] = p
	}
	return res
}

func (ps AggregatedOrderPayments) OrderTotal() int64 {
	var total int64
	for _, p := range ps {
		total += p.OrderCount
	}
	return total
}

// AggregatedOrderPromotion - プロモーションコード利用履歴集計情報
type AggregatedOrderPromotion struct {
	PromotionID   string // プロモーションID
	OrderCount    int64  // 利用合計回数
	DiscountTotal int64  // 割引合計金額
}

type AggregatedOrderPromotions []*AggregatedOrderPromotion

func (os AggregatedOrderPromotions) Map() map[string]*AggregatedOrderPromotion {
	res := make(map[string]*AggregatedOrderPromotion, len(os))
	for _, o := range os {
		res[o.PromotionID] = o
	}
	return res
}

// AggregatedPeriodOrder - 注文履歴集計情報（期間別）
type AggregatedPeriodOrder struct {
	Period        time.Time // 期間
	OrderCount    int64     // 注文合計回数
	UserCount     int64     // 注文ユーザー数
	SalesTotal    int64     // 購入合計金額
	DiscountTotal int64     // 割引合計金額
}

type AggregatedPeriodOrders []*AggregatedPeriodOrder

func (os AggregatedPeriodOrders) MapByPeriod() map[time.Time]*AggregatedPeriodOrder {
	res := make(map[time.Time]*AggregatedPeriodOrder, len(os))
	for _, o := range os {
		res[o.Period] = o
	}
	return res
}
