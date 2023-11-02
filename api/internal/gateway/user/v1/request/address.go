package request

type CreateAddressRequest struct {
	Lastname     string `json:"lastname,omitempty"`     // 姓
	Firstname    string `json:"firstname,omitempty"`    // 名
	PostalCode   string `json:"postalCode,omitempty"`   // 郵便番号
	Prefecture   int64  `json:"prefecture,omitempty"`   // 都道府県
	City         string `json:"city,omitempty"`         // 市区町村
	AddressLine1 string `json:"addressLine1,omitempty"` // 町名・番地
	AddressLine2 string `json:"addressLine2,omitempty"` // ビル名・号室など
	PhoneNumber  string `json:"phoneNumber,omitempty"`  // 電話番号
	IsDefault    bool   `json:"isDefault,omitempty"`    // デフォルト設定フラグ
}

type UpdateAddressRequest struct {
	Lastname     string `json:"lastname,omitempty"`     // 姓
	Firstname    string `json:"firstname,omitempty"`    // 名
	PostalCode   string `json:"postalCode,omitempty"`   // 郵便番号
	Prefecture   int64  `json:"prefecture,omitempty"`   // 都道府県
	City         string `json:"city,omitempty"`         // 市区町村
	AddressLine1 string `json:"addressLine1,omitempty"` // 町名・番地
	AddressLine2 string `json:"addressLine2,omitempty"` // ビル名・号室など
	PhoneNumber  string `json:"phoneNumber,omitempty"`  // 電話番号
	IsDefault    bool   `json:"isDefault,omitempty"`    // デフォルト設定フラグ
}
