package types

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

type CheckoutResponse struct {
	URL string `json:"url"` // 遷移先URL
}

type CheckoutStateResponse struct {
	OrderID string `json:"orderId"` // 注文履歴ID
	Status  int32  `json:"status"`  // 注文ステータス
}

type PreCheckoutExperienceResponse struct {
	RequestID  string      `json:"requestId"`  // 支払い時にAPIへ送信するキー(重複判定用)
	Experience *Experience `json:"experience"` // 体験情報
	Promotion  *Promotion  `json:"promotion"`  // プロモーション情報
	SubTotal   int64       `json:"subtotal"`   // 購入金額(税込)
	Discount   int64       `json:"discount"`   // 割引金額(税込)
	Total      int64       `json:"total"`      // 合計金額(税込)
}

type GuestCheckoutProductRequest struct {
	RequestID       string                `json:"requestId" validate:"required"`                                            // 支払いキー(重複判定用)
	CoordinatorID   string                `json:"coordinatorId" validate:"required"`                                        // コーディネータID
	BoxNumber       int64                 `json:"boxNumber" validate:"min=0"`                                               // 箱の通番（箱単位で購入する場合）
	PromotionCode   string                `json:"promotionCode" validate:"omitempty,len=8"`                                 // プロモーションコード
	PaymentMethod   int32                 `json:"paymentMethod" validate:"required"`                                        // 支払い方法
	CreditCard      *CheckoutCreditCard   `json:"creditCard" validate:"omitempty,dive"`                                     // クレジットカード決済情報
	CallbackURL     string                `json:"callbackUrl" validate:"required,http_url"`                                 // 決済完了後のリダイレクト先URL
	Total           int64                 `json:"total" validate:"min=0"`                                                   // 支払い合計金額（誤り検出用）
	Email           string                `json:"email" validate:"required,email"`                                          // メールアドレス
	IsSameAddress   bool                  `json:"isSameAddress"`                                                            // 配送先住所を請求先住所と同一にする
	BillingAddress  *GuestCheckoutAddress `json:"billingAddress" validate:"required,dive"`                                  // 請求先住所
	ShippingAddress *GuestCheckoutAddress `json:"shippingAddress" validate:"required_without=IsSameAddress,omitempty,dive"` // 配送先住所
}

type GuestCheckoutExperienceRequest struct {
	RequestID             string                `json:"requestId" validate:"required"`                 // 支払いキー(重複判定用)
	PromotionCode         string                `json:"promotionCode" validate:"omitempty,len=8"`      // プロモーションコード
	AdultCount            int64                 `json:"adultCount" validate:"min=0,max=99"`            // 大人人数
	JuniorHighSchoolCount int64                 `json:"juniorHighSchoolCount" validate:"min=0,max=99"` // 中学生人数
	ElementarySchoolCount int64                 `json:"elementarySchoolCount" validate:"min=0,max=99"` // 小学生人数
	PreschoolCount        int64                 `json:"preschoolCount" validate:"min=0,max=99"`        // 幼児人数
	SeniorCount           int64                 `json:"seniorCount" validate:"min=0,max=99"`           // シニア人数
	Transportation        string                `json:"transportation" validate:"omitempty,max=256"`   // 交通手段
	RequestedDate         string                `json:"requestedDate" validate:"omitempty,date"`       // 体験希望日
	RequestedTime         string                `json:"requestedTime" validate:"omitempty,time"`       // 体験希望時間
	PaymentMethod         int32                 `json:"paymentMethod" validate:"required"`             // 支払い方法
	CreditCard            *CheckoutCreditCard   `json:"creditCard" validate:"omitempty,dive"`          // クレジットカード決済情報
	CallbackURL           string                `json:"callbackUrl" validate:"required,http_url"`      // 決済完了後のリダイレクト先URL
	Total                 int64                 `json:"total" validate:"min=0"`                        // 支払い合計金額（誤り検出用）
	Email                 string                `json:"email" validate:"required,email"`               // メールアドレス
	BillingAddress        *GuestCheckoutAddress `json:"billingAddress" validate:"required,dive"`       // 請求先住所
}

type GuestCheckoutAddress struct {
	Lastname       string `json:"lastname" validate:"required,max=16"`               // 姓
	Firstname      string `json:"firstname" validate:"required,max=16"`              // 名
	LastnameKana   string `json:"lastnameKana" validate:"required,max=32,hiragana"`  // 姓（かな）
	FirstnameKana  string `json:"firstnameKana" validate:"required,max=32,hiragana"` // 名（かな）
	PostalCode     string `json:"postalCode" validate:"required,max=16,numeric"`     // 郵便番号
	PrefectureCode int32  `json:"prefectureCode" validate:"required,min=1,max=47"`   // 都道府県
	City           string `json:"city" validate:"required,max=32"`                   // 市区町村
	AddressLine1   string `json:"addressLine1" validate:"required,max=64"`           // 町名・番地
	AddressLine2   string `json:"addressLine2" validate:"omitempty,max=64"`          // ビル名・号室など
	PhoneNumber    string `json:"phoneNumber" validate:"required,phone_number"`      // 電話番号
}

type GuestCheckoutResponse struct {
	URL string `json:"url"` // 遷移先URL
}

type GuestCheckoutStateResponse struct {
	OrderID string `json:"orderId"` // 注文履歴ID
	Status  int32  `json:"status"`  // 注文ステータス
}
