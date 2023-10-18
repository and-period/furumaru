package entity

import (
	"time"

	"gorm.io/gorm"
)

// 支払いステータス
type PaymentStatus int32

const (
	PaymentStatusUnknown    PaymentStatus = 0
	PaymentStatusPending    PaymentStatus = 1 // 保留中
	PaymentStatusAuthorized PaymentStatus = 2 // 仮売上・オーソリ
	PaymentStatusCaptured   PaymentStatus = 3 // 実売上・キャプチャ
	PaymentStatusRefunded   PaymentStatus = 4 // 返金
	PaymentStatusFailed     PaymentStatus = 5 // 失敗/期限切れ
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
	Payment           `gorm:"-"`
	Fulfillment       `gorm:"-"`
	Activities        `gorm:"-"`
	OrderItems        `gorm:"-"`
	ID                string            `gorm:"primaryKey;<-:create"` // 注文履歴ID
	UserID            string            `gorm:""`                     // ユーザーID
	CoordinatorID     string            `gorm:""`                     // 注文受付担当者ID
	ScheduleID        string            `gorm:"default:null"`         // マルシェ開催スケジュールID
	PromotionID       string            `gorm:"default:null"`         // プロモーションID
	PaymentStatus     PaymentStatus     `gorm:""`                     // 支払いステータス
	FulfillmentStatus FulfillmentStatus `gorm:""`                     // 配送ステータス
	RefundReason      string            `gorm:""`                     // 注文キャンセル理由
	CreatedAt         time.Time         `gorm:"<-:create"`            // 登録日時
	UpdatedAt         time.Time         `gorm:""`                     // 更新日時
	OrderedAt         time.Time         `gorm:"default:null"`         // 決済要求日時
	PaidAt            time.Time         `gorm:"default:null"`         // 決済承認日時(仮売上)
	CapturedAt        time.Time         `gorm:"default:null"`         // 決済確定日時(実売上)
	FailedAt          time.Time         `gorm:"default:null"`         // 決済失敗日時
	RefundedAt        time.Time         `gorm:"default:null"`         // キャンセル日時(返金)
	ShippedAt         time.Time         `gorm:"default:null"`         // 配送日時
	DeletedAt         gorm.DeletedAt    `gorm:"default:null"`         // 削除日時
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

func (o *Order) Fill(payment *Payment, fulfillment *Fulfillment, activities Activities, items OrderItems) {
	o.Payment = *payment
	o.Fulfillment = *fulfillment
	o.Activities = activities
	o.OrderItems = items
}

func (o *Order) IsCanceled() bool {
	return o.PaymentStatus == PaymentStatusRefunded
}

func (os Orders) IDs() []string {
	res := make([]string, len(os))
	for i := range os {
		res[i] = os[i].ID
	}
	return res
}

func (os AggregatedOrders) Map() map[string]*AggregatedOrder {
	res := make(map[string]*AggregatedOrder, len(os))
	for _, o := range os {
		res[o.UserID] = o
	}
	return res
}
