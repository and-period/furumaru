package types

type CheckoutRequest struct {
	RequestID      string              `json:"requestId" validate:"required"`            // 支払いキー(重複判定用)
	CoordinatorID  string              `json:"coordinatorId" validate:"required"`        // コーディネータID
	BoxNumber      int64               `json:"boxNumber" validate:"min=0"`               // 箱の通番（箱単位で購入する場合）
	PromotionCode  string              `json:"promotionCode" validate:"omitempty,len=8"` // プロモーションコード
	PickupAt       int64               `json:"pickupAt" validate:"min=0"`                // 受け取り日時（UNIX時間）
	PickupLocation string              `json:"pickupLocation" validate:"max=256"`        // 受け取り場所
	OrderRequest   string              `json:"orderRequest" validate:"max=256"`          // 注文時リクエスト
	PaymentMethod  PaymentMethodType   `json:"paymentMethod" validate:"required"`        // 支払い方法
	CreditCard     *CheckoutCreditCard `json:"creditCard" validate:"omitempty,dive"`     // クレジットカード決済情報
	CallbackURL    string              `json:"callbackUrl" validate:"required,http_url"` // 決済完了後のリダイレクト先URL
	Total          int64               `json:"total" validate:"min=0"`                   // 支払い合計金額（誤り検出用）
}

type CheckoutCreditCard struct {
	Token             string `json:"token"`                                                                  // カードトークン（KOMOJU Tokens API で取得）
	Name              string `json:"name" validate:"required"`                                               // カード名義
	Number            string `json:"number" validate:"required_without=Token,omitempty,credit_card"`         // カード番号
	Month             int64  `json:"month" validate:"required_without=Token,omitempty,min=1,max=12"`         // 有効期限（月）
	Year              int64  `json:"year" validate:"required_without=Token,omitempty,min=2000,max=2100"`     // 有効期限（年）
	VerificationValue string `json:"verificationValue" validate:"required_without=Token,omitempty,min=3,max=4,numeric"` // セキュリティコード
}

type CheckoutResponse struct {
	URL string `json:"url"` // 遷移先URL
}

type CheckoutStateResponse struct {
	OrderID string        `json:"orderId"` // 注文履歴ID
	Status  PaymentStatus `json:"status"`  // 注文ステータス
}
