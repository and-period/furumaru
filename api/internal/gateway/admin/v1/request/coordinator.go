package request

import "time"

type CreateCoordinatorRequest struct {
	Lastname          string         `json:"lastname"`          // 姓
	Firstname         string         `json:"firstname"`         // 名
	LastnameKana      string         `json:"lastnameKana"`      // 姓(かな)
	FirstnameKana     string         `json:"firstnameKana"`     // 名(かな)
	MarcheName        string         `json:"marcheName"`        // マルシェ名
	Username          string         `json:"username"`          // 表示名
	Profile           string         `json:"profile"`           // 紹介文
	ProductTypeIDs    []string       `json:"productTypeIds"`    // 取り扱い品目一覧
	ThumbnailURL      string         `json:"thumbnailUrl"`      // サムネイルURL
	HeaderURL         string         `json:"headerUrl"`         // ヘッダー画像URL
	PromotionVideoURL string         `json:"promotionVideoUrl"` // 紹介映像URL
	BonusVideoURL     string         `json:"bonusVideoUrl"`     // 購入特典映像URL
	InstagramID       string         `json:"instagramId"`       // Instagramアカウント
	FacebookID        string         `json:"facebookId"`        // Facebookアカウント
	Email             string         `json:"email"`             // メールアドレス
	PhoneNumber       string         `json:"phoneNumber"`       // 電話番号
	PostalCode        string         `json:"postalCode"`        // 郵便番号
	PrefectureCode    int32          `json:"prefectureCode"`    // 都道府県
	City              string         `json:"city"`              // 市区町村
	AddressLine1      string         `json:"addressLine1"`      // 町名・番地
	AddressLine2      string         `json:"addressLine2"`      // ビル名・号室など
	BusinessDays      []time.Weekday `json:"businessDays"`      // 営業曜日(発送可能日)
}

type UpdateCoordinatorRequest struct {
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
	PhoneNumber       string `json:"phoneNumber"`       // 電話番号
	PostalCode        string `json:"postalCode"`        // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode"`    // 都道府県
	City              string `json:"city"`              // 市区町村
	AddressLine1      string `json:"addressLine1"`      // 町名・番地
	AddressLine2      string `json:"addressLine2"`      // ビル名・号室など
}

type UpdateCoordinatorEmailRequest struct {
	Email string `json:"email"` // メールアドレス
}
