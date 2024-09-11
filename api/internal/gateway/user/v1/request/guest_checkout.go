package request

type GuestCheckoutProductRequest struct {
	RequestID       string                `json:"requestId,omitempty"`       // 支払いキー(重複判定用)
	CoordinatorID   string                `json:"coordinatorId,omitempty"`   // コーディネータID
	BoxNumber       int64                 `json:"boxNumber,omitempty"`       // 箱の通番（箱単位で購入する場合）
	PromotionCode   string                `json:"promotionCode,omitempty"`   // プロモーションコード
	PaymentMethod   int32                 `json:"paymentMethod,omitempty"`   // 支払い方法
	CreditCard      *CheckoutCreditCard   `json:"creditCard,omitempty"`      // クレジットカード決済情報
	CallbackURL     string                `json:"callbackUrl,omitempty"`     // 決済完了後のリダイレクト先URL
	Total           int64                 `json:"total,omitempty"`           // 支払い合計金額（誤り検出用）
	Email           string                `json:"email,omitempty"`           // メールアドレス
	IsSameAddress   bool                  `json:"isSameAddress,omitempty"`   // 配送先住所を請求先住所と同一にする
	BillingAddress  *GuestCheckoutAddress `json:"billingAddress,omitempty"`  // 請求先住所
	ShippingAddress *GuestCheckoutAddress `json:"shippingAddress,omitempty"` // 配送先住所
}

type GuestCheckoutExperienceRequest struct {
	RequestID             string                `json:"requestId,omitempty"`             // 支払いキー(重複判定用)
	PromotionCode         string                `json:"promotionCode,omitempty"`         // プロモーションコード
	AdultCount            int64                 `json:"adultCount,omitempty"`            // 大人人数
	JuniorHighSchoolCount int64                 `json:"juniorHighSchoolCount,omitempty"` // 中学生人数
	ElementarySchoolCount int64                 `json:"elementarySchoolCount,omitempty"` // 小学生人数
	PreschoolCount        int64                 `json:"preschoolCount,omitempty"`        // 幼児人数
	SeniorCount           int64                 `json:"seniorCount,omitempty"`           // シニア人数
	Transportation        string                `json:"transportation,omitempty"`        // 交通手段
	RequestedDate         string                `json:"requestedDate,omitempty"`         // 体験希望日
	RequestedTime         string                `json:"requestedTime,omitempty"`         // 体験希望時間
	PaymentMethod         int32                 `json:"paymentMethod,omitempty"`         // 支払い方法
	CreditCard            *CheckoutCreditCard   `json:"creditCard,omitempty"`            // クレジットカード決済情報
	CallbackURL           string                `json:"callbackUrl,omitempty"`           // 決済完了後のリダイレクト先URL
	Total                 int64                 `json:"total,omitempty"`                 // 支払い合計金額（誤り検出用）
	Email                 string                `json:"email,omitempty"`                 // メールアドレス
	BillingAddress        *GuestCheckoutAddress `json:"billingAddress,omitempty"`        // 請求先住所
}

type GuestCheckoutAddress struct {
	Lastname       string `json:"lastname,omitempty"`       // 姓
	Firstname      string `json:"firstname,omitempty"`      // 名
	LastnameKana   string `json:"lastnameKana,omitempty"`   // 姓（かな）
	FirstnameKana  string `json:"firstnameKana,omitempty"`  // 名（かな）
	PostalCode     string `json:"postalCode,omitempty"`     // 郵便番号
	PrefectureCode int32  `json:"prefectureCode,omitempty"` // 都道府県
	City           string `json:"city,omitempty"`           // 市区町村
	AddressLine1   string `json:"addressLine1,omitempty"`   // 町名・番地
	AddressLine2   string `json:"addressLine2,omitempty"`   // ビル名・号室など
	PhoneNumber    string `json:"phoneNumber,omitempty"`    // 電話番号
}
