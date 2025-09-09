package types

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

type CreateAddressRequest struct {
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
	IsDefault      bool   `json:"isDefault"`                                         // デフォルト設定フラグ
}

type UpdateAddressRequest struct {
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
	IsDefault      bool   `json:"isDefault"`                                         // デフォルト設定フラグ
}

type AddressResponse struct {
	Address *Address `json:"address"` // アドレス情報
}

type AddressesResponse struct {
	Addresses []*Address `json:"addresses"` // アドレス一覧
	Total     int64      `json:"total"`     // 合計数
}
