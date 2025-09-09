package types

import "time"

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
