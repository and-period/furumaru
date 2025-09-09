package types

type CheckoutRequest struct {
	RequestID     string              `json:"requestId" validate:"required"`            // 支払いキー(重複判定用)
	CoordinatorID string              `json:"coordinatorId" validate:"required"`        // コーディネータID
	BoxNumber     int64               `json:"boxNumber" validate:"min=0"`               // 箱の通番（箱単位で購入する場合）
	PromotionCode string              `json:"promotionCode" validate:"omitempty,len=8"` // プロモーションコード
	PaymentMethod int32               `json:"paymentMethod" validate:"required"`        // 支払い方法
	CreditCard    *CheckoutCreditCard `json:"creditCard" validate:"omitempty,dive"`     // クレジットカード決済情報
	CallbackURL   string              `json:"callbackUrl" validate:"required,http_url"` // 決済完了後のリダイレクト先URL
	Total         int64               `json:"total" validate:"min=0"`                   // 支払い合計金額（誤り検出用）
}

type CheckoutCreditCard struct {
	Name              string `json:"name" validate:"required"`                         // カード名義
	Number            string `json:"number" validate:"required,credit_card"`           // カード番号
	Month             int64  `json:"month" validate:"min=1,max=12"`                    // 有効期限（月）
	Year              int64  `json:"year" validate:"min=2000,max=2100"`                // 有効期限（年）
	VerificationValue string `json:"verificationValue" validate:"min=3,max=4,numeric"` // セキュリティコード
}

type CheckoutResponse struct {
	URL string `json:"url"` // 遷移先URL
}

type CheckoutStateResponse struct {
	OrderID string `json:"orderId"` // 注文履歴ID
	Status  int32  `json:"status"`  // 注文ステータス
}
