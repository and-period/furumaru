package request

type CheckoutProductRequest struct {
	RequestID         string              `json:"requestId,omitempty"`         // 支払いキー(重複判定用)
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

type CheckoutExperienceRequest struct {
	RequestID             string              `json:"requestId,omitempty"`             // 支払いキー(重複判定用)
	BillingAddressID      string              `json:"billingAddressId,omitempty"`      // 請求先住所ID
	PromotionCode         string              `json:"promotionCode,omitempty"`         // プロモーションコード
	AdultCount            int64               `json:"adultCount,omitempty"`            // 大人人数
	JuniorHighSchoolCount int64               `json:"juniorHighSchoolCount,omitempty"` // 中学生人数
	ElementarySchoolCount int64               `json:"elementarySchoolCount,omitempty"` // 小学生人数
	PreschoolCount        int64               `json:"preschoolCount,omitempty"`        // 幼児人数
	SeniorCount           int64               `json:"seniorCount,omitempty"`           // シニア人数
	Transportation        string              `json:"transportation,omitempty"`        // 交通手段
	RequestedDate         string              `json:"requestedDate,omitempty"`         // 体験希望日
	RequestedTime         string              `json:"requestedTime,omitempty"`         // 体験希望時間
	PaymentMethod         int32               `json:"paymentMethod,omitempty"`         // 支払い方法
	CreditCard            *CheckoutCreditCard `json:"creditCard,omitempty"`            // クレジットカード決済情報
	CallbackURL           string              `json:"callbackUrl,omitempty"`           // 決済完了後のリダイレクト先URL
	Total                 int64               `json:"total,omitempty"`                 // 支払い合計金額（誤り検出用）
}

type CheckoutCreditCard struct {
	Name              string `json:"name,omitempty"`              // カード名義
	Number            string `json:"number,omitempty"`            // カード番号
	Month             int64  `json:"month,omitempty"`             // 有効期限（月）
	Year              int64  `json:"year,omitempty"`              // 有効期限（年）
	VerificationValue string `json:"verificationValue,omitempty"` // セキュリティコード
}
