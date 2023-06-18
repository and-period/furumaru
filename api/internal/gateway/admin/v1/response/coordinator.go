package response

import "github.com/and-period/furumaru/api/internal/user/entity"

// Coordinator - コーディネータ情報
type Coordinator struct {
	ID                string             `json:"id"`                // コーディネータID
	Status            entity.AdminStatus `json:"status"`            // コーディネータの状態
	Lastname          string             `json:"lastname"`          // 姓
	Firstname         string             `json:"firstname"`         // 名
	LastnameKana      string             `json:"lastnameKana"`      // 姓(かな)
	FirstnameKana     string             `json:"firstnameKana"`     // 名(かな)
	MarcheName        string             `json:"marcheName"`        // マルシェ名
	Username          string             `json:"username"`          // 表示名
	Profile           string             `json:"profile"`           // 紹介文
	ProductTypeIDs    []string           `json:"productTypeIds"`    // 取り扱い品目一覧
	ThumbnailURL      string             `json:"thumbnailUrl"`      // サムネイルURL
	Thumbnails        []*Image           `json:"thumbnails"`        // サムネイルURL(リサイズ済み)一覧
	HeaderURL         string             `json:"headerUrl"`         // ヘッダー画像URL
	Headers           []*Image           `json:"headers"`           // ヘッダー画像URL(リサイズ済み)一覧
	PromotionVideoURL string             `json:"promotionVideoUrl"` // 紹介映像URL
	BonusVideoURL     string             `json:"bonusVideoUrl"`     // 購入特典映像URL
	InstagramID       string             `json:"instagramId"`       // Instagramアカウント
	FacebookID        string             `json:"facebookId"`        // Facebookアカウント
	Email             string             `json:"email"`             // メールアドレス
	PhoneNumber       string             `json:"phoneNumber"`       // 電話番号
	PostalCode        string             `json:"postalCode"`        // 郵便番号
	Prefecture        string             `json:"prefecture"`        // 都道府県
	City              string             `json:"city"`              // 市区町村
	AddressLine1      string             `json:"addressLine1"`      // 町名・番地
	AddressLine2      string             `json:"addressLine2"`      // ビル名・号室など
	CreatedAt         int64              `json:"createdAt"`         // 登録日時
	UpdatedAt         int64              `json:"updatedAt"`         // 更新日時
}

type CoordinatorResponse struct {
	*Coordinator
}

type CoordinatorsResponse struct {
	Coordinators []*Coordinator `json:"coordinators"` // 生産者一覧
	Total        int64          `json:"total"`        // 合計数
}
