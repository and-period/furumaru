package request

type CreateAddressRequest struct {
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
	IsDefault      bool   `json:"isDefault"`                                        // デフォルト設定フラグ
}

type UpdateAddressRequest struct {
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
	IsDefault      bool   `json:"isDefault"`                                        // デフォルト設定フラグ
}
