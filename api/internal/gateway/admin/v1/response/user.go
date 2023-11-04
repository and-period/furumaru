package response

// User - 購入者情報
type User struct {
	ID             string `json:"id"`             // 購入者ID
	Lastname       string `json:"lastname"`       // 姓
	Firstname      string `json:"firstname"`      // 名
	LastnameKana   string `json:"lastnameKana"`   // 姓(かな)
	FirstnameKana  string `json:"firstnameKana"`  // 名(かな)
	Registered     bool   `json:"registered"`     // 会員登録フラグ
	Email          string `json:"email"`          // メールアドレス
	PhoneNumber    string `json:"phoneNumber"`    // 電話番号
	PostalCode     string `json:"postalCode"`     // 郵便番号
	PrefectureCode int32  `json:"prefectureCode"` // 都道府県
	City           string `json:"city"`           // 市区町村
	AddressLine1   string `json:"addressLine1"`   // 町名・番地
	AddressLine2   string `json:"addressLine2"`   // ビル名・号室など
	CreatedAt      int64  `json:"createdAt"`      // 登録日時
	UpdatedAt      int64  `json:"updateAt"`       // 更新日時
}

// UserSummary - 購入者一覧情報
type UserSummary struct {
	ID             string `json:"id"`             // 購入者ID
	Lastname       string `json:"lastname"`       // 姓
	Firstname      string `json:"firstname"`      // 名
	Registered     bool   `json:"registered"`     // 会員登録フラグ
	PrefectureCode int32  `json:"prefectureCode"` // 都道府県
	City           string `json:"city"`           // 市区町村
	TotalOrder     int64  `json:"totalOrder"`     // 購入回数
	TotalAmount    int64  `json:"totalAmount"`    // 購入金額
}

type UserResponse struct {
	User *User `json:"user"` // 購入者情報
}

type UsersResponse struct {
	Users []*UserSummary `json:"users"` // 購入者一覧
	Total int64          `json:"total"` // 購入者合計数
}
