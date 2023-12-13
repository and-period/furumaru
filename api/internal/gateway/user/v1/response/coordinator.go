package response

import "time"

// Coordinator - コーディネータ情報
type Coordinator struct {
	ID                string         `json:"id"`                // コーディネータID
	MarcheName        string         `json:"marcheName"`        // マルシェ名
	Username          string         `json:"username"`          // 表示名
	Profile           string         `json:"profile"`           // 紹介文
	ProductTypeIDs    []string       `json:"productTypeIds"`    // 取り扱い品目一覧
	BusinessDays      []time.Weekday `json:"businessDays"`      // 営業曜日(発送可能日)
	ThumbnailURL      string         `json:"thumbnailUrl"`      // サムネイルURL
	Thumbnails        []*Image       `json:"thumbnails"`        // サムネイルURL(リサイズ済み)一覧
	HeaderURL         string         `json:"headerUrl"`         // ヘッダー画像URL
	Headers           []*Image       `json:"headers"`           // ヘッダー画像URL(リサイズ済み)一覧
	PromotionVideoURL string         `json:"promotionVideoUrl"` // 紹介映像URL
	InstagramID       string         `json:"instagramId"`       // Instagramアカウント
	FacebookID        string         `json:"facebookId"`        // Facebookアカウント
	Prefecture        string         `json:"prefecture"`        // 都道府県
	City              string         `json:"city"`              // 市区町村
}

type CoordinatorResponse struct {
	Coordinator  *Coordinator      `json:"coordinator"`  // コーディネータ情報
	Lives        []*LiveSummary    `json:"lives"`        // 配信中・配信予定のマルシェ一覧
	Archives     []*ArchiveSummary `json:"archives"`     // 過去のマルシェ一覧
	ProductTypes []*ProductType    `json:"productTypes"` // 品目一覧
	Producers    []*Producer       `json:"producers"`    // 生産者一覧
	Products     []*Product        `json:"products"`     // 生産者に関連づく商品一覧
}

type CoordinatorsResponse struct {
	Coordinator *Coordinator `json:"coordinator"` // コーディネータ情報
	Total       int64        `json:"total"`       // コーディネータ合計数
}
