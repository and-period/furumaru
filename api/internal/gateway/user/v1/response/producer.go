package response

// Producer - 生産者情報
type Producer struct {
	ID                string `json:"id"`                // 生産者ID
	CoordinatorID     string `json:"coordinatorId"`     // 担当コーディネータID
	Username          string `json:"username"`          // 生産者名
	Profile           string `json:"profile"`           // 紹介文
	ThumbnailURL      string `json:"thumbnailUrl"`      // サムネイルURL
	HeaderURL         string `json:"headerUrl"`         // ヘッダー画像URL
	PromotionVideoURL string `json:"promotionVideoUrl"` // 紹介映像URL
	InstagramID       string `json:"instagramId"`       // Instagramアカウント
	FacebookID        string `json:"facebookId"`        // Facebookアカウント
	Prefecture        string `json:"prefecture"`        // 都道府県
	City              string `json:"city"`              // 市区町村
}

type ProducerResponse struct {
	Producer    *Producer         `json:"producer"`    // 生産者情報
	Lives       []*LiveSummary    `json:"lives"`       // 配信中・配信予定のマルシェ一覧
	Archives    []*ArchiveSummary `json:"archives"`    // 過去のマルシェ一覧
	Products    []*Product        `json:"products"`    // 商品一覧
	Experiences []*Experience     `json:"experiences"` // 体験一覧
}

type ProducersResponse struct {
	Producers []*Producer `json:"producers"` // 生産者一覧
	Total     int64       `json:"total"`     // 生産者合計数
}
