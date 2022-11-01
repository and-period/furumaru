package response

// Coordinator - コーディネータ情報
type Coordinator struct {
	ID               string `json:"id"`               // 生産者ID
	Lastname         string `json:"lastname"`         // 姓
	Firstname        string `json:"firstname"`        // 名
	LastnameKana     string `json:"lastnameKana"`     // 姓(かな)
	FirstnameKana    string `json:"firstnameKana"`    // 名(かな)
	CompanyName      string `json:"companyName"`      // 会社名
	StoreName        string `json:"storeName"`        // 店舗名
	ThumbnailURL     string `json:"thumbnailUrl"`     // サムネイルURL
	HeaderURL        string `json:"headerUrl"`        // ヘッダー画像URL
	TwitterAccount   string `json:"twitterAccount"`   // Twitterアカウント
	InstagramAccount string `json:"instagramAccount"` // Instagramアカウント
	FacebookAccount  string `json:"facebookAccount"`  // Facebookアカウント
	Email            string `json:"email"`            // メールアドレス
	PhoneNumber      string `json:"phoneNumber"`      // 電話番号
	PostalCode       string `json:"postalCode"`       // 郵便番号
	Prefecture       string `json:"prefecture"`       // 都道府県
	City             string `json:"city"`             // 市区町村
	AddressLine1     string `json:"addressLine1"`     // 町名・番地
	AddressLine2     string `json:"addressLine2"`     // ビル名・号室など
	CreatedAt        int64  `json:"createdAt"`        // 登録日時
	UpdatedAt        int64  `json:"updatedAt"`        // 更新日時
}

type CoordinatorResponse struct {
	*Coordinator
}

type CoordinatorsResponse struct {
	Coordinators []*Coordinator `json:"coordinators"` // 生産者一覧
	Total        int64          `json:"total"`        // 合計数
}
