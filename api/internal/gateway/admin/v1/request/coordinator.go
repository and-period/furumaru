package request

import "time"

type CreateCoordinatorRequest struct {
	Lastname          string         `json:"lastname,omitempty"`          // 姓
	Firstname         string         `json:"firstname,omitempty"`         // 名
	LastnameKana      string         `json:"lastnameKana,omitempty"`      // 姓(かな)
	FirstnameKana     string         `json:"firstnameKana,omitempty"`     // 名(かな)
	MarcheName        string         `json:"marcheName,omitempty"`        // マルシェ名
	Username          string         `json:"username,omitempty"`          // 表示名
	Profile           string         `json:"profile,omitempty"`           // 紹介文
	ProductTypeIDs    []string       `json:"productTypeIds,omitempty"`    // 取り扱い品目一覧
	ThumbnailURL      string         `json:"thumbnailUrl,omitempty"`      // サムネイルURL
	HeaderURL         string         `json:"headerUrl,omitempty"`         // ヘッダー画像URL
	PromotionVideoURL string         `json:"promotionVideoUrl,omitempty"` // 紹介映像URL
	BonusVideoURL     string         `json:"bonusVideoUrl,omitempty"`     // 購入特典映像URL
	InstagramID       string         `json:"instagramId,omitempty"`       // Instagramアカウント
	FacebookID        string         `json:"facebookId,omitempty"`        // Facebookアカウント
	Email             string         `json:"email,omitempty"`             // メールアドレス
	PhoneNumber       string         `json:"phoneNumber,omitempty"`       // 電話番号
	PostalCode        string         `json:"postalCode,omitempty"`        // 郵便番号
	PrefectureCode    int32          `json:"prefectureCode,omitempty"`    // 都道府県
	City              string         `json:"city,omitempty"`              // 市区町村
	AddressLine1      string         `json:"addressLine1,omitempty"`      // 町名・番地
	AddressLine2      string         `json:"addressLine2,omitempty"`      // ビル名・号室など
	BusinessDays      []time.Weekday `json:"businessDays,omitempty"`      // 営業曜日(発送可能日)
}

type UpdateCoordinatorRequest struct {
	Lastname          string `json:"lastname,omitempty"`          // 姓
	Firstname         string `json:"firstname,omitempty"`         // 名
	LastnameKana      string `json:"lastnameKana,omitempty"`      // 姓(かな)
	FirstnameKana     string `json:"firstnameKana,omitempty"`     // 名(かな)
	Username          string `json:"username,omitempty"`          // 表示名
	Profile           string `json:"profile,omitempty"`           // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl,omitempty"`      // サムネイルURL
	HeaderURL         string `json:"headerUrl,omitempty"`         // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl,omitempty"` // 紹介映像URL
	BonusVideoURL     string `json:"bonusVideoUrl,omitempty"`     // 購入特典映像URL
	InstagramID       string `json:"instagramId,omitempty"`       // Instagramアカウント
	FacebookID        string `json:"facebookId,omitempty"`        // Facebookアカウント
	PhoneNumber       string `json:"phoneNumber,omitempty"`       // 電話番号
	PostalCode        string `json:"postalCode,omitempty"`        // 郵便番号
	PrefectureCode    int32  `json:"prefectureCode,omitempty"`    // 都道府県
	City              string `json:"city,omitempty"`              // 市区町村
	AddressLine1      string `json:"addressLine1,omitempty"`      // 町名・番地
	AddressLine2      string `json:"addressLine2,omitempty"`      // ビル名・号室など
}

type UpdateCoordinatorEmailRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}
