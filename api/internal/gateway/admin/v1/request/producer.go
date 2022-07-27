package request

type CreateProducerRequest struct {
	Lastname      string `json:"lastname,omitempty"`      // 姓
	Firstname     string `json:"firstname,omitempty"`     // 名
	LastnameKana  string `json:"lastnameKana,omitempty"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana,omitempty"` // 名(かな)
	StoreName     string `json:"storeName,omitempty"`     // 店舗名
	ThumbnailURL  string `json:"thumbnailUrl,omitempty"`  // サムネイルURL
	HeaderURL     string `json:"headerUrl,omitempty"`     // ヘッダー画像URL
	Email         string `json:"email,omitempty"`         // メールアドレス
	PhoneNumber   string `json:"phoneNumber,omitempty"`   // 電話番号
	PostalCode    string `json:"postalCode,omitempty"`    // 郵便番号
	Prefecture    string `json:"prefecture,omitempty"`    // 都道府県
	City          string `json:"city,omitempty"`          // 市区町村
	AddressLine1  string `json:"addressLine1,omitempty"`  // 町名・番地
	AddressLine2  string `json:"addressLine2,omitempty"`  // ビル名・号室など
}

type UpdateProducerRequest struct {
	Lastname      string `json:"lastname,omitempty"`      // 姓
	Firstname     string `json:"firstname,omitempty"`     // 名
	LastnameKana  string `json:"lastnameKana,omitempty"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana,omitempty"` // 名(かな)
	StoreName     string `json:"storeName,omitempty"`     // 店舗名
	ThumbnailURL  string `json:"thumbnailUrl,omitempty"`  // サムネイルURL
	HeaderURL     string `json:"headerUrl,omitempty"`     // ヘッダー画像URL
	PhoneNumber   string `json:"phoneNumber,omitempty"`   // 電話番号
	PostalCode    string `json:"postalCode,omitempty"`    // 郵便番号
	Prefecture    string `json:"prefecture,omitempty"`    // 都道府県
	City          string `json:"city,omitempty"`          // 市区町村
	AddressLine1  string `json:"addressLine1,omitempty"`  // 町名・番地
	AddressLine2  string `json:"addressLine2,omitempty"`  // ビル名・号室など
}

type UpdateProducerEmailRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}
