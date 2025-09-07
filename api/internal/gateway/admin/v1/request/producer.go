package request

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
