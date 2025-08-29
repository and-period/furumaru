package request

type CheckoutProductRequest struct {
	RequestID         string              `json:"requestId"`         // 支払いキー(重複判定用)
	CoordinatorID     string              `json:"coordinatorId"`     // コーディネータID
	BoxNumber         int64               `json:"boxNumber"`         // 箱の通番（箱単位で購入する場合）
	BillingAddressID  string              `json:"billingAddressId"`  // 請求先住所ID
	ShippingAddressID string              `json:"shippingAddressId"` // 配送先住所ID
	PromotionCode     string              `json:"promotionCode"`     // プロモーションコード
	PaymentMethod     int32               `json:"paymentMethod"`     // 支払い方法
	CreditCard        *CheckoutCreditCard `json:"creditCard"`        // クレジットカード決済情報
	CallbackURL       string              `json:"callbackUrl"`       // 決済完了後のリダイレクト先URL
	Total             int64               `json:"total"`             // 支払い合計金額（誤り検出用）
}

type CheckoutExperienceRequest struct {
	RequestID             string              `json:"requestId"`             // 支払いキー(重複判定用)
	BillingAddressID      string              `json:"billingAddressId"`      // 請求先住所ID
	PromotionCode         string              `json:"promotionCode"`         // プロモーションコード
	AdultCount            int64               `json:"adultCount"`            // 大人人数
	JuniorHighSchoolCount int64               `json:"juniorHighSchoolCount"` // 中学生人数
	ElementarySchoolCount int64               `json:"elementarySchoolCount"` // 小学生人数
	PreschoolCount        int64               `json:"preschoolCount"`        // 幼児人数
	SeniorCount           int64               `json:"seniorCount"`           // シニア人数
	Transportation        string              `json:"transportation"`        // 交通手段
	RequestedDate         string              `json:"requestedDate"`         // 体験希望日
	RequestedTime         string              `json:"requestedTime"`         // 体験希望時間
	PaymentMethod         int32               `json:"paymentMethod"`         // 支払い方法
	CreditCard            *CheckoutCreditCard `json:"creditCard"`            // クレジットカード決済情報
	CallbackURL           string              `json:"callbackUrl"`           // 決済完了後のリダイレクト先URL
	Total                 int64               `json:"total"`                 // 支払い合計金額（誤り検出用）
}

type CheckoutCreditCard struct {
	Name              string `json:"name"`              // カード名義
	Number            string `json:"number"`            // カード番号
	Month             int64  `json:"month"`             // 有効期限（月）
	Year              int64  `json:"year"`              // 有効期限（年）
	VerificationValue string `json:"verificationValue"` // セキュリティコード
}
