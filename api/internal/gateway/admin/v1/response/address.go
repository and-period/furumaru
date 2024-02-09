package response

// Address - 請求・配送先情報
type Address struct {
	AddressID      string `json:"addressId"`      // 住所ID
	Lastname       string `json:"lastname"`       // 氏名（姓）
	Firstname      string `json:"firstname"`      // 氏名（名）
	LastnameKana   string `json:"lastnameKana"`   // 氏名(姓:かな)
	FirstnameKana  string `json:"firstnameKana"`  // 氏名(名:かな)
	PostalCode     string `json:"postalCode"`     // 郵便番号
	Prefecture     string `json:"prefecture"`     // 都道府県
	PrefectureCode int32  `json:"prefectureCode"` // 都道府県コード
	City           string `json:"city"`           // 市区町村
	AddressLine1   string `json:"addressLine1"`   // 町名・番地
	AddressLine2   string `json:"addressLine2"`   // ビル名・号室など
	PhoneNumber    string `json:"phoneNumber"`    // 電話番号
}
