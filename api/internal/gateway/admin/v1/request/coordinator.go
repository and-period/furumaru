package request

import "time"

type CreateCoordinatorRequest struct {
	Lastname          string         `json:"lastname" binding:"required,max=16"`               // 姓
	Firstname         string         `json:"firstname" binding:"required,max=16"`              // 名
	LastnameKana      string         `json:"lastnameKana" binding:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana     string         `json:"firstnameKana" binding:"required,max=32,hiragana"` // 名(かな)
	MarcheName        string         `json:"marcheName" binding:"required,max=64"`             // マルシェ名
	Username          string         `json:"username" binding:"required,max=32"`               // 表示名
	Profile           string         `json:"profile" binding:"required,max=2000"`              // 紹介文
	ProductTypeIDs    []string       `json:"productTypeIds" binding:"dive,required"`           // 取り扱い品目一覧
	ThumbnailURL      string         `json:"thumbnailUrl" binding:"omitempty,url"`             // サムネイルURL
	HeaderURL         string         `json:"headerUrl" binding:"omitempty,url"`                // ヘッダー画像URL
	PromotionVideoURL string         `json:"promotionVideoUrl" binding:"omitempty,url"`        // 紹介映像URL
	BonusVideoURL     string         `json:"bonusVideoUrl" binding:"omitempty,url"`            // 購入特典映像URL
	InstagramID       string         `json:"instagramId" binding:"omitempty,max=30"`           // Instagramアカウント
	FacebookID        string         `json:"facebookId" binding:"omitempty,max=50"`            // Facebookアカウント
	Email             string         `json:"email" binding:"required,email"`                   // メールアドレス
	PhoneNumber       string         `json:"phoneNumber" binding:"required,e164"`              // 電話番号
	PostalCode        string         `json:"postalCode" binding:"numeric,max=16"`              // 郵便番号
	PrefectureCode    int32          `json:"prefectureCode" binding:"required,min=1,max=47"`   // 都道府県
	City              string         `json:"city" binding:"required,max=32"`                   // 市区町村
	AddressLine1      string         `json:"addressLine1" binding:"required,max=64"`           // 町名・番地
	AddressLine2      string         `json:"addressLine2" binding:"omitempty,max=64"`          // ビル名・号室など
	BusinessDays      []time.Weekday `json:"businessDays" binding:"max=7,unique"`              // 営業曜日(発送可能日)
}

type UpdateCoordinatorRequest struct {
	Lastname          string `json:"lastname" binding:"required,max=16"`               // 姓
	Firstname         string `json:"firstname" binding:"required,max=16"`              // 名
	LastnameKana      string `json:"lastnameKana" binding:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana     string `json:"firstnameKana" binding:"required,max=32,hiragana"` // 名(かな)
	Username          string `json:"username" binding:"required,max=32"`               // 表示名
	Profile           string `json:"profile" binding:"required,max=2000"`              // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl" binding:"omitempty,url"`             // サムネイルURL
	HeaderURL         string `json:"headerUrl" binding:"omitempty,url"`                // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl" binding:"omitempty,url"`        // 紹介映像URL
	BonusVideoURL     string `json:"bonusVideoUrl" binding:"omitempty,url"`            // 購入特典映像URL
	InstagramID       string `json:"instagramId" binding:"omitempty,max=32"`           // Instagramアカウント
	FacebookID        string `json:"facebookId" binding:"omitempty,max=50"`            // Facebookアカウント
	PhoneNumber       string `json:"phoneNumber" binding:"required,e164"`              // 電話番号
	PostalCode        string `json:"postalCode" binding:"numeric,max=16"`              // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode" binding:"required,min=1,max=47"`   // 都道府県
	City              string `json:"city" binding:"required,max=32"`                   // 市区町村
	AddressLine1      string `json:"addressLine1" binding:"required,max=64"`           // 町名・番地
	AddressLine2      string `json:"addressLine2" binding:"omitempty,max=64"`          // ビル名・号室など
}

type UpdateCoordinatorEmailRequest struct {
	Email string `json:"email" binding:"required,email"` // メールアドレス
}
