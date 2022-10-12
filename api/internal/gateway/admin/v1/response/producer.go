package response

// Producer - 生産者情報
type Producer struct {
	ID            string `json:"id"`            // 生産者ID
	CoordinatorID string `json:"coordinatorId"` // 担当仲介者ID
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	StoreName     string `json:"storeName"`     // 店舗名
	ThumbnailURL  string `json:"thumbnailUrl"`  // サムネイルURL
	HeaderURL     string `json:"headerUrl"`     // ヘッダー画像URL
	Email         string `json:"email"`         // メールアドレス
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
	PostalCode    string `json:"postalCode"`    // 郵便番号
	Prefecture    string `json:"prefecture"`    // 都道府県
	City          string `json:"city"`          // 市区町村
	AddressLine1  string `json:"addressLine1"`  // 町名・番地
	AddressLine2  string `json:"addressLine2"`  // ビル名・号室など
	CreatedAt     int64  `json:"createdAt"`     // 登録日時
	UpdatedAt     int64  `json:"updatedAt"`     // 更新日時
}

type ProducerResponse struct {
	*Producer
}

type ProducersResponse struct {
	Producers []*Producer `json:"producers"` // 生産者一覧
	Total     int64       `json:"total"`     // 合計数
}
