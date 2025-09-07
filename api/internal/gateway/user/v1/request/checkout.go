package request

type CheckoutProductRequest struct {
	RequestID         string              `json:"requestId" validate:"required"`            // 支払いキー(重複判定用)
	CoordinatorID     string              `json:"coordinatorId" validate:"required"`        // コーディネータID
	BoxNumber         int64               `json:"boxNumber" validate:"min=0"`               // 箱の通番（箱単位で購入する場合）
	BillingAddressID  string              `json:"billingAddressId" validate:"required"`     // 請求先住所ID
	ShippingAddressID string              `json:"shippingAddressId" validate:"required"`    // 配送先住所ID
	PromotionCode     string              `json:"promotionCode" validate:"omitempty,len=8"` // プロモーションコード
	PaymentMethod     int32               `json:"paymentMethod" validate:"required"`        // 支払い方法
	CreditCard        *CheckoutCreditCard `json:"creditCard" validate:"omitempty,dive"`     // クレジットカード決済情報
	CallbackURL       string              `json:"callbackUrl" validate:"required,http_url"` // 決済完了後のリダイレクト先URL
	Total             int64               `json:"total" validate:"min=0"`                   // 支払い合計金額（誤り検出用）
}

type CheckoutExperienceRequest struct {
	RequestID             string              `json:"requestId" validate:"required"`                 // 支払いキー(重複判定用)
	BillingAddressID      string              `json:"billingAddressId" validate:"required"`          // 請求先住所ID
	PromotionCode         string              `json:"promotionCode" validate:"omitempty,len=8"`      // プロモーションコード
	AdultCount            int64               `json:"adultCount" validate:"min=0,max=99"`            // 大人人数
	JuniorHighSchoolCount int64               `json:"juniorHighSchoolCount" validate:"min=0,max=99"` // 中学生人数
	ElementarySchoolCount int64               `json:"elementarySchoolCount" validate:"min=0,max=99"` // 小学生人数
	PreschoolCount        int64               `json:"preschoolCount" validate:"min=0,max=99"`        // 幼児人数
	SeniorCount           int64               `json:"seniorCount" validate:"min=0,max=99"`           // シニア人数
	Transportation        string              `json:"transportation" validate:"omitempty,max=256"`   // 交通手段
	RequestedDate         string              `json:"requestedDate" validate:"omitempty,date"`       // 体験希望日
	RequestedTime         string              `json:"requestedTime" validate:"omitempty,time"`       // 体験希望時間
	PaymentMethod         int32               `json:"paymentMethod" validate:"required"`             // 支払い方法
	CreditCard            *CheckoutCreditCard `json:"creditCard" validate:"omitempty,dive"`          // クレジットカード決済情報
	CallbackURL           string              `json:"callbackUrl" validate:"required,http_url"`      // 決済完了後のリダイレクト先URL
	Total                 int64               `json:"total" validate:"min=0"`                        // 支払い合計金額（誤り検出用）
}

type CheckoutCreditCard struct {
	Name              string `json:"name" validate:"required"`                         // カード名義
	Number            string `json:"number" validate:"required,credit_card"`           // カード番号
	Month             int64  `json:"month" validate:"min=1,max=12"`                    // 有効期限（月）
	Year              int64  `json:"year" validate:"min=2000,max=2100"`                // 有効期限（年）
	VerificationValue string `json:"verificationValue" validate:"min=3,max=4,numeric"` // セキュリティコード
}
