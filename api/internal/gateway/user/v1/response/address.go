package response

// Address - アドレス情報
type Address struct {
	ID             string `json:"id"`             // アドレス帳ID
	IsDefault      bool   `json:"isDefault"`      // デフォルト設定フラグ
	Lastname       string `json:"lastname"`       // 姓
	Firstname      string `json:"firstname"`      // 名
	LastnameKana   string `json:"lastnameKana"`   // 姓（かな）
	FirstnameKana  string `json:"firstnameKana"`  // 名（かな）
	PostalCode     string `json:"postalCode"`     // 郵便番号
	Prefecture     string `json:"prefecture"`     // 都道府県
	PrefectureCode int32  `json:"prefectureCode"` // 都道府県コード
	City           string `json:"city"`           // 市区町村
	AddressLine1   string `json:"addressLine1"`   // 町名・番地
	AddressLine2   string `json:"addressLine2"`   // ビル名・号室など
	PhoneNumber    string `json:"phoneNumber"`    // 電話番号
}

type AddressResponse struct {
	Address *Address `json:"address"` // アドレス情報
}

type AddressesResponse struct {
	Addresses []*Address `json:"addresses"` // アドレス一覧
	Total     int64      `json:"total"`     // 合計数
}
