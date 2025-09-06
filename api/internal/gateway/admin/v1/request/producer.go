package request

type CreateProducerRequest struct {
	CoordinatorID     string `json:"coordinatorId" binding:"required"`                 // 担当コーディネータ名
	Lastname          string `json:"lastname" binding:"required,max=16"`               // 姓
	Firstname         string `json:"firstname" binding:"required,max=16"`              // 名
	LastnameKana      string `json:"lastnameKana" binding:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana     string `json:"firstnameKana" binding:"required,max=32,hiragana"` // 名(かな)
	Username          string `json:"username" binding:"required,max=64"`               // 表示名
	Profile           string `json:"profile" binding:"max=2000"`                       // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl" binding:"omitempty,url"`             // サムネイルURL
	HeaderURL         string `json:"headerUrl" binding:"omitempty,url"`                // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl" binding:"omitempty,url"`        // 紹介映像URL
	BonusVideoURL     string `json:"bonusVideoUrl" binding:"omitempty,url"`            // 購入特典映像URL
	InstagramID       string `json:"instagramId" binding:"omitempty,max=30"`           // Instagramアカウント
	FacebookID        string `json:"facebookId" binding:"omitempty,max=50"`            // Facebookアカウント
	Email             string `json:"email" binding:"omitempty,email"`                  // メールアドレス
	PhoneNumber       string `json:"phoneNumber" binding:"omitempty,e164"`             // 電話番号
	PostalCode        string `json:"postalCode" binding:"omitempty,numeric,max=16"`    // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode" binding:"min=0,max=47"`            // 都道府県
	City              string `json:"city" binding:"max=32"`                            // 市区町村
	AddressLine1      string `json:"addressLine1" binding:"max=64"`                    // 町名・番地
	AddressLine2      string `json:"addressLine2" binding:"max=64"`                    // ビル名・号室など
}

type UpdateProducerRequest struct {
	Lastname          string `json:"lastname" binding:"required,max=16"`               // 姓
	Firstname         string `json:"firstname" binding:"required,max=16"`              // 名
	LastnameKana      string `json:"lastnameKana" binding:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana     string `json:"firstnameKana" binding:"required,max=32,hiragana"` // 名(かな)
	Username          string `json:"username" binding:"required,max=64"`               // 表示名
	Profile           string `json:"profile" binding:"max=2000"`                       // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl" binding:"omitempty,url"`             // サムネイルURL
	HeaderURL         string `json:"headerUrl" binding:"omitempty,url"`                // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl" binding:"omitempty,url"`        // 紹介映像URL
	BonusVideoURL     string `json:"bonusVideoUrl" binding:"omitempty,url"`            // 購入特典映像URL
	InstagramID       string `json:"instagramId" binding:"omitempty,max=30"`           // Instagramアカウント
	FacebookID        string `json:"facebookId" binding:"omitempty,max=50"`            // Facebookアカウント
	Email             string `json:"email" binding:"omitempty,email"`                  // メールアドレス
	PhoneNumber       string `json:"phoneNumber" binding:"omitempty,e164"`             // 電話番号
	PostalCode        string `json:"postalCode" binding:"omitempty,numeric,max=16"`    // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode" binding:"min=0,max=47"`            // 都道府県
	City              string `json:"city" binding:"max=32"`                            // 市区町村
	AddressLine1      string `json:"addressLine1" binding:"max=64"`                    // 町名・番地
	AddressLine2      string `json:"addressLine2" binding:"max=64"`                    // ビル名・号室など
}
