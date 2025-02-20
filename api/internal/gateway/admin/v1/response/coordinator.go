package response

// Coordinator - コーディネータ情報
type Coordinator struct {
	ID                string `json:"id"`                // コーディネータID
	ShopID            string `json:"shopId"`            // 店舗ID
	Status            int32  `json:"status"`            // コーディネータの状態
	Lastname          string `json:"lastname"`          // 姓
	Firstname         string `json:"firstname"`         // 名
	LastnameKana      string `json:"lastnameKana"`      // 姓(かな)
	FirstnameKana     string `json:"firstnameKana"`     // 名(かな)
	Username          string `json:"username"`          // 表示名
	Profile           string `json:"profile"`           // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl"`      // サムネイルURL
	HeaderURL         string `json:"headerUrl"`         // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl"` // 紹介映像URL
	BonusVideoURL     string `json:"bonusVideoUrl"`     // 購入特典映像URL
	InstagramID       string `json:"instagramId"`       // Instagramアカウント
	FacebookID        string `json:"facebookId"`        // Facebookアカウント
	Email             string `json:"email"`             // メールアドレス
	PhoneNumber       string `json:"phoneNumber"`       // 電話番号
	PostalCode        string `json:"postalCode"`        // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode"`    // 都道府県
	City              string `json:"city"`              // 市区町村
	AddressLine1      string `json:"addressLine1"`      // 町名・番地
	AddressLine2      string `json:"addressLine2"`      // ビル名・号室など
	ProducerTotal     int64  `json:"producerTotal"`     // 担当する生産者数
	CreatedAt         int64  `json:"createdAt"`         // 登録日時
	UpdatedAt         int64  `json:"updatedAt"`         // 更新日時
}

type CoordinatorResponse struct {
	Coordinator  *Coordinator   `json:"coordinator"`  // コーディネータ情報
	Shop         *Shop          `json:"shop"`         // 店舗情報
	ProductTypes []*ProductType `json:"productTypes"` // 品目一覧
	Password     string         `json:"password"`     // パスワード（登録時のみ）
}

type CoordinatorsResponse struct {
	Coordinators []*Coordinator `json:"coordinators"` // 生産者一覧
	Shops        []*Shop        `json:"shops"`        // 店舗一覧
	ProductTypes []*ProductType `json:"productTypes"` // 品目一覧
	Total        int64          `json:"total"`        // 合計数
}
