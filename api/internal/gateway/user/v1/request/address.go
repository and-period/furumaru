package request

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
