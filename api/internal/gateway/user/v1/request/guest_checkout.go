package request

type GuestCheckoutProductRequest struct {
	RequestID       string                `json:"requestId" binding:"required"`                                            // 支払いキー(重複判定用)
	CoordinatorID   string                `json:"coordinatorId" binding:"required"`                                        // コーディネータID
	BoxNumber       int64                 `json:"boxNumber" binding:"min=0"`                                               // 箱の通番（箱単位で購入する場合）
	PromotionCode   string                `json:"promotionCode" binding:"omitempty,len=8"`                                 // プロモーションコード
	PaymentMethod   int32                 `json:"paymentMethod" binding:"required"`                                        // 支払い方法
	CreditCard      *CheckoutCreditCard   `json:"creditCard" binding:"omitempty,dive"`                                     // クレジットカード決済情報
	CallbackURL     string                `json:"callbackUrl" binding:"required,http_url"`                                 // 決済完了後のリダイレクト先URL
	Total           int64                 `json:"total" binding:"min=0"`                                                   // 支払い合計金額（誤り検出用）
	Email           string                `json:"email" binding:"required,email"`                                          // メールアドレス
	IsSameAddress   bool                  `json:"isSameAddress"`                                                           // 配送先住所を請求先住所と同一にする
	BillingAddress  *GuestCheckoutAddress `json:"billingAddress" binding:"required,dive"`                                  // 請求先住所
	ShippingAddress *GuestCheckoutAddress `json:"shippingAddress" binding:"required_without=IsSameAddress,omitempty,dive"` // 配送先住所
}

type GuestCheckoutExperienceRequest struct {
	RequestID             string                `json:"requestId" binding:"required"`                 // 支払いキー(重複判定用)
	PromotionCode         string                `json:"promotionCode" binding:"omitempty,len=8"`      // プロモーションコード
	AdultCount            int64                 `json:"adultCount" binding:"min=0,max=99"`            // 大人人数
	JuniorHighSchoolCount int64                 `json:"juniorHighSchoolCount" binding:"min=0,max=99"` // 中学生人数
	ElementarySchoolCount int64                 `json:"elementarySchoolCount" binding:"min=0,max=99"` // 小学生人数
	PreschoolCount        int64                 `json:"preschoolCount" binding:"min=0,max=99"`        // 幼児人数
	SeniorCount           int64                 `json:"seniorCount" binding:"min=0,max=99"`           // シニア人数
	Transportation        string                `json:"transportation" binding:"omitempty,max=256"`   // 交通手段
	RequestedDate         string                `json:"requestedDate" binding:"omitempty,date"`       // 体験希望日
	RequestedTime         string                `json:"requestedTime" binding:"omitempty,time"`       // 体験希望時間
	PaymentMethod         int32                 `json:"paymentMethod" binding:"required"`             // 支払い方法
	CreditCard            *CheckoutCreditCard   `json:"creditCard" binding:"omitempty,dive"`          // クレジットカード決済情報
	CallbackURL           string                `json:"callbackUrl" binding:"required,http_url"`      // 決済完了後のリダイレクト先URL
	Total                 int64                 `json:"total" binding:"min=0"`                        // 支払い合計金額（誤り検出用）
	Email                 string                `json:"email" binding:"required,email"`               // メールアドレス
	BillingAddress        *GuestCheckoutAddress `json:"billingAddress" binding:"required,dive"`       // 請求先住所
}

type GuestCheckoutAddress struct {
	Lastname       string `json:"lastname" binding:"required,max=16"`               // 姓
	Firstname      string `json:"firstname" binding:"required,max=16"`              // 名
	LastnameKana   string `json:"lastnameKana" binding:"required,max=32,hiragana"`  // 姓（かな）
	FirstnameKana  string `json:"firstnameKana" binding:"required,max=32,hiragana"` // 名（かな）
	PostalCode     string `json:"postalCode" binding:"required,max=16,numeric"`     // 郵便番号
	PrefectureCode int32  `json:"prefectureCode" binding:"required,min=1,max=47"`   // 都道府県
	City           string `json:"city" binding:"required,max=32"`                   // 市区町村
	AddressLine1   string `json:"addressLine1" binding:"required,max=64"`           // 町名・番地
	AddressLine2   string `json:"addressLine2" binding:"omitempty,max=64"`          // ビル名・号室など
	PhoneNumber    string `json:"phoneNumber" binding:"required,phone_number"`      // 電話番号
}
