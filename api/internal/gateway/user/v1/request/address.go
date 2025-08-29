package request

type CreateAddressRequest struct {
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
	IsDefault      bool   `json:"isDefault"`      // デフォルト設定フラグ
}

type UpdateAddressRequest struct {
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
	IsDefault      bool   `json:"isDefault"`      // デフォルト設定フラグ
}
