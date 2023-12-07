package request

type CheckoutRequest struct {
	CoordinatorID     string              `json:"coordinatorId,omitempty"`     // コーディネータID
	BoxNumber         int64               `json:"boxNumber,omitempty"`         // 箱の通番（箱単位で購入する場合）
	BillingAddressID  string              `json:"billingAddressId,omitempty"`  // 請求先住所ID
	ShippingAddressID string              `json:"shippingAddressId,omitempty"` // 配送先住所ID
	PromotionCode     string              `json:"promotionCode,omitempty"`     // プロモーションコード
	PaymentMethod     int32               `json:"paymentMethod,omitempty"`     // 支払い方法
	CreditCard        *CheckoutCreditCard `json:"creditCard,omitempty"`        // クレジットカード決済情報
	CallbackURL       string              `json:"callbackUrl,omitempty"`       // 決済完了後のリダイレクト先URL
	Total             int64               `json:"total,omitempty"`             // 支払い合計金額（誤り検出用）
}

type CheckoutCreditCard struct {
	Number            string `json:"number,omitempty"`            // カード番号
	Month             int64  `json:"month,omitempty"`             // 有効期限（月）
	Year              int64  `json:"year,omitempty"`              // 有効期限（年）
	VerificationValue string `json:"verificationValue,omitempty"` // セキュリティコード
}
