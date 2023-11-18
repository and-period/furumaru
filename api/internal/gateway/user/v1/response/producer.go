package response

// Producer - 生産者情報
type Producer struct {
	ID                string   `json:"id"`                // 生産者ID
	CoordinatorID     string   `json:"coordinatorId"`     // 担当コーディネータID
	Username          string   `json:"username"`          // 生産者名
	Profile           string   `json:"profile"`           // 紹介文
	ThumbnailURL      string   `json:"thumbnailUrl"`      // サムネイルURL
	Thumbnails        []*Image `json:"thumbnails"`        // サムネイルURL(リサイズ済み)一覧
	HeaderURL         string   `json:"headerUrl"`         // ヘッダー画像URL
	Headers           []*Image `json:"headers"`           // ヘッダー画像URL(リサイズ済み)一覧
	PromotionVideoURL string   `json:"promotionVideoUrl"` // 紹介映像URL
	InstagramID       string   `json:"instagramId"`       // Instagramアカウント
	FacebookID        string   `json:"facebookId"`        // Facebookアカウント
	Prefecture        string   `json:"prefecture"`        // 都道府県
	City              string   `json:"city"`              // 市区町村
}
