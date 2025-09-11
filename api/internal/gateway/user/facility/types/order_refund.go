package types

// RefundType - 注文キャンセル種別
type RefundType int32

const (
	RefundTypeNone     RefundType = 0 // キャンセルなし
	RefundTypeCanceled RefundType = 1 // キャンセル
	RefundTypeRefunded RefundType = 2 // 返金
)

// OrderRefund - 注文キャンセル情報
type OrderRefund struct {
	Total      int64      `json:"total"`      // 返金金額
	Type       RefundType `json:"type"`       // 注文キャンセル種別
	Reason     string     `json:"reason"`     // 注文キャンセル理由
	Canceled   bool       `json:"canceled"`   // 注文キャンセルフラグ
	CanceledAt int64      `json:"canceledAt"` // 注文キャンセル日時
}
