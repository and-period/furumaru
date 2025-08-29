package request

type GuestCheckoutProductRequest struct {
	RequestID       string                `json:"requestId"`       // 支払いキー(重複判定用)
	CoordinatorID   string                `json:"coordinatorId"`   // コーディネータID
	BoxNumber       int64                 `json:"boxNumber"`       // 箱の通番（箱単位で購入する場合）
	PromotionCode   string                `json:"promotionCode"`   // プロモーションコード
	PaymentMethod   int32                 `json:"paymentMethod"`   // 支払い方法
	CreditCard      *CheckoutCreditCard   `json:"creditCard"`      // クレジットカード決済情報
	CallbackURL     string                `json:"callbackUrl"`     // 決済完了後のリダイレクト先URL
	Total           int64                 `json:"total"`           // 支払い合計金額（誤り検出用）
	Email           string                `json:"email"`           // メールアドレス
	IsSameAddress   bool                  `json:"isSameAddress"`   // 配送先住所を請求先住所と同一にする
	BillingAddress  *GuestCheckoutAddress `json:"billingAddress"`  // 請求先住所
	ShippingAddress *GuestCheckoutAddress `json:"shippingAddress"` // 配送先住所
}

type GuestCheckoutExperienceRequest struct {
	RequestID             string                `json:"requestId"`             // 支払いキー(重複判定用)
	PromotionCode         string                `json:"promotionCode"`         // プロモーションコード
	AdultCount            int64                 `json:"adultCount"`            // 大人人数
	JuniorHighSchoolCount int64                 `json:"juniorHighSchoolCount"` // 中学生人数
	ElementarySchoolCount int64                 `json:"elementarySchoolCount"` // 小学生人数
	PreschoolCount        int64                 `json:"preschoolCount"`        // 幼児人数
	SeniorCount           int64                 `json:"seniorCount"`           // シニア人数
	Transportation        string                `json:"transportation"`        // 交通手段
	RequestedDate         string                `json:"requestedDate"`         // 体験希望日
	RequestedTime         string                `json:"requestedTime"`         // 体験希望時間
	PaymentMethod         int32                 `json:"paymentMethod"`         // 支払い方法
	CreditCard            *CheckoutCreditCard   `json:"creditCard"`            // クレジットカード決済情報
	CallbackURL           string                `json:"callbackUrl"`           // 決済完了後のリダイレクト先URL
	Total                 int64                 `json:"total"`                 // 支払い合計金額（誤り検出用）
	Email                 string                `json:"email"`                 // メールアドレス
	BillingAddress        *GuestCheckoutAddress `json:"billingAddress"`        // 請求先住所
}

type GuestCheckoutAddress struct {
	Lastname       string `json:"lastname"`       // 姓
	Firstname      string `json:"firstname"`      // 名
	LastnameKana   string `json:"lastnameKana"`   // 姓（かな）
	FirstnameKana  string `json:"firstnameKana"`  // 名（かな）
	PostalCode     string `json:"postalCode"`     // 郵便番号
	PrefectureCode int32  `json:"prefectureCode"` // 都道府県
	City           string `json:"city"`           // 市区町村
	AddressLine1   string `json:"addressLine1"`   // 町名・番地
	AddressLine2   string `json:"addressLine2"`   // ビル名・号室など
	PhoneNumber    string `json:"phoneNumber"`    // 電話番号
}
