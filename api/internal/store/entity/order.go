package entity

import "time"

// 支払いステータス
type PaymentStatus int32

const (
	PaymentStatusUnknown     PaymentStatus = 0
	PaymentStatusInitialized PaymentStatus = 1 // 未完了
	PaymentStatusPending     PaymentStatus = 2 // 保留中
	PaymentStatusAuthorized  PaymentStatus = 3 // 仮売上・オーソリ
	PaymentStatusCaptured    PaymentStatus = 4 // 実売上・キャプチャ
	PaymentStatusCanceled    PaymentStatus = 5 // キャンセル
	PaymentStatusFailed      PaymentStatus = 6 // 失敗/期限切れ
)

// 配送ステータス
type FulfillmentStatus int32

const (
	FulfillmentStatusUnknown     FulfillmentStatus = 0
	FulfillmentStatusUnfulfilled FulfillmentStatus = 1 // 未発送
	FulfillmentStatusFulfilled   FulfillmentStatus = 2 // 発送済み
)

// 注文キャンセル種別
type CancelType int32

const (
	CancelTypeUnknown CancelType = 0
)

type OrderOrderBy string

const (
	OrderOrderByPaymentStatus     OrderOrderBy = "payment_status"
	OrderOrderByFulfillmentStatus OrderOrderBy = "fulfillment_status"
	OrderOrderByCanceledAt        OrderOrderBy = "canceled_at"
	OrderOrderByOrderedAt         OrderOrderBy = "ordered_at"
	OrderOrderByConfirmedAt       OrderOrderBy = "confirmed_at"
	OrderOrderByCapturedAt        OrderOrderBy = "captured_at"
	OrderOrderByDeliveredAt       OrderOrderBy = "delivered_at"
	OrderOrderByCreatedAt         OrderOrderBy = "created_at"
	OrderOrderByUpdatedAt         OrderOrderBy = "updated_at"
)

// Order - 注文履歴情報
type Order struct {
	OrderItems        `gorm:"-"`
	OrderPayment      `gorm:"-"`
	OrderFulfillment  `gorm:"-"`
	OrderActivities   `gorm:"-"`
	ID                string            `gorm:"primaryKey;<-:create"` // 注文履歴ID
	UserID            string            `gorm:""`                     // ユーザーID
	ScheduleID        string            `gorm:""`                     // 購入対象マルシェID
	CoordinatorID     string            `gorm:""`                     // 購入対象仲介者ID
	PaymentStatus     PaymentStatus     `gorm:""`                     // 支払いステータス
	FulfillmentStatus FulfillmentStatus `gorm:""`                     // 配送ステータス
	CancelType        CancelType        `gorm:""`                     // 注文キャンセル種別
	CancelReason      string            `gorm:""`                     // 注文キャンセル理由
	CanceledAt        time.Time         `gorm:"default:null"`         // 注文キャンセル日時
	OrderedAt         time.Time         `gorm:"default:null"`         // 決済要求日時
	ConfirmedAt       time.Time         `gorm:"default:null"`         // 決済実行日時
	CapturedAt        time.Time         `gorm:"default:null"`         // 決済確定日時
	DeliveredAt       time.Time         `gorm:"default:null"`         // 配送日時
	CreatedAt         time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time         `gorm:""`                     // 更新日時
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

func (o *Order) Fill(
	items OrderItems, payment *OrderPayment, fulfillment *OrderFulfillment, activities OrderActivities,
) {
	o.OrderItems = items
	o.OrderPayment = *payment
	o.OrderFulfillment = *fulfillment
	o.OrderActivities = activities
}

func (os Orders) IDs() []string {
	res := make([]string, len(os))
	for i := range os {
		res[i] = os[i].ID
	}
	return res
}
