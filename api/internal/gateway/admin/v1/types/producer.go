package types

// Producer - 生産者情報
type Producer struct {
	ID                string      `json:"id"`                // 生産者ID
	Status            AdminStatus `json:"status"`            // 生産者の状態
	Lastname          string      `json:"lastname"`          // 姓
	Firstname         string      `json:"firstname"`         // 名
	LastnameKana      string      `json:"lastnameKana"`      // 姓(かな)
	FirstnameKana     string      `json:"firstnameKana"`     // 名(かな)
	Username          string      `json:"username"`          // 生産者名
	Profile           string      `json:"profile"`           // 紹介文
	ThumbnailURL      string      `json:"thumbnailUrl"`      // サムネイルURL
	HeaderURL         string      `json:"headerUrl"`         // ヘッダー画像URL
	PromotionVideoURL string      `json:"promotionVideoUrl"` // 紹介映像URL
	BonusVideoURL     string      `json:"bonusVideoUrl"`     // 購入特典映像URL
	InstagramID       string      `json:"instagramId"`       // Instagramアカウント
	FacebookID        string      `json:"facebookId"`        // Facebookアカウント
	Email             string      `json:"email"`             // メールアドレス
	PhoneNumber       string      `json:"phoneNumber"`       // 電話番号
	PostalCode        string      `json:"postalCode"`        // 郵便番号
	PrefectureCode    int32       `json:"prefectureCode"`    // 都道府県
	City              string      `json:"city"`              // 市区町村
	AddressLine1      string      `json:"addressLine1"`      // 町名・番地
	AddressLine2      string      `json:"addressLine2"`      // ビル名・号室など
	CreatedAt         int64       `json:"createdAt"`         // 登録日時
	UpdatedAt         int64       `json:"updatedAt"`         // 更新日時
}

type CreateProducerRequest struct {
	CoordinatorID     string `json:"coordinatorId" validate:"required"`                 // 担当コーディネータ名
	Lastname          string `json:"lastname" validate:"required,max=16"`               // 姓
	Firstname         string `json:"firstname" validate:"required,max=16"`              // 名
	LastnameKana      string `json:"lastnameKana" validate:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana     string `json:"firstnameKana" validate:"required,max=32,hiragana"` // 名(かな)
	Username          string `json:"username" validate:"required,max=64"`               // 表示名
	Profile           string `json:"profile" validate:"max=2000"`                       // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl" validate:"omitempty,url"`             // サムネイルURL
	HeaderURL         string `json:"headerUrl" validate:"omitempty,url"`                // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl" validate:"omitempty,url"`        // 紹介映像URL
	BonusVideoURL     string `json:"bonusVideoUrl" validate:"omitempty,url"`            // 購入特典映像URL
	InstagramID       string `json:"instagramId" validate:"omitempty,max=30"`           // Instagramアカウント
	FacebookID        string `json:"facebookId" validate:"omitempty,max=50"`            // Facebookアカウント
	Email             string `json:"email" validate:"omitempty,email"`                  // メールアドレス
	PhoneNumber       string `json:"phoneNumber" validate:"omitempty,e164"`             // 電話番号
	PostalCode        string `json:"postalCode" validate:"omitempty,numeric,max=16"`    // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode" validate:"min=0,max=47"`            // 都道府県
	City              string `json:"city" validate:"max=32"`                            // 市区町村
	AddressLine1      string `json:"addressLine1" validate:"max=64"`                    // 町名・番地
	AddressLine2      string `json:"addressLine2" validate:"max=64"`                    // ビル名・号室など
}

type UpdateProducerRequest struct {
	Lastname          string `json:"lastname" validate:"required,max=16"`               // 姓
	Firstname         string `json:"firstname" validate:"required,max=16"`              // 名
	LastnameKana      string `json:"lastnameKana" validate:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana     string `json:"firstnameKana" validate:"required,max=32,hiragana"` // 名(かな)
	Username          string `json:"username" validate:"required,max=64"`               // 表示名
	Profile           string `json:"profile" validate:"max=2000"`                       // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl" validate:"omitempty,url"`             // サムネイルURL
	HeaderURL         string `json:"headerUrl" validate:"omitempty,url"`                // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl" validate:"omitempty,url"`        // 紹介映像URL
	BonusVideoURL     string `json:"bonusVideoUrl" validate:"omitempty,url"`            // 購入特典映像URL
	InstagramID       string `json:"instagramId" validate:"omitempty,max=30"`           // Instagramアカウント
	FacebookID        string `json:"facebookId" validate:"omitempty,max=50"`            // Facebookアカウント
	Email             string `json:"email" validate:"omitempty,email"`                  // メールアドレス
	PhoneNumber       string `json:"phoneNumber" validate:"omitempty,e164"`             // 電話番号
	PostalCode        string `json:"postalCode" validate:"omitempty,numeric,max=16"`    // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode" validate:"min=0,max=47"`            // 都道府県
	City              string `json:"city" validate:"max=32"`                            // 市区町村
	AddressLine1      string `json:"addressLine1" validate:"max=64"`                    // 町名・番地
	AddressLine2      string `json:"addressLine2" validate:"max=64"`                    // ビル名・号室など
}

type ProducerResponse struct {
	Producer     *Producer      `json:"producer"`     // 生産者情報
	Shops        []*Shop        `json:"shops"`        // 店舗情報
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
}

type ProducersResponse struct {
	Producers    []*Producer    `json:"producers"`    // 生産者一覧
	Shops        []*Shop        `json:"shops"`        // 店舗一覧
	Coordinators []*Coordinator `json:"coordinators"` // コーディネータ一覧
	Total        int64          `json:"total"`        // 合計数
}
