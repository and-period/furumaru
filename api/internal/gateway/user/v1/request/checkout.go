package request

type CheckoutProductRequest struct {
	RequestID         string              `json:"requestId" binding:"required"`            // 支払いキー(重複判定用)
	CoordinatorID     string              `json:"coordinatorId" binding:"required"`        // コーディネータID
	BoxNumber         int64               `json:"boxNumber" binding:"min=0"`               // 箱の通番（箱単位で購入する場合）
	BillingAddressID  string              `json:"billingAddressId" binding:"required"`     // 請求先住所ID
	ShippingAddressID string              `json:"shippingAddressId" binding:"required"`    // 配送先住所ID
	PromotionCode     string              `json:"promotionCode" binding:"omitempty,len=8"` // プロモーションコード
	PaymentMethod     int32               `json:"paymentMethod" binding:"required"`        // 支払い方法
	CreditCard        *CheckoutCreditCard `json:"creditCard" binding:"omitempty,dive"`     // クレジットカード決済情報
	CallbackURL       string              `json:"callbackUrl" binding:"required,http_url"` // 決済完了後のリダイレクト先URL
	Total             int64               `json:"total" binding:"min=0"`                   // 支払い合計金額（誤り検出用）
}

type CheckoutExperienceRequest struct {
	RequestID             string              `json:"requestId" binding:"required"`                 // 支払いキー(重複判定用)
	BillingAddressID      string              `json:"billingAddressId" binding:"required"`          // 請求先住所ID
	PromotionCode         string              `json:"promotionCode" binding:"omitempty,len=8"`      // プロモーションコード
	AdultCount            int64               `json:"adultCount" binding:"min=0,max=99"`            // 大人人数
	JuniorHighSchoolCount int64               `json:"juniorHighSchoolCount" binding:"min=0,max=99"` // 中学生人数
	ElementarySchoolCount int64               `json:"elementarySchoolCount" binding:"min=0,max=99"` // 小学生人数
	PreschoolCount        int64               `json:"preschoolCount" binding:"min=0,max=99"`        // 幼児人数
	SeniorCount           int64               `json:"seniorCount" binding:"min=0,max=99"`           // シニア人数
	Transportation        string              `json:"transportation" binding:"omitempty,max=256"`   // 交通手段
	RequestedDate         string              `json:"requestedDate" binding:"omitempty,date"`       // 体験希望日
	RequestedTime         string              `json:"requestedTime" binding:"omitempty,time"`       // 体験希望時間
	PaymentMethod         int32               `json:"paymentMethod" binding:"required"`             // 支払い方法
	CreditCard            *CheckoutCreditCard `json:"creditCard" binding:"omitempty,dive"`          // クレジットカード決済情報
	CallbackURL           string              `json:"callbackUrl" binding:"required,http_url"`      // 決済完了後のリダイレクト先URL
	Total                 int64               `json:"total" binding:"min=0"`                        // 支払い合計金額（誤り検出用）
}

type CheckoutCreditCard struct {
	Name              string `json:"name" binding:"required"`                         // カード名義
	Number            string `json:"number" binding:"required,credit_card"`           // カード番号
	Month             int64  `json:"month" binding:"min=1,max=12"`                    // 有効期限（月）
	Year              int64  `json:"year" binding:"min=2000,max=2100"`                // 有効期限（年）
	VerificationValue string `json:"verificationValue" binding:"min=3,max=4,numeric"` // セキュリティコード
}
