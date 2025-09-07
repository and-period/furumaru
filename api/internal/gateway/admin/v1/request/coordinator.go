package request

import "time"

type CreateCoordinatorRequest struct {
	Lastname          string         `json:"lastname" validate:"required,max=16"`               // 姓
	Firstname         string         `json:"firstname" validate:"required,max=16"`              // 名
	LastnameKana      string         `json:"lastnameKana" validate:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana     string         `json:"firstnameKana" validate:"required,max=32,hiragana"` // 名(かな)
	MarcheName        string         `json:"marcheName" validate:"required,max=64"`             // マルシェ名
	Username          string         `json:"username" validate:"required,max=32"`               // 表示名
	Profile           string         `json:"profile" validate:"required,max=2000"`              // 紹介文
	ProductTypeIDs    []string       `json:"productTypeIds" validate:"dive,required"`           // 取り扱い品目一覧
	ThumbnailURL      string         `json:"thumbnailUrl" validate:"omitempty,url"`             // サムネイルURL
	HeaderURL         string         `json:"headerUrl" validate:"omitempty,url"`                // ヘッダー画像URL
	PromotionVideoURL string         `json:"promotionVideoUrl" validate:"omitempty,url"`        // 紹介映像URL
	BonusVideoURL     string         `json:"bonusVideoUrl" validate:"omitempty,url"`            // 購入特典映像URL
	InstagramID       string         `json:"instagramId" validate:"omitempty,max=30"`           // Instagramアカウント
	FacebookID        string         `json:"facebookId" validate:"omitempty,max=50"`            // Facebookアカウント
	Email             string         `json:"email" validate:"required,email"`                   // メールアドレス
	PhoneNumber       string         `json:"phoneNumber" validate:"required,e164"`              // 電話番号
	PostalCode        string         `json:"postalCode" validate:"numeric,max=16"`              // 郵便番号
	PrefectureCode    int32          `json:"prefectureCode" validate:"required,min=1,max=47"`   // 都道府県
	City              string         `json:"city" validate:"required,max=32"`                   // 市区町村
	AddressLine1      string         `json:"addressLine1" validate:"required,max=64"`           // 町名・番地
	AddressLine2      string         `json:"addressLine2" validate:"omitempty,max=64"`          // ビル名・号室など
	BusinessDays      []time.Weekday `json:"businessDays" validate:"max=7,unique"`              // 営業曜日(発送可能日)
}

type UpdateCoordinatorRequest struct {
	Lastname          string `json:"lastname" validate:"required,max=16"`               // 姓
	Firstname         string `json:"firstname" validate:"required,max=16"`              // 名
	LastnameKana      string `json:"lastnameKana" validate:"required,max=32,hiragana"`  // 姓(かな)
	FirstnameKana     string `json:"firstnameKana" validate:"required,max=32,hiragana"` // 名(かな)
	Username          string `json:"username" validate:"required,max=32"`               // 表示名
	Profile           string `json:"profile" validate:"required,max=2000"`              // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl" validate:"omitempty,url"`             // サムネイルURL
	HeaderURL         string `json:"headerUrl" validate:"omitempty,url"`                // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl" validate:"omitempty,url"`        // 紹介映像URL
	BonusVideoURL     string `json:"bonusVideoUrl" validate:"omitempty,url"`            // 購入特典映像URL
	InstagramID       string `json:"instagramId" validate:"omitempty,max=32"`           // Instagramアカウント
	FacebookID        string `json:"facebookId" validate:"omitempty,max=50"`            // Facebookアカウント
	PhoneNumber       string `json:"phoneNumber" validate:"required,e164"`              // 電話番号
	PostalCode        string `json:"postalCode" validate:"numeric,max=16"`              // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode" validate:"required,min=1,max=47"`   // 都道府県
	City              string `json:"city" validate:"required,max=32"`                   // 市区町村
	AddressLine1      string `json:"addressLine1" validate:"required,max=64"`           // 町名・番地
	AddressLine2      string `json:"addressLine2" validate:"omitempty,max=64"`          // ビル名・号室など
}

type UpdateCoordinatorEmailRequest struct {
	Email string `json:"email" validate:"required,email"` // メールアドレス
}
